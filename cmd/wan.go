package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ljw/tplink-cli/internal/model"
	"github.com/spf13/cobra"
)

var wanCmd = &cobra.Command{
	Use:   "wan",
	Short: "WAN接口管理",
	Long:  `查看和配置 TP-Link 路由器的 WAN 接口（支持 static/dhcp/pppoe 三种模式）。`,
}

// wan list
var wanListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看WAN接口配置",
	Long:  `列出所有 WAN 接口的配置信息，包括 IP、子网、网关、DNS、连接状态等。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		baseName, _ := cmd.Flags().GetString("base-name")
		items, err := client.GetWanInfo(baseName)
		if err != nil {
			return err
		}

		type WanRow struct {
			IfName     string `json:"if_name"`
			Proto      string `json:"proto"`
			IP         string `json:"ip"`
			Netmask    string `json:"netmask"`
			Gateway    string `json:"gateway"`
			PriDNS     string `json:"pri_dns"`
			LinkStatus string `json:"link_status"`
		}

		rows := make([]WanRow, 0, len(items))
		for _, item := range items {
			for _, data := range item {
				rows = append(rows, WanRow{
					IfName:     data.IfName,
					Proto:      protoLabel(data.Proto),
					IP:         data.IP,
					Netmask:    data.Netmask,
					Gateway:    data.Gateway,
					PriDNS:     data.PriDNS,
					LinkStatus: linkLabel(data.LinkStatus),
				})
			}
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- if_name: %s\n  proto: %s\n  ip: %s\n  netmask: %s\n  gateway: %s\n  pri_dns: %s\n  link: %s\n",
					r.IfName, r.Proto, r.IP, r.Netmask, r.Gateway, r.PriDNS, r.LinkStatus)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("没有找到WAN接口信息")
				return nil
			}
			headers := []string{"IF_NAME", "PROTO", "IP", "NETMASK", "GATEWAY", "DNS", "LINK"}
			printTable(headers, rows, func(r WanRow) []string {
				return []string{r.IfName, r.Proto, r.IP, r.Netmask, r.Gateway, r.PriDNS, r.LinkStatus}
			})
			return nil
		}
	},
}

// wan set
var wanSetCmd = &cobra.Command{
	Use:   "set",
	Short: "设置WAN接口",
	Long: `修改 WAN 接口的连接方式。支持 static、dhcp、pppoe 三种模式。

示例:
  tplink wan set --if-name wan1_eth --mode static --ip 192.168.1.3 --netmask 255.255.255.0 --gateway 192.168.1.1
  tplink wan set --if-name wan1_eth --mode dhcp
  tplink wan set --if-name wan1_pppoe --mode pppoe --username admin --password 123456`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ifName, _ := cmd.Flags().GetString("if-name")
		mode, _ := cmd.Flags().GetString("mode")

		if ifName == "" || mode == "" {
			return fmt.Errorf("--if-name 和 --mode 为必选参数")
		}

		client, err := newAPIClient()
		if err != nil {
			return err
		}

		// 先获取当前配置
		items, err := client.GetWanInfo("")
		if err != nil {
			return err
		}

		// 查找目标接口的当前配置作为基础
		var current *model.WanIfData
		for _, item := range items {
			for name, data := range item {
				if name == ifName || data.IfName == ifName {
					d := data
					current = &d
					break
				}
			}
			if current != nil {
				break
			}
		}
		if current == nil {
			return fmt.Errorf("未找到接口: %s", ifName)
		}

		para := *current
		para.Proto = mode
		para.IfName = ifName

		switch mode {
		case "static":
			if v, _ := cmd.Flags().GetString("ip"); v != "" {
				para.IP = v
			}
			if v, _ := cmd.Flags().GetString("netmask"); v != "" {
				para.Netmask = v
			}
			if v, _ := cmd.Flags().GetString("gateway"); v != "" {
				para.Gateway = v
			}
			if v, _ := cmd.Flags().GetString("dns1"); v != "" {
				para.PriDNS = v
			}
			if v, _ := cmd.Flags().GetString("dns2"); v != "" {
				para.SndDNS = v
			}

		case "dhcp":
			if v, _ := cmd.Flags().GetString("hostname"); v != "" {
				para.Hostname = v
			}

		case "pppoe":
			para.IfType = "pppoe"
			if strings.HasSuffix(ifName, "_pppoe") {
				para.TPPPoEEnable = "1"
			}
			if v, _ := cmd.Flags().GetString("username"); v != "" {
				para.Username = v
			}
			if v, _ := cmd.Flags().GetString("password"); v != "" {
				para.Password = v
			}
			if v, _ := cmd.Flags().GetString("service-name"); v != "" {
				para.Service = v
			}

		default:
			return fmt.Errorf("不支持的连接模式: %s (仅支持 static/dhcp/pppoe)", mode)
		}

		if err := client.SetWan(ifName, para); err != nil {
			return err
		}

		fmt.Printf("WAN接口 %s 配置成功, 模式: %s\n", ifName, mode)
		return nil
	},
}

func init() {
	wanListCmd.Flags().String("base-name", "wan1_eth", "WAN 基础名称过滤")

	wanSetCmd.Flags().StringP("if-name", "i", "", "接口名称, 如 wan1_eth/wan1_pppoe (必选)")
	wanSetCmd.Flags().StringP("mode", "m", "", "连接模式: static|dhcp|pppoe (必选)")
	// static 模式参数
	wanSetCmd.Flags().String("ip", "", "IP地址 (static模式)")
	wanSetCmd.Flags().String("netmask", "", "子网掩码 (static模式)")
	wanSetCmd.Flags().String("gateway", "", "网关 (static模式)")
	wanSetCmd.Flags().String("dns1", "", "首选DNS (static/pppoe模式)")
	wanSetCmd.Flags().String("dns2", "", "备用DNS (static/pppoe模式)")
	// dhcp 模式参数
	wanSetCmd.Flags().String("hostname", "", "DHCP主机名 (dhcp模式)")
	// pppoe 模式参数
	wanSetCmd.Flags().StringP("username", "u", "", "PPPoE用户名 (pppoe模式)")
	wanSetCmd.Flags().StringP("password", "p", "", "PPPoE密码 (pppoe模式)")
	wanSetCmd.Flags().String("service-name", "", "PPPoE服务名 (pppoe模式)")

	wanCmd.AddCommand(wanListCmd)
	wanCmd.AddCommand(wanSetCmd)
	rootCmd.AddCommand(wanCmd)
}
