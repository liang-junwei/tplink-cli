package api

import (
	"fmt"

	"github.com/ljw/tplink-cli/internal/model"
)

// ========== wireless config ==========

// GetWirelessConfig 查看无线配置（2.4G和5G）
func (c *Client) GetWirelessConfig() (*model.WirelessConfigResult, error) {
	req := model.WirelessConfigRequest{
		Method: "get",
		Wireless: model.WirelessConfigName{
			Name: []string{"wlan_host_2g", "wlan_host_5g"},
		},
	}

	var resp model.WirelessConfigResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("查询无线配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	return &resp.Wireless, nil
}

// SetWirelessConfig 设置无线配置（2.4G或5G）
// band: "wlan_host_2g" 或 "wlan_host_5g", cfg: 需要设置的字段
func (c *Client) SetWirelessConfig(band string, cfg map[string]interface{}) error {
	req := model.WirelessSetRequest{
		Method:   "set",
		Wireless: map[string]interface{}{band: cfg},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("设置无线配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("设置失败, error_code: %d", resp.ErrorCode)
	}

	return nil
}

// ========== guest network ==========

// GetGuestNetwork 查看访客网络配置
func (c *Client) GetGuestNetwork() (*model.GuestNetworkResult, error) {
	req := model.GuestNetworkRequest{
		Method: "get",
		GuestNetwork: model.GuestNetworkName{
			Name: "guest_2g",
		},
	}

	var resp model.GuestNetworkResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("查询访客网络失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	return &resp.GuestNetwork, nil
}

// SetGuestNetwork 设置访客网络配置
func (c *Client) SetGuestNetwork(cfg map[string]interface{}) error {
	req := model.GuestNetworkSetRequest{
		Method:       "set",
		GuestNetwork: map[string]interface{}{"guest_2g": cfg},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("设置访客网络失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("设置失败, error_code: %d", resp.ErrorCode)
	}

	return nil
}

// ========== wlan access config ==========

// GetWlanAccessConfig 查看MAC地址过滤配置
func (c *Client) GetWlanAccessConfig() (*model.WlanAccessConfig, error) {
	req := model.WlanAccessConfigRequest{
		Method: "get",
		WlanAccess: model.WlanAccessNameReq{
			Name: "config",
		},
	}

	var resp model.WlanAccessConfigResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("查询MAC过滤配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	return &resp.WlanAccess.Config, nil
}

// SetWlanAccessConfig 设置MAC地址过滤配置
func (c *Client) SetWlanAccessConfig(cfg map[string]interface{}) error {
	req := model.WlanAccessConfigSetRequest{
		Method:     "set",
		WlanAccess: map[string]interface{}{"config": cfg},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("设置MAC过滤配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("设置失败, error_code: %d", resp.ErrorCode)
	}

	return nil
}

// ========== wlan access white list ==========

// GetWlanAccessWhiteList 查看白名单
func (c *Client) GetWlanAccessWhiteList() ([]model.WlanAccessListItem, int, error) {
	req := model.WlanAccessListRequest{
		Method: "get",
		WlanAccess: model.WlanAccessTableReq{
			Table: "white_list",
		},
	}

	var resp model.WlanAccessListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("查询白名单失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	total := 0
	for _, v := range resp.WlanAccess.Count {
		total = v
		break
	}
	return resp.WlanAccess.WhiteList, total, nil
}

// AddWlanAccessWhite 添加白名单条目
func (c *Client) AddWlanAccessWhite(mac, name string) (string, error) {
	req := model.WlanAccessAddRequest{
		Method: "add",
		WlanAccess: model.WlanAccessAddIn{
			Table: "white_list",
			Para:  map[string]interface{}{
				"mac":  mac,
				"name": name,
			},
		},
	}

	var resp model.WlanAccessAddResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return "", fmt.Errorf("添加白名单失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("添加失败, error_code: %d", resp.ErrorCode)
	}
	if len(resp.WlanAccess.Name) > 0 {
		return resp.WlanAccess.Name[0], nil
	}
	return "", nil
}

// DelWlanAccessWhite 删除白名单条目（用完整key名如 white_list_1782875167）
func (c *Client) DelWlanAccessWhite(listName string) error {
	req := model.WlanAccessDelRequest{
		Method: "delete",
		WlanAccess: model.WlanAccessDelIn{
			Name: listName,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("删除白名单失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除失败, error_code: %d", resp.ErrorCode)
	}

	return nil
}

// ========== wlan access black list ==========

// GetWlanAccessBlackList 查看黑名单
func (c *Client) GetWlanAccessBlackList() ([]model.WlanAccessListItem, int, error) {
	req := model.WlanAccessListRequest{
		Method: "get",
		WlanAccess: model.WlanAccessTableReq{
			Table: "black_list",
		},
	}

	var resp model.WlanAccessListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("查询黑名单失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	total := 0
	for _, v := range resp.WlanAccess.Count {
		total = v
		break
	}
	return resp.WlanAccess.BlackList, total, nil
}

// AddWlanAccessBlack 添加黑名单条目
func (c *Client) AddWlanAccessBlack(mac, name string) (string, error) {
	req := model.WlanAccessAddRequest{
		Method: "add",
		WlanAccess: model.WlanAccessAddIn{
			Table: "black_list",
			Para:  map[string]interface{}{
				"mac":  mac,
				"name": name,
			},
		},
	}

	var resp model.WlanAccessAddResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return "", fmt.Errorf("添加黑名单失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("添加失败, error_code: %d", resp.ErrorCode)
	}
	if len(resp.WlanAccess.Name) > 0 {
		return resp.WlanAccess.Name[0], nil
	}
	return "", nil
}

// DelWlanAccessBlack 删除黑名单条目（用完整key名如 black_list_1782875167）
func (c *Client) DelWlanAccessBlack(listName string) error {
	req := model.WlanAccessDelRequest{
		Method: "delete",
		WlanAccess: model.WlanAccessDelIn{
			Name: listName,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("删除黑名单失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除失败, error_code: %d", resp.ErrorCode)
	}

	return nil
}

// ========== wlan service list ==========

// GetWlanServList 查看无线Wlan服务列表
func (c *Client) GetWlanServList() ([]model.WlanServItem, int, error) {
	req := model.WlanServRequest{
		Method: "get",
		Wireless: model.WlanServTable{
			Table: "wlan_serv",
		},
	}

	var resp model.WlanServResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("查询Wlan服务列表失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	total := 0
	for _, v := range resp.Wireless.Count {
		total = v
		break
	}
	return resp.Wireless.WlanServ, total, nil
}

// ========== wireless client list ==========

// GetWirelessClients 查看无线客户端列表，可选过滤
func (c *Client) GetWirelessClients(radioID, servID string) ([]model.WirelessClientItem, int, error) {
	req := model.WirelessClientRequest{
		Method: "get",
		Wireless: model.WirelessClientTableReq{
			Table: "sta_list",
		},
	}
	if radioID != "" || servID != "" {
		req.Wireless.Filter = &model.WirelessClientFilter{
			RadioID: radioID,
			ServID:  servID,
		}
	}

	var resp model.WirelessClientResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("查询无线客户端失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	total := 0
	for _, v := range resp.Wireless.Count {
		total = v
		break
	}
	return resp.Wireless.StaList, total, nil
}
