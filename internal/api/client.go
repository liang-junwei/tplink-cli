package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ljw/tplink-cli/internal/config"
	"github.com/ljw/tplink-cli/internal/model"
)

// Client TP-Link 路由器 API 客户端
type Client struct {
	Server     *config.ServerConfig
	ServerName string
	DryRun     bool
	httpClient *http.Client
}

// NewClient 从 ServerConfig 创建客户端
func NewClient(srv *config.ServerConfig, serverName string) *Client {
	return &Client{
		Server:     srv,
		ServerName: serverName,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// EnsureAuth 确保已认证（优先使用缓存的 stok，无效则重新登录）
func (c *Client) EnsureAuth() error {
	if c.Server.IsStokValid() {
		return nil
	}

	return c.relogin()
}

// relogin 强制重新登录并持久化 stok
func (c *Client) relogin() error {
	if err := c.Login(); err != nil {
		return err
	}

	// 持久化 stok 到配置
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	if err := cfg.UpdateStok(c.ServerName, c.Server.Stok); err != nil {
		// stok 持久化失败不影响本次执行
		fmt.Fprintf(os.Stderr, "警告: 持久化 stok 失败: %v\n", err)
	}

	return nil
}

// stokExpiredError 表示 stok 过期需要重新登录的 error_code
const stokExpiredCode = -40401

// Login 登录获取 stok
// 如果 Server.DynamicAuth 为 true，优先使用预编码的密码（EncodedPassword）
// 如果预编码密码为空，则动态获取认证密钥后对密码编码
func (c *Client) Login() error {
	password := c.Server.Password

	if c.Server.DynamicAuth {
		if c.Server.EncodedPassword != "" {
			// 使用预编码的密码
			password = c.Server.EncodedPassword
		} else {
			// 预编码密码为空，尝试动态获取（兼容旧配置或首次登录）
			fmt.Fprintf(os.Stderr,
				"提示: 动态认证已启用但未找到预编码密码，正在从路由器获取密钥...\n")
			keys, err := FetchAuthKeys(c.Server.ServerURL)
			if err != nil {
				return fmt.Errorf(
					"动态认证失败: %w\n"+
						"  提示: 请重新开启动态认证：\n"+
						"        tplink context update %s --dynamic-auth=true",
					err, c.ServerName,
				)
			}
			password = EncodePasswordWithKeys(password, keys)
		}
	}

	reqBody := model.LoginRequest{
		Method: "do",
		Login: model.LoginParams{
			Username: c.Server.Username,
			Password: password,
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化登录请求失败: %w", err)
	}

	if c.DryRun {
		c.printDryRun("POST", "/", body)
		return nil
	}

	url := fmt.Sprintf("%s/", c.Server.ServerURL)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("登录请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取登录响应失败: %w", err)
	}

	var loginResp model.LoginResponse
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return fmt.Errorf("解析登录响应失败: %w", err)
	}

	if loginResp.ErrorCode != 0 {
		return fmt.Errorf("登录失败, error_code: %d", loginResp.ErrorCode)
	}

	if loginResp.Stok == "" {
		return fmt.Errorf("登录成功但未获取到 stok")
	}

	c.Server.Stok = loginResp.Stok
	return nil
}

// DoRequest 发送带 stok 的通用 API 请求
// 当响应 error_code 为 -40401（stok 过期）时自动重新登录并重试一次
func (c *Client) DoRequest(method, path string, reqBody interface{}, result interface{}) error {
	return c.doRequest(method, path, reqBody, result, true)
}

func (c *Client) doRequest(method, path string, reqBody interface{}, result interface{}, canRetry bool) error {
	if err := c.EnsureAuth(); err != nil {
		return err
	}

	var body []byte
	var err error
	if reqBody != nil {
		body, err = json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("序列化请求失败: %w", err)
		}
	}

	if c.DryRun {
		c.printDryRun(method, path, body)
		return nil
	}

	url := fmt.Sprintf("%s/stok=%s/ds", c.Server.ServerURL, c.Server.Stok)
	if path != "" && path != "/" {
		url = fmt.Sprintf("%s/stok=%s/ds%s", c.Server.ServerURL, c.Server.Stok, path)
	}

	var resp *http.Response
	switch strings.ToUpper(method) {
	case "GET":
		resp, err = c.httpClient.Get(url)
	case "POST":
		resp, err = c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	case "PUT":
		req, _ := http.NewRequest("PUT", url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err = c.httpClient.Do(req)
	case "DELETE":
		req, _ := http.NewRequest("DELETE", url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err = c.httpClient.Do(req)
	default:
		return fmt.Errorf("不支持的 HTTP 方法: %s", method)
	}

	if err != nil {
		return fmt.Errorf("API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查是否 stok 过期，自动重试一次
	if canRetry {
		var ec struct {
			ErrorCode int `json:"error_code"`
		}
		if json.Unmarshal(respBody, &ec) == nil && ec.ErrorCode == stokExpiredCode {
			if err := c.relogin(); err != nil {
				return fmt.Errorf("stok 过期, 重新登录失败: %w", err)
			}
			return c.doRequest(method, path, reqBody, result, false)
		}
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("解析响应失败: %w", err)
		}
	}

	return nil
}

// DoRawRequest 发送原始 API 请求（为 api 命令设计）
// method: GET/POST/PUT/DELETE, path: 不含 host 的 API 路径，body: 原始 JSON 字节
// 当响应 error_code 为 -40401（stok 过期）时自动重新登录并重试一次
func (c *Client) DoRawRequest(method, path string, body []byte) ([]byte, error) {
	return c.doRawRequest(method, path, body, true)
}

func (c *Client) doRawRequest(method, path string, body []byte, canRetry bool) ([]byte, error) {
	if err := c.EnsureAuth(); err != nil {
		return nil, err
	}

	if c.DryRun {
		c.printDryRun(method, path, body)
		return nil, nil
	}

	url := fmt.Sprintf("%s/stok=%s/ds%s", c.Server.ServerURL, c.Server.Stok, path)

	var resp *http.Response
	var err error

	switch strings.ToUpper(method) {
	case "GET":
		resp, err = c.httpClient.Get(url)
	case "DELETE":
		req, _ := http.NewRequest("DELETE", url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err = c.httpClient.Do(req)
	default:
		// POST, PUT 等
		req, _ := http.NewRequest(strings.ToUpper(method), url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err = c.httpClient.Do(req)
	}

	if err != nil {
		return nil, fmt.Errorf("API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查是否 stok 过期，自动重试一次
	if canRetry {
		var ec struct {
			ErrorCode int `json:"error_code"`
		}
		if json.Unmarshal(respBody, &ec) == nil && ec.ErrorCode == stokExpiredCode {
			if err := c.relogin(); err != nil {
				return nil, fmt.Errorf("stok 过期, 重新登录失败: %w", err)
			}
			return c.doRawRequest(method, path, body, false)
		}
	}

	return respBody, nil
}

// printDryRun 输出 dry-run 请求信息
func (c *Client) printDryRun(method, path string, body []byte) {
	url := c.Server.ServerURL
	if !strings.Contains(path, "/stok=") {
		url = fmt.Sprintf("%s/stok=<stok>/ds%s", url, path)
	} else {
		url = fmt.Sprintf("%s%s", url, path)
	}

	fmt.Printf("[DRY RUN] %s %s\n", strings.ToUpper(method), url)
	if len(body) > 0 {
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, body, "", "  "); err == nil {
			fmt.Printf("[DRY RUN] Body:\n%s\n", pretty.String())
		} else {
			fmt.Printf("[DRY RUN] Body: %s\n", string(body))
		}
	}
}
