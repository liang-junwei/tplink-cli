package model

// LoginRequest 登录请求
type LoginRequest struct {
	Method string      `json:"method"`
	Login  LoginParams `json:"login"`
}

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Stok      string `json:"stok"`
	ErrorCode int    `json:"error_code"`
}

// RedirectRule 端口映射规则
type RedirectRule struct {
	Name           string `json:"name"`
	If             string `json:"if"`
	SrcDport       string `json:"src_dport"`
	DestPort       string `json:"dest_port"`
	DestIP         string `json:"dest_ip"`
	Proto          string `json:"proto"`
	SrcDportStart  string `json:"src_dport_start"`
	SrcDportEnd    string `json:"src_dport_end"`
	DestPortStart  string `json:"dest_port_start"`
	DestPortEnd    string `json:"dest_port_end"`
	Enable         string `json:"enable"`
}

// RedirectRuleFull 包含元数据的完整规则（查询返回）
type RedirectRuleFull struct {
	Name           string `json:"name"`
	If             string `json:"if"`
	SrcDport       string `json:"src_dport"`
	DestPort       string `json:"dest_port"`
	DestIP         string `json:"dest_ip"`
	Proto          string `json:"proto"`
	SrcDportStart  string `json:"src_dport_start"`
	SrcDportEnd    string `json:"src_dport_end"`
	DestPortStart  string `json:"dest_port_start"`
	DestPortEnd    string `json:"dest_port_end"`
	Enable         string `json:"enable"`
	InternalName   string `json:".name"`
	Index          int    `json:".index"`
	Type           string `json:".type"`
	Anonymous      bool   `json:".anonymous"`
}

// RedirectItem 查询结果中的单条规则（key 为规则ID）
type RedirectItem map[string]RedirectRuleFull

// GetRedirectsRequest 查询端口映射规则请求
type GetRedirectsRequest struct {
	Method   string      `json:"method"`
	Firewall FirewallGet `json:"firewall"`
}

type FirewallGet struct {
	Table string `json:"table"`
}

// GetRedirectsResponse 查询端口映射规则响应
type GetRedirectsResponse struct {
	Firewall FirewallResult `json:"firewall"`
	ErrorCode int           `json:"error_code"`
}

type FirewallResult struct {
	Redirect []RedirectItem   `json:"redirect"`
	Count    map[string]int   `json:"count"`
}

// AddRedirectRequest 添加端口映射规则请求
type AddRedirectRequest struct {
	Method   string      `json:"method"`
	Firewall FirewallAdd `json:"firewall"`
}

type FirewallAdd struct {
	Table string       `json:"table"`
	Para  RedirectRule `json:"para"`
}

// AddRedirectResponse 添加端口映射规则响应
type AddRedirectResponse struct {
	Firewall struct {
		Name []string `json:"name"`
	} `json:"firewall"`
	ErrorCode int `json:"error_code"`
}

// SetRedirectRequest 修改/启用/禁用端口映射规则请求
type SetRedirectRequest struct {
	Method   string               `json:"method"`
	Firewall map[string]RedirectRule `json:"firewall"`
}

// DeleteRedirectRequest 删除端口映射规则请求
type DeleteRedirectRequest struct {
	Method   string         `json:"method"`
	Firewall FirewallDelete `json:"firewall"`
}

type FirewallDelete struct {
	Name string `json:"name"`
}

// APIResponse 通用API响应（仅含error_code）
type APIResponse struct {
	ErrorCode int `json:"error_code"`
}

// Config CLI 配置文件结构
type Config struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
