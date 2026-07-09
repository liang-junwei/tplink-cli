package api

import (
	"fmt"

	"github.com/ljw/tplink-cli/internal/model"
)

// ========== 设备信息 ==========

// GetDeviceInfo 获取设备信息
func (c *Client) GetDeviceInfo() (*model.DeviceInfo, error) {
	req := map[string]interface{}{
		"method": "get",
		"device_info": map[string]interface{}{
			"name": "info",
		},
	}

	var resp model.DeviceInfoResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("获取设备信息失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("获取设备信息失败, error_code: %d", resp.ErrorCode)
	}

	// dry-run 返回模拟数据
	if c.DryRun {
		return &model.DeviceInfo{}, nil
	}
	return &resp.DeviceInfo.Info, nil
}

// ========== 重启 ==========

// RebootSystem 重启设备
func (c *Client) RebootSystem() error {
	req := map[string]interface{}{
		"method": "do",
		"system": map[string]interface{}{
			"reboot": nil,
		},
	}

	var resp model.SystemRebootResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("重启设备失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("重启设备失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}
