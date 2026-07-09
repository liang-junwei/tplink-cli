package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// route
	routeCmd.AddCommand(routeSysCmd)
	routeSysCmd.AddCommand(routeSysListCmd)
	routeCmd.AddCommand(routePolicyCmd)
	routePolicyCmd.AddCommand(routePolicyListCmd)
	routePolicyCmd.AddCommand(routePolicyAddCmd)
	routePolicyCmd.AddCommand(routePolicyDelCmd)
	routeCmd.AddCommand(routeStaticCmd)
	routeStaticCmd.AddCommand(routeStaticListCmd)
	routeStaticCmd.AddCommand(routeStaticAddCmd)
	routeStaticCmd.AddCommand(routeStaticDelCmd)
	rootCmd.AddCommand(routeCmd)

	// napt
	naptCmd.AddCommand(naptRuleCmd)
	naptRuleCmd.AddCommand(naptRuleListCmd)
	naptRuleCmd.AddCommand(naptRuleAddCmd)
	naptRuleCmd.AddCommand(naptRuleDelCmd)
	rootCmd.AddCommand(naptCmd)

	// alg
	algCmd.AddCommand(algListCmd)
	algCmd.AddCommand(algSetCmd)
	rootCmd.AddCommand(algCmd)

	// phddns
	phddnsCmd.AddCommand(phddnsListCmd)
	phddnsCmd.AddCommand(phddnsAddCmd)
	phddnsCmd.AddCommand(phddnsSetCmd)
	phddnsCmd.AddCommand(phddnsDelCmd)
	rootCmd.AddCommand(phddnsCmd)
}

// ========== 路由 - 系统路由 ==========

var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "路由管理",
}

var routeSysCmd = &cobra.Command{
	Use:   "sys",
	Short: "系统路由",
}

var routeSysListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看系统路由表",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		items, total, err := client.GetSysRoutes()
		if err != nil {
			return err
		}

		type Row struct {
			ID      string `json:"id"`
			Target  string `json:"target"`
			Netmask string `json:"netmask"`
			Gateway string `json:"gateway"`
			If      string `json:"if"`
			Metric  string `json:"metric"`
		}

		rows := make([]Row, 0, len(items))
		for _, item := range items {
			rows = append(rows, Row{
				ID:      item.DotName,
				Target:  item.Target,
				Netmask: item.Netmask,
				Gateway: item.Gateway,
				If:      item.If,
				Metric:  item.Metric,
			})
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- id: %s\n  target: %s\n  netmask: %s\n  gateway: %s\n  if: %s\n  metric: %s\n",
					r.ID, r.Target, r.Netmask, r.Gateway, r.If, r.Metric)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("系统路由表为空")
				return nil
			}
			headers := []string{"ID", "TARGET", "NETMASK", "GATEWAY", "IF", "METRIC"}
			printTable(headers, rows, func(r Row) []string {
				return []string{r.ID, r.Target, r.Netmask, r.Gateway, r.If, r.Metric}
			})
			fmt.Printf("\n# total: %d\n", total)
			return nil
		}
	},
}

// ========== 路由 - 策略路由 ==========

var routePolicyCmd = &cobra.Command{
	Use:   "policy",
	Short: "策略路由管理",
}

var routePolicyListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看策略路由规则",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		items, total, err := client.GetPolicyRoutes()
		if err != nil {
			return err
		}

		type Row struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Enable   string `json:"enable"`
			SrcGroup string `json:"src_ipgroup"`
			DstGroup string `json:"dst_ipgroup"`
			If       string `json:"if"`
			TimeObj  string `json:"timeobj"`
		}

		rows := make([]Row, 0, len(items))
		for _, item := range items {
			rows = append(rows, Row{
				ID:       item.DotName,
				Name:     item.Name,
				Enable:   statusLabel(item.Enable),
				SrcGroup: item.SrcIPGroup,
				DstGroup: item.DstIPGroup,
				If:       item.If,
				TimeObj:  item.TimeObj,
			})
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- id: %s\n  name: %s\n  enable: %s\n  src_ipgroup: %s\n  dst_ipgroup: %s\n  if: %s\n  timeobj: %s\n",
					r.ID, r.Name, r.Enable, r.SrcGroup, r.DstGroup, r.If, r.TimeObj)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("策略路由为空")
				return nil
			}
			headers := []string{"ID", "NAME", "ENABLE", "SRC_IPGROUP", "DST_IPGROUP", "IF", "TIMEOBJ"}
			printTable(headers, rows, func(r Row) []string {
				return []string{r.ID, r.Name, r.Enable, r.SrcGroup, r.DstGroup, r.If, r.TimeObj}
			})
			fmt.Printf("\n# total: %d\n", total)
			return nil
		}
	},
}

var routePolicyAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加策略路由",
	Long: `添加策略路由规则（两步: 先添加 service，再添加 rule）。

示例:
  tplink route policy add --name my_policy --proto tcp-udp --sport "1-65535" --dport "5000-5200"
  tplink route policy add --name my_policy --proto tcp-udp --sport "1-65535" --dport "80" --src-ipgroup ipgroup1 --dst-ipgroup ipgroup2 --if WAN1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		proto, _ := cmd.Flags().GetString("proto")
		sport, _ := cmd.Flags().GetString("sport")
		dport, _ := cmd.Flags().GetString("dport")
		ifName, _ := cmd.Flags().GetString("if")
		srcIPGroup, _ := cmd.Flags().GetString("src-ipgroup")
		dstIPGroup, _ := cmd.Flags().GetString("dst-ipgroup")
		timeObj, _ := cmd.Flags().GetString("timeobj")
		enable, _ := cmd.Flags().GetString("enable")

		if name == "" {
			return fmt.Errorf("--name 参数必填")
		}
		if proto == "" {
			proto = "tcp-udp"
		}
		if sport == "" {
			sport = "1-65535"
		}
		if dport == "" {
			return fmt.Errorf("--dport 参数必填")
		}
		if ifName == "" {
			ifName = "WAN1"
		}
		if enable == "" {
			enable = "on"
		}

		// 生成内部 service 名称
		serviceName := "_pr_" + timeBasedSuffix()

		id, err := client.AddPolicyRoute(name, serviceName, proto, sport, dport, ifName, srcIPGroup, dstIPGroup, timeObj, enable)
		if err != nil {
			return err
		}
		fmt.Printf("策略路由添加成功, ID: %s, service: %s\n", id, serviceName)
		return nil
	},
}

var routePolicyDelCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "删除策略路由",
	Long: `根据ID删除策略路由规则及关联的 service。ID 来自 list 输出的 ID 列。

删除流程: 先删除 rule，再删除关联的 service。

