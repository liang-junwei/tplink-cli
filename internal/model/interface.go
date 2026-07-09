package model

// ========== 接口模式 (ifmode) ==========

// IfModeGetRequest 查看接口模式请求
type IfModeGetRequest struct {
	Method  string            `json:"method"`
	Network map[string]string `json:"network"`
}

// IfModeGetResponse 查看接口模式响应
type IfModeGetResponse struct {
	Network struct {
		IfMode IfModeData `json:"if_mode"`
	} `json:"network"`
	ErrorCode int `json:"error_code"`
}

// IfModeData 接口模式数据
type IfModeData struct {
	WanMode   string `json:"wan_mode"`
	SingleWan string `json:"single_wan"`
}

// IfModeSetRequest 设置接口模式请求
type IfModeSetRequest struct {
	Method  string                     `json:"method"`
	Network map[string]*IfModeSetData  `json:"network"`
}

// IfModeSetData 设置接口模式数据
type IfModeSetData struct {
	WanMode string `json:"wan_mode"`
}

// ========== IPv6桥模式 (brv6mode) ==========

// BridgeV6GetRequest 查看IPv6桥模式请求
type BridgeV6GetRequest struct {
	Method  string            `json:"method"`
	Network map[string]string `json:"network"`
}

// BridgeV6GetResponse 查看IPv6桥模式响应
type BridgeV6GetResponse struct {
	Network struct {
		BridgeV6 BridgeV6Data `json:"bridge_v6"`
	} `json:"network"`
	ErrorCode int `json:"error_code"`
}

// BridgeV6Data IPv6桥模式数据
type BridgeV6Data struct {
	Enable  string   `json:"enable"`
	BindIf  []string `json:"bindif"`
	Type    string   `json:"type"`
	IfName  string   `json:"if_name"`
	Stp     string   `json:"stp,omitempty"`
}

// BridgeV6SetRequest 设置IPv6桥模式请求
type BridgeV6SetRequest struct {
	Method  string                   `json:"method"`
	Network map[string]*BridgeV6Data `json:"network"`
}

// ========== 有线接口状态 (port) ==========

// PortGetRequest 查看有线接口状态请求
type PortGetRequest struct {
	Method string            `json:"method"`
	Port   map[string]string `json:"port"`
}

// PortGetResponse 查看有线接口状态响应
type PortGetResponse struct {
	ErrorCode int       `json:"error_code"`
	Port      PortResult `json:"port"`
}

// PortResult 有线接口查询结果
type PortResult struct {
	Count map[string]int `json:"count"`
	Ports []PortItem     `json:"port"`
}

// PortItem 单个端口（key为port_N）
type PortItem map[string]PortData

// PortData 端口详细数据
type PortData struct {
	PortID            string `json:"port_id"`
	PortState         string `json:"port_state"`
	LinkSpeed         string `json:"link_speed"`
	LinkDuplex        string `json:"link_duplex"`
	Enable            string `json:"enable"`
	TxAll             string `json:"tx_all"`
	RxAll             string `json:"rx_all"`
	TxUnicast         string `json:"tx_unicast"`
	RxUnicast         string `json:"rx_unicast"`
	TxBroadcast       string `json:"tx_broadcast"`
	RxBroadcast       string `json:"rx_broadcast"`
	TxMulticast       string `json:"tx_multicast"`
	RxMulticast       string `json:"rx_multicast"`
	TxPause           string `json:"tx_pause"`
	RxPause           string `json:"rx_pause"`
	FlowcontrolEnable string `json:"flowcontrol_enable"`
	Pvid              string `json:"pvid"`
	CfgSpeed          string `json:"cfg_speed"`
	CfgDuplex         string `json:"cfg_duplex"`
}

// ========== WAN 配置 (wan) ==========

// WanListRequest 查看WAN配置请求
type WanListRequest struct {
	Method  string   `json:"method"`
	Network WanQuery `json:"network"`
}

// WanQuery WAN查询参数
type WanQuery struct {
	Table  string          `json:"table"`
	Filter []WanIfFilter   `json:"filter,omitempty"`
}

// WanIfFilter WAN接口过滤条件
type WanIfFilter struct {
	BaseName []string `json:"base_name,omitempty"`
	IfName   []string `json:"if_name,omitempty"`
}

