package model

// ========== wireless config (wlan_host_2g / wlan_host_5g) ==========

// WlanHostConfig 无线主机配置（2.4G/5G 共有字段）
type WlanHostConfig struct {
	Enable        string `json:"enable"`
	SSID          string `json:"ssid"`
	SSIDCodeType  string `json:"ssid_code_type"`
	SSIDBrd       string `json:"ssidbrd"`
	Isolate       string `json:"isolate"`
	Encryption    string `json:"encryption,omitempty"`
	Auth          string `json:"auth,omitempty"`
	Cipher        string `json:"cipher,omitempty"`
	Key           string `json:"key,omitempty"`
	KeyUpdateIntv string `json:"key_update_intv,omitempty"`
	Bandwidth     string `json:"bandwidth,omitempty"`
	Channel       string `json:"channel,omitempty"`
	Power         string `json:"power,omitempty"`
	Mode          string `json:"mode,omitempty"`
	WMM           string `json:"wmm,omitempty"`
	BeaconIntv    string `json:"beacon_intv,omitempty"`
	RadioIsolate  string `json:"radio_isolate,omitempty"`
}

// WirelessConfigRequest 查看无线配置请求
type WirelessConfigRequest struct {
	Method   string             `json:"method"`
	Wireless WirelessConfigName `json:"wireless"`
}

// WirelessConfigName 请求中指定要获取的配置名
type WirelessConfigName struct {
	Name []string `json:"name"`
}

// WirelessConfigResponse 查看无线配置响应
type WirelessConfigResponse struct {
	Wireless  WirelessConfigResult `json:"wireless"`
	ErrorCode int                  `json:"error_code"`
}

// WirelessConfigResult 无线配置结果
type WirelessConfigResult struct {
	WlanHost2G *WlanHostConfig `json:"wlan_host_2g"`
	WlanHost5G *WlanHostConfig `json:"wlan_host_5g"`
}

// WirelessSetRequest 设置无线配置请求
type WirelessSetRequest struct {
	Method   string                 `json:"method"`
	Wireless map[string]interface{} `json:"wireless"`
}

// ========== guest network ==========

// GuestNetworkConfig 访客网络配置
type GuestNetworkConfig struct {
	Enable        string `json:"enable"`
	SSID          string `json:"ssid"`
	SSIDCodeType  string `json:"ssid_code_type"`
	Encrypt       string `json:"encrypt"`
	Key           string `json:"key,omitempty"`
	Upload        string `json:"upload"`
	Download      string `json:"download"`
	RadioMaxSta   string `json:"radio_max_sta"`
}

// GuestNetworkRequest 查看访客网络请求
type GuestNetworkRequest struct {
	Method       string            `json:"method"`
	GuestNetwork GuestNetworkName  `json:"guest_network"`
}

// GuestNetworkName 请求中指定获取的访客网络名
type GuestNetworkName struct {
	Name string `json:"name"`
}

// GuestNetworkResponse 查看访客网络响应
type GuestNetworkResponse struct {
	GuestNetwork GuestNetworkResult `json:"guest_network"`
	ErrorCode    int                `json:"error_code"`
}

// GuestNetworkResult 访客网络结果
type GuestNetworkResult struct {
	Guest2G GuestNetworkConfig `json:"guest_2g"`
}

// GuestNetworkSetRequest 设置访客网络请求
type GuestNetworkSetRequest struct {
	Method       string                       `json:"method"`
	GuestNetwork map[string]interface{}       `json:"guest_network"`
}

// ========== wlan access config ==========

// WlanAccessConfig MAC地址过滤配置
type WlanAccessConfig struct {
	Name      string   `json:".name"`
	Type      string   `json:".type"`
	SSIDList  []string `json:"ssid_list"`
	Enable    string   `json:"enable"`
	Mode      string   `json:"mode"` // 1=白名单, 2=黑名单
	InitFlag  string   `json:"init_flag,omitempty"`
	Anonymous bool     `json:".anonymous"`
}

// WlanAccessConfigRequest 查看MAC过滤配置请求
type WlanAccessConfigRequest struct {
	Method     string             `json:"method"`
	WlanAccess WlanAccessNameReq  `json:"wlan_access"`
}

// WlanAccessNameReq MAC过滤配置名请求
type WlanAccessNameReq struct {
	Name string `json:"name"`
}

// WlanAccessConfigResponse 查看MAC过滤配置响应
type WlanAccessConfigResponse struct {
	WlanAccess WlanAccessConfigResult `json:"wlan_access"`
	ErrorCode  int                    `json:"error_code"`
}

// WlanAccessConfigResult MAC过滤配置结果
type WlanAccessConfigResult struct {
	Config WlanAccessConfig `json:"config"`
}

// WlanAccessConfigSetRequest 设置MAC过滤配置请求
type WlanAccessConfigSetRequest struct {
	Method     string           `json:"method"`
	WlanAccess map[string]interface{} `json:"wlan_access"`
}

// ========== wlan access white/black list ==========

// WlanAccessListItem 白/黑名单条目（key为 white_list_N / black_list_N）
type WlanAccessListItem map[string]WlanAccessItemData

// WlanAccessItemData 白/黑名单条目数据
type WlanAccessItemData struct {
	Name      string `json:".name"`
	Type      string `json:".type"`
	Mac       string `json:"mac"`
	NameField string `json:"name"` // 用户设置的名称
	Anonymous bool   `json:".anonymous"`
	Index     int    `json:".index"`
}

