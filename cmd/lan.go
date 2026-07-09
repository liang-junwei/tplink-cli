package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var lanCmd = &cobra.Command{
	Use:   "lan",
	Short: "LAN接口管理",
	Long:  `查看和修改 TP-Link 路由器 LAN 接口的网络配置。`,
}

// lan list
var lanListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看LAN接口信息",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		data, err := client.GetLanInfo()
		if err != nil {
			return err
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(data)
		case "yaml":
			fmt.Printf("proto: %s\nipaddr: %s\nnetmask: %s\nip_mode: %s\ntype: %s\nmacaddr: %s\nifname: %s\n",
				data.Proto, data.IPAddr, data.Netmask, data.IPMode, data.Type, data.MacAddr, data.IfName)
			return nil
		default:
			fmt.Printf("协议:     %s\n", protoLabel(data.Proto))
			fmt.Printf("IP 地址:  %s\n", data.IPAddr)
			fmt.Printf("子网掩码: %s\n", data.Netmask)
			fmt.Printf("IP 模式:  %s\n", data.IPMode)
			fmt.Printf("接口类型: %s\n", data.Type)
			fmt.Printf("MAC 地址: %s\n", data.MacAddr)
			fmt.Printf("接口名称: %s\n", data.IfName)
			return nil
		}
	},
}

// lan set
var lanSetCmd = &cobra.Command{
	Use:   "set",
	Short: "修改LAN接口信息",
	Long:  `修改 LAN 接口的 IP 地址和子网掩码等配置。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		// 先获取当前配置
		current, err := client.GetLanInfo()
		if err != nil {
			return err
		}

		// 合并用户指定的字段
		if v, _ := cmd.Flags().GetString("ip"); v != "" {
			current.IPAddr = v
		}
		if v, _ := cmd.Flags().GetString("netmask"); v != "" {
			current.Netmask = v
		}
		if v, _ := cmd.Flags().GetString("ip-mode"); v != "" {
			current.IPMode = v
		}

		if err := client.SetLan(current); err != nil {
			return err
		}

		fmt.Printf("LAN接口配置成功, IP: %s/%s\n", current.IPAddr, current.Netmask)
		return nil
	},
}

func init() {
	lanSetCmd.Flags().String("ip", "", "LAN IP 地址")
	lanSetCmd.Flags().String("netmask", "", "子网掩码")
	lanSetCmd.Flags().String("ip-mode", "", "IP模式: manual|auto")

	lanCmd.AddCommand(lanListCmd)
	lanCmd.AddCommand(lanSetCmd)
	rootCmd.AddCommand(lanCmd)
}
