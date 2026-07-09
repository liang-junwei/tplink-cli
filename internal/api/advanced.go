package api

import (
	"fmt"

	"github.com/ljw/tplink-cli/internal/model"
)

// ========== 路由 - 系统路由 ==========

// GetSysRoutes 获取系统路由列表
func (c *Client) GetSysRoutes() ([]model.SysRouteItem, int, error) {
	req := map[string]interface{}{
		"method": "get",
		"network": map[string]interface{}{
			"table": "sys_route",
		},
	}

	var resp model.SysRouteListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("获取系统路由失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("获取系统路由失败, error_code: %d", resp.ErrorCode)
	}

	items := make([]model.SysRouteItem, 0)
	for _, m := range resp.Network.SysRoute {
		for key, item := range m {
			item.DotName = key
			items = append(items, item)
		}
	}
	return items, resp.Network.Count.SysRoute, nil
}

// ========== 路由 - 策略路由 ==========

// GetPolicyRoutes 获取策略路由列表
func (c *Client) GetPolicyRoutes() ([]model.PolicyRouteItem, int, error) {
	req := map[string]interface{}{
		"method": "get",
		"policy_route": map[string]interface{}{
			"table": "policy_rule",
		},
	}

	var resp model.PolicyRouteListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("获取策略路由失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("获取策略路由失败, error_code: %d", resp.ErrorCode)
	}

	items := make([]model.PolicyRouteItem, 0)
	for _, m := range resp.PolicyRoute.PolicyRule {
		for key, item := range m {
			item.DotName = key
			items = append(items, item)
		}
	}
	return items, resp.PolicyRoute.Count.PolicyRule, nil
}

// AddPolicyRoute 添加策略路由（两步：先添加 service，再添加 rule）
func (c *Client) AddPolicyRoute(name, serviceName, proto, sport, dport, ifName, srcIPGroup, dstIPGroup, timeObj, enable string) (string, error) {
	// 第一步: 添加 service
	serviceReq := map[string]interface{}{
		"method": "add",
		"service": map[string]interface{}{
			"table": "service",
			"para": []map[string]interface{}{
				{
					"name":  serviceName,
					"proto": proto,
					"sport": sport,
					"dport": dport,
					"flag":  "inner",
				},
			},
		},
	}

	if c.DryRun {
		var sResp model.PolicyRouteServiceAddResponse
		c.DoRequest("POST", "", serviceReq, &sResp)
	} else {
		var sResp model.PolicyRouteServiceAddResponse
		if err := c.DoRequest("POST", "", serviceReq, &sResp); err != nil {
			return "", fmt.Errorf("添加策略路由service失败: %w", err)
		}
		if sResp.ErrorCode != 0 {
			return "", fmt.Errorf("添加策略路由service失败, error_code: %d", sResp.ErrorCode)
		}
	}

	// 第二步: 添加 rule
	ruleReq := map[string]interface{}{
		"method": "add",
		"policy_route": map[string]interface{}{
			"table": "policy_rule",
			"para": map[string]interface{}{
				"name":         name,
				"service_type": serviceName,
				"src_ipgroup":  srcIPGroup,
				"dst_ipgroup":  dstIPGroup,
				"if":           ifName,
				"timeobj":      timeObj,
				"enable":       enable,
				"index":        "",
			},
		},
	}

	if c.DryRun {
		var rResp model.PolicyRouteRuleAddResponse
		c.DoRequest("POST", "", ruleReq, &rResp)
		return "policy_rule_dryrun", nil
	}

	var rResp model.PolicyRouteRuleAddResponse
	if err := c.DoRequest("POST", "", ruleReq, &rResp); err != nil {
		return "", fmt.Errorf("添加策略路由rule失败: %w", err)
	}
	if rResp.ErrorCode != 0 {
		return "", fmt.Errorf("添加策略路由rule失败, error_code: %d", rResp.ErrorCode)
	}
	if len(rResp.PolicyRoute.Name) > 0 {
		return rResp.PolicyRoute.Name[0], nil
	}
	return "", fmt.Errorf("添加策略路由失败: 未返回名称")
}

