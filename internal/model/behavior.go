package model

// ========== 通用响应/请求 ==========

// GenericNameResponse 通用 {"xxx":{"name":["..."]}} 格式响应（add 操作）
type GenericNameResponse struct {
	ErrorCode int      `json:"error_code"`
	Name      []string `json:"-"` // 解析时动态填充
}

// GenericCount 通用 count 结构 {"count": {"table_name": N}}
type GenericCount map[string]int

// ========== IP地址组 (ipgroup) ==========

// IPGroupItem IP地址组条目（来自 rule_ipgroup 表）
type IPGroupItem struct {
	DotName string `json:".name"` // 内部名称（map key），如 rule_ipgroup_xxx
	Name    string `json:"name"`  // 显示名称
	Comment string `json:"comment"`
	Flag    string `json:"flag"` // system | user
	Ref     string `json:"ref"`
}

// IPGroupListResponse IP地址组列表响应
type IPGroupListResponse struct {
	IPGroup struct {
		RuleIPGroup  []map[string]IPGroupItem `json:"rule_ipgroup"`
		Count        GenericCount             `json:"count"`
	} `json:"ipgroup"`
	ErrorCode int `json:"error_code"`
}

// IPGroupAddScopeRequest 添加IP范围请求（第一步：添加到 rule_ipscope）
type IPGroupAddScopeRequest struct {
	Method   string                   `json:"method"`
	IPGroup  IPGroupAddScopeIn        `json:"ipgroup"`
}

type IPGroupAddScopeIn struct {
	Table string                   `json:"table"`
	Para  []map[string]interface{} `json:"para"`
}

// IPGroupAddGroupRequest 添加IP地址组请求（第二步：添加到 rule_ipgroup）
type IPGroupAddGroupRequest struct {
	Method  string            `json:"method"`
	IPGroup IPGroupAddGroupIn `json:"ipgroup"`
}

type IPGroupAddGroupIn struct {
	Table string                 `json:"table"`
	Para  map[string]interface{} `json:"para"`
}

// IPGroupAddResponse IP地址组添加响应
type IPGroupAddResponse struct {
	IPGroup struct {
		Name []string `json:"name"`
	} `json:"ipgroup"`
	ErrorCode int `json:"error_code"`
}

// IPGroupDelRequest IP地址组删除请求
type IPGroupDelRequest struct {
	Method  string         `json:"method"`
	IPGroup IPGroupDelIn   `json:"ipgroup"`
}

type IPGroupDelIn struct {
	Name string `json:"name"`
}

// ========== 时间段 (timerange) ==========

// TimeRangeItem 时间段条目
type TimeRangeItem struct {
	DotName     string   `json:".name"`     // 内部名称（map key），如 time_obj_xxx
	Name        string   `json:"name"`      // 显示名称
	Comment     string   `json:"comment"`
	Flag        string   `json:"flag"` // system | user
	Ref         string   `json:"ref"`
	Mode        string   `json:"mode"`        // manual
	Weekday     string   `json:"weekday"`     // 位掩码
	TimeSection []string `json:"time_section"` // "0100,0359"
}

// TimeRangeListResponse 时间段列表响应
type TimeRangeListResponse struct {
	TimeMngt struct {
		TimeObj []map[string]TimeRangeItem `json:"time_obj"`
		Count   GenericCount               `json:"count"`
	} `json:"time_mngt"`
	ErrorCode int `json:"error_code"`
}

// TimeRangeAddRequest 添加时间段请求
type TimeRangeAddRequest struct {
	Method   string            `json:"method"`
	TimeMngt TimeRangeAddIn    `json:"time_mngt"`
}

type TimeRangeAddIn struct {
	Table string                 `json:"table"`
	Para  map[string]interface{} `json:"para"`
}

// TimeRangeAddResponse 添加时间段响应
type TimeRangeAddResponse struct {
	TimeMngt struct {
		Name []string `json:"name"`
	} `json:"time_mngt"`
	ErrorCode int `json:"error_code"`
}

// TimeRangeDelRequest 删除时间段请求
type TimeRangeDelRequest struct {
	Method   string         `json:"method"`
	TimeMngt TimeRangeDelIn `json:"time_mngt"`
}

type TimeRangeDelIn struct {
	Name string `json:"name"`
}

// ========== 带宽控制 (Qos) ==========

// QosSetting Qos配置
type QosSetting struct {
	Name             string   `json:".name"`
	Type             string   `json:".type"`
	QosEnable        string   `json:"qos_enable"`
	ThresholdEnable  string   `json:"threshold_enable"`
	QosThreshold     string   `json:"qos_threshold"`
	Interface        []string `json:"interface"`
}

// QosConfigResponse Qos配置响应
type QosConfigResponse struct {
	Qos struct {
		Setting QosSetting `json:"setting"`
	} `json:"qos"`
	ErrorCode int `json:"error_code"`
}

// QosConfigSetRequest 设置Qos配置请求
type QosConfigSetRequest struct {
	Method string               `json:"method"`
	Qos    map[string]interface{} `json:"qos"`
}

