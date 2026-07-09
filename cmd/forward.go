package cmd

import (
	"fmt"
	"strings"

	"github.com/ljw/tplink-cli/internal/api"
	"github.com/ljw/tplink-cli/internal/config"
	"github.com/ljw/tplink-cli/internal/format"
	"github.com/ljw/tplink-cli/internal/model"
	"github.com/spf13/cobra"
)

var forwardCmd = &cobra.Command{
	Use:   "forward",
	Short: "端口映射(NAT)规则管理",
	Long:  `管理 TP-Link 路由器的端口映射规则，支持查询、添加、修改、删除和启停操作。`,
}

// list 子命令
var forwardListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询端口映射规则",
	Long:  `查询所有端口映射规则，支持按名称、协议、端口、目标IP、WAN接口、状态过滤，以及分页。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newForwardClient()
		if err != nil {
			return err
		}

		items, err := client.GetRedirects()
		if err != nil {
			return err
		}

		// 过滤
		items = filterRedirects(items, cmd)

		// 分页
		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		totalCount := len(items)
		pagination := format.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    totalCount,
		}

		if page > 0 && pageSize > 0 {
			start := (page - 1) * pageSize
			end := start + pageSize
			if start >= len(items) {
				items = nil
			} else {
				if end > len(items) {
					end = len(items)
				}
				items = items[start:end]
			}
		}

		return format.PrintRedirects(items, format.OutputFormat(output), &pagination)
	},
}

// add 子命令
var forwardAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加端口映射规则",
	Long:  `添加一条新的端口映射规则。--port 支持 "33001" 或 "33001-33010" 范围格式。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		port, _ := cmd.Flags().GetString("port")
		destIP, _ := cmd.Flags().GetString("dest-ip")
		destPort, _ := cmd.Flags().GetString("dest-port")
		proto, _ := cmd.Flags().GetString("proto")
		wan, _ := cmd.Flags().GetString("wan")

		if name == "" || port == "" || destIP == "" {
			return fmt.Errorf("--name, --port, --dest-ip 为必选参数")
		}

		rule, err := buildRedirectRule(name, port, destIP, destPort, proto, wan, "on")
		if err != nil {
			return err
		}

		client, err := newForwardClient()
		if err != nil {
			return err
		}

		ruleID, err := client.AddRedirect(rule)
		if err != nil {
			return err
		}

		fmt.Printf("添加成功, 规则ID: %s\n", ruleID)
		return nil
	},
}

