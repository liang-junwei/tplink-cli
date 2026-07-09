package api

import (
	"fmt"

	"github.com/ljw/tplink-cli/internal/model"
)

// GetRedirects 查询所有端口映射规则
func (c *Client) GetRedirects() ([]model.RedirectItem, error) {
	req := model.GetRedirectsRequest{
		Method: "get",
		Firewall: model.FirewallGet{
			Table: "redirect",
		},
	}

	var resp model.GetRedirectsResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return nil, fmt.Errorf("查询端口映射规则失败: %w", err)
	}

	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("查询失败, error_code: %d", resp.ErrorCode)
	}

	return resp.Firewall.Redirect, nil
}

// AddRedirect 添加端口映射规则
func (c *Client) AddRedirect(rule model.RedirectRule) (string, error) {
	req := model.AddRedirectRequest{
		Method: "add",
		Firewall: model.FirewallAdd{
			Table: "redirect",
			Para:  rule,
		},
	}

	var resp model.AddRedirectResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return "", fmt.Errorf("添加端口映射规则失败: %w", err)
	}

	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("添加失败, error_code: %d", resp.ErrorCode)
	}

	if len(resp.Firewall.Name) > 0 {
		return resp.Firewall.Name[0], nil
	}

	return "", nil
}

// SetRedirect 修改端口映射规则
func (c *Client) SetRedirect(ruleID string, rule model.RedirectRule) error {
	req := model.SetRedirectRequest{
		Method: "set",
		Firewall: map[string]model.RedirectRule{
			ruleID: rule,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("修改端口映射规则失败: %w", err)
	}

	if resp.ErrorCode != 0 {
		return fmt.Errorf("修改失败, error_code: %d", resp.ErrorCode)
	}

	return nil
}

// EnableRedirect 启用端口映射规则
func (c *Client) EnableRedirect(ruleID string, rule model.RedirectRule) error {
	rule.Enable = "on"
	return c.SetRedirect(ruleID, rule)
}

// DisableRedirect 禁用端口映射规则
func (c *Client) DisableRedirect(ruleID string, rule model.RedirectRule) error {
	rule.Enable = "off"
	return c.SetRedirect(ruleID, rule)
}

// DeleteRedirect 删除端口映射规则
func (c *Client) DeleteRedirect(ruleID string) error {
	req := model.DeleteRedirectRequest{
		Method: "delete",
		Firewall: model.FirewallDelete{
			Name: ruleID,
		},
	}

	var resp model.APIResponse
	if err := c.DoRequest("POST", "", req, &resp); err != nil {
		return fmt.Errorf("删除端口映射规则失败: %w", err)
	}

	if resp.ErrorCode != 0 {
		return fmt.Errorf("删除失败, error_code: %d", resp.ErrorCode)
	}

	return nil
}

// FindRedirectByID 根据 ruleID 查找单条规则
func (c *Client) FindRedirectByID(ruleID string) (*model.RedirectRuleFull, error) {
	items, err := c.GetRedirects()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if rule, ok := item[ruleID]; ok {
			return &rule, nil
		}
	}

	return nil, fmt.Errorf("未找到规则: %s", ruleID)
}
