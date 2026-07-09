package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ljw/tplink-cli/internal/model"
	"github.com/spf13/cobra"
)

var dhcpCmd = &cobra.Command{
	Use:   "dhcp",
	Short: "DHCP服务管理",
	Long:  `管理 TP-Link 路由器的 DHCP 服务，包括配置、客户端列表和静态地址分配。`,
}

// ========== dhcp config list ==========
var dhcpConfigListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看DHCP配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		data, err := client.GetDhcpConfig()
		if err != nil {
			return err
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(data)
		case "yaml":
			fmt.Printf("enable: %s\npool_start: %s\npool_end: %s\nlease_time: %s\ngateway: %s\npri_dns: %s\nsnd_dns: %s\n",
				data.Enable, data.PoolStart, data.PoolEnd, data.LeaseTime, data.Gateway, data.PriDNS, data.SndDNS)
			return nil
		default:
			fmt.Printf("状态:       %s\n", statusLabel(data.Enable))
			fmt.Printf("地址池起始: %s\n", data.PoolStart)
			fmt.Printf("地址池结束: %s\n", data.PoolEnd)
			fmt.Printf("租约时间:   %s 秒\n", data.LeaseTime)
			if data.Gateway != "" {
				fmt.Printf("网关:       %s\n", data.Gateway)
			}
			if data.PriDNS != "" {
				fmt.Printf("首选DNS:    %s\n", data.PriDNS)
			}
			if data.SndDNS != "" {
				fmt.Printf("备用DNS:    %s\n", data.SndDNS)
			}
			return nil
		}
	},
}

// ========== dhcp config set ==========
var dhcpConfigSetCmd = &cobra.Command{
	Use:   "set",
	Short: "修改DHCP配置",
	Long:  `修改 LAN 接口的 DHCP 服务配置。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		current, err := client.GetDhcpConfig()
		if err != nil {
			return err
		}

		if v, _ := cmd.Flags().GetString("enable"); v != "" {
			current.Enable = v
		}
		if v, _ := cmd.Flags().GetString("pool-start"); v != "" {
			current.PoolStart = v
		}
		if v, _ := cmd.Flags().GetString("pool-end"); v != "" {
			current.PoolEnd = v
		}
		if v, _ := cmd.Flags().GetString("lease-time"); v != "" {
			current.LeaseTime = v
		}
		if v, _ := cmd.Flags().GetString("gateway"); v != "" {
			current.Gateway = v
		}
		if v, _ := cmd.Flags().GetString("dns1"); v != "" {
			current.PriDNS = v
		}
		if v, _ := cmd.Flags().GetString("dns2"); v != "" {
			current.SndDNS = v
		}
		if v, _ := cmd.Flags().GetString("domain"); v != "" {
			current.Domain = v
		}

		if err := client.SetDhcpConfig(current); err != nil {
			return err
		}

		fmt.Println("DHCP配置修改成功")
		return nil
	},
}

// ========== dhcp client list ==========
var dhcpClientListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看DHCP客户端列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 500
		}
		start := (page - 1) * pageSize
		end := page*pageSize - 1

		items, total, err := client.GetDhcpClients(start, end)
		if err != nil {
			return err
		}

		type ClientRow struct {
			IP        string `json:"ip"`
			Mac       string `json:"mac"`
			Hostname  string `json:"hostname"`
			Expires   string `json:"expires"`
			Interface string `json:"interface"`
		}

		rows := make([]ClientRow, 0, len(items))
		for _, item := range items {
			for _, data := range item {
				rows = append(rows, ClientRow{
					IP:        data.IPAddr,
					Mac:       data.MacAddr,
					Hostname:  decodeURL(data.Hostname),
					Expires:   formatExpires(data.Expires),
					Interface: data.Interface,
				})
			}
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- ip: %s\n  mac: %s\n  hostname: %s\n  expires: %s\n  interface: %s\n",
					r.IP, r.Mac, r.Hostname, r.Expires, r.Interface)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d, page: %d, page_size: %d\n", total, page, pageSize)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("没有找到DHCP客户端")
				return nil
			}
			headers := []string{"IP", "MAC", "HOSTNAME", "EXPIRES", "IF"}
			printTable(headers, rows, func(r ClientRow) []string {
				return []string{r.IP, r.Mac, r.Hostname, r.Expires, r.Interface}
			})
			if total > 0 {
				fmt.Printf("\n[总记录: %d, 第 %d 页 (每页 %d 条)]\n", total, page, pageSize)
			}
			return nil
		}
	},
}

// ========== dhcp static list ==========
var dhcpStaticListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看静态地址分配列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 100
		}
		start := (page - 1) * pageSize
		end := page*pageSize - 1

		items, total, err := client.GetDhcpStatic(start, end)
		if err != nil {
			return err
		}

		type StaticRow struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Mac    string `json:"mac"`
			IP     string `json:"ip"`
			Note   string `json:"note"`
			Enable string `json:"enable"`
		}

		rows := make([]StaticRow, 0, len(items))
		for _, item := range items {
			for _, data := range item {
				rows = append(rows, StaticRow{
					ID:     data.DhcpStaticID,
					Name:   decodeURL(data.Name),
					Mac:    data.Mac,
					IP:     data.IP,
					Note:   decodeURL(data.Note),
					Enable: data.Enable,
				})
			}
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- id: %s\n  name: %s\n  mac: %s\n  ip: %s\n  note: %s\n  enable: %s\n",
					r.ID, r.Name, r.Mac, r.IP, r.Note, r.Enable)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d, page: %d, page_size: %d\n", total, page, pageSize)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("没有找到静态地址分配规则")
				return nil
			}
			headers := []string{"ID", "MAC", "IP", "NAME", "NOTE", "ENABLE"}
			printTable(headers, rows, func(r StaticRow) []string {
				return []string{r.ID, r.Mac, r.IP, r.Name, r.Note, r.Enable}
			})
			if total > 0 {
				fmt.Printf("\n[总记录: %d, 第 %d 页 (每页 %d 条)]\n", total, page, pageSize)
			}
			return nil
		}
	},
}

// ========== dhcp static add ==========
var dhcpStaticAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加静态地址分配",
	Long: `为指定 MAC 地址绑定固定 IP。

示例:
  tplink dhcp static add -m 90-E2-BA-7F-6D-58 -i 192.168.0.250 -n "备注信息"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mac, _ := cmd.Flags().GetString("mac")
		ip, _ := cmd.Flags().GetString("ip")
		note, _ := cmd.Flags().GetString("note")

		if mac == "" || ip == "" {
			return fmt.Errorf("--mac 和 --ip 为必选参数")
		}

		client, err := newAPIClient()
		if err != nil {
			return err
		}

		para := model.DhcpStaticPara{
			Mac:    mac,
			IP:     ip,
			Note:   note,
			Enable: "on",
		}

		name, err := client.AddDhcpStatic(para)
		if err != nil {
			return err
		}

		fmt.Printf("静态地址分配添加成功, ID: %s\n", strings.TrimPrefix(name, "dhcp_static_"))
		return nil
	},
}

