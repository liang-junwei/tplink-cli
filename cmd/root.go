package cmd

import (
	"fmt"
	"os"

	"github.com/ljw/tplink-cli/internal/api"
	"github.com/ljw/tplink-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	output     string
	serverName string
	dryRun     bool
	appConfig  *config.AppConfig
)

var rootCmd = &cobra.Command{
	Use:   "tplink",
	Short: "TP-Link 路由器 CLI 管理工具",
	Long:  `tplink-cli 是一个用于管理 TP-Link 路由器的命令行工具，支持端口映射(NAT)规则的查询、添加、修改、删除和启停，以及原始 API 请求。`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 加载配置（某些命令不需要配置，如 context add 首次添加）
		if cfg, err := config.Load(); err == nil {
			appConfig = cfg
		} else {
			appConfig = &config.AppConfig{Servers: make(map[string]config.ServerConfig)}
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// newAPIClient 从配置创建 API 客户端（所有命令文件共用）
func newAPIClient() (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	srvCfg, srvName, err := cfg.GetServer(serverName)
	if err != nil {
		return nil, err
	}

	client := api.NewClient(srvCfg, srvName)
	client.DryRun = dryRun
	return client, nil
}

func init() {
	// 配置文件路径（可选覆盖默认 ~/.tplink.json）
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (默认 ~/.tplink.json)")

	// 目标 server
	rootCmd.PersistentFlags().StringVarP(&serverName, "server", "S", "", "指定目标 server 名称")

	// 输出格式
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "输出格式: table|json|yaml")

	// dry-run
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "模拟执行，仅输出请求信息不实际发送")

	// 初始化自定义配置路径
	cobra.OnInitialize(func() {
		if cfgFile != "" {
			config.SetConfigPath(cfgFile)
		}
	})
}