// WlanAccessListRequest 查看白/黑名单请求
type WlanAccessListRequest struct {
	Method     string              `json:"method"`
	WlanAccess WlanAccessTableReq  `json:"wlan_access"`
}

// WlanAccessTableReq 带 table 的请求
type WlanAccessTableReq struct {
	Table string `json:"table"`
}

// WlanAccessListResponse 查看白/黑名单响应
type WlanAccessListResponse struct {
	WlanAccess WlanAccessListResult `json:"wlan_access"`
	ErrorCode  int                  `json:"error_code"`
}

// WlanAccessListResult 白/黑名单结果
type WlanAccessListResult struct {
	WhiteList []WlanAccessListItem `json:"white_list,omitempty"`
	BlackList []WlanAccessListItem `json:"black_list,omitempty"`
	Count     map[string]int       `json:"count,omitempty"`
}

// WlanAccessAddRequest 添加白/黑名单请求
type WlanAccessAddRequest struct {
	Method     string             `json:"method"`
	WlanAccess WlanAccessAddIn    `json:"wlan_access"`
}

// WlanAccessAddIn 添加白/黑名单输入
type WlanAccessAddIn struct {
	Table string                 `json:"table"`
	Para  map[string]interface{} `json:"para"`
}

// WlanAccessAddResponse 添加白/黑名单响应
type WlanAccessAddResponse struct {
	WlanAccess struct {
		Name []string `json:"name"`
	} `json:"wlan_access"`
	ErrorCode int `json:"error_code"`
}

// WlanAccessDelRequest 删除白/黑名单请求
type WlanAccessDelRequest struct {
	Method     string           `json:"method"`
	WlanAccess WlanAccessDelIn  `json:"wlan_access"`
}

// WlanAccessDelIn 删除白/黑名单输入
type WlanAccessDelIn struct {
	Name string `json:"name"` // 完整key名如 white_list_1782875167
}

// ========== wlan service ==========

// WlanServItem 无线服务条目（key为 wlan_serv_N）
type WlanServItem map[string]WlanServData

// WlanServData 无线服务数据
type WlanServData struct {
	ServID       string `json:"serv_id"`
	RadioID      string `json:"radio_id"`
	SSID         string `json:"ssid"`
	SSIDCodeType string `json:"ssid_code_type"`
	SSIDBrd      string `json:"ssidbrd"`
	Enable       string `json:"enable"`
	Isolate      string `json:"isolate"`
	Encryption   string `json:"encryption"`
	Auth         string `json:"auth,omitempty"`
	Cipher       string `json:"cipher,omitempty"`
	Key          string `json:"key,omitempty"`
	KeyUpdateIntv string `json:"key_update_intv,omitempty"`
	NetworkType  string `json:"network_type"`
}

// WlanServRequest 查看服务列表请求
type WlanServRequest struct {
	Method   string         `json:"method"`
	Wireless WlanServTable  `json:"wireless"`
}

// WlanServTable 服务表请求
type WlanServTable struct {
	Table string `json:"table"`
}

// WlanServResponse 查看服务列表响应
type WlanServResponse struct {
	Wireless  WlanServResult `json:"wireless"`
	ErrorCode int            `json:"error_code"`
}

// WlanServResult 服务列表结果
type WlanServResult struct {
	WlanServ []WlanServItem  `json:"wlan_serv"`
	Count    map[string]int  `json:"count,omitempty"`
}

// ========== wireless client (sta_list) ==========

// WirelessClientItem 无线客户端条目（key为 sta_list_N）
type WirelessClientItem map[string]WirelessClientData

// WirelessClientData 无线客户端数据
type WirelessClientData struct {
	Mac          string `json:"mac"`
	ServID       string `json:"serv_id"`
	RadioID      string `json:"radio_id"`
	Name         string `json:"name"`
	RSSI         string `json:"rssi"`
	RxFlow       string `json:"rx_flow"`
	IP           string `json:"ip"`
	TxRate       string `json:"tx_rate"`
	SSIDCodeType string `json:"ssid_code_type"`
	TxFlow       string `json:"tx_flow"`
	EntryID      int    `json:"entry_id"`
	Status       string `json:"status"`
	SSID         string `json:"ssid"`
	RxRate       string `json:"rx_rate"`
	NetworkType  string `json:"network_type"`
}

// WirelessClientRequest 查看无线客户端请求
type WirelessClientRequest struct {
	Method   string                  `json:"method"`
	Wireless WirelessClientTableReq  `json:"wireless"`
}

// WirelessClientTableReq 无线客户端表请求
type WirelessClientTableReq struct {
	Table  string                       `json:"table"`
	Filter *WirelessClientFilter        `json:"filter,omitempty"`
}

// WirelessClientFilter 无线客户端过滤
type WirelessClientFilter struct {
	RadioID string `json:"radio_id,omitempty"`
	ServID  string `json:"serv_id,omitempty"`
}

// WirelessClientResponse 查看无线客户端响应
type WirelessClientResponse struct {
	Wireless  WirelessClientResult `json:"wireless"`
	ErrorCode int                  `json:"error_code"`
}

// WirelessClientResult 无线客户端结果
type WirelessClientResult struct {
	StaList []WirelessClientItem `json:"sta_list"`
	Count   map[string]int       `json:"count,omitempty"`
}