// DelPolicyRoute 删除策略路由（两步：先删 rule，再删关联的 service）
func (c *Client) DelPolicyRoute(ruleID, serviceID string) error {
	// 第一步: 删除 rule
	delRuleReq := map[string]interface{}{
		"method": "delete",
		"policy_route": map[string]interface{}{
			"name": ruleID,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", delRuleReq, &resp); err != nil {
		return fmt.Errorf("删除策略路由rule失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除策略路由rule失败, error_code: %d", resp.ErrorCode)
	}

	// 第二步: 删除关联的 service
	delSvcReq := map[string]interface{}{
		"method": "delete",
		"service": map[string]interface{}{
			"name": serviceID,
		},
	}
	if err := c.DoRequest("POST", "", delSvcReq, &resp); err != nil {
		return fmt.Errorf("删除策略路由service失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除策略路由service失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// ========== 路由 - 静态路由 ==========

// GetStaticRoutes 获取静态路由列表
func (c *Client) GetStaticRoutes() ([]model.StaticRouteItem, int, error) {
	req := map[string]interface{}{
		"method": "get",
		"network": map[string]interface{}{
			"table": "user_route",
		},
	}

	var resp model.StaticRouteListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("获取静态路由失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("获取静态路由失败, error_code: %d", resp.ErrorCode)
	}

	items := make([]model.StaticRouteItem, 0)
	for _, m := range resp.Network.UserRoute {
		for key, item := range m {
			item.DotName = key
			items = append(items, item)
		}
	}
	return items, resp.Network.Count.UserRoute, nil
}

// AddStaticRoute 添加静态路由
func (c *Client) AddStaticRoute(name, target, netmask, gateway, ifName, metric, note, enable string) (string, error) {
	req := map[string]interface{}{
		"method": "add",
		"network": map[string]interface{}{
			"table": "user_route",
			"para": map[string]interface{}{
				"name":    name,
				"target":  target,
				"netmask": netmask,
				"gateway": gateway,
				"if":      ifName,
				"metric":  metric,
				"note":    note,
				"enable":  enable,
			},
		},
	}

	if c.DryRun {
		var resp model.StaticRouteAddResponse
		c.DoRequest("POST", "", req, &resp)
		return "user_route_dryrun", nil
	}

	var resp model.StaticRouteAddResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return "", fmt.Errorf("添加静态路由失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("添加静态路由失败, error_code: %d", resp.ErrorCode)
	}
	if len(resp.Network.Name) > 0 {
		return resp.Network.Name[0], nil
	}
	return "", fmt.Errorf("添加静态路由失败: 未返回名称")
}

// DelStaticRoute 删除静态路由
func (c *Client) DelStaticRoute(id string) error {
	req := map[string]interface{}{
		"method": "delete",
		"network": map[string]interface{}{
			"name": id,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("删除静态路由失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除静态路由失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// ========== NAPT ==========

// GetNaptRules 获取NAPT规则列表
func (c *Client) GetNaptRules() ([]model.NaptRuleItem, int, error) {
	req := map[string]interface{}{
		"method": "get",
		"nat": map[string]interface{}{
			"table": "rule_napt",
		},
	}

	var resp model.NaptRuleListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("获取NAPT规则失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("获取NAPT规则失败, error_code: %d", resp.ErrorCode)
	}

	items := make([]model.NaptRuleItem, 0)
	for _, m := range resp.Nat.RuleNapt {
		for key, item := range m {
			item.DotName = key
			items = append(items, item)
		}
	}
	return items, resp.Nat.Count.RuleNapt, nil
}

// AddNaptRule 添加NAPT规则
func (c *Client) AddNaptRule(name, ip, ifName, enable string) (string, error) {
	req := map[string]interface{}{
		"method": "add",
		"nat": map[string]interface{}{
			"table": "rule_napt",
			"para": map[string]interface{}{
				"name":   name,
				"ip":     ip,
				"if":     ifName,
				"enable": enable,
			},
		},
	}

	if c.DryRun {
		var resp model.NaptRuleAddResponse
		c.DoRequest("POST", "", req, &resp)
		return "rule_napt_dryrun", nil
	}

	var resp model.NaptRuleAddResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return "", fmt.Errorf("添加NAPT规则失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("添加NAPT规则失败, error_code: %d", resp.ErrorCode)
	}
	if len(resp.Nat.Name) > 0 {
		return resp.Nat.Name[0], nil
	}
	return "", fmt.Errorf("添加NAPT规则失败: 未返回名称")
}

// DelNaptRule 删除NAPT规则
func (c *Client) DelNaptRule(id string) error {
	req := map[string]interface{}{
		"method": "delete",
		"nat": map[string]interface{}{
			"name": id,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("删除NAPT规则失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除NAPT规则失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// ========== ALG ==========

// GetAlgConfig 获取ALG配置
func (c *Client) GetAlgConfig() (*model.AlgConfig, error) {
	req := map[string]interface{}{
		"method": "get",
		"nat": map[string]interface{}{
			"name": "alg_glb",
		},
	}

	var resp model.AlgConfigResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("获取ALG配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("获取ALG配置失败, error_code: %d", resp.ErrorCode)
	}

	// dry-run 返回模拟数据
	if c.DryRun {
		return &model.AlgConfig{}, nil
	}
	return &resp.Nat.AlgGlb, nil
}

// SetAlgConfig 设置ALG配置（只发送用户指定的字段）
func (c *Client) SetAlgConfig(cfg map[string]interface{}) error {
	req := map[string]interface{}{
		"method": "set",
		"nat": map[string]interface{}{
			"alg_glb": cfg,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("设置ALG配置失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("设置ALG配置失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// ========== Phddns ==========

// GetPhddnsList 获取Phddns列表
func (c *Client) GetPhddnsList() ([]model.PhddnsItem, int, error) {
	req := map[string]interface{}{
		"method": "get",
		"phddns": map[string]interface{}{
			"table": "phddns",
		},
	}

	var resp model.PhddnsListResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, 0, fmt.Errorf("获取Phddns列表失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, 0, fmt.Errorf("获取Phddns列表失败, error_code: %d", resp.ErrorCode)
	}

	items := make([]model.PhddnsItem, 0)
	for _, m := range resp.Phddns.Phddns {
		for key, item := range m {
			item.DotName = key
			items = append(items, item)
		}
	}
	return items, resp.Phddns.Count.Phddns, nil
}

// AddPhddns 添加Phddns配置
func (c *Client) AddPhddns(username, password, ifName, enable string) (string, error) {
	req := map[string]interface{}{
		"method": "add",
		"phddns": map[string]interface{}{
			"table": "phddns",
			"para": map[string]interface{}{
				"interface": ifName,
				"username":  username,
				"password":  password,
				"enable":    enable,
			},
		},
	}

	if c.DryRun {
		var resp model.PhddnsAddResponse
		c.DoRequest("POST", "", req, &resp)
		return "phddns_dryrun", nil
	}

	var resp model.PhddnsAddResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return "", fmt.Errorf("添加Phddns失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("添加Phddns失败, error_code: %d", resp.ErrorCode)
	}
	if len(resp.Phddns.Name) > 0 {
		return resp.Phddns.Name[0], nil
	}
	return "", fmt.Errorf("添加Phddns失败: 未返回名称")
}

// SetPhddns 修改Phddns配置（需要传入条目的 DotName）
func (c *Client) SetPhddns(dotName string, cfg map[string]interface{}) error {
	req := map[string]interface{}{
		"method": "set",
		"phddns": map[string]interface{}{
			dotName: cfg,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("设置Phddns失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("设置Phddns失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}

// DelPhddns 删除Phddns配置
func (c *Client) DelPhddns(id string) error {
	req := map[string]interface{}{
		"method": "delete",
		"phddns": map[string]interface{}{
			"name": id,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("删除Phddns失败: %w", err)
	}
	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除Phddns失败, error_code: %d", resp.ErrorCode)
	}
	return nil
}
