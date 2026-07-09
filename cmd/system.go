package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	systemCmd.AddCommand(systemInfoCmd)
	systemCmd.AddCommand(systemRebootCmd)
	rootCmd.AddCommand(systemCmd)
}

// ========== system ==========

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "系统工具",
}

// ========== system info ==========

var systemInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "查看设备信息",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		info, err := client.GetDeviceInfo()
		if err != nil {
			return err
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(info)
		case "yaml":
			fmt.Printf("device_model: %s\n", decodeURL(info.DeviceModel))
			fmt.Printf("device_name: %s\n", decodeURL(info.DeviceName))
			fmt.Printf("device_type: %s\n", info.DeviceType)
			fmt.Printf("fw_description: %s\n", info.FwDescription)
			fmt.Printf("hw_version: %s\n", info.HwVersion)
			fmt.Printf("sw_version: %s\n", decodeURL(info.SwVersion))
			fmt.Printf("manufacturer: %s\n", info.ManufacturerName)
			fmt.Printf("manufacturer_url: %s\n", info.ManufacturerURL)
			fmt.Printf("product_id: %s\n", info.ProductID)
			fmt.Printf("language: %s\n", info.Language)
			fmt.Printf("domain_name: %s\n", info.DomainName)
			fmt.Printf("radio_count: %s\n", info.RadioCount)
			fmt.Printf("sys_software_revision: %s\n", info.SysSoftwareRevision)
			fmt.Printf("sys_software_revision_minor: %s\n", info.SysSoftwareRevisionMinor)
			fmt.Printf("vendor_id: %s\n", info.VendorID)
			fmt.Printf("zone_code: %s\n", info.ZoneCode)
			fmt.Printf("enable_dns: %s\n", boolLabel(info.EnableDNS))
			return nil
		default:
			type Row struct {
				Field string `json:"field"`
				Value string `json:"value"`
			}

			rows := []Row{
				{"设备型号", decodeURL(info.DeviceModel)},
				{"设备名称", decodeURL(info.DeviceName)},
				{"设备类型", info.DeviceType},
				{"固件描述", info.FwDescription},
				{"硬件版本", info.HwVersion},
				{"软件版本", decodeURL(info.SwVersion)},
				{"制造商", info.ManufacturerName},
				{"制造商URL", info.ManufacturerURL},
				{"产品ID", info.ProductID},
				{"语言", info.Language},
				{"域名", info.DomainName},
				{"无线频段数", info.RadioCount},
				{"系统软件版本", info.SysSoftwareRevision},
				{"系统软件次版本", info.SysSoftwareRevisionMinor},
				{"厂商ID", info.VendorID},
				{"区域代码", info.ZoneCode},
				{"DNS状态", boolLabel(info.EnableDNS)},
			}

			headers := []string{"FIELD", "VALUE"}
			printTable(headers, rows, func(r Row) []string {
				return []string{r.Field, r.Value}
			})
			return nil
		}
	},
}

// ========== system reboot ==========

var systemRebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "重启设备",
	Long: `重启TP-Link路由器。

警告: 此操作会导致网络暂时中断，请谨慎使用。

示例:
  tplink system reboot`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		if err := client.RebootSystem(); err != nil {
			return err
		}
		fmt.Println("重启指令已发送，设备正在重启...")
		return nil
	},
}

// ========== 辅助函数 ==========

// boolLabel 将 1/0 转为 enabled/disabled
func boolLabel(s string) string {
	switch s {
	case "1":
		return "enabled"
	case "0":
		return "disabled"
	case "on":
		return "enabled"
	case "off":
		return "disabled"
	default:
		return s
	}
}

// decodeURL 解码 URL 编码的字符串