示例:
  tplink route policy del policy_rule_1782980672`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		ruleID := args[0]

		// 先查询规则列表找到关联的 service 名称
		var serviceID string
		if !client.DryRun {
			rules, _, err := client.GetPolicyRoutes()
			if err != nil {
				return err
			}
			found := false
			for _, rule := range rules {
				if rule.DotName == ruleID {
					serviceID = rule.ServiceType
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("未找到ID为 %s 的策略路由规则", ruleID)
			}
		} else {
			// dry-run: 仍打印 GET 请求
			client.GetPolicyRoutes()
			serviceID = "_pr_dryrun"
		}

		if err := client.DelPolicyRoute(ruleID, serviceID); err != nil {
			return err
		}
		fmt.Printf("策略路由删除成功, ID: %s\n", ruleID)
		fmt.Printf("关联服务删除成功, service: %s\n", serviceID)
		return nil
	},
}

// ========== 路由 - 静态路由 ==========

var routeStaticCmd = &cobra.Command{
	Use:   "static",
	Short: "静态路由管理",
}

var routeStaticListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看静态路由规则",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		items, total, err := client.GetStaticRoutes()
		if err != nil {
			return err
		}

		type Row struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Enable  string `json:"enable"`
			Target  string `json:"target"`
			Netmask string `json:"netmask"`
			Gateway string `json:"gateway"`
			If      string `json:"if"`
			Metric  string `json:"metric"`
			Note    string `json:"note"`
		}

		rows := make([]Row, 0, len(items))
		for _, item := range items {
			rows = append(rows, Row{
				ID:      item.DotName,
				Name:    item.Name,
				Enable:  statusLabel(item.Enable),
				Target:  item.Target,
				Netmask: item.Netmask,
				Gateway: item.Gateway,
				If:      item.If,
				Metric:  item.Metric,
				Note:    decodeURL(item.Note),
			})
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- id: %s\n  name: %s\n  enable: %s\n  target: %s\n  netmask: %s\n  gateway: %s\n  if: %s\n  metric: %s\n  note: %s\n",
					r.ID, r.Name, r.Enable, r.Target, r.Netmask, r.Gateway, r.If, r.Metric, r.Note)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("静态路由为空")
				return nil
			}
			headers := []string{"ID", "NAME", "ENABLE", "TARGET", "NETMASK", "GATEWAY", "IF", "METRIC", "NOTE"}
			printTable(headers, rows, func(r Row) []string {
				return []string{r.ID, r.Name, r.Enable, r.Target, r.Netmask, r.Gateway, r.If, r.Metric, r.Note}
			})
			fmt.Printf("\n# total: %d\n", total)
			return nil
		}
	},
}

var routeStaticAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加静态路由",
	Long: `添加静态路由规则。

示例:
  tplink route static add --name my_route --target 192.168.0.0 --netmask 255.255.255.0 --gateway 192.168.0.1 --if LAN --note "测试"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		target, _ := cmd.Flags().GetString("target")
		netmask, _ := cmd.Flags().GetString("netmask")
		gateway, _ := cmd.Flags().GetString("gateway")
		ifName, _ := cmd.Flags().GetString("if")
		metric, _ := cmd.Flags().GetString("metric")
		note, _ := cmd.Flags().GetString("note")
		enable, _ := cmd.Flags().GetString("enable")

		if name == "" || target == "" || netmask == "" || gateway == "" {
			return fmt.Errorf("--name, --target, --netmask, --gateway 参数必填")
		}
		if ifName == "" {
			ifName = "LAN"
		}
		if metric == "" {
			metric = "0"
		}
		if enable == "" {
			enable = "on"
		}

		id, err := client.AddStaticRoute(name, target, netmask, gateway, ifName, metric, note, enable)
		if err != nil {
			return err
		}
		fmt.Printf("静态路由添加成功, ID: %s\n", id)
		return nil
	},
}

var routeStaticDelCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "删除静态路由",
	Long: `根据ID删除静态路由。ID 来自 list 输出的 ID 列。

示例:
  tplink route static del user_route_1782981222`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		if err := client.DelStaticRoute(args[0]); err != nil {
			return err
		}
		fmt.Printf("静态路由删除成功, ID: %s\n", args[0])
		return nil
	},
}

// ========== NAPT ==========

var naptCmd = &cobra.Command{
	Use:   "napt",
	Short: "NAPT规则管理",
}

var naptRuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "NAPT规则",
}

var naptRuleListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看NAPT规则",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		items, total, err := client.GetNaptRules()
		if err != nil {
			return err
		}

		type Row struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Enable string `json:"enable"`
			If     string `json:"if"`
			IP     string `json:"ip"`
		}

		rows := make([]Row, 0, len(items))
		for _, item := range items {
			rows = append(rows, Row{
				ID:     item.DotName,
				Name:   decodeURL(item.Name),
				Enable: statusLabel(item.Enable),
				If:     item.If,
				IP:     decodeURL(item.IP),
			})
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- id: %s\n  name: %s\n  enable: %s\n  if: %s\n  ip: %s\n",
					r.ID, r.Name, r.Enable, r.If, r.IP)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("NAPT规则为空")
				return nil
			}
			headers := []string{"ID", "NAME", "ENABLE", "IF", "IP"}
			printTable(headers, rows, func(r Row) []string {
				return []string{r.ID, r.Name, r.Enable, r.If, r.IP}
			})
			fmt.Printf("\n# total: %d\n", total)
			return nil
		}
	},
}

var naptRuleAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加NAPT规则",
	Long: `添加NAPT规则。

示例:
  tplink napt rule add --name my_napt --ip "192.168.2.0/24" --if WAN1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		ip, _ := cmd.Flags().GetString("ip")
		ifName, _ := cmd.Flags().GetString("if")
		enable, _ := cmd.Flags().GetString("enable")

		if name == "" || ip == "" {
			return fmt.Errorf("--name 和 --ip 参数必填")
		}
		if ifName == "" {
			ifName = "WAN1"
		}
		if enable == "" {
			enable = "on"
		}

		id, err := client.AddNaptRule(name, ip, ifName, enable)
		if err != nil {
			return err
		}
		fmt.Printf("NAPT规则添加成功, ID: %s\n", id)
		return nil
	},
}

var naptRuleDelCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "删除NAPT规则",
	Long: `根据ID删除NAPT规则。ID 来自 list 输出的 ID 列。

示例:
  tplink napt rule del rule_napt_1782982451`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		if err := client.DelNaptRule(args[0]); err != nil {
			return err
		}
		fmt.Printf("NAPT规则删除成功, ID: %s\n", args[0])
		return nil
	},
}

// ========== ALG ==========

var algCmd = &cobra.Command{
	Use:   "alg",
	Short: "ALG配置管理",
}

var algListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看ALG配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		cfg, err := client.GetAlgConfig()
		if err != nil {
			return err
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(cfg)
		case "yaml":
			fmt.Printf("ftp: %s\n", statusLabel(cfg.Ftp))
			fmt.Printf("h323: %s\n", statusLabel(cfg.H323))
			fmt.Printf("pptp: %s\n", statusLabel(cfg.Pptp))
			fmt.Printf("sip: %s\n", statusLabel(cfg.Sip))
			fmt.Printf("l2tp: %s\n", statusLabel(cfg.L2tp))
			fmt.Printf("tftp: %s\n", statusLabel(cfg.Tftp))
			fmt.Printf("ipsec: %s\n", statusLabel(cfg.Ipsec))
			return nil
		default:
			type Row struct {
				Field string `json:"field"`
				Value string `json:"value"`
			}
			rows := []Row{
				{"ftp", statusLabel(cfg.Ftp)},
				{"h323", statusLabel(cfg.H323)},
				{"pptp", statusLabel(cfg.Pptp)},
				{"sip", statusLabel(cfg.Sip)},
				{"l2tp", statusLabel(cfg.L2tp)},
				{"tftp", statusLabel(cfg.Tftp)},
				{"ipsec", statusLabel(cfg.Ipsec)},
			}
			headers := []string{"SERVICE", "ENABLE"}
			printTable(headers, rows, func(r Row) []string {
				return []string{r.Field, r.Value}
			})
			return nil
		}
	},
}

var algSetCmd = &cobra.Command{
	Use:   "set",
	Short: "修改ALG配置",
	Long: `修改ALG配置。只会发送用户显式指定的字段。

