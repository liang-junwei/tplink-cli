package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ljw/tplink-cli/internal/api"
	"github.com/ljw/tplink-cli/internal/config"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "管理 server 配置",
	Long:  `管理多个 TP-Link 路由器 server 配置，支持添加、删除、修改、列表和默认 server 切换。`,
}

// context list
var contextListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有 server 配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if len(cfg.Servers) == 0 {
			fmt.Println("没有配置任何 server，使用 'tplink context add' 添加")
			return nil
		}

		printServerTable(cfg)
		return nil
	},
}

// context add
	var contextAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "添加 server 配置",
	Long:  `添加一个新的 server 配置。使用 --default 将其设为默认 server。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		url, _ := cmd.Flags().GetString("url")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		setDefault, _ := cmd.Flags().GetBool("default")
		dynamicAuth, _ := cmd.Flags().GetBool("dynamic-auth")

		if url == "" || username == "" || password == "" {
			return fmt.Errorf("--url, --username, --password 为必选参数")
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		srv := config.ServerConfig{
			ServerURL:   url,
			Username:    username,
			Password:    password,
			DynamicAuth: dynamicAuth,
		}

		// 如果开启动态认证，立即获取密钥并编码密码
		if dynamicAuth {
			encodedPwd, err := api.FetchAndEncodePassword(url, password)
			if err != nil {
				return fmt.Errorf(
					"动态认证开启失败（无法从路由器获取密钥）\n"+
						"  错误: %w\n"+
						"  提示: 请确认路由器 URL 正确且可访问，或关闭动态认证：\n"+
						"        tplink context update %s --dynamic-auth=false", name, err,
				)
			}
			srv.EncodedPassword = encodedPwd
			fmt.Println("  动态认证密钥获取成功，密码已编码并持久化")
		}

		if err := cfg.AddServer(name, srv, setDefault); err != nil {
			return err
		}

		fmt.Printf("server '%s' 添加成功", name)
		if dynamicAuth {
			fmt.Print(" (已启用动态认证)")
		}
		if setDefault || cfg.Current == name {
			fmt.Print(" (已设为默认)")
		}
		fmt.Println()
		return nil
	},
}

// context delete
var contextDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "删除 server 配置",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if err := cfg.RemoveServer(name); err != nil {
			return err
		}

		fmt.Printf("server '%s' 已删除\n", name)
		return nil
	},
}

// context use
var contextUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "切换默认 server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if err := cfg.SetCurrent(name); err != nil {
			return err
		}

		fmt.Printf("默认 server 已切换为 '%s'\n", name)
		return nil
	},
}

// context update
	var contextUpdateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "更新 server 配置",
	Long:  `更新指定 server 的配置，仅更新指定的字段。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		srv, ok := cfg.Servers[name]
		if !ok {
			return fmt.Errorf("server '%s' 不存在", name)
		}

		if v, _ := cmd.Flags().GetString("url"); v != "" {
			srv.ServerURL = v
		}
		if v, _ := cmd.Flags().GetString("username"); v != "" {
			srv.Username = v
		}
		if v, _ := cmd.Flags().GetString("password"); v != "" {
			srv.Password = v
		}

		// 处理 dynamic-auth 变更
		if cmd.Flags().Changed("dynamic-auth") {
			v, _ := cmd.Flags().GetBool("dynamic-auth")
			srv.DynamicAuth = v
			if v {
				// 开启动态认证，立即获取密钥并编码密码
				encodedPwd, err := api.FetchAndEncodePassword(srv.ServerURL, srv.Password)
				if err != nil {
					return fmt.Errorf(
						"动态认证开启失败（无法从路由器获取密钥）\n" +
							"  错误: %w\n" +
							"  提示: 请确认路由器 URL 正确且可访问，或关闭动态认证：\n" +
							"        tplink context update %s --dynamic-auth=false",
						name, err,
					)
				}
				srv.EncodedPassword = encodedPwd
				fmt.Println("  动态认证密钥获取成功，密码已编码并持久化")
			} else {
				// 关闭动态认证，清除已编码的密码
				srv.EncodedPassword = ""
				fmt.Println("  动态认证已关闭，已清除编码密码")
			}
		} else if srv.DynamicAuth {
			// dynamic-auth 未变更，但检查密码是否变更
			if cmd.Flags().Changed("password") {
				// 密码已变更，重新编码
				encodedPwd, err := api.FetchAndEncodePassword(srv.ServerURL, srv.Password)
				if err != nil {
					return fmt.Errorf(
						"重新编码密码失败\n"+
							"  错误: %w\n"+
							"  提示: 请重新开启动态认证：\n"+
							"        tplink context update %s --dynamic-auth=true",
						name, err,
					)
				}
				srv.EncodedPassword = encodedPwd
				fmt.Println("  密码已变更，重新编码并持久化")
			}
		}

		// 修改认证信息后清除缓存的 stok
		srv.Stok = ""
		srv.StokExpiresAt = 0

		if err := cfg.UpdateServer(name, srv); err != nil {
			return err
		}

		fmt.Printf("server '%s' 更新成功\n", name)
		return nil
	},
}

