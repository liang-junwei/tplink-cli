package model

// ========== VPN (vpn_wan) ==========

// VpnConfig VPN配置条目
type VpnConfig struct {
	DotName      string `json:".name"`
	Protocol     string `json:"protocol"`     // auto
	RouteMode    string `json:"route_mode"`   // all | manual
	RemoteSubnet string `json:"remotesubnet"` // URL编码, 仅在route_mode=manual时有效
	Username     string `json:"username"`
	Server       string `json:"server"`
	Interface    string `json:"interface"`   // WAN1
	Password     string `json:"password"`
	ForwardMode  string `json:"forward_mode"` // nat | route
	ConnectMode  string `json:"connect_mode"` // auto | manual
}

// VpnListResponse VPN列表响应
type VpnListResponse struct {
	ErrorCode int         `json:"error_code"`
	Vpn       VpnListData `json:"vpn"`
	TimeMngt  map[string]any `json:"time_mngt"` // _vpn_wan_1: []
}

// VpnListData VPN列表数据
type VpnListData struct {
	Count  VpnCount               `json:"count"`
	VpnWan []map[string]VpnConfig `json:"vpn_wan"`
}

// VpnCount VPN计数
type VpnCount struct {
	VpnWan int `json:"vpn_wan"`
}