示例:
  tplink alg set --ftp on --pptp off
  tplink alg set --sip on --ipsec on`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		current, err := client.GetAlgConfig()
		if err != nil {
			return err
		}

		cfg := map[string]interface{}{}
		mergeFieldBoolMapped(cmd, "ftp", "ftp", cfg)
		mergeFieldBoolMapped(cmd, "h323", "h323", cfg)
		mergeFieldBoolMapped(cmd, "pptp", "pptp", cfg)
		mergeFieldBoolMapped(cmd, "sip", "sip", cfg)
		mergeFieldBoolMapped(cmd, "l2tp", "l2tp", cfg)
		mergeFieldBoolMapped(cmd, "tftp", "tftp", cfg)
		mergeFieldBoolMapped(cmd, "ipsec", "ipsec", cfg)
		_ = current

		if len(cfg) == 0 {
			return fmt.Errorf("没有指定要修改的字段")
		}

		return client.SetAlgConfig(cfg)
	},
}

// ========== Phddns ==========

var phddnsCmd = &cobra.Command{
	Use:   "phddns",
	Short: "花生壳DDNS管理",
}

var phddnsListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看花生壳DDNS配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		items, total, err := client.GetPhddnsList()
		if err != nil {
			return err
		}

		type Row struct {
			ID        string `json:"id"`
			Domain    string `json:"domain"`
			Enable    string `json:"enable"`
			ConnState string `json:"connstate"`
			Username  string `json:"username"`
			Interface string `json:"interface"`
		}

		rows := make([]Row, 0, len(items))
		for _, item := range items {
			rows = append(rows, Row{
				ID:        item.DotName,
				Domain:    item.Domain,
				Enable:    statusLabel(item.Enable),
				ConnState: phddnsStateLabel(item.ConnState),
				Username:  item.Username,
				Interface: item.Interface,
			})
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- id: %s\n  domain: %s\n  enable: %s\n  connstate: %s\n  username: %s\n  interface: %s\n",
					r.ID, r.Domain, r.Enable, r.ConnState, r.Username, r.Interface)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("Phddns配置为空")
				return nil
			}
			headers := []string{"ID", "DOMAIN", "ENABLE", "CONN_STATE", "USERNAME", "INTERFACE"}
			printTable(headers, rows, func(r Row) []string {
				return []string{r.ID, r.Domain, r.Enable, r.ConnState, r.Username, r.Interface}
			})
			fmt.Printf("\n# total: %d\n", total)
			return nil
		}
	},
}

var phddnsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加花生壳DDNS配置",
	Long: `添加花生壳DDNS配置。

示例:
  tplink phddns add --username myuser --password mypass --if WAN1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		ifName, _ := cmd.Flags().GetString("if")
		enable, _ := cmd.Flags().GetString("enable")

		if username == "" || password == "" {
			return fmt.Errorf("--username 和 --password 参数必填")
		}
		if ifName == "" {
			ifName = "WAN1"
		}
		if enable == "" {
			enable = "on"
		}

		id, err := client.AddPhddns(username, password, ifName, enable)
		if err != nil {
			return err
		}
		fmt.Printf("Phddns添加成功, ID: %s\n", id)
		return nil
	},
}

var phddnsSetCmd = &cobra.Command{
	Use:   "set <id>",
	Short: "修改花生壳DDNS配置",
	Long: `修改Phddns配置。需要指定条目ID，只发送用户显式指定的字段。

示例:
  tplink phddns set phddns_1751252413 --username newuser --password newpass`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		dotName := args[0]

		// 先获取当前配置（dry-run 时会打印 GET 请求）
		_, _, err = client.GetPhddnsList()
		if err != nil {
			return err
		}

		cfg := map[string]interface{}{}
		mergeFieldStr(cmd, "interface", "", cfg)
		mergeFieldStr(cmd, "username", "", cfg)
		mergeFieldStr(cmd, "password", "", cfg)
		mergeFieldStr(cmd, "enable", "", cfg)

		if len(cfg) == 0 {
			return fmt.Errorf("没有指定要修改的字段")
		}

		return client.SetPhddns(dotName, cfg)
	},
}

var phddnsDelCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "删除花生壳DDNS配置",
	Long: `根据ID删除Phddns配置。ID 来自 list 输出的 ID 列。

示例:
  tplink phddns del phddns_1751252413`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		if err := client.DelPhddns(args[0]); err != nil {
			return err
		}
		fmt.Printf("Phddns删除成功, ID: %s\n", args[0])
		return nil
	},
}