// ========== dhcp static del ==========
var dhcpStaticDelCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "删除静态地址分配",
	Long:  `根据 dhcp_static_id 删除静态地址分配规则。ID 为 list 输出的编号。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		staticID := args[0]
		if err := client.DelDhcpStatic(staticID); err != nil {
			return err
		}

		fmt.Printf("静态地址分配删除成功, ID: %s\n", staticID)
		return nil
	},
}

func init() {
	// dhcp config 子命令
	var dhcpConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "DHCP服务配置",
	}
	dhcpConfigCmd.AddCommand(dhcpConfigListCmd)
	dhcpConfigCmd.AddCommand(dhcpConfigSetCmd)
	dhcpConfigSetCmd.Flags().String("enable", "", "DHCP启用状态: on|off")
	dhcpConfigSetCmd.Flags().String("pool-start", "", "地址池起始IP")
	dhcpConfigSetCmd.Flags().String("pool-end", "", "地址池结束IP")
	dhcpConfigSetCmd.Flags().String("lease-time", "", "租约时间(秒)")
	dhcpConfigSetCmd.Flags().String("gateway", "", "网关地址")
	dhcpConfigSetCmd.Flags().String("dns1", "", "首选DNS")
	dhcpConfigSetCmd.Flags().String("dns2", "", "备用DNS")
	dhcpConfigSetCmd.Flags().String("domain", "", "域名")

	// dhcp client 子命令
	var dhcpClientCmd = &cobra.Command{
		Use:   "client",
		Short: "DHCP客户端管理",
	}
	dhcpClientListCmd.Flags().IntP("page", "", 1, "页码 (从1开始)")
	dhcpClientListCmd.Flags().IntP("page-size", "", 500, "每页条数")
	dhcpClientCmd.AddCommand(dhcpClientListCmd)

	// dhcp static 子命令
	var dhcpStaticCmd = &cobra.Command{
		Use:   "static",
		Short: "静态地址分配管理",
	}
	dhcpStaticAddCmd.Flags().StringP("mac", "m", "", "MAC 地址 (必选)")
	dhcpStaticAddCmd.Flags().StringP("ip", "i", "", "分配的IP地址 (必选)")
	dhcpStaticAddCmd.Flags().StringP("note", "n", "", "备注信息")
	dhcpStaticListCmd.Flags().IntP("page", "", 1, "页码 (从1开始)")
	dhcpStaticListCmd.Flags().IntP("page-size", "", 100, "每页条数")
	dhcpStaticCmd.AddCommand(dhcpStaticListCmd)
	dhcpStaticCmd.AddCommand(dhcpStaticAddCmd)
	dhcpStaticCmd.AddCommand(dhcpStaticDelCmd)

	dhcpCmd.AddCommand(dhcpConfigCmd)
	dhcpCmd.AddCommand(dhcpClientCmd)
	dhcpCmd.AddCommand(dhcpStaticCmd)

	rootCmd.AddCommand(dhcpCmd)
}
