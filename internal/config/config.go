package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ServerConfig 单个 server 配置
type ServerConfig struct {
	ServerURL       string `json:"server_url"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	DynamicAuth     bool   `json:"dynamic_auth,omitempty"`
	EncodedPassword string `json:"encoded_password,omitempty"` // 预编码的密码（dynamic_auth=true 时生成）
	Stok            string `json:"stok,omitempty"`
	StokExpiresAt   int64  `json:"stok_expires_at,omitempty"`
}

// AppConfig 根配置结构
type AppConfig struct {
	Current string                  `json:"current"`
	Servers map[string]ServerConfig `json:"servers"`
}

// customPath 自定义配置文件路径（空则用默认）
var customPath string

// SetConfigPath 设置自定义配置文件路径
func SetConfigPath(path string) {
	customPath = path
}

// configPath 配置文件路径
func configPath() (string, error) {
	if customPath != "" {
		return customPath, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".tplink.json"), nil
}

// Load 加载配置文件
func Load() (*AppConfig, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &AppConfig{Servers: make(map[string]ServerConfig)}, nil
		}
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	if cfg.Servers == nil {
		cfg.Servers = make(map[string]ServerConfig)
	}

	return &cfg, nil
}

// Save 保存配置到文件
func Save(cfg *AppConfig) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// GetServer 获取当前或指定 server 配置
// 优先使用 serverName，为空则使用 Current
func (cfg *AppConfig) GetServer(serverName string) (*ServerConfig, string, error) {
	name := serverName
	if name == "" {
		name = cfg.Current
	}
	if name == "" {
		return nil, "", fmt.Errorf("未指定 server，请使用 --server 或先设置默认 server (context use <name>)")
	}

	srv, ok := cfg.Servers[name]
	if !ok {
		return nil, "", fmt.Errorf("server '%s' 不存在", name)
	}

	return &srv, name, nil
}

// AddServer 添加 server
func (cfg *AppConfig) AddServer(name string, srv ServerConfig, setDefault bool) error {
	if _, exists := cfg.Servers[name]; exists {
		return fmt.Errorf("server '%s' 已存在", name)
	}

	cfg.Servers[name] = srv

	if setDefault || cfg.Current == "" {
		cfg.Current = name
	}

	return Save(cfg)
}

// RemoveServer 删除 server
func (cfg *AppConfig) RemoveServer(name string) error {
	if _, exists := cfg.Servers[name]; !exists {
		return fmt.Errorf("server '%s' 不存在", name)
	}

	delete(cfg.Servers, name)

	if cfg.Current == name {
		cfg.Current = ""
	}

	return Save(cfg)
}

// SetCurrent 切换默认 server
func (cfg *AppConfig) SetCurrent(name string) error {
	if _, exists := cfg.Servers[name]; !exists {
		return fmt.Errorf("server '%s' 不存在", name)
	}

	cfg.Current = name
	return Save(cfg)
}

// UpdateServer 更新 server 配置
func (cfg *AppConfig) UpdateServer(name string, srv ServerConfig) error {
	if _, exists := cfg.Servers[name]; !exists {
		return fmt.Errorf("server '%s' 不存在", name)
	}

	cfg.Servers[name] = srv
	return Save(cfg)
}

// UpdateStok 更新指定 server 的 stok
func (cfg *AppConfig) UpdateStok(serverName string, stok string) error {
	srv, ok := cfg.Servers[serverName]
	if !ok {
		return fmt.Errorf("server '%s' 不存在", serverName)
	}

	srv.Stok = stok
	srv.StokExpiresAt = time.Now().Add(30 * time.Minute).Unix()
	cfg.Servers[serverName] = srv
	return Save(cfg)
}

// IsStokValid 检查 stok 是否仍然有效
func (srv *ServerConfig) IsStokValid() bool {
	if srv.Stok == "" {
		return false
	}
	return time.Now().Unix() < srv.StokExpiresAt
}