// ========== 辅助函数 ==========

// phddnsStateLabel 将花生壳连接状态转为中文
func phddnsStateLabel(s string) string {
	switch s {
	case "2":
		return "已连接"
	case "1":
		return "连接中"
	case "0":
		return "未连接"
	default:
		return s
	}
}

func init() {
	// route policy add flags
	routePolicyAddCmd.Flags().StringP("name", "n", "", "规则名称")
	routePolicyAddCmd.Flags().String("proto", "tcp-udp", "协议: tcp-udp|tcp|udp")
	routePolicyAddCmd.Flags().String("sport", "1-65535", "源端口范围")
	routePolicyAddCmd.Flags().StringP("dport", "d", "", "目标端口范围")
	routePolicyAddCmd.Flags().StringP("if", "i", "WAN1", "出口接口")
	routePolicyAddCmd.Flags().String("src-ipgroup", "", "源IP地址组名称")
	routePolicyAddCmd.Flags().String("dst-ipgroup", "", "目标IP地址组名称")
	routePolicyAddCmd.Flags().String("timeobj", "", "时间对象名称")
	routePolicyAddCmd.Flags().String("enable", "on", "启用: on|off")

	// route static add flags
	routeStaticAddCmd.Flags().StringP("name", "n", "", "规则名称")
	routeStaticAddCmd.Flags().StringP("target", "t", "", "目标网络")
	routeStaticAddCmd.Flags().StringP("netmask", "m", "", "子网掩码")
	routeStaticAddCmd.Flags().StringP("gateway", "g", "", "网关地址")
	routeStaticAddCmd.Flags().StringP("if", "i", "LAN", "接口: LAN|WAN1")
	routeStaticAddCmd.Flags().String("metric", "0", "跃点数")
	routeStaticAddCmd.Flags().String("note", "", "备注")
	routeStaticAddCmd.Flags().String("enable", "on", "启用: on|off")

	// napt rule add flags
	naptRuleAddCmd.Flags().StringP("name", "n", "", "规则名称")
	naptRuleAddCmd.Flags().StringP("ip", "i", "", "IP地址段 (如 192.168.2.0/24)")
	naptRuleAddCmd.Flags().String("if", "WAN1", "出口接口: WAN1")
	naptRuleAddCmd.Flags().String("enable", "on", "启用: on|off")

	// alg set flags
	algSetCmd.Flags().String("ftp", "", "启用FTP ALG: on|off")
	algSetCmd.Flags().String("h323", "", "启用H323 ALG: on|off")
	algSetCmd.Flags().String("pptp", "", "启用PPTP ALG: on|off")
	algSetCmd.Flags().String("sip", "", "启用SIP ALG: on|off")
	algSetCmd.Flags().String("l2tp", "", "启用L2TP ALG: on|off")
	algSetCmd.Flags().String("tftp", "", "启用TFTP ALG: on|off")
	algSetCmd.Flags().String("ipsec", "", "启用IPSEC ALG: on|off")

	// phddns add flags
	phddnsAddCmd.Flags().StringP("username", "u", "", "花生壳用户名")
	phddnsAddCmd.Flags().StringP("password", "p", "", "花生壳密码")
	phddnsAddCmd.Flags().StringP("if", "i", "WAN1", "接口名称: WAN1")
	phddnsAddCmd.Flags().String("enable", "on", "启用: on|off")

	// phddns set flags
	phddnsSetCmd.Flags().String("interface", "", "接口名称: WAN1")
	phddnsSetCmd.Flags().StringP("username", "u", "", "花生壳用户名")
	phddnsSetCmd.Flags().StringP("password", "p", "", "花生壳密码")
	phddnsSetCmd.Flags().String("enable", "", "启用: on|off|1|0")
}
