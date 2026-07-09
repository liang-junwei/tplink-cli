package model

import "encoding/json"

// ========== 路由 - 系统路由 (sys_route) ==========

// SysRouteItem 系统路由条目
type SysRouteItem struct {
	DotName string `json:".name"`
	Metric  string `json:"metric"`
	Target  string `json:"target"`
	Netmask string `json:"netmask"`
	If      string `json:"if"` // WAN1 | LOOPBACK
	Gateway string `json:"gateway"`
}

// SysRouteListResponse 系统路由列表响应
type SysRouteListResponse struct {
	ErrorCode int              `json:"error_code"`
	Network   SysRouteListData `json:"network"`
}

// SysRouteListData 系统路由列表数据
type SysRouteListData struct {
	Count    SysRouteCount              `json:"count"`
	SysRoute []map[string]SysRouteItem  `json:"sys_route"`
}

// SysRouteCount 系统路由计数
type SysRouteCount struct {
	SysRoute int `json:"sys_route"`
}

// ========== 路由 - 策略路由 (policy_rule) ==========

// PolicyRouteItem 策略路由条目
type PolicyRouteItem struct {
	DotName     string `json:".name"`
	Name        string `json:"name"`
	ServiceType string `json:"service_type"` // 关联的 service 名称 (inner name)
	If          string `json:"if"`           // WAN1
	SrcIPGroup  string `json:"src_ipgroup"`
	DstIPGroup  string `json:"dst_ipgroup"`
	TimeObj     string `json:"timeobj"`
	Enable      string `json:"enable"` // on | off
	Forced      string `json:"forced"`
	Index       json.Number `json:"index"`
}

// PolicyRouteListResponse 策略路由列表响应
type PolicyRouteListResponse struct {
	ErrorCode   int                  `json:"error_code"`
	PolicyRoute PolicyRouteListData  `json:"policy_route"`
}

// PolicyRouteListData 策略路由列表数据
type PolicyRouteListData struct {
	Count      PolicyRouteCount               `json:"count"`
	PolicyRule []map[string]PolicyRouteItem   `json:"policy_rule"`
}

// PolicyRouteCount 策略路由计数
type PolicyRouteCount struct {
	PolicyRule int `json:"policy_rule"`
}

// PolicyRouteServiceAddResponse 策略路由 service add 响应
type PolicyRouteServiceAddResponse struct {
	Service struct {
		Name []string `json:"name"`
	} `json:"service"`
	ErrorCode int `json:"error_code"`
}

// PolicyRouteRuleAddResponse 策略路由 rule add 响应
type PolicyRouteRuleAddResponse struct {
	PolicyRoute struct {
		Name []string `json:"name"`
	} `json:"policy_route"`
	ErrorCode int `json:"error_code"`
}

// ========== 路由 - 静态路由 (user_route) ==========

// StaticRouteItem 静态路由条目
type StaticRouteItem struct {
	DotName      string `json:".name"`
	Name         string `json:"name"`
	Target       string `json:"target"`
	Netmask      string `json:"netmask"`
	Gateway      string `json:"gateway"`
	If           string `json:"if"` // LAN | WAN1
	Metric       string `json:"metric"`
	Note         string `json:"note"`
	Enable       string `json:"enable"` // on | off
	Reachability string `json:"reachability"`
}

// StaticRouteListResponse 静态路由列表响应
type StaticRouteListResponse struct {
	ErrorCode int                 `json:"error_code"`
	Network   StaticRouteListData `json:"network"`
}

// StaticRouteListData 静态路由列表数据
type StaticRouteListData struct {
	Count     StaticRouteCount              `json:"count"`
	UserRoute []map[string]StaticRouteItem  `json:"user_route"`
}

// StaticRouteCount 静态路由计数
type StaticRouteCount struct {
	UserRoute int `json:"user_route"`
}

// StaticRouteAddResponse 静态路由 add 响应
type StaticRouteAddResponse struct {
	Network struct {
		Name []string `json:"name"`
	} `json:"network"`
	ErrorCode int `json:"error_code"`
}

// ========== NAPT (rule_napt) ==========

// NaptRuleItem NAPT规则条目
type NaptRuleItem struct {
	DotName string `json:".name"`
	Name    string `json:"name"`
	Enable  string `json:"enable"` // on | off
	If      string `json:"if"`     // WAN1
	IP      string `json:"ip"`     // URL编码
	SysRule string `json:"sysrule"` // 1=系统规则
}

// NaptRuleListResponse NAPT规则列表响应
type NaptRuleListResponse struct {
	ErrorCode int              `json:"error_code"`
	Nat       NaptRuleListData `json:"nat"`
}

// NaptRuleListData NAPT规则列表数据
type NaptRuleListData struct {
	Count    NaptRuleCount              `json:"count"`
	RuleNapt []map[string]NaptRuleItem  `json:"rule_napt"`
}

// NaptRuleCount NAPT规则计数
type NaptRuleCount struct {
	RuleNapt int `json:"rule_napt"`
}

// NaptRuleAddResponse NAPT规则 add 响应
type NaptRuleAddResponse struct {
	Nat struct {
		Name []string `json:"name"`
	} `json:"nat"`
	ErrorCode int `json:"error_code"`
}

// ========== ALG (alg_glb) ==========

// AlgConfig ALG配置
type AlgConfig struct {
	Ftp   string `json:"ftp"`
	H323  string `json:"h323"`
	Pptp  string `json:"pptp"`
	Sip   string `json:"sip"`
	L2tp  string `json:"l2tp"`
	Tftp  string `json:"tftp"`
	Ipsec string `json:"ipsec"`
}

// AlgConfigResponse ALG配置响应
type AlgConfigResponse struct {
	ErrorCode int           `json:"error_code"`
	Nat       AlgConfigData `json:"nat"`
}

// AlgConfigData ALG配置数据
type AlgConfigData struct {
	AlgGlb AlgConfig `json:"alg_glb"`
}

// ========== Phddns (phddns) ==========

// PhddnsItem Phddns条目
type PhddnsItem struct {
	DotName   string `json:".name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Enable    string `json:"enable"`
	Interface string `json:"interface"`
	Domain    string `json:"domain"`
	ConnState string `json:"connstate"`
	DomainNum string `json:"domain_num"`
	UserType  string `json:"usertype"`
}

// PhddnsListResponse Phddns列表响应
type PhddnsListResponse struct {
	ErrorCode int            `json:"error_code"`
	Phddns    PhddnsListData `json:"phddns"`
}

// PhddnsListData Phddns列表数据
type PhddnsListData struct {
	Count  PhddnsCount              `json:"count"`
	Phddns []map[string]PhddnsItem  `json:"phddns"`
}

// PhddnsCount Phddns计数
type PhddnsCount struct {
	Phddns int `json:"phddns"`
}

// PhddnsAddResponse Phddns add 响应
type PhddnsAddResponse struct {
	Phddns struct {
		Name []string `json:"name"`
	} `json:"phddns"`
	ErrorCode int `json:"error_code"`
}
