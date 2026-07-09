package api

import (
	"fmt"

	"github.com/ljw/tplink-cli/internal/model"
)

// ========== ARP 防护 (arp) ==========

// GetArpConfig 获取ARP防护配置
func (c *Client) GetArpConfig() (*model.ArpConfig, error) {
	req := map[string]interface{}{
		"method": "get",
		"arp_defense": map[string]interface{}{
			"name": "global",
		},
	}

	var resp model.ArpConfigResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("获取ARP配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("获取ARP配置失败, error_code: %d", resp.ErrorCode)
	}
	return &resp.ArpDefense.Global, nil
}

// SetArpConfig 设置ARP防护配置（仅发送用户指定的字段）
func (c *Client) SetArpConfig(cfg map[string]interface{}) error {
	req := model.ArpConfigSetRequest{
		Method: "set",
	}
	req.ArpDefense.Global = cfg

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("设置ARP配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("设置ARP配置失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// GetArpBindList 获取ARP绑定列表，支持分页
func (c *Client) GetArpBindList(start, end int) ([]model.ArpBindItem, int, error) {
	req := map[string]interface{}{
		"method": "get",
		"ip_mac_bind": map[string]interface{}{
			"table": "sys_arp",
			"para": map[string]interface{}{
				"start": start,
				"end":   end,
			},
		},
	}

	var resp model.ArpBindListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("获取ARP列表失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("获取ARP列表失败, error_code: %d", resp.ErrorCode)
	}

	items := make([]model.ArpBindItem, 0)
	for _, m := range resp.IPMacBind.SysArp {
		for key, item := range m {
			item.DotName = key
			items = append(items, item)
		}
	}
	return items, resp.IPMacBind.Count.SysArp, nil
}

// ========== MAC 地址过滤 (macfilter) ==========

// GetMacFilterConfig 获取MAC过滤配置
func (c *Client) GetMacFilterConfig() (*model.MacFilterConfig, error) {
	req := map[string]interface{}{
		"method": "get",
		"mac_filter": map[string]interface{}{
			"name": "global",
		},
	}

	var resp model.MacFilterConfigResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("获取MAC过滤配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("获取MAC过滤配置失败, error_code: %d", resp.ErrorCode)
	}
	return &resp.MacFilter.Global, nil
}

// SetMacFilterConfig 设置MAC过滤配置（仅发送用户指定的字段）
func (c *Client) SetMacFilterConfig(cfg map[string]interface{}) error {
	req := model.MacFilterConfigSetRequest{
		Method: "set",
	}
	req.MacFilter.Global = cfg

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("设置MAC过滤配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("设置MAC过滤配置失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// GetMacFilterRules 获取MAC过滤规则列表，支持分页
func (c *Client) GetMacFilterRules(start, end int) ([]model.MacFilterRuleItem, int, error) {
	req := map[string]interface{}{
		"method": "get",
		"mac_filter": map[string]interface{}{
			"table": "mac_filter_list",
			"para": map[string]interface{}{
				"start": start,
				"end":   end,
			},
		},
	}

	var resp model.MacFilterRuleListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("获取MAC过滤规则列表失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("获取MAC过滤规则列表失败, error_code: %d", resp.ErrorCode)
	}

	items := make([]model.MacFilterRuleItem, 0)
	for _, m := range resp.MacFilter.MacFilterList {
		for key, item := range m {
			item.DotName = key
			items = append(items, item)
		}
	}
	return items, resp.MacFilter.Count.MacFilterList, nil
}

// AddMacFilterRule 添加MAC过滤规则
func (c *Client) AddMacFilterRule(name, mac string) (string, error) {
	req := model.MacFilterRuleAddRequest{
		Method: "add",
	}
	req.MacFilter.Table = "mac_filter_list"
	req.MacFilter.Para = map[string]interface{}{
		"name": name,
		"mac":  mac,
	}

	// dry-run 模式: 先调用 DoRequest 打印请求，再返回模拟ID
	if c.DryRun {
		var resp model.MacFilterRuleAddResponse
		c.DoRequest("POST", "", req, &resp)
		return "mac_filter_list_dryrun", nil
	}

	var resp model.MacFilterRuleAddResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return "", fmt.Errorf("添加MAC过滤规则失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("添加失败, error_code: %d", resp.ErrorCode)
	}
	if len(resp.MacFilter.Name) > 0 {
		return resp.MacFilter.Name[0], nil
	}
	return "", fmt.Errorf("添加MAC过滤规则失败: 未返回名称")
}

// DelMacFilterRule 删除MAC过滤规则
func (c *Client) DelMacFilterRule(id string) error {
	req := model.MacFilterRuleDelRequest{
		Method: "delete",
	}
	req.MacFilter.Name = id

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("删除MAC过滤规则失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// ========== DoS 攻击防护 (dos) ==========

// GetDosConfig 获取DoS攻击防护配置
func (c *Client) GetDosConfig() (*model.DosConfig, error) {
	req := map[string]interface{}{
		"method": "get",
		"dos_defense": map[string]interface{}{
			"name": "global",
		},
	}

	var resp model.DosConfigResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("获取DoS配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("获取DoS配置失败, error_code: %d", resp.ErrorCode)
	}
	return &resp.DosDefense.Global, nil
}

// SetDosConfig 设置DoS攻击防护配置（仅发送用户指定的字段）
func (c *Client) SetDosConfig(cfg map[string]interface{}) error {
	req := model.DosConfigSetRequest{
		Method: "set",
	}
	req.DosDefense.Global = cfg

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("设置DoS配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("设置DoS配置失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// ========== Flood 攻击防护 (flood) ==========

// GetFloodConfig 获取Flood攻击防护配置（含 global + threshold）
func (c *Client) GetFloodConfig() (*model.FloodGlobal, *model.FloodThreshold, error) {
	req := map[string]interface{}{
		"method": "get",
		"flood_defense": map[string]interface{}{
			"name": []string{"global", "threshold"},
		},
	}

	var resp model.FloodConfigResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, nil, fmt.Errorf("获取Flood配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, nil, fmt.Errorf("获取Flood配置失败, error_code: %d", resp.ErrorCode)
	}
	return &resp.FloodDefense.Global, &resp.FloodDefense.Threshold, nil
}

// SetFloodConfig 设置Flood攻击防护配置（同时设置 global 和 threshold）
func (c *Client) SetFloodConfig(globalCfg, thresholdCfg map[string]interface{}) error {
	req := model.FloodConfigSetRequest{
		Method: "set",
	}
	req.FloodDefense.Global = globalCfg
	req.FloodDefense.Threshold = thresholdCfg

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("设置Flood配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("设置Flood配置失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}
