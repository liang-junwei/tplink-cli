package api

import (
	"fmt"

	"github.com/ljw/tplink-cli/internal/model"
)

// ========== VPN (vpn_wan) ==========

// GetVpnConfig 获取VPN配置
func (c *Client) GetVpnConfig() (*model.VpnConfig, error) {
	req := map[string]interface{}{
		"method": "get",
		"vpn": map[string]interface{}{
			"table": "vpn_wan",
			"filter": []map[string]interface{}{
				{"interface": "WAN1"},
			},
		},
		"time_mngt": map[string]interface{}{
			"name": "_vpn_wan_1",
		},
	}

	var resp model.VpnListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("获取VPN配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("获取VPN配置失败, error_code: %d", resp.ErrorCode)
	}

	// dry-run 返回模拟数据
	if c.DryRun {
		return &model.VpnConfig{}, nil
	}

	items := make([]model.VpnConfig, 0)
	for _, m := range resp.Vpn.VpnWan {
		for key, item := range m {
			item.DotName = key
			items = append(items, item)
		}
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("未找到VPN配置")
	}
	return &items[0], nil
}

// SetVpnConfig 设置VPN配置（只发送用户指定的字段）
func (c *Client) SetVpnConfig(cfg map[string]interface{}) error {
	req := map[string]interface{}{
		"method": "set",
		"vpn": map[string]interface{}{
			"table": "vpn_wan",
			"para":  cfg,
			"filter": []map[string]interface{}{
				{"interface": "WAN1"},
			},
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("设置VPN配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("设置VPN配置失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}