// update 子命令
var forwardUpdateCmd = &cobra.Command{
	Use:   "update <rule-id>",
	Short: "修改端口映射规则",
	Long:  `根据规则ID修改端口映射规则。仅更新指定的字段，未指定的字段保持原值。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ruleID := args[0]

		client, err := newForwardClient()
		if err != nil {
			return err
		}

		// 先查询当前规则
		current, err := client.FindRedirectByID(ruleID)
		if err != nil {
			return err
		}

		// 合并更新字段
		rule := model.RedirectRule{
			Name:          current.Name,
			If:            current.If,
			SrcDport:      current.SrcDport,
			DestPort:      current.DestPort,
			DestIP:        current.DestIP,
			Proto:         current.Proto,
			SrcDportStart: current.SrcDportStart,
			SrcDportEnd:   current.SrcDportEnd,
			DestPortStart: current.DestPortStart,
			DestPortEnd:   current.DestPortEnd,
			Enable:        current.Enable,
		}

		// 覆盖指定字段
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			rule.Name = v
		}
		if v, _ := cmd.Flags().GetString("port"); v != "" {
			parsed, err := parsePort(v)
			if err != nil {
				return err
			}
			rule.SrcDport = parsed.Display
			rule.SrcDportStart = parsed.Start
			rule.SrcDportEnd = parsed.End
		}
		if v, _ := cmd.Flags().GetString("dest-ip"); v != "" {
			rule.DestIP = v
		}
		if v, _ := cmd.Flags().GetString("dest-port"); v != "" {
			parsed, err := parsePort(v)
			if err != nil {
				return err
			}
			rule.DestPort = parsed.Display
			rule.DestPortStart = parsed.Start
			rule.DestPortEnd = parsed.End
		}
		if v, _ := cmd.Flags().GetString("proto"); v != "" {
			rule.Proto = v
		}
		if v, _ := cmd.Flags().GetString("wan"); v != "" {
			rule.If = v
		}

		if err := client.SetRedirect(ruleID, rule); err != nil {
			return err
		}

		fmt.Printf("修改成功, 规则ID: %s\n", ruleID)
		return nil
	},
}

// enable 子命令
var forwardEnableCmd = &cobra.Command{
	Use:   "enable <rule-id>",
	Short: "启用端口映射规则",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return toggleRedirect(args[0], "on")
	},
}

// disable 子命令
var forwardDisableCmd = &cobra.Command{
	Use:   "disable <rule-id>",
	Short: "禁用端口映射规则",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return toggleRedirect(args[0], "off")
	},
}

// delete 子命令
var forwardDeleteCmd = &cobra.Command{
	Use:   "delete <rule-id>",
	Short: "删除端口映射规则",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newForwardClient()
		if err != nil {
			return err
		}

		if err := client.DeleteRedirect(args[0]); err != nil {
			return err
		}

		fmt.Printf("删除成功, 规则ID: %s\n", args[0])
		return nil
	},
}

func init() {
	// forward list 过滤参数
	forwardListCmd.Flags().StringP("name", "n", "", "按规则名称过滤(模糊匹配)")
	forwardListCmd.Flags().String("proto", "", "按协议过滤: ALL|TCP|UDP")
	forwardListCmd.Flags().StringP("port", "p", "", "按外部端口过滤")
	forwardListCmd.Flags().StringP("dest-ip", "d", "", "按目标IP过滤")
	forwardListCmd.Flags().String("wan", "", "按WAN接口过滤")
	forwardListCmd.Flags().String("status", "", "按状态过滤: on|off")
	forwardListCmd.Flags().Int("page", 0, "页码 (从1开始)")
	forwardListCmd.Flags().Int("page-size", 0, "每页条数")

	// forward add 参数
	forwardAddCmd.Flags().StringP("name", "n", "", "规则名称 (必选)")
	forwardAddCmd.Flags().StringP("port", "p", "", "外部端口, 支持 33001 或 33001-33010 格式 (必选)")
	forwardAddCmd.Flags().StringP("dest-ip", "d", "", "目标IP地址 (必选)")
	forwardAddCmd.Flags().StringP("dest-port", "P", "", "目标端口, 默认与外部端口一致")
	forwardAddCmd.Flags().String("proto", "ALL", "协议: ALL|TCP|UDP")
	forwardAddCmd.Flags().String("wan", "WAN1", "WAN 接口")

	// forward update 参数
	forwardUpdateCmd.Flags().StringP("name", "n", "", "规则名称")
	forwardUpdateCmd.Flags().StringP("port", "p", "", "外部端口")
	forwardUpdateCmd.Flags().StringP("dest-ip", "d", "", "目标IP地址")
	forwardUpdateCmd.Flags().StringP("dest-port", "P", "", "目标端口")
	forwardUpdateCmd.Flags().String("proto", "", "协议: ALL|TCP|UDP")
	forwardUpdateCmd.Flags().String("wan", "", "WAN 接口")

	// 注册子命令
	forwardCmd.AddCommand(forwardListCmd)
	forwardCmd.AddCommand(forwardAddCmd)
	forwardCmd.AddCommand(forwardUpdateCmd)
	forwardCmd.AddCommand(forwardEnableCmd)
	forwardCmd.AddCommand(forwardDisableCmd)
	forwardCmd.AddCommand(forwardDeleteCmd)

	// 注册到根命令
	rootCmd.AddCommand(forwardCmd)
}

// newForwardClient 从配置创建 API 客户端
func newForwardClient() (*api.Client, error) {
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

// portRange 解析后的端口范围
type portRange struct {
	Start   string
	End     string
	Display string
}

// parsePort 解析端口参数，支持 "33001" 和 "33001-33010" 两种格式
func parsePort(port string) (*portRange, error) {
	if strings.Contains(port, "-") {
		parts := strings.SplitN(port, "-", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("端口范围格式错误: %s, 正确格式: 33001-33010", port)
		}
		return &portRange{Start: parts[0], End: parts[1], Display: port}, nil
	}
	return &portRange{Start: port, End: port, Display: port}, nil
}

// buildRedirectRule 构建重定向规则
func buildRedirectRule(name, port, destIP, destPort, proto, wan, enable string) (model.RedirectRule, error) {
	srcPort, err := parsePort(port)
	if err != nil {
		return model.RedirectRule{}, err
	}

	// 目标端口默认与外部端口一致
	if destPort == "" {
		destPort = port
	}
	dstPort, err := parsePort(destPort)
	if err != nil {
		return model.RedirectRule{}, err
	}

	return model.RedirectRule{
		Name:          name,
		If:            wan,
		SrcDport:      srcPort.Display,
		DestPort:      dstPort.Display,
		DestIP:        destIP,
		Proto:         proto,
		SrcDportStart: srcPort.Start,
		SrcDportEnd:   srcPort.End,
		DestPortStart: dstPort.Start,
		DestPortEnd:   dstPort.End,
		Enable:        enable,
	}, nil
}

// toggleRedirect 启用或禁用规则
func toggleRedirect(ruleID, status string) error {
	client, err := newForwardClient()
	if err != nil {
		return err
	}

	rule, err := client.FindRedirectByID(ruleID)
	if err != nil {
		return err
	}

	updateRule := model.RedirectRule{
		Name:          rule.Name,
		If:            rule.If,
		SrcDport:      rule.SrcDport,
		DestPort:      rule.DestPort,
		DestIP:        rule.DestIP,
		Proto:         rule.Proto,
		SrcDportStart: rule.SrcDportStart,
		SrcDportEnd:   rule.SrcDportEnd,
		DestPortStart: rule.DestPortStart,
		DestPortEnd:   rule.DestPortEnd,
		Enable:        status,
	}

	if err := client.SetRedirect(ruleID, updateRule); err != nil {
		return err
	}

	action := "启用"
	if status == "off" {
		action = "禁用"
	}
	fmt.Printf("%s成功, 规则ID: %s\n", action, ruleID)
	return nil
}

// filterRedirects 客户端过滤规则列表
func filterRedirects(items []model.RedirectItem, cmd *cobra.Command) []model.RedirectItem {
	name, _ := cmd.Flags().GetString("name")
	proto, _ := cmd.Flags().GetString("proto")
	port, _ := cmd.Flags().GetString("port")
	destIP, _ := cmd.Flags().GetString("dest-ip")
	wan, _ := cmd.Flags().GetString("wan")
	status, _ := cmd.Flags().GetString("status")

	if name == "" && proto == "" && port == "" && destIP == "" && wan == "" && status == "" {
		return items
	}

	var result []model.RedirectItem
	for _, item := range items {
		for id, rule := range item {
			if name != "" && !strings.Contains(strings.ToLower(rule.Name), strings.ToLower(name)) {
				continue
			}
			if proto != "" && !strings.EqualFold(rule.Proto, proto) {
				continue
			}
			if port != "" && rule.SrcDport != port && rule.SrcDportStart != port {
				continue
			}
			if destIP != "" && rule.DestIP != destIP {
				continue
			}
			if wan != "" && !strings.EqualFold(rule.If, wan) {
				continue
			}
			if status != "" && !strings.EqualFold(rule.Enable, status) {
				continue
			}
			filtered := make(model.RedirectItem)
			filtered[id] = rule
			result = append(result, filtered)
		}
	}

	return result
}
