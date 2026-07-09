package api

import (
	"fmt"

	"github.com/ljw/tplink-cli/internal/model"
)

// ========== IP地址组 (ipgroup) ==========

// GetIPGroups 获取IP地址组列表，支持分页
func (c *Client) GetIPGroups(start, end int) ([]model.IPGroupItem, int, error) {
	req := map[string]interface{}{
		"method": "get",
		"ipgroup": map[string]interface{}{
			"table": "rule_ipgroup",
			"filter": []map[string]interface{}{
				{"flag": "system"},
				{"flag": "user"},
			},
			"para": map[string]interface{}{
				"start": start,
				"end":   end,
			},
		},
	}

	var resp model.IPGroupListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("查询IP地址组失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	items := make([]model.IPGroupItem, 0)
	for _, m := range resp.IPGroup.RuleIPGroup {
		for k, item := range m {
			item.DotName = k
			items = append(items, item)
		}
	}

	total := 0
	for _, v := range resp.IPGroup.Count {
		total = v
		break
	}
	return items, total, nil
}

// AddIPGroupScope 添加IP地址范围（第一步）
func (c *Client) AddIPGroupScope(name, scopeType, scope string) (string, error) {
	req := model.IPGroupAddScopeRequest{
		Method: "add",
		IPGroup: model.IPGroupAddScopeIn{
			Table: "rule_ipscope",
			Para: []map[string]interface{}{
				{
					"flag":  "inner",
					"name":  name,
					"type":  scopeType,
					"scope": scope,
				},
			},
		},
	}

	var resp model.IPGroupAddResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return "", fmt.Errorf("添加IP范围失败: %w", err)
	}
	if c.DryRun {
		return "rule_ipscope_dryrun", nil
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("添加失败, error_code: %d", resp.ErrorCode)
	}
	if len(resp.IPGroup.Name) > 0 {
		return resp.IPGroup.Name[0], nil
	}
	return "", fmt.Errorf("添加IP范围失败: 未返回名称")
}

// AddIPGroup 添加IP地址组（第二步）
func (c *Client) AddIPGroup(name, scopeName string) (string, error) {
	req := model.IPGroupAddGroupRequest{
		Method: "add",
		IPGroup: model.IPGroupAddGroupIn{
			Table: "rule_ipgroup",
			Para: map[string]interface{}{
				"flag":       "user",
				"name":       name,
				"rule_scope": scopeName,
				"ref":        "0",
			},
		},
	}

	var resp model.IPGroupAddResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return "", fmt.Errorf("添加IP地址组失败: %w", err)
	}
	if c.DryRun {
		return "rule_ipgroup_dryrun", nil
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("添加失败, error_code: %d", resp.ErrorCode)
	}
	if len(resp.IPGroup.Name) > 0 {
		return resp.IPGroup.Name[0], nil
	}
	return "", fmt.Errorf("添加IP地址组失败: 未返回名称")
}

