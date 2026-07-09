package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// vpn
	vpnCmd.AddCommand(vpnListCmd)
	vpnCmd.AddCommand(vpnSetCmd)
	rootCmd.AddCommand(vpnCmd)
}

// ========== VPN ==========

var vpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "VPN管理",
}

var vpnListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看VPN配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		cfg, err := client.GetVpnConfig()
		if err != nil {
			return err
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(cfg)
		case "yaml":
			fmt.Printf("server: %s\n", cfg.Server)
			fmt.Printf("username: %s\n", cfg.Username)
			fmt.Printf("password: %s\n", cfg.Password)
			fmt.Printf("protocol: %s\n", cfg.Protocol)
			fmt.Printf("route_mode: %s\n", routeModeLabel(cfg.RouteMode))
			if cfg.RouteMode == "manual" {
				fmt.Printf("remotesubnet: %s\n", decodeURL(cfg.RemoteSubnet))
			}
			fmt.Printf("forward_mode: %s\n", forwardModeLabel(cfg.ForwardMode))
			fmt.Printf("connect_mode: %s\n", connectModeLabel(cfg.ConnectMode))
			fmt.Printf("interface: %s\n", cfg.Interface)
			return nil
		default:
			type VpnRow struct {
				Server      string `json:"server"`
				Username    string `json:"username"`
				Protocol    string `json:"protocol"`
				RouteMode   string `json:"route_mode"`
				ForwardMode string `json:"forward_mode"`
				ConnectMode string `json:"connect_mode"`
				Interface   string `json:"interface"`
			}
			rows := []VpnRow{{
				Server:      cfg.Server,
				Username:    cfg.Username,
				Protocol:    cfg.Protocol,
				RouteMode:   routeModeLabel(cfg.RouteMode),
				ForwardMode: forwardModeLabel(cfg.ForwardMode),
				ConnectMode: connectModeLabel(cfg.ConnectMode),
				Interface:   cfg.Interface,
			}}
			headers := []string{"SERVER", "USERNAME", "PROTOCOL", "ROUTE", "FORWARD", "CONNECT", "INTERFACE"}
			printTable(headers, rows, func(r VpnRow) []string {
				return []string{r.Server, r.Username, r.Protocol, r.RouteMode, r.ForwardMode, r.ConnectMode, r.Interface}
			})
			return nil
		}
	},
}

var vpnSetCmd = &cobra.Command{
	Use:   "set",
	Short: "修改VPN配置",
	Long: `修改VPN配置。只会发送用户显式指定的字段。

示例:
  tplink vpn set --server vpn.example.com --username user1
  tplink vpn set --route-mode manual --remotesubnet "192.168.1.0/24"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		current, err := client.GetVpnConfig()
		if err != nil {
			return err
		}

		cfg := map[string]interface{}{}
		mergeFieldStr(cmd, "server", current.Server, cfg)
		mergeFieldStr(cmd, "username", current.Username, cfg)
		mergeFieldStr(cmd, "password", current.Password, cfg)
		mergeFieldStr(cmd, "protocol", current.Protocol, cfg)
		mergeFieldStr(cmd, "route-mode", current.RouteMode, cfg)
		mergeFieldStr(cmd, "remotesubnet", current.RemoteSubnet, cfg)
		mergeFieldStr(cmd, "forward-mode", current.ForwardMode, cfg)
		mergeFieldStr(cmd, "connect-mode", current.ConnectMode, cfg)
		mergeFieldStr(cmd, "interface", current.Interface, cfg)

		if len(cfg) == 0 {
			return fmt.Errorf("没有指定要修改的字段")
		}

		return client.SetVpnConfig(cfg)
	},
}

// ========== VPN 辅助函数 ==========

func routeModeLabel(mode string) string {
	switch mode {
	case "all":
		return "全网段(all)"
	case "manual":
		return "手动(manual)"
	default:
		return mode
	}
}

func forwardModeLabel(mode string) string {
	switch mode {
	case "nat":
		return "NAT模式"
	case "route":
		return "Route模式"
	default:
		return mode
	}
}

func connectModeLabel(mode string) string {
	switch mode {
	case "auto":
		return "自动(auto)"
	case "manual":
		return "手动(manual)"
	default:
		return mode
	}
}

func init() {
	// vpn set flags
	vpnSetCmd.Flags().String("server", "", "VPN服务器地址")
	vpnSetCmd.Flags().StringP("username", "u", "", "用户名")
	vpnSetCmd.Flags().StringP("password", "p", "", "密码")
	vpnSetCmd.Flags().String("protocol", "", "协议: auto|l2tp|pptp")
	vpnSetCmd.Flags().String("route-mode", "", "路由模式: all|manual")
	vpnSetCmd.Flags().String("remotesubnet", "", "远端子网 (route_mode=manual时有效)")
	vpnSetCmd.Flags().String("forward-mode", "", "转发模式: nat|route")
	vpnSetCmd.Flags().String("connect-mode", "", "连接模式: auto|manual")
	vpnSetCmd.Flags().StringP("interface", "i", "", "接口名称: WAN1")
}
