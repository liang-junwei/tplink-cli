package model

// ========== ARP 防护 ==========

// ArpConfig ARP防护全局配置
type ArpConfig struct {
	DotName   string   `json:".name"`
	Type      string   `json:".type"`
	Garp      string   `json:"garp"`
	Enable    string   `json:"enable"`
	LogEnable string   `json:"log_enable"`
	ImbPass   string   `json:"imb_pass"`
	Interval  string   `json:"interval"`
	Interface []string `json:"interface"`
}

// ArpConfigResponse ARP配置查询响应
type ArpConfigResponse struct {
	ArpDefense struct {
		Global ArpConfig `json:"global"`
	} `json:"arp_defense"`
	ErrorCode int `json:"error_code"`
}

// ArpConfigSetRequest ARP配置设置请求
type ArpConfigSetRequest struct {
	Method     string `json:"method"`
	ArpDefense struct {
		Global map[string]interface{} `json:"global"`
	} `json:"arp_defense"`
}

// ArpBindItem ARP绑定条目
type ArpBindItem struct {
	DotName   string `json:"-"`
	Mac       string `json:"mac"`
	Hostname  string `json:"hostname"`
	Status    string `json:"status"`
	Interface string `json:"interface"`
	IP        string `json:"ip"`
}

// ArpBindListResponse ARP绑定列表响应
type ArpBindListResponse struct {
	IPMacBind struct {
		SysArp []map[string]ArpBindItem `json:"sys_arp"`
		Count  struct {
			SysArp int `json:"sys_arp"`
		} `json:"count"`
	} `json:"ip_mac_bind"`
	ErrorCode int `json:"error_code"`
}

// ========== MAC 地址过滤 ==========

// MacFilterConfig MAC过滤全局配置
type MacFilterConfig struct {
	DotName    string `json:".name"`
	Type       string `json:".type"`
	Interfaces string `json:"interfaces"`
	Enable     string `json:"enable"`
	FilterMode string `json:"filter_mode"` // black | white
}

// MacFilterConfigResponse MAC过滤配置查询响应
type MacFilterConfigResponse struct {
	MacFilter struct {
		Global MacFilterConfig `json:"global"`
	} `json:"mac_filter"`
	ErrorCode int `json:"error_code"`
}

// MacFilterConfigSetRequest MAC过滤配置设置请求
type MacFilterConfigSetRequest struct {
	Method    string `json:"method"`
	MacFilter struct {
		Global map[string]interface{} `json:"global"`
	} `json:"mac_filter"`
}

// MacFilterRuleItem MAC过滤规则条目
type MacFilterRuleItem struct {
	DotName string `json:"-"`
	Name    string `json:"name"`
	Mac     string `json:"mac"`
}

// MacFilterRuleListResponse MAC过滤规则列表响应
type MacFilterRuleListResponse struct {
	MacFilter struct {
		MacFilterList []map[string]MacFilterRuleItem `json:"mac_filter_list"`
		Count         struct {
			MacFilterList int `json:"mac_filter_list"`
		} `json:"count"`
	} `json:"mac_filter"`
	ErrorCode int `json:"error_code"`
}

// MacFilterRuleAddRequest MAC过滤规则添加请求
type MacFilterRuleAddRequest struct {
	Method    string `json:"method"`
	MacFilter struct {
		Table string                 `json:"table"`
		Para  map[string]interface{} `json:"para"`
	} `json:"mac_filter"`
}

// MacFilterRuleAddResponse MAC过滤规则添加响应
type MacFilterRuleAddResponse struct {
	MacFilter struct {
		Name []string `json:"name"`
	} `json:"mac_filter"`
	ErrorCode int `json:"error_code"`
}

// MacFilterRuleDelRequest MAC过滤规则删除请求
type MacFilterRuleDelRequest struct {
	Method    string `json:"method"`
	MacFilter struct {
		Name string `json:"name"`
	} `json:"mac_filter"`
}

// ========== DoS 攻击防护 ==========

// DosConfig DoS攻击防护配置
type DosConfig struct {
	DotName          string `json:".name"`
	Type             string `json:".type"`
	TcpWinNuke       string `json:"tcp_winnuke"`
	IpoptTimestamp   string `json:"ipopt_timestamp"`
	IpoptStream      string `json:"ipopt_stream"`
	PingDeath        string `json:"ping_death"`
	IpoptNoop        string `json:"ipopt_noop"`
	TcpNoflag        string `json:"tcp_noflag"`
	PingLarge        string `json:"ping_large"`
	IpoptSecure      string `json:"ipopt_secure"`
	IpOption         string `json:"ip_option"`
	TcpFinSyn        string `json:"tcp_fin_syn"`
	IpoptRecordRoute string `json:"ipopt_record_route"`
	IpFrag           string `json:"ip_frag"`
	IpoptLooseRoute  string `json:"ipopt_loose_route"`
	TcpFinNoack      string `json:"tcp_fin_noack"`
	IpoptStrictRoute string `json:"ipopt_strict_route"`
}

// DosConfigResponse DoS配置查询响应
type DosConfigResponse struct {
	DosDefense struct {
		Global DosConfig `json:"global"`
	} `json:"dos_defense"`
	ErrorCode int `json:"error_code"`
}

// DosConfigSetRequest DoS配置设置请求
type DosConfigSetRequest struct {
	Method     string `json:"method"`
	DosDefense struct {
		Global map[string]interface{} `json:"global"`
	} `json:"dos_defense"`
}

// ========== Flood 攻击防护 ==========

// FloodGlobal Flood全局开关配置
type FloodGlobal struct {
	DotName    string `json:".name"`
	Type       string `json:".type"`
	UdpConnEn  string `json:"udp_conn_en"`
	UdpSrcEn   string `json:"udp_src_en"`
	IcmpConnEn string `json:"icmp_conn_en"`
	IcmpSrcEn  string `json:"icmp_src_en"`
	TcpSrcEn   string `json:"tcp_src_en"`
	TcpConnEn  string `json:"tcp_conn_en"`
}

// FloodThreshold Flood阈值配置
type FloodThreshold struct {
	DotName     string `json:".name"`
	Type        string `json:".type"`
	TcpConnBst  string `json:"tcp_conn_bst"`
	TcpConnLim  string `json:"tcp_conn_lim"`
	UdpConnLim  string `json:"udp_conn_lim"`
	UdpConnBst  string `json:"udp_conn_bst"`
	IcmpConnLim string `json:"icmp_conn_lim"`
	IcmpConnBst string `json:"icmp_conn_bst"`
	TcpSrcLim   string `json:"tcp_src_lim"`
	TcpSrcBst   string `json:"tcp_src_bst"`
	UdpSrcBst   string `json:"udp_src_bst"`
	UdpSrcLim   string `json:"udp_src_lim"`
	IcmpSrcBst  string `json:"icmp_src_bst"`
	IcmpSrcLim  string `json:"icmp_src_lim"`
}

// FloodConfigResponse Flood配置查询响应
type FloodConfigResponse struct {
	FloodDefense struct {
		Global    FloodGlobal    `json:"global"`
		Threshold FloodThreshold `json:"threshold"`
	} `json:"flood_defense"`
	ErrorCode int `json:"error_code"`
}

// FloodConfigSetRequest Flood配置设置请求
type FloodConfigSetRequest struct {
	Method       string `json:"method"`
	FloodDefense struct {
		Global    map[string]interface{} `json:"global"`
		Threshold map[string]interface{} `json:"threshold"`
	} `json:"flood_defense"`
}
