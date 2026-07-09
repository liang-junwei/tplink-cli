package model

// ========== 设备信息 ==========

// DeviceInfo 设备信息
type DeviceInfo struct {
	DotName                      string      `json:".name"`
	RadioCount                   string      `json:"radio_count"`
	ZoneCode                     string      `json:"zone_code"`
	ManufacturerName             string      `json:"manufacturer_name"`
	SwVersion                    string      `json:"sw_version"`
	ManufacturerURL              string      `json:"manufacturer_url"`
	Language                     string      `json:"language"`
	DomainName                   string      `json:"domain_name"`
	DotType                      string      `json:".type"`
	SysSoftwareRevision          string      `json:"sys_software_revision"`
	ProductID                    string      `json:"product_id"`
	FwDescription                string      `json:"fw_description"`
	HwVersion                    string      `json:"hw_version"`
	DotAnonymous                 bool        `json:".anonymous"`
	DeviceName                   string      `json:"device_name"`
	VendorID                     string      `json:"vendor_id"`
	DeviceModel                  string      `json:"device_model"`
	EnableDNS                    string      `json:"enable_dns"`
	DeviceInfoField              string      `json:"device_info"`
	SysSoftwareRevisionMinor     string      `json:"sys_software_revision_minor"`
	DeviceType                   string      `json:"device_type"`
}

// DeviceInfoResponse 设备信息响应
type DeviceInfoResponse struct {
	DeviceInfo DeviceInfoData `json:"device_info"`
	ErrorCode  int            `json:"error_code"`
}

// DeviceInfoData 设备信息数据
type DeviceInfoData struct {
	Info DeviceInfo `json:"info"`
}

// ========== 重启 ==========

// SystemRebootResponse 重启响应（通常无body）
type SystemRebootResponse struct {
	ErrorCode int `json:"error_code"`
}