// QosRuleItem Qos规则条目
type QosRuleItem struct {
	DotName     string `json:"-"` // map key：内部名称，如 rule_xxx
	Name        string `json:"name"`
	IfPing      string `json:"if_ping"`
	IfPong      string `json:"if_pong"`
	IPGroup     string `json:"ip_group"`
	RateMax     string `json:"rate_max"`
	RateMaxMate string `json:"rate_max_mate"`
	Mode        string `json:"mode"` // share | priv
	Time        string `json:"time"`
	Comment     string `json:"comment"`
	Position    string `json:"position"`
	Enable      string `json:"enable"`
	IPType      string `json:"ip_type"`
}

// QosRuleListResponse Qos规则列表响应
type QosRuleListResponse struct {
	Qos struct {
		Rule  []map[string]QosRuleItem `json:"rule"`
		Count GenericCount             `json:"count"`
	} `json:"qos"`
	ErrorCode int `json:"error_code"`
}

// QosRuleAddRequest 添加Qos规则请求
type QosRuleAddRequest struct {
	Method string            `json:"method"`
	Qos    QosRuleAddIn      `json:"qos"`
}

type QosRuleAddIn struct {
	Table string                 `json:"table"`
	Para  map[string]interface{} `json:"para"`
}

// QosRuleAddResponse 添加Qos规则响应
type QosRuleAddResponse struct {
	Qos struct {
		Name []string `json:"name"`
	} `json:"qos"`
	ErrorCode int `json:"error_code"`
}

// QosRuleDelRequest 删除Qos规则请求
type QosRuleDelRequest struct {
	Method string       `json:"method"`
	Qos    QosRuleDelIn `json:"qos"`
}

type QosRuleDelIn struct {
	Name string `json:"name"`
}

// ========== 访问控制 (ACL) ==========

// ACLRuleItem ACL规则条目
type ACLRuleItem struct {
	DotName  string `json:"-"` // map key：内部名称，如 rule_acl_inner_xxx
	Name     string `json:"name"`
	Policy   string `json:"policy"` // DROP | ACCEPT
	Service  string `json:"service"`
	Zone     string `json:"zone"`
	Src      string `json:"src"`
	Dest     string `json:"dest"`
	Time     string `json:"time"`
	User     string `json:"user"`
	Position int    `json:"position"`
}

// ACLRuleListResponse ACL规则列表响应
type ACLRuleListResponse struct {
	AccessCtl struct {
		RuleACLInner []map[string]ACLRuleItem `json:"rule_acl_inner"`
		Count        GenericCount             `json:"count"`
	} `json:"access_ctl"`
	ErrorCode int `json:"error_code"`
}

// ACLServiceItem ACL服务条目
type ACLServiceItem struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Flag    string `json:"flag"` // system | inner
	Proto   string `json:"proto"`
	SPort   string `json:"sport"`
	DPort   string `json:"dport"`
	Ref     string `json:"ref"`
}

// ACLServiceListResponse ACL服务列表响应
type ACLServiceListResponse struct {
	Service struct {
		Service []map[string]ACLServiceItem `json:"service"`
		Count   GenericCount                `json:"count"`
	} `json:"service"`
	ErrorCode int `json:"error_code"`
}

// ACLServiceAddRequest 添加ACL服务请求
type ACLServiceAddRequest struct {
	Method  string            `json:"method"`
	Service ACLServiceAddIn   `json:"service"`
}

type ACLServiceAddIn struct {
	Table string                   `json:"table"`
	Para  []map[string]interface{} `json:"para"`
}

// ACLServiceAddResponse 添加ACL服务响应
type ACLServiceAddResponse struct {
	Service struct {
		Name []string `json:"name"`
	} `json:"service"`
	ErrorCode int `json:"error_code"`
}

// ACLRuleAddRequest 添加ACL规则请求
type ACLRuleAddRequest struct {
	Method    string            `json:"method"`
	AccessCtl ACLRuleAddIn      `json:"access_ctl"`
}

type ACLRuleAddIn struct {
	Table string                 `json:"table"`
	Para  map[string]interface{} `json:"para"`
}

// ACLRuleAddResponse 添加ACL规则响应
type ACLRuleAddResponse struct {
	AccessCtl struct {
		Name []string `json:"name"`
	} `json:"access_ctl"`
	ErrorCode int `json:"error_code"`
}

// ACLRuleDelRequest 删除ACL规则请求
type ACLRuleDelRequest struct {
	Method    string        `json:"method"`
	AccessCtl ACLRuleDelIn  `json:"access_ctl"`
}

type ACLRuleDelIn struct {
	Name string `json:"name"`
}

// ACLServiceDelRequest 删除ACL服务请求
type ACLServiceDelRequest struct {
	Method  string          `json:"method"`
	Service ACLServiceDelIn `json:"service"`
}

type ACLServiceDelIn struct {
	Name string `json:"name"`
}
