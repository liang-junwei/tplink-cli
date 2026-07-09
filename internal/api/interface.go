package api

import (
	"fmt"

	"github.com/ljw/tplink-cli/internal/model"
)

// ========== 接口模式 (ifmode) ==========

// GetIfMode 查看接口模式
func (c *Client) GetIfMode() (*model.IfModeData, error) {
	req := model.IfModeGetRequest{
		Method:  "get",
		Network: map[string]string{"name": "if_mode"},
	}

	var resp model.IfModeGetResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("查询接口模式失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	return &resp.Network.IfMode, nil
}

// SetIfMode 设置接口模式
func (c *Client) SetIfMode(wanMode string) error {
	req := model.IfModeSetRequest{
		Method: "set",
		Network: map[string]*model.IfModeSetData{
			"if_mode": {WanMode: wanMode},
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("设置接口模式失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("设置失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// ========== IPv6桥模式 (brv6mode) ==========

// GetBridgeV6 查看IPv6桥模式
func (c *Client) GetBridgeV6() (*model.BridgeV6Data, error) {
	req := model.BridgeV6GetRequest{
		Method:  "get",
		Network: map[string]string{"name": "bridge_v6"},
	}

	var resp model.BridgeV6GetResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("查询IPv6桥模式失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	return &resp.Network.BridgeV6, nil
}

// SetBridgeV6 设置IPv6桥模式
func (c *Client) SetBridgeV6(data *model.BridgeV6Data) error {
	req := model.BridgeV6SetRequest{
		Method: "set",
		Network: map[string]*model.BridgeV6Data{
			"bridge_v6": data,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("设置IPv6桥模式失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("设置失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// ========== 有线接口状态 (port) ==========

// GetPorts 查看有线接口状态
func (c *Client) GetPorts() ([]model.PortItem, error) {
	req := model.PortGetRequest{
		Method: "get",
		Port:   map[string]string{"table": "port"},
	}

	var resp model.PortGetResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("查询有线接口状态失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	return resp.Port.Ports, nil
}

// ========== WAN 配置 (wan) ==========

// GetWanInfo 查看WAN配置
func (c *Client) GetWanInfo(baseName string) ([]model.WanIfItem, error) {
	req := model.WanListRequest{
		Method: "get",
		Network: model.WanQuery{
			Table: "if_info",
		},
	}

	// 指定 base_name 过滤
	if baseName != "" {
		req.Network.Filter = []model.WanIfFilter{
			{BaseName: []string{baseName}},
		}
	}

	var resp model.WanListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("查询WAN配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	return resp.Network.IfInfo, nil
}

// SetWan 设置WAN配置(合并现有配置后发送)
func (c *Client) SetWan(ifName string, para model.WanIfData) error {
	req := model.WanSetRequest{
		Method: "set",
		Network: model.WanSetParams{
			Table: "if_info",
			Para:  para,
			Filter: []model.WanIfFilter{
				{IfName: []string{ifName}},
			},
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("设置WAN配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("设置失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// ========== LAN 接口 (lan) ==========

// GetLanInfo 查看LAN接口信息
func (c *Client) GetLanInfo() (*model.LanData, error) {
	req := model.LanGetRequest{
		Method:  "get",
		Network: map[string]string{"name": "lan"},
	}

	var resp model.LanGetResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("查询LAN接口信息失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	return &resp.Network.Lan, nil
}

// SetLan 修改LAN接口信息
func (c *Client) SetLan(data *model.LanData) error {
	req := model.LanSetRequest{
		Method: "set",
		Network: map[string]*model.LanData{
			"lan": data,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("修改LAN接口信息失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("修改失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// ========== DHCP 管理 ==========

// GetDhcpConfig 查看DHCP配置
func (c *Client) GetDhcpConfig() (*model.DhcpConfigData, error) {
	req := model.DhcpConfigGetRequest{
		Method: "get",
		Dhcpd:  map[string]string{"name": "lan"},
	}

	var resp model.DhcpConfigGetResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("查询DHCP配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	return &resp.Dhcpd.Lan, nil
}

// SetDhcpConfig 修改DHCP配置
func (c *Client) SetDhcpConfig(data *model.DhcpConfigData) error {
	req := model.DhcpConfigSetRequest{
		Method: "set",
		Dhcpd: map[string]*model.DhcpConfigData{
			"lan": data,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("修改DHCP配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("修改失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// GetDhcpClients 查看DHCP客户端列表
// start/end 对应 API 的 para.start/para.end（含首尾，如 0~499 为第1页）
func (c *Client) GetDhcpClients(start, end int) ([]model.DhcpClientItem, int, error) {
	req := model.DhcpClientListRequest{
		Method: "get",
		Dhcpd: model.DhcpTableReq{
			Table: "dhcp_clients",
			Para:  model.DhcpRangeParams{Start: start, End: end},
		},
	}

	var resp model.DhcpClientListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("查询DHCP客户端失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	total := 0
	for _, v := range resp.Dhcpd.Count {
		total = v
		break
	}
	return resp.Dhcpd.DhcpClients, total, nil
}

// GetDhcpStatic 查看静态地址分配列表
// start/end 对应 API 的 para.start/para.end（含首尾，如 0~99 为第1页）
func (c *Client) GetDhcpStatic(start, end int) ([]model.DhcpStaticItem, int, error) {
	req := model.DhcpStaticListRequest{
		Method: "get",
		Dhcpd: model.DhcpTableReq{
			Table: "dhcp_static",
			Para:  model.DhcpRangeParams{Start: start, End: end},
		},
	}

	var resp model.DhcpStaticListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("查询静态地址分配失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	total := 0
	for _, v := range resp.Dhcpd.Count {
		total = v
		break
	}
	return resp.Dhcpd.DhcpStatic, total, nil
}

// AddDhcpStatic 添加静态地址分配
func (c *Client) AddDhcpStatic(para model.DhcpStaticPara) (string, error) {
	req := model.DhcpStaticAddRequest{
		Method: "add",
		Dhcpd: model.DhcpStaticAddIn{
			Table: "dhcp_static",
			Para:  para,
		},
	}

	var resp model.DhcpStaticAddResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return "", fmt.Errorf("添加静态地址分配失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("添加失败, error_code: %d", resp.ErrorCode)
	}
	if len(resp.Dhcpd.Name) > 0 {
		return resp.Dhcpd.Name[0], nil
	}
	return "", nil
}

// DelDhcpStatic 删除静态地址分配
func (c *Client) DelDhcpStatic(staticID string) error {
	req := model.DhcpStaticDelRequest{
		Method: "delete",
		Dhcpd: model.DhcpStaticDelIn{
			Table: "dhcp_static",
			Filter: []model.DhcpStaticIDFilter{
				{DhcpStaticID: staticID},
			},
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("删除静态地址分配失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}