func init() {
	// context add flags
	contextAddCmd.Flags().StringP("url", "U", "", "server URL (必选)")
	contextAddCmd.Flags().StringP("username", "u", "", "用户名 (必选)")
	contextAddCmd.Flags().StringP("password", "p", "", "密码 (必选)")
	contextAddCmd.Flags().Bool("default", false, "设为默认 server")
	contextAddCmd.Flags().Bool("dynamic-auth", false, "启用动态认证（自动从路由器获取加密密钥）")

	// context update flags
	contextUpdateCmd.Flags().StringP("url", "U", "", "server URL")
	contextUpdateCmd.Flags().StringP("username", "u", "", "用户名")
	contextUpdateCmd.Flags().StringP("password", "p", "", "密码")
	contextUpdateCmd.Flags().Bool("dynamic-auth", false, "启用/禁用动态认证")

	// 注册子命令
	contextCmd.AddCommand(contextListCmd)
	contextCmd.AddCommand(contextAddCmd)
	contextCmd.AddCommand(contextDeleteCmd)
	contextCmd.AddCommand(contextUseCmd)
	contextCmd.AddCommand(contextUpdateCmd)

	// 注册到根命令
	rootCmd.AddCommand(contextCmd)
}

// printServerTable 打印 server 列表表格
func printServerTable(cfg *config.AppConfig) {
	// 获取排序后的 server 名称
	names := make([]string, 0, len(cfg.Servers))
	for name := range cfg.Servers {
		names = append(names, name)
	}
	sort.Strings(names)

	// 计算列宽
	headers := []string{"NAME", "URL", "USERNAME", "DEFAULT"}
	widths := []int{4, 3, 8, 7}

	rows := make([][]string, 0, len(names))
	for _, name := range names {
		srv := cfg.Servers[name]
		isDefault := ""
		if name == cfg.Current {
			isDefault = "*"
		}
		row := []string{name, srv.ServerURL, srv.Username, isDefault}
		rows = append(rows, row)
		for i, v := range row {
			if len(v) > widths[i] {
				widths[i] = len(v)
			}
		}
	}

	// 打印表头
	printTableRow(headers, widths)

	// 分隔线
	seps := make([]string, len(headers))
	for i, w := range widths {
		seps[i] = strings.Repeat("-", w)
	}
	fmt.Println("  " + strings.Join(seps, "  "))

	// 打印数据行
	for _, row := range rows {
		printTableRow(row, widths)
	}
}

func printTableRow(cols []string, widths []int) {
	padded := make([]string, len(cols))
	for i, col := range cols {
		if len(col) < widths[i] {
			padded[i] = col + strings.Repeat(" ", widths[i]-len(col))
		} else {
			padded[i] = col
		}
	}
	fmt.Println("  " + strings.Join(padded, "  "))
}
