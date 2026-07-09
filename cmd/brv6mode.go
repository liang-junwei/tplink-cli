package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var brv6modeCmd = &cobra.Command{
	Use:   "brv6mode",
	Short: "IPv6桥模式管理",
	Long:  `查看和设置 TP-Link 路由器的 IPv6 桥模式配置。`,
}

// brv6mode list
var brv6modeListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看IPv6桥模式",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		data, err := client.GetBridgeV6()
		if err != nil {
			return err
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(data)
		case "yaml":
			fmt.Printf("enable: %s\nbindif: %s\nif_name: %s\ntype: %s\n",
				data.Enable, strings.Join(data.BindIf, ","), data.IfName, data.Type)
			return nil
		default:
			fmt.Printf("状态:     %s\n", statusLabel(data.Enable))
			fmt.Printf("绑定接口: %s\n", strings.Join(data.BindIf, ", "))
			fmt.Printf("接口名称: %s\n", data.IfName)
			fmt.Printf("类型:     %s\n", data.Type)
			return nil
		}
	},
}

// brv6mode set
var brv6modeSetCmd = &cobra.Command{
	Use:   "set",
	Short: "设置IPv6桥模式",
	Long:  `更新 IPv6 桥模式配置。需要指定 --enable、--bindif 等参数。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		// 先获取当前配置作为基础
		current, err := client.GetBridgeV6()
		if err != nil {
			return err
		}

		// 合并用户指定的字段
		if v, _ := cmd.Flags().GetString("enable"); v != "" {
			current.Enable = v
		}
		if v, _ := cmd.Flags().GetString("bindif"); v != "" {
			current.BindIf = strings.Split(v, ",")
		}
		if v, _ := cmd.Flags().GetString("stp"); v != "" {
			current.Stp = v
		}

		if current.IfName == "" {
			current.IfName = "bridge_v6"
		}
		current.Type = "bridgev6"

		if err := client.SetBridgeV6(current); err != nil {
			return err
		}

		fmt.Println("IPv6桥模式设置成功")
		return nil
	},
}

func init() {
	brv6modeSetCmd.Flags().String("enable", "", "启用状态: on|off")
	brv6modeSetCmd.Flags().String("bindif", "", "绑定接口, 逗号分隔, 如: wan1_eth,lan")
	brv6modeSetCmd.Flags().String("stp", "", "STP: 0|1")

	brv6modeCmd.AddCommand(brv6modeListCmd)
	brv6modeCmd.AddCommand(brv6modeSetCmd)
	rootCmd.AddCommand(brv6modeCmd)
}