// DelIPGroup 删除IP地址组或IP范围
func (c *Client) DelIPGroup(name string) error {
	req := model.IPGroupDelRequest{
		Method: "delete",
		IPGroup: model.IPGroupDelIn{
			Name: name,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// ========== 时间段 (timerange) ==========

// GetTimeRanges 获取时间段列表，支持分页
func (c *Client) GetTimeRanges(start, end int) ([]model.TimeRangeItem, int, error) {
	req := map[string]interface{}{
		"method": "get",
		"time_mngt": map[string]interface{}{
			"table": "time_obj",
			"filter": map[string]interface{}{
				"flag": []string{"system", "user"},
			},
			"para": map[string]interface{}{
				"start": start,
				"end":   end,
			},
		},
	}

	var resp model.TimeRangeListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("查询时间段失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	items := make([]model.TimeRangeItem, 0)
	for _, m := range resp.TimeMngt.TimeObj {
		for k, item := range m {
			item.DotName = k
			items = append(items, item)
		}
	}

	total := 0
	for _, v := range resp.TimeMngt.Count {
		total = v
		break
	}
	return items, total, nil
}

// AddTimeRange 添加时间段
func (c *Client) AddTimeRange(name, mode string, timeSection []string, weekday, comment string) (string, error) {
	req := model.TimeRangeAddRequest{
		Method: "add",
		TimeMngt: model.TimeRangeAddIn{
			Table: "time_obj",
			Para: map[string]interface{}{
				"name":         name,
				"mode":         mode,
				"time_section": timeSection,
				"weekday":      weekday,
				"comment":      comment,
				"flag":         "user",
			},
		},
	}

	var resp model.TimeRangeAddResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return "", fmt.Errorf("添加时间段失败: %w", err)
	}
	if c.DryRun {
		return "time_obj_dryrun", nil
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("添加失败, error_code: %d", resp.ErrorCode)
	}
	if len(resp.TimeMngt.Name) > 0 {
		return resp.TimeMngt.Name[0], nil
	}
	return "", fmt.Errorf("添加时间段失败: 未返回名称")
}

// DelTimeRange 删除时间段
func (c *Client) DelTimeRange(name string) error {
	req := model.TimeRangeDelRequest{
		Method: "delete",
		TimeMngt: model.TimeRangeDelIn{
			Name: name,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// ========== 带宽控制 (Qos) ==========

// GetQosConfig 获取Qos配置
func (c *Client) GetQosConfig() (*model.QosSetting, error) {
	req := map[string]interface{}{
		"method": "get",
		"qos": map[string]interface{}{
			"name": "setting",
		},
	}

	var resp model.QosConfigResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("查询Qos配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	return &resp.Qos.Setting, nil
}

// SetQosConfig 设置Qos配置
func (c *Client) SetQosConfig(cfg map[string]interface{}) error {
	req := model.QosConfigSetRequest{
		Method: "set",
		Qos: map[string]interface{}{
			"setting": cfg,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("设置Qos配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("设置失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// GetQosRules 获取Qos规则列表，支持分页
func (c *Client) GetQosRules(start, end int) ([]model.QosRuleItem, int, error) {
	req := map[string]interface{}{
		"method": "get",
		"qos": map[string]interface{}{
			"table": "rule",
			"para": map[string]interface{}{
				"start": start,
				"end":   end,
			},
		},
	}

	var resp model.QosRuleListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("查询Qos规则失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	items := make([]model.QosRuleItem, 0)
	for _, m := range resp.Qos.Rule {
		for key, item := range m {
			item.DotName = key // 捕获 map key 作为内部 ID
			items = append(items, item)
		}
	}

	total := 0
	for _, v := range resp.Qos.Count {
		total = v
		break
	}
	return items, total, nil
}

// AddQosRule 添加Qos规则
func (c *Client) AddQosRule(para map[string]interface{}) (string, error) {
	req := model.QosRuleAddRequest{
		Method: "add",
		Qos: model.QosRuleAddIn{
			Table: "rule",
			Para:  para,
		},
	}

	var resp model.QosRuleAddResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return "", fmt.Errorf("添加Qos规则失败: %w", err)
	}
	if c.DryRun {
		return "rule_dryrun", nil
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("添加失败, error_code: %d", resp.ErrorCode)
	}
	if len(resp.Qos.Name) > 0 {
		return resp.Qos.Name[0], nil
	}
	return "", fmt.Errorf("添加Qos规则失败: 未返回名称")
}

// DelQosRule 删除Qos规则
func (c *Client) DelQosRule(name string) error {
	req := model.QosRuleDelRequest{
		Method: "delete",
		Qos: model.QosRuleDelIn{
			Name: name,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// ========== 访问控制 (ACL) ==========

// GetACLRules 获取ACL规则列表，支持分页
func (c *Client) GetACLRules(start, end int) ([]model.ACLRuleItem, int, error) {
	req := map[string]interface{}{
		"method": "get",
		"access_ctl": map[string]interface{}{
			"table": "rule_acl_inner",
			"para": map[string]interface{}{
				"start": start,
				"end":   end,
			},
		},
	}

	var resp model.ACLRuleListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("查询ACL规则失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	items := make([]model.ACLRuleItem, 0)
	for _, m := range resp.AccessCtl.RuleACLInner {
		for key, item := range m {
			item.DotName = key // 捕获 map key 作为内部 ID
			items = append(items, item)
		}
	}

	total := 0
	for _, v := range resp.AccessCtl.Count {
		total = v
		break
	}
	return items, total, nil
}

// GetACLServices 获取ACL服务列表
func (c *Client) GetACLServices() ([]model.ACLServiceItem, int, error) {
	req := map[string]interface{}{
		"method": "get",
		"service": map[string]interface{}{
			"table": "service",
		},
	}

	var resp model.ACLServiceListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("查询ACL服务失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	items := make([]model.ACLServiceItem, 0)
	for _, m := range resp.Service.Service {
		for _, item := range m {
			items = append(items, item)
		}
	}

	total := 0
	for _, v := range resp.Service.Count {
		total = v
		break
	}
	return items, total, nil
}

// AddACLService 添加ACL服务（第一步：先建服务）
func (c *Client) AddACLService(name, proto, sport, dport string) (string, error) {
	req := model.ACLServiceAddRequest{
		Method: "add",
		Service: model.ACLServiceAddIn{
			Table: "service",
			Para: []map[string]interface{}{
				{
					"name":  name,
					"proto": proto,
					"sport": sport,
					"dport": dport,
					"ref":   "0",
					"flag":  "inner",
				},
			},
		},
	}

	var resp model.ACLServiceAddResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return "", fmt.Errorf("添加ACL服务失败: %w", err)
	}
	if c.DryRun {
		return "service_dryrun", nil
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("添加失败, error_code: %d", resp.ErrorCode)
	}
	if len(resp.Service.Name) > 0 {
		return resp.Service.Name[0], nil
	}
	return "", fmt.Errorf("添加ACL服务失败: 未返回名称")
}

// AddACLRule 添加ACL规则（第二步：建规则引用服务）
func (c *Client) AddACLRule(para map[string]interface{}) (string, error) {
	req := model.ACLRuleAddRequest{
		Method: "add",
		AccessCtl: model.ACLRuleAddIn{
			Table: "rule_acl_inner",
			Para:  para,
		},
	}

	var resp model.ACLRuleAddResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return "", fmt.Errorf("添加ACL规则失败: %w", err)
	}
	if c.DryRun {
		return "rule_acl_inner_dryrun", nil
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("添加失败, error_code: %d", resp.ErrorCode)
	}
	if len(resp.AccessCtl.Name) > 0 {
		return resp.AccessCtl.Name[0], nil
	}
	return "", fmt.Errorf("添加ACL规则失败: 未返回名称")
}

// DelACLRule 删除ACL规则
func (c *Client) DelACLRule(name string) error {
	req := model.ACLRuleDelRequest{
		Method: "delete",
		AccessCtl: model.ACLRuleDelIn{
			Name: name,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// DelACLService 删除ACL服务
func (c *Client) DelACLService(name string) error {
	req := model.ACLServiceDelRequest{
		Method: "delete",
		Service: model.ACLServiceDelIn{
			Name: name,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}
