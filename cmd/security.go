package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ljw/tplink-cli/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	// arp
	arpCmd.AddCommand(arpConfigCmd)
	arpConfigCmd.AddCommand(arpConfigListCmd)
	arpConfigCmd.AddCommand(arpConfigSetCmd)
	arpBindCmd.AddCommand(arpBindListCmd)
	arpCmd.AddCommand(arpBindCmd)
	rootCmd.AddCommand(arpCmd)

	// macfilter
	macFilterCmd.AddCommand(macFilterConfigCmd)
	macFilterConfigCmd.AddCommand(macFilterConfigListCmd)
	macFilterConfigCmd.AddCommand(macFilterConfigSetCmd)
	macFilterCmd.AddCommand(macFilterRuleCmd)
	macFilterRuleCmd.AddCommand(macFilterRuleListCmd)
	macFilterRuleCmd.AddCommand(macFilterRuleAddCmd)
	macFilterRuleCmd.AddCommand(macFilterRuleDelCmd)
	rootCmd.AddCommand(macFilterCmd)

	// dos
	dosCmd.AddCommand(dosListCmd)
	dosCmd.AddCommand(dosSetCmd)
	rootCmd.AddCommand(dosCmd)

	// flood
	floodCmd.AddCommand(floodListCmd)
	floodCmd.AddCommand(floodSetCmd)
	rootCmd.AddCommand(floodCmd)
}

// ========== ARP 防护 ==========

var arpCmd = &cobra.Command{
	Use:   "arp",
	Short: "ARP防护管理",
}

var arpConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "ARP防护配置",
}

var arpConfigListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看ARP防护配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		cfg, err := client.GetArpConfig()
		if err != nil {
			return err
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(cfg)
		case "yaml":
			fmt.Printf("enable: %s\n", statusLabel(cfg.Enable))
			fmt.Printf("garp: %s\n", statusLabel(cfg.Garp))
			fmt.Printf("imb_pass: %s\n", statusLabel(cfg.ImbPass))
			fmt.Printf("log_enable: %s\n", statusLabel(cfg.LogEnable))
			fmt.Printf("interval: %s\n", cfg.Interval)
			fmt.Printf("interface: %v\n", cfg.Interface)
			return nil
		default:
			type ArpCfgRow struct {
				Enable    string `json:"enable"`
				Garp      string `json:"garp"`
				ImbPass   string `json:"imb_pass"`
				LogEnable string `json:"log_enable"`
				Interval  string `json:"interval"`
				Interface string `json:"interface"`
			}
			rows := []ArpCfgRow{{
				Enable:    statusLabel(cfg.Enable),
				Garp:      statusLabel(cfg.Garp),
				ImbPass:   statusLabel(cfg.ImbPass),
				LogEnable: statusLabel(cfg.LogEnable),
				Interval:  cfg.Interval,
				Interface: fmt.Sprintf("%v", cfg.Interface),
			}}
			headers := []string{"ENABLE", "GARP", "IMB_PASS", "LOG", "INTERVAL(ms)", "INTERFACE"}
			printTable(headers, rows, func(r ArpCfgRow) []string {
				return []string{r.Enable, r.Garp, r.ImbPass, r.LogEnable, r.Interval, r.Interface}
			})
			return nil
		}
	},
}

var arpConfigSetCmd = &cobra.Command{
	Use:   "set",
	Short: "修改ARP防护配置",
	Long: `修改ARP防护配置。只会发送用户显式指定的字段。

示例:
  tplink arp config set --enable on
  tplink arp config set --garp off --imb-pass on`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		current, err := client.GetArpConfig()
		if err != nil {
			return err
		}

		cfg := map[string]interface{}{}
		mergeFieldBool(cmd, "enable", cfg)
		mergeFieldBool(cmd, "garp", cfg)
		mergeFieldBool(cmd, "imb-pass", cfg)

		if len(cfg) == 0 {
			return fmt.Errorf("没有指定要修改的字段")
		}

		_ = current // 仅用于 dry-run 展示 GET 请求

		return client.SetArpConfig(cfg)
	},
}

var arpBindCmd = &cobra.Command{
	Use:   "bind",
	Short: "ARP绑定管理",
}

var arpBindListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看ARP绑定列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		start := (page - 1) * pageSize
		end := page*pageSize - 1

		items, total, err := client.GetArpBindList(start, end)
		if err != nil {
			return err
		}

		type BindRow struct {
			ID        string `json:"id"`
			IP        string `json:"ip"`
			Mac       string `json:"mac"`
			Hostname  string `json:"hostname"`
			Status    string `json:"status"`
			Interface string `json:"interface"`
		}

		rows := make([]BindRow, 0, len(items))
		for _, item := range items {
			rows = append(rows, BindRow{
				ID:        item.DotName,
				IP:        item.IP,
				Mac:       item.Mac,
				Hostname:  item.Hostname,
				Status:    arpStatusLabel(item.Status),
				Interface: item.Interface,
			})
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- id: %s\n  ip: %s\n  mac: %s\n  hostname: %s\n  status: %s\n  interface: %s\n",
					r.ID, r.IP, r.Mac, r.Hostname, r.Status, r.Interface)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("ARP列表为空")
				return nil
			}
			headers := []string{"ID", "IP", "MAC", "HOSTNAME", "STATUS", "INTERFACE"}
			printTable(headers, rows, func(r BindRow) []string {
				return []string{r.ID, r.IP, r.Mac, r.Hostname, r.Status, r.Interface}
			})
			if total > 0 && len(rows) < total {
				fmt.Printf("\n# page: %d/%d, total: %d\n", page, (total+pageSize-1)/pageSize, total)
			} else {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		}
	},
}

// ========== MAC 地址过滤 ==========

var macFilterCmd = &cobra.Command{
	Use:   "macfilter",
	Short: "MAC地址过滤管理",
}

var macFilterConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "MAC过滤配置",
}

var macFilterConfigListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看MAC地址过滤配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		cfg, err := client.GetMacFilterConfig()
		if err != nil {
			return err
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(cfg)
		case "yaml":
			fmt.Printf("enable: %s\n", statusLabel(cfg.Enable))
			fmt.Printf("filter_mode: %s\n", modeNameLabel(cfg.FilterMode))
			fmt.Printf("interfaces: %s\n", cfg.Interfaces)
			return nil
		default:
			type McRow struct {
				Enable     string `json:"enable"`
				FilterMode string `json:"filter_mode"`
				Interfaces string `json:"interfaces"`
			}
			rows := []McRow{{
				Enable:     statusLabel(cfg.Enable),
				FilterMode: modeNameLabel(cfg.FilterMode),
				Interfaces: cfg.Interfaces,
			}}
			headers := []string{"ENABLE", "MODE", "INTERFACES"}
			printTable(headers, rows, func(r McRow) []string {
				return []string{r.Enable, r.FilterMode, r.Interfaces}
			})
			return nil
		}
	},
}

var macFilterConfigSetCmd = &cobra.Command{
	Use:   "set",
	Short: "修改MAC地址过滤配置",
	Long: `修改MAC地址过滤配置。只会发送用户显式指定的字段。

示例:
  tplink macfilter config set --enable on --filter-mode white
  tplink macfilter config set --interfaces LAN`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		current, err := client.GetMacFilterConfig()
		if err != nil {
			return err
		}

		cfg := map[string]interface{}{}
		mergeFieldBool(cmd, "enable", cfg)
		mergeFieldStr(cmd, "filter-mode", current.FilterMode, cfg)
		mergeFieldStr(cmd, "interfaces", current.Interfaces, cfg)

		if len(cfg) == 0 {
			return fmt.Errorf("没有指定要修改的字段")
		}

		return client.SetMacFilterConfig(cfg)
	},
}

var macFilterRuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "MAC过滤规则管理",
}

var macFilterRuleListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看MAC过滤规则列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		start := (page - 1) * pageSize
		end := page*pageSize - 1

		items, total, err := client.GetMacFilterRules(start, end)
		if err != nil {
			return err
		}

		type RuleRow struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Mac  string `json:"mac"`
		}

		rows := make([]RuleRow, 0, len(items))
		for _, item := range items {
			rows = append(rows, RuleRow{
				ID:   item.DotName,
				Name: decodeURL(item.Name),
				Mac:  item.Mac,
			})
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- id: %s\n  name: %s\n  mac: %s\n", r.ID, r.Name, r.Mac)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("MAC过滤规则为空")
				return nil
			}
			headers := []string{"ID", "NAME", "MAC"}
			printTable(headers, rows, func(r RuleRow) []string {
				return []string{r.ID, r.Name, r.Mac}
			})
			if total > 0 && len(rows) < total {
				fmt.Printf("\n# page: %d/%d, total: %d\n", page, (total+pageSize-1)/pageSize, total)
			} else {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		}
	},
}

var macFilterRuleAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加MAC过滤规则",
	Long: `添加MAC地址过滤规则。

示例:
  tplink macfilter rule add --name NoAccess --mac "4C-CC-6A-C0-6F-C1"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		mac, _ := cmd.Flags().GetString("mac")

		if name == "" || mac == "" {
			return fmt.Errorf("--name 和 --mac 参数必填")
		}

		id, err := client.AddMacFilterRule(name, mac)
		if err != nil {
			return err
		}
		fmt.Printf("MAC过滤规则添加成功, ID: %s\n", id)
		return nil
	},
}

var macFilterRuleDelCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "删除MAC过滤规则",
	Long: `根据ID删除MAC过滤规则。ID 来自 list 输出的 ID 列。

示例:
  tplink macfilter rule del mac_filter_list_1782977516`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		if err := client.DelMacFilterRule(args[0]); err != nil {
			return err
		}
		fmt.Printf("MAC过滤规则删除成功, ID: %s\n", args[0])
		return nil
	},
}

// ========== DoS 攻击防护 ==========

var dosCmd = &cobra.Command{
	Use:   "dos",
	Short: "DoS攻击防护管理",
}

var dosListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看DoS攻击防护配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		cfg, err := client.GetDosConfig()
		if err != nil {
			return err
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(cfg)
		case "yaml":
			yamlDosConfig(cfg)
			return nil
		default:
			tableDosConfig(cfg)
			return nil
		}
	},
}

var dosSetCmd = &cobra.Command{
	Use:   "set",
	Short: "修改DoS攻击防护配置",
	Long: `修改DoS攻击防护配置。只会发送用户显式指定的字段。

所有字段值为 1(启用) 或 0(禁用)。
ipopt_* 系列字段依赖 ip_option=1 时才生效。

示例:
  tplink dos set --ip-frag 1 --tcp-noflag 1
  tplink dos set --ip-option 1 --ipopt-secure 1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		current, err := client.GetDosConfig()
		if err != nil {
			return err
		}
		_ = current // 仅用于 dry-run 展示 GET 请求

		cfg := map[string]interface{}{}
		dosMergeFlags(cmd, cfg)

		if len(cfg) == 0 {
			return fmt.Errorf("没有指定要修改的字段")
		}

		return client.SetDosConfig(cfg)
	},
}

// ========== Flood 攻击防护 ==========

var floodCmd = &cobra.Command{
	Use:   "flood",
	Short: "Flood攻击防护管理",
}

var floodListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看Flood攻击防护配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		global, threshold, err := client.GetFloodConfig()
		if err != nil {
			return err
		}

		switch output {
		case "json":
			result := map[string]interface{}{
				"global":    global,
				"threshold": threshold,
			}
			return json.NewEncoder(os.Stdout).Encode(result)
		case "yaml":
			yamlFloodConfig(global, threshold)
			return nil
		default:
			tableFloodConfig(global, threshold)
			return nil
		}
	},
}

