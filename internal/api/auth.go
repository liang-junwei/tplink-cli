package api

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// AuthKeys 动态获取到的认证密钥
type AuthKeys struct {
	ShortKey string
	LongKey  string
}

// securityEncode 对应前端 $.su.securityEncode(a, b, e)
// a=密码, b=shortKey, e=longKey
// 算法：对密码和短key逐字符做 XOR，结果对长key长度取模，取长key对应字符拼接
func securityEncode(a, b, e string) string {
	h := len(a)
	m := len(b)
	d := len(e)
	g := h
	if m > g {
		g = m
	}

	result := make([]byte, 0, g)
	for l := 0; l < g; l++ {
		k := 187
		u := 187
		if l >= h {
			u = int(b[l])
		} else if l >= m {
			k = int(a[l])
		} else {
			k = int(a[l])
			u = int(b[l])
		}
		idx := (k ^ u) % d
		result = append(result, e[idx])
	}
	return string(result)
}

// OrgAuthPwd 对密码进行 TP-Link 专用编码（使用硬编码的默认密钥）
// 对应前端 $.su.orgAuthPwd(password)
// 注意：密钥因固件版本而异，建议优先使用动态获取方式（FetchAuthKeys + EncodePasswordWithKeys）
func OrgAuthPwd(password string) string {
	return securityEncode(password, defaultShortKey, defaultLongKey)
}

const (
	defaultShortKey = "RDpbLfCPsJZ7fiv"
	defaultLongKey  = "yLwVl0zKqws7LgKPRQ84Mdt708T1qQ3Ha7xv3H7NyU84p21BriUWBU43odz3iP4rBL3cD02KZciXTysVXiV8ngg6vL48rPJyAUw0HurW20xqxv9aYb4M9wK1Ae0wlro510qXeU07kV57fQMc8L6aLgMLwygtc0F10a0Dg70TOoouyFhdysuRMO51yY5ZlOZZLEal1h0t9YQW0Ko7oBwmCAHoic4HYbUyVeU3sfQ1xtXcPcf1aT303wAQhv66qzW"
)

// FetchAuthKeys 从路由器动态获取 shortKey 和 longKey
// JS 文件路径固定为 /web-static/js/su/su.js
// 如果获取或解析失败，返回详细错误指引用户回退
func FetchAuthKeys(serverURL string) (*AuthKeys, error) {
	// 构造 JS 文件 URL（固定路径）
	jsURL := strings.TrimSuffix(serverURL, "/") + "/web-static/js/su/su.js"

	// 获取 JS 文件内容
	jsContent, err := fetchJSContent(jsURL)
	if err != nil {
		return nil, fmt.Errorf(
			"动态获取认证密钥失败（获取 JS 文件）\n"+
				"  JS 路径: %s\n"+
				"  错误: %v\n"+
				"  提示: 请检查路由器 URL 是否正确，或关闭动态认证：\n"+
				"        tplink context update <name> --dynamic-auth=false", jsURL, err,
		)
	}

	// 从 JS 中提取 shortKey 和 longKey
	keys, err := extractAuthKeys(jsContent, jsURL)
	if err != nil {
		return nil, fmt.Errorf(
			"动态获取认证密钥失败（解析 JS 内容）\n"+
				"  JS 路径: %s\n"+
				"  %v\n"+
				"  提示: 路由器固件版本可能已变更，JS 文件格式不匹配，请关闭动态认证：\n"+
				"        tplink context update <name> --dynamic-auth=false", jsURL, err,
		)
	}

	return keys, nil
}

// fetchJSContent 获取 JS 文件内容
func fetchJSContent(jsURL string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(jsURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// extractAuthKeys 从 JS 内容中提取 shortKey 和 longKey
// 匹配模式: securityEncode(a, "短key", "长key")
func extractAuthKeys(jsContent string, jsURL string) (*AuthKeys, error) {
	// 先检查 JS 文件是否包含关键函数名
	hasOrgAuthPwd := strings.Contains(jsContent, "orgAuthPwd")
	hasSecurityEncode := strings.Contains(jsContent, "securityEncode")

	// 诊断信息
	diag := fmt.Sprintf("JS 文件大小: %d 字节", len(jsContent))
	if hasOrgAuthPwd {
		diag += ", 包含 orgAuthPwd: 是"
	}
	if hasSecurityEncode {
		diag += ", 包含 securityEncode: 是"
	}

	// 模式1: securityEncode(a, "短key", "长key")
	re := regexp.MustCompile(`securityEncode\s*\(\s*[a-z]+\s*,\s*"([^"]*)"\s*,\s*"([^"]*)"\s*\)`)
	matches := re.FindStringSubmatch(jsContent)
	if matches != nil && matches[1] != "" && matches[2] != "" {
		return &AuthKeys{
			ShortKey: matches[1],
			LongKey:  matches[2],
		}, nil
	}

	// 模式2: orgAuthPwd = function(a) { ... securityEncode(..., "短key", "长key") }
	re2 := regexp.MustCompile(`orgAuthPwd\s*=\s*function\s*\([^)]*\)\s*\{[^}]*securityEncode\s*\([^,]+,\s*"([^"]*)"\s*,\s*"([^"]*)"\s*\)`)
	matches2 := re2.FindStringSubmatch(jsContent)
	if matches2 != nil && matches2[1] != "" && matches2[2] != "" {
		return &AuthKeys{
			ShortKey: matches2[1],
			LongKey:  matches2[2],
		}, nil
	}

	// 构建详细错误信息
	errMsg := "JS 文件中未找到 securityEncode 调用的密钥参数"
	if !hasOrgAuthPwd && !hasSecurityEncode {
		errMsg += "\n  原因: JS 文件中不包含 orgAuthPwd 或 securityEncode 关键字"
		errMsg += "\n  说明: 这可能不是 TP-Link 路由器的 su.js 文件，或固件版本差异较大"
	} else if hasOrgAuthPwd || hasSecurityEncode {
		errMsg += "\n  原因: JS 文件中包含关键函数，但正则表达式未能提取密钥"
		errMsg += "\n  说明: 函数格式可能与预期不符，需要更新正则匹配规则"
		// 尝试提取相关代码段供调试
		if hasOrgAuthPwd {
			idx := strings.Index(jsContent, "orgAuthPwd")
			if idx != -1 {
				start := idx - 20
				if start < 0 {
					start = 0
				}
				end := idx + 300
				if end > len(jsContent) {
					end = len(jsContent)
				}
				errMsg += fmt.Sprintf("\n  相关代码段: %s", jsContent[start:end])
			}
		}
	}
	errMsg += "\n  诊断信息: " + diag
	errMsg += fmt.Sprintf("\n  JS 文件 URL: %s", jsURL)

	return nil, fmt.Errorf("%s", errMsg)
}

// FetchAndEncodePassword 从路由器获取密钥并编码密码
// 返回编码后的密码，错误时返回详细提示
func FetchAndEncodePassword(serverURL, password string) (string, error) {
	keys, err := FetchAuthKeys(serverURL)
	if err != nil {
		return "", err
	}
	return EncodePasswordWithKeys(password, keys), nil
}

// EncodePasswordWithKeys 使用动态获取的 key 对密码进行编码
func EncodePasswordWithKeys(password string, keys *AuthKeys) string {
	return securityEncode(password, keys.ShortKey, keys.LongKey)
}