// WanListResponse 查看WAN配置响应
type WanListResponse struct {
	Network struct {
		IfInfo []WanIfItem   `json:"if_info"`
		Count  map[string]int `json:"count"`
	} `json:"network"`
	ErrorCode int `json:"error_code"`
}

// WanIfItem 单个WAN接口（key为接口名）
type WanIfItem map[string]WanIfData

// WanIfData WAN接口详细数据
type WanIfData struct {
	IfName    string   `json:"if_name"`
	IfType    string   `json:"if_type"`
	Proto     string   `json:"proto"`
	IP        string   `json:"ip"`
	Netmask   string   `json:"netmask"`
	Gateway   string   `json:"gateway"`
	Mac       string   `json:"mac"`
	FactoryMac string  `json:"factory_mac"`
	PriDNS    string   `json:"pri_dns"`
	SndDNS    string   `json:"snd_dns"`
	LinkStatus string  `json:"link_status"`
	UpSpeed   int64    `json:"up_speed"`
	DownSpeed int64    `json:"down_speed"`
	Mtu       string   `json:"mtu"`
	Mtu6      string   `json:"mtu6"`
	Uptime    int      `json:"uptime"`
	BindIf    []string `json:"bindif"`
	VlanID    string   `json:"vlanid"`
	Untag     string   `json:"untag"`
	IsBridged int      `json:"isbridged"`
	IPv6Enable string  `json:"ipv6_enable"`
	IsSys     string   `json:"issys"`
	MngtEnable string  `json:"mngt_enable"`
	IsUserDNS string   `json:"is_user_dns"`
	// PPPoE 特有
	TPPPoEEnable string `json:"t_pppoe_enable,omitempty"`
	LinkType     string `json:"linktype,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	Service      string `json:"service,omitempty"`
	DupIPCheck   bool   `json:"dup_ip_check,omitempty"`
	// DHCP 特有
	Hostname string `json:"hostname,omitempty"`
}

// WanSetRequest 设置WAN配置请求
type WanSetRequest struct {
	Method  string  `json:"method"`
	Network WanSetParams `json:"network"`
}

// WanSetParams 设置WAN参数
type WanSetParams struct {
	Table  string        `json:"table"`
	Para   WanIfData     `json:"para"`
	Filter []WanIfFilter `json:"filter"`
}

// ========== LAN 接口 (lan) ==========

// LanGetRequest 查看LAN请求
type LanGetRequest struct {
	Method  string            `json:"method"`
	Network map[string]string `json:"network"`
}

// LanGetResponse 查看LAN响应
type LanGetResponse struct {
	Network struct {
		Lan LanData `json:"lan"`
	} `json:"network"`
	ErrorCode int `json:"error_code"`
}

// LanData LAN接口数据
type LanData struct {
	Proto     string `json:"proto"`
	IPAddr    string `json:"ipaddr"`
	Netmask   string `json:"netmask"`
	IPMode    string `json:"ip_mode"`
	FacNetmask string `json:"fac_netmask"`
	Type      string `json:"type"`
	MacAddr   string `json:"macaddr"`
	FacIPAddr string `json:"fac_ipaddr"`
	IfName    string `json:"ifname"`
	IPv6Enable string `json:"ipv6_enable,omitempty"`
}

// LanSetRequest 设置LAN请求
type LanSetRequest struct {
	Method  string               `json:"method"`
	Network map[string]*LanData  `json:"network"`
}

// ========== DHCP (dhcp) ==========

// DhcpConfigGetRequest 查看DHCP配置请求
type DhcpConfigGetRequest struct {
	Method string            `json:"method"`
	Dhcpd  map[string]string `json:"dhcpd"`
}

// DhcpConfigGetResponse 查看DHCP配置响应
type DhcpConfigGetResponse struct {
	Dhcpd struct {
		Lan DhcpConfigData `json:"lan"`
	} `json:"dhcpd"`
	ErrorCode int `json:"error_code"`
}

// DhcpConfigData DHCP配置数据
type DhcpConfigData struct {
	Enable    string `json:"enable"`
	PoolStart string `json:"pool_start"`
	PoolEnd   string `json:"pool_end"`
	LeaseTime string `json:"lease_time"`
	Gateway   string `json:"gateway,omitempty"`
	Domain    string `json:"domain,omitempty"`
	PriDNS    string `json:"pri_dns,omitempty"`
	SndDNS    string `json:"snd_dns,omitempty"`
	Option60  string `json:"option60,omitempty"`
	Option138 string `json:"option138,omitempty"`
	Interface string `json:"interface,omitempty"`
	AddrType  string `json:"addrtype,omitempty"`
}

// DhcpConfigSetRequest 设置DHCP配置请求
type DhcpConfigSetRequest struct {
	Method string                  `json:"method"`
	Dhcpd  map[string]*DhcpConfigData `json:"dhcpd"`
}

// DhcpClientListRequest 查看DHCP客户端列表请求
type DhcpClientListRequest struct {
	Method string       `json:"method"`
	Dhcpd  DhcpTableReq `json:"dhcpd"`
}

// DhcpTableReq DHCP表查询参数
type DhcpTableReq struct {
	Table string          `json:"table"`
	Para  DhcpRangeParams `json:"para,omitempty"`
}

// DhcpRangeParams 范围查询参数
type DhcpRangeParams struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// DhcpClientListResponse 查看DHCP客户端列表响应
type DhcpClientListResponse struct {
	Dhcpd struct {
		DhcpClients []DhcpClientItem `json:"dhcp_clients"`
		Count       map[string]int   `json:"count,omitempty"`
	} `json:"dhcpd"`
	ErrorCode int `json:"error_code"`
}

// DhcpClientItem 单个DHCP客户端（key为dhcp_client_N）
type DhcpClientItem map[string]DhcpClientData

// DhcpClientData DHCP客户端数据
type DhcpClientData struct {
	Expires   string `json:"expires"`
	IPAddr    string `json:"ipaddr"`
	Hostname  string `json:"hostname"`
	MacAddr   string `json:"macaddr"`
	Interface string `json:"interface"`
}

// DhcpStaticListRequest 查看静态地址分配请求
type DhcpStaticListRequest struct {
	Method string       `json:"method"`
	Dhcpd  DhcpTableReq `json:"dhcpd"`
}

// DhcpStaticListResponse 查看静态地址分配响应
type DhcpStaticListResponse struct {
	Dhcpd struct {
		DhcpStatic []DhcpStaticItem `json:"dhcp_static"`
		Count      map[string]int   `json:"count,omitempty"`
	} `json:"dhcpd"`
	ErrorCode int `json:"error_code"`
}

// DhcpStaticItem 单条静态地址分配（key为dhcp_static_N）
type DhcpStaticItem map[string]DhcpStaticData

// DhcpStaticData 静态地址分配数据
type DhcpStaticData struct {
	Mac           string `json:"mac"`
	IP            string `json:"ip"`
	Name          string `json:"name"`
	Note          string `json:"note"`
	Enable        string `json:"enable"`
	DhcpStaticID  string `json:"dhcp_static_id"`
}

// DhcpStaticAddRequest 添加静态地址分配请求
type DhcpStaticAddRequest struct {
	Method string           `json:"method"`
	Dhcpd  DhcpStaticAddIn  `json:"dhcpd"`
}

// DhcpStaticAddIn 添加静态地址分配输入
type DhcpStaticAddIn struct {
	Table string          `json:"table"`
	Para  DhcpStaticPara  `json:"para"`
}

// DhcpStaticPara 静态地址分配参数
type DhcpStaticPara struct {
	Mac    string `json:"mac"`
	IP     string `json:"ip"`
	Note   string `json:"note"`
	Enable string `json:"enable"`
}

// DhcpStaticDelRequest 删除静态地址分配请求
type DhcpStaticDelRequest struct {
	Method string          `json:"method"`
	Dhcpd  DhcpStaticDelIn `json:"dhcpd"`
}

// DhcpStaticDelIn 删除静态地址分配输入
type DhcpStaticDelIn struct {
	Table  string              `json:"table"`
	Filter []DhcpStaticIDFilter `json:"filter"`
}

// DhcpStaticIDFilter 静态地址分配ID过滤
type DhcpStaticIDFilter struct {
	DhcpStaticID string `json:"dhcp_static_id"`
}

// NameResult 添加操作返回的名称结果
type NameResult struct {
	Name      []string `json:"name"`
	ErrorCode int      `json:"error_code"`
}

// DhcpStaticAddResponse 添加静态地址响应
type DhcpStaticAddResponse struct {
	Dhcpd struct {
		Name []string `json:"name"`
	} `json:"dhcpd"`
	ErrorCode int `json:"error_code"`
}