var floodSetCmd = &cobra.Command{
	Use:   "set",
	Short: "修改Flood攻击防护配置",
	Long: `修改Flood攻击防护配置。只会发送用户显式指定的字段。

分为两部分:
  --xxx-en 系列: 攻击防护开关 (global)，值为 1(启用) 或 0(禁用)
  --xxx-lim / --xxx-bst 系列: 阈值 (threshold)

示例:
  tplink flood set --tcp-conn-en 1 --tcp-conn-lim 3000
  tplink flood set --udp-conn-en 0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		// 先获取当前配置
		_, _, err = client.GetFloodConfig()
		if err != nil {
			return err
		}
		// GET 结果仅用于 dry-run 展示

		globalCfg := map[string]interface{}{}
		thresholdCfg := map[string]interface{}{}

		floodMergeFlags(cmd, globalCfg, thresholdCfg)

		if len(globalCfg) == 0 && len(thresholdCfg) == 0 {
			return fmt.Errorf("没有指定要修改的字段")
		}

		return client.SetFloodConfig(globalCfg, thresholdCfg)
	},
}

// ========== 辅助函数 ==========

// arpStatusLabel ARP 状态标签 (0=未绑定, 1=已绑定, 非标准字段用 "")
func arpStatusLabel(s string) string {
	switch s {
	case "0":
		return "未绑定"
	case "1":
		return "已绑定"
	default:
		return s
	}
}

// modeNameLabel 过滤模式标签
func modeNameLabel(mode string) string {
	switch mode {
	case "white":
		return "白名单"
	case "black":
		return "黑名单"
	default:
		return mode
	}
}

// yamlDosConfig 以 YAML 格式输出 DoS 配置
func yamlDosConfig(cfg *model.DosConfig) {
	fmt.Printf("tcp_winnuke: %s\n", statusLabel(cfg.TcpWinNuke))
	fmt.Printf("ip_frag: %s\n", statusLabel(cfg.IpFrag))
	fmt.Printf("tcp_noflag: %s\n", statusLabel(cfg.TcpNoflag))
	fmt.Printf("ping_death: %s\n", statusLabel(cfg.PingDeath))
	fmt.Printf("ping_large: %s\n", statusLabel(cfg.PingLarge))
	fmt.Printf("tcp_fin_syn: %s\n", statusLabel(cfg.TcpFinSyn))
	fmt.Printf("tcp_fin_noack: %s\n", statusLabel(cfg.TcpFinNoack))
	fmt.Printf("ip_option: %s\n", statusLabel(cfg.IpOption))
	fmt.Printf("ipopt_secure: %s\n", statusLabel(cfg.IpoptSecure))
	fmt.Printf("ipopt_loose_route: %s\n", statusLabel(cfg.IpoptLooseRoute))
	fmt.Printf("ipopt_strict_route: %s\n", statusLabel(cfg.IpoptStrictRoute))
	fmt.Printf("ipopt_record_route: %s\n", statusLabel(cfg.IpoptRecordRoute))
	fmt.Printf("ipopt_stream: %s\n", statusLabel(cfg.IpoptStream))
	fmt.Printf("ipopt_timestamp: %s\n", statusLabel(cfg.IpoptTimestamp))
	fmt.Printf("ipopt_noop: %s\n", statusLabel(cfg.IpoptNoop))
}

// tableDosConfig 以 Table 格式输出 DoS 配置
func tableDosConfig(cfg *model.DosConfig) {
	type DosRow struct {
		Field  string `json:"field"`
		Value  string `json:"value"`
		Remark string `json:"remark"`
	}
	rows := []DosRow{
		{"tcp_winnuke", statusLabel(cfg.TcpWinNuke), "防WinNuke攻击"},
		{"ip_frag", statusLabel(cfg.IpFrag), "防碎片包攻击"},
		{"tcp_noflag", statusLabel(cfg.TcpNoflag), "防TCP Scan"},
		{"ping_death", statusLabel(cfg.PingDeath), "防Ping of Death"},
		{"ping_large", statusLabel(cfg.PingLarge), "防Large Ping"},
		{"tcp_fin_syn", statusLabel(cfg.TcpFinSyn), "阻止FIN+SYN的TCP包"},
		{"tcp_fin_noack", statusLabel(cfg.TcpFinNoack), "阻止仅FIN无ACK的TCP包"},
		{"ip_option", statusLabel(cfg.IpOption), "阻止带选项的包(总开关)"},
		{"ipopt_secure", statusLabel(cfg.IpoptSecure), "安全限制"},
		{"ipopt_loose_route", statusLabel(cfg.IpoptLooseRoute), "宽松选路"},
		{"ipopt_strict_route", statusLabel(cfg.IpoptStrictRoute), "严格选路"},
		{"ipopt_record_route", statusLabel(cfg.IpoptRecordRoute), "记录路径"},
		{"ipopt_stream", statusLabel(cfg.IpoptStream), "流标记"},
		{"ipopt_timestamp", statusLabel(cfg.IpoptTimestamp), "时间戳"},
		{"ipopt_noop", statusLabel(cfg.IpoptNoop), "空标记"},
	}
	headers := []string{"FIELD", "VALUE", "REMARK"}
	printTable(headers, rows, func(r DosRow) []string {
		return []string{r.Field, r.Value, r.Remark}
	})
}

// dosMergeFlags 合并 DoS set 的所有 flag 到 cfg map
// flag 名使用 - 分隔，API 字段名使用 _ 分隔
func dosMergeFlags(cmd *cobra.Command, cfg map[string]interface{}) {
	dosFlags := [][2]string{
		{"tcp-winnuke", "tcp_winnuke"},
		{"ip-frag", "ip_frag"},
		{"tcp-noflag", "tcp_noflag"},
		{"ping-death", "ping_death"},
		{"ping-large", "ping_large"},
		{"tcp-fin-syn", "tcp_fin_syn"},
		{"tcp-fin-noack", "tcp_fin_noack"},
		{"ip-option", "ip_option"},
		{"ipopt-secure", "ipopt_secure"},
		{"ipopt-loose-route", "ipopt_loose_route"},
		{"ipopt-strict-route", "ipopt_strict_route"},
		{"ipopt-record-route", "ipopt_record_route"},
		{"ipopt-stream", "ipopt_stream"},
		{"ipopt-timestamp", "ipopt_timestamp"},
		{"ipopt-noop", "ipopt_noop"},
	}
	for _, pair := range dosFlags {
		mergeFieldBoolMapped(cmd, pair[0], pair[1], cfg)
	}
}

// yamlFloodConfig 以 YAML 格式输出 Flood 配置
func yamlFloodConfig(global *model.FloodGlobal, threshold *model.FloodThreshold) {
	fmt.Println("# 开关配置")
	fmt.Printf("tcp_conn_en: %s\n", statusLabel(global.TcpConnEn))
	fmt.Printf("udp_conn_en: %s\n", statusLabel(global.UdpConnEn))
	fmt.Printf("icmp_conn_en: %s\n", statusLabel(global.IcmpConnEn))
	fmt.Printf("tcp_src_en: %s\n", statusLabel(global.TcpSrcEn))
	fmt.Printf("udp_src_en: %s\n", statusLabel(global.UdpSrcEn))
	fmt.Printf("icmp_src_en: %s\n", statusLabel(global.IcmpSrcEn))
	fmt.Println("# 阈值配置")
	fmt.Printf("tcp_conn_lim: %s\n", threshold.TcpConnLim)
	fmt.Printf("tcp_conn_bst: %s\n", threshold.TcpConnBst)
	fmt.Printf("udp_conn_lim: %s\n", threshold.UdpConnLim)
	fmt.Printf("udp_conn_bst: %s\n", threshold.UdpConnBst)
	fmt.Printf("icmp_conn_lim: %s\n", threshold.IcmpConnLim)
	fmt.Printf("icmp_conn_bst: %s\n", threshold.IcmpConnBst)
	fmt.Printf("tcp_src_lim: %s\n", threshold.TcpSrcLim)
	fmt.Printf("tcp_src_bst: %s\n", threshold.TcpSrcBst)
	fmt.Printf("udp_src_lim: %s\n", threshold.UdpSrcLim)
	fmt.Printf("udp_src_bst: %s\n", threshold.UdpSrcBst)
	fmt.Printf("icmp_src_lim: %s\n", threshold.IcmpSrcLim)
	fmt.Printf("icmp_src_bst: %s\n", threshold.IcmpSrcBst)
}

// tableFloodConfig 以 Table 格式输出 Flood 配置
func tableFloodConfig(global *model.FloodGlobal, threshold *model.FloodThreshold) {
	type FloodRow struct {
		Field  string `json:"field"`
		Value  string `json:"value"`
		Remark string `json:"remark"`
	}
	rows := []FloodRow{
		{"tcp_conn_en", statusLabel(global.TcpConnEn), "防多连接TCP SYN Flood"},
		{"udp_conn_en", statusLabel(global.UdpConnEn), "防多连接UDP Flood"},
		{"icmp_conn_en", statusLabel(global.IcmpConnEn), "防多连接ICMP Flood"},
		{"tcp_src_en", statusLabel(global.TcpSrcEn), "防固定源TCP SYN Flood"},
		{"udp_src_en", statusLabel(global.UdpSrcEn), "防固定源UDP Flood"},
		{"icmp_src_en", statusLabel(global.IcmpSrcEn), "防固定源ICMP Flood"},
	}
	headers := []string{"FIELD", "VALUE", "REMARK"}
	fmt.Println("--- 攻击防护开关 ---")
	printTable(headers, rows, func(r FloodRow) []string {
		return []string{r.Field, r.Value, r.Remark}
	})

	rows2 := []FloodRow{
		{"tcp_conn_lim", threshold.TcpConnLim, "多连接TCP SYN Flood 限制"},
		{"tcp_conn_bst", threshold.TcpConnBst, "多连接TCP SYN Flood 突发"},
		{"udp_conn_lim", threshold.UdpConnLim, "多连接UDP Flood 限制"},
		{"udp_conn_bst", threshold.UdpConnBst, "多连接UDP Flood 突发"},
		{"icmp_conn_lim", threshold.IcmpConnLim, "多连接ICMP Flood 限制"},
		{"icmp_conn_bst", threshold.IcmpConnBst, "多连接ICMP Flood 突发"},
		{"tcp_src_lim", threshold.TcpSrcLim, "固定源TCP SYN Flood 限制"},
		{"tcp_src_bst", threshold.TcpSrcBst, "固定源TCP SYN Flood 突发"},
		{"udp_src_lim", threshold.UdpSrcLim, "固定源UDP Flood 限制"},
		{"udp_src_bst", threshold.UdpSrcBst, "固定源UDP Flood 突发"},
		{"icmp_src_lim", threshold.IcmpSrcLim, "固定源ICMP Flood 限制"},
		{"icmp_src_bst", threshold.IcmpSrcBst, "固定源ICMP Flood 突发"},
	}
	fmt.Println("\n--- 攻击阈值 ---")
	printTable(headers, rows2, func(r FloodRow) []string {
		return []string{r.Field, r.Value, r.Remark}
	})
}

// floodMergeFlags 合并 Flood set 的所有 flag
func floodMergeFlags(cmd *cobra.Command, globalCfg, thresholdCfg map[string]interface{}) {
	// global 开关字段（on/off → 1/0），flag名 → API字段名映射
	globalFlags := [][2]string{
		{"tcp-conn-en", "tcp_conn_en"},
		{"udp-conn-en", "udp_conn_en"},
		{"icmp-conn-en", "icmp_conn_en"},
		{"tcp-src-en", "tcp_src_en"},
		{"udp-src-en", "udp_src_en"},
		{"icmp-src-en", "icmp_src_en"},
	}
	for _, pair := range globalFlags {
		mergeFieldBoolMapped(cmd, pair[0], pair[1], globalCfg)
	}

	// threshold 阈值字段（字符串值）
	thresholdFlags := [][2]string{
		{"tcp-conn-lim", "tcp_conn_lim"}, {"tcp-conn-bst", "tcp_conn_bst"},
		{"udp-conn-lim", "udp_conn_lim"}, {"udp-conn-bst", "udp_conn_bst"},
		{"icmp-conn-lim", "icmp_conn_lim"}, {"icmp-conn-bst", "icmp_conn_bst"},
		{"tcp-src-lim", "tcp_src_lim"}, {"tcp-src-bst", "tcp_src_bst"},
		{"udp-src-lim", "udp_src_lim"}, {"udp-src-bst", "udp_src_bst"},
		{"icmp-src-lim", "icmp_src_lim"}, {"icmp-src-bst", "icmp_src_bst"},
	}
	for _, pair := range thresholdFlags {
		mergeFieldStrMapped(cmd, pair[0], pair[1], "", thresholdCfg)
	}
}

func init() {
	// ARP bind list 支持分页
	arpBindListCmd.Flags().Int("page", 1, "页码")
	arpBindListCmd.Flags().Int("page-size", 50, "每页数量")

	// ARP config set
	arpConfigSetCmd.Flags().String("enable", "", "启用: on|off|1|0")
	arpConfigSetCmd.Flags().String("garp", "", "免费ARP: on|off|1|0")
	arpConfigSetCmd.Flags().String("imb-pass", "", "IMB直通: on|off|1|0")

	// MAC filter config set
	macFilterConfigSetCmd.Flags().String("enable", "", "启用: on|off|1|0")
	macFilterConfigSetCmd.Flags().String("filter-mode", "", "过滤模式: white|black")
	macFilterConfigSetCmd.Flags().String("interfaces", "", "接口名称: LAN")

	// MAC filter rule 分页
	macFilterRuleListCmd.Flags().Int("page", 1, "页码")
	macFilterRuleListCmd.Flags().Int("page-size", 10, "每页数量")

	// MAC filter rule add
	macFilterRuleAddCmd.Flags().StringP("name", "n", "", "规则名称")
	macFilterRuleAddCmd.Flags().StringP("mac", "m", "", "MAC地址 (格式: XX-XX-XX-XX-XX-XX)")

	// DoS set flags
	dosSetCmd.Flags().String("tcp-winnuke", "", "防WinNuke攻击: on|off|1|0")
	dosSetCmd.Flags().String("ip-frag", "", "防碎片包攻击: on|off|1|0")
	dosSetCmd.Flags().String("tcp-noflag", "", "防TCP Scan: on|off|1|0")
	dosSetCmd.Flags().String("ping-death", "", "防Ping of Death: on|off|1|0")
	dosSetCmd.Flags().String("ping-large", "", "防Large Ping: on|off|1|0")
	dosSetCmd.Flags().String("tcp-fin-syn", "", "阻止FIN+SYN的TCP包: on|off|1|0")
	dosSetCmd.Flags().String("tcp-fin-noack", "", "阻止仅FIN无ACK的TCP包: on|off|1|0")
	dosSetCmd.Flags().String("ip-option", "", "阻止带IP选项的包(总开关): on|off|1|0")
	dosSetCmd.Flags().String("ipopt-secure", "", "安全限制(依赖ip-option): on|off|1|0")
	dosSetCmd.Flags().String("ipopt-loose-route", "", "宽松选路(依赖ip-option): on|off|1|0")
	dosSetCmd.Flags().String("ipopt-strict-route", "", "严格选路(依赖ip-option): on|off|1|0")
	dosSetCmd.Flags().String("ipopt-record-route", "", "记录路径(依赖ip-option): on|off|1|0")
	dosSetCmd.Flags().String("ipopt-stream", "", "流标记(依赖ip-option): on|off|1|0")
	dosSetCmd.Flags().String("ipopt-timestamp", "", "时间戳(依赖ip-option): on|off|1|0")
	dosSetCmd.Flags().String("ipopt-noop", "", "空标记(依赖ip-option): on|off|1|0")

	// Flood set flags - global 开关
	floodSetCmd.Flags().String("tcp-conn-en", "", "防多连接TCP SYN Flood: on|off|1|0")
	floodSetCmd.Flags().String("udp-conn-en", "", "防多连接UDP Flood: on|off|1|0")
	floodSetCmd.Flags().String("icmp-conn-en", "", "防多连接ICMP Flood: on|off|1|0")
	floodSetCmd.Flags().String("tcp-src-en", "", "防固定源TCP SYN Flood: on|off|1|0")
	floodSetCmd.Flags().String("udp-src-en", "", "防固定源UDP Flood: on|off|1|0")
	floodSetCmd.Flags().String("icmp-src-en", "", "防固定源ICMP Flood: on|off|1|0")
	// Flood set flags - threshold 阈值
	floodSetCmd.Flags().String("tcp-conn-lim", "", "多连接TCP SYN Flood 限制")
	floodSetCmd.Flags().String("tcp-conn-bst", "", "多连接TCP SYN Flood 突发")
	floodSetCmd.Flags().String("udp-conn-lim", "", "多连接UDP Flood 限制")
	floodSetCmd.Flags().String("udp-conn-bst", "", "多连接UDP Flood 突发")
	floodSetCmd.Flags().String("icmp-conn-lim", "", "多连接ICMP Flood 限制")
	floodSetCmd.Flags().String("icmp-conn-bst", "", "多连接ICMP Flood 突发")
	floodSetCmd.Flags().String("tcp-src-lim", "", "固定源TCP SYN Flood 限制")
	floodSetCmd.Flags().String("tcp-src-bst", "", "固定源TCP SYN Flood 突发")
	floodSetCmd.Flags().String("udp-src-lim", "", "固定源UDP Flood 限制")
	floodSetCmd.Flags().String("udp-src-bst", "", "固定源UDP Flood 突发")
	floodSetCmd.Flags().String("icmp-src-lim", "", "固定源ICMP Flood 限制")
	floodSetCmd.Flags().String("icmp-src-bst", "", "固定源ICMP Flood 突发")
}
