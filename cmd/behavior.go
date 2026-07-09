package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// ============================================================
// ipgroup 命令
// ============================================================

var ipgroupCmd = &cobra.Command{
	Use:   "ipgroup",
	Short: "IP地址组管理",
}

var ipgroupListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看IP地址组列表",
	Long: `查看IP地址组列表，支持分页。

ID 列为 API 内部名称（如 rule_ipgroup_xxx），用于 del 操作。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		start := (page - 1) * pageSize
		end := page*pageSize - 1

		items, total, err := client.GetIPGroups(start, end)
		if err != nil {
			return err
		}

		type IPGroupRow struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Flag    string `json:"flag"`
			Comment string `json:"comment"`
		}

		rows := make([]IPGroupRow, 0, len(items))
		for _, item := range items {
			rows = append(rows, IPGroupRow{
				ID:      item.DotName, // API 内部名称如 rule_ipgroup_xxx
				Name:    item.Name,    // 显示名称如 IPGROUP_LAN
				Flag:    item.Flag,
				Comment: decodeURL(item.Comment),
			})
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- id: %s\n  name: %s\n  flag: %s\n  comment: %s\n", r.ID, r.Name, r.Flag, r.Comment)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("IP地址组为空")
				return nil
			}
			headers := []string{"ID", "NAME", "FLAG", "COMMENT"}
			printTable(headers, rows, func(r IPGroupRow) []string {
				return []string{r.ID, r.Name, flagLabel(r.Flag), r.Comment}
			})
			fmt.Printf("\n第 %d 页, 共 %d 条\n", page, total)
			return nil
		}
	},
}

var ipgroupAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加IP地址组",
	Long: `添加IP地址组（两步操作：先添加IP范围，再创建IP组引用该范围）。

示例:
  tplink ipgroup add --name "test_grp" --scope-type range --scope "192.168.0.20-192.168.0.25"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		groupName, _ := cmd.Flags().GetString("name")
		scopeType, _ := cmd.Flags().GetString("scope-type")
		scope, _ := cmd.Flags().GetString("scope")

		if groupName == "" || scopeType == "" || scope == "" {
			return fmt.Errorf("--name, --scope-type, --scope 均为必填")
		}

		// 第一步：添加IP范围（内部名称用 scope_xxx 格式）
		scopeInternalName := "scope_" + timeBasedSuffix()
		scopeID, err := client.AddIPGroupScope(scopeInternalName, scopeType, scope)
		if err != nil {
			return err
		}
		fmt.Printf("IP范围创建成功, ID: %s\n", scopeID)

		// 第二步：添加IP地址组，引用该范围
		groupID, err := client.AddIPGroup(groupName, scopeInternalName)
		if err != nil {
			return fmt.Errorf("IP范围已创建(%s), 但添加IP组失败: %w", scopeID, err)
		}
		fmt.Printf("IP地址组创建成功, ID: %s\n", groupID)

		return nil
	},
}

var ipgroupDelCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "删除IP地址组或IP范围",
	Long: `根据ID删除IP地址组或IP范围。ID 来自 list 输出。

支持的ID格式:
  rule_ipgroup_xxx   (删除IP地址组)
  rule_ipscope_xxx   (删除IP范围)

示例:
  tplink ipgroup del rule_ipgroup_1782972586
  tplink ipgroup del rule_ipscope_17829725841`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		if err := client.DelIPGroup(args[0]); err != nil {
			return err
		}
		fmt.Printf("删除成功, ID: %s\n", args[0])
		return nil
	},
}

// ============================================================
// timerange 命令
// ============================================================

var timerangeCmd = &cobra.Command{
	Use:   "timerange",
	Short: "时间段管理",
}

var timerangeListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看时间段列表",
	Long: `查看时间段列表，支持分页。

ID 列为 API 内部名称（如 time_obj_xxx），用于 del 操作。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		start := (page - 1) * pageSize
		end := page*pageSize - 1

		items, total, err := client.GetTimeRanges(start, end)
		if err != nil {
			return err
		}

		type TimeRangeRow struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Flag        string `json:"flag"`
			Mode        string `json:"mode"`
			Weekday     string `json:"weekday"`
			TimeSection string `json:"time_section"`
			Comment     string `json:"comment"`
		}

		rows := make([]TimeRangeRow, 0, len(items))
		for _, item := range items {
			var timeSecStr string
			for i, s := range item.TimeSection {
				if i > 0 {
					timeSecStr += "; "
				}
				timeSecStr += decodeURL(s)
			}
			rows = append(rows, TimeRangeRow{
				ID:          item.DotName, // API 内部名称如 time_obj_xxx
				Name:        item.Name,    // 显示名称
				Flag:        item.Flag,
				Mode:        item.Mode,
				Weekday:     weekdayLabel(item.Weekday),
				TimeSection: timeSecStr,
				Comment:     decodeURL(item.Comment),
			})
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- id: %s\n  name: %s\n  flag: %s\n  mode: %s\n  weekday: %s\n  time_section: %s\n  comment: %s\n",
					r.ID, r.Name, r.Flag, r.Mode, r.Weekday, r.TimeSection, r.Comment)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("时间段列表为空")
				return nil
			}
			headers := []string{"ID", "NAME", "MODE", "WEEKDAY", "TIME_SECTION", "COMMENT"}
			printTable(headers, rows, func(r TimeRangeRow) []string {
				return []string{r.ID, r.Name, r.Mode, r.Weekday, r.TimeSection, r.Comment}
			})
			fmt.Printf("\n第 %d 页, 共 %d 条\n", page, total)
			return nil
		}
	},
}

var timerangeAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加时间段",
	Long: `添加时间段。

weekday 为位掩码表示法:
  1=周日, 2=周一, 4=周二, 8=周三, 16=周四, 32=周五, 64=周六
  多天叠加: 3=周日+周一, 127=所有天

time-section 格式: "HHMM,HHMM" (如 "0100,0359" 表示 01:00-03:59)

示例:
  tplink timerange add --name "工作时间" --time-section "0900,1800" --weekday 2 --comment "工作日"
  tplink timerange add --name "晚间" --time-section "2200,2359" --time-section "0000,0600" --weekday 3 --comment "周末晚间"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		mode, _ := cmd.Flags().GetString("mode")
		weekday, _ := cmd.Flags().GetString("weekday")
		comment, _ := cmd.Flags().GetString("comment")
		timeSection, _ := cmd.Flags().GetStringSlice("time-section")

		if name == "" || weekday == "" || len(timeSection) == 0 {
			return fmt.Errorf("--name, --weekday, --time-section 均为必填")
		}
		if mode == "" {
			mode = "manual"
		}

		id, err := client.AddTimeRange(name, mode, timeSection, weekday, comment)
		if err != nil {
			return err
		}
		fmt.Printf("时间段创建成功, ID: %s\n", id)
		return nil
	},
}

var timerangeDelCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "删除时间段",
	Long: `根据ID删除时间段。ID 来自 list 输出。

示例:
  tplink timerange del time_obj_1782960315`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		if err := client.DelTimeRange(args[0]); err != nil {
			return err
		}
		fmt.Printf("删除成功, ID: %s\n", args[0])
		return nil
	},
}

// ============================================================
// qos 命令
// ============================================================

var qosCmd = &cobra.Command{
	Use:   "qos",
	Short: "带宽控制(Qos)管理",
}

// --- qos config ---

var qosConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Qos配置管理",
}

var qosConfigListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看Qos配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		setting, err := client.GetQosConfig()
		if err != nil {
			return err
		}

		type QosConfigRow struct {
			QosEnable       string `json:"qos_enable"`
			ThresholdEnable string `json:"threshold_enable"`
			QosThreshold    string `json:"qos_threshold"`
			Interface       string `json:"interface"`
		}

		iface := ""
		for i, s := range setting.Interface {
			if i > 0 {
				iface += ", "
			}
			iface += s
		}

		row := QosConfigRow{
			QosEnable:       statusLabel(setting.QosEnable),
			ThresholdEnable: statusLabel(setting.ThresholdEnable),
			QosThreshold:    setting.QosThreshold,
			Interface:       iface,
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(row)
		case "yaml":
			fmt.Printf("qos_enable: %s\nthreshold_enable: %s\nqos_threshold: %s\ninterface: %s\n",
				row.QosEnable, row.ThresholdEnable, row.QosThreshold, row.Interface)
			return nil
		default:
			headers := []string{"QOS_ENABLE", "THRESHOLD_ENABLE", "THRESHOLD", "INTERFACE"}
			printTable(headers, []QosConfigRow{row}, func(r QosConfigRow) []string {
				return []string{r.QosEnable, r.ThresholdEnable, r.QosThreshold, r.Interface}
			})
			return nil
		}
	},
}

var qosConfigSetCmd = &cobra.Command{
	Use:   "set",
	Short: "设置Qos配置",
	Long: `设置带宽控制(Qos)配置。只发送用户指定的字段。

示例:
  tplink qos config set --qos-enable on
  tplink qos config set --qos-enable on --threshold-enable off --qos-threshold 80`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		cfg := map[string]interface{}{}

		f := cmd.Flags().Lookup("qos-enable")
		if f != nil && f.Changed {
			cfg["qos_enable"] = boolToOnOff(f.Value.String())
		}

		f = cmd.Flags().Lookup("threshold-enable")
		if f != nil && f.Changed {
			cfg["threshold_enable"] = boolToOnOff(f.Value.String())
		}

		f = cmd.Flags().Lookup("qos-threshold")
		if f != nil && f.Changed {
			cfg["qos_threshold"] = f.Value.String()
		}

		if len(cfg) == 0 {
			return fmt.Errorf("请至少指定一个要修改的参数")
		}

		if err := client.SetQosConfig(cfg); err != nil {
			return err
		}
		fmt.Println("Qos配置更新成功")
		return nil
	},
}

// --- qos rule ---

var qosRuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Qos规则管理",
}

var qosRuleListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看Qos规则列表",
	Long: `查看带宽控制规则列表，支持分页。

ID 列为 API 内部名称（如 rule_xxx），用于 del 操作。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		start := (page - 1) * pageSize
		end := page*pageSize - 1

		items, total, err := client.GetQosRules(start, end)
		if err != nil {
			return err
		}

		type QosRuleRow struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Enable      string `json:"enable"`
			Mode        string `json:"mode"`
			IPGroup     string `json:"ip_group"`
			RateMax     string `json:"rate_max"`
			RateUp      string `json:"rate_max_mate"`
			Time        string `json:"time"`
			IfPing      string `json:"if_ping"`
			IfPong      string `json:"if_pong"`
		}

		rows := make([]QosRuleRow, 0, len(items))
		for _, item := range items {
			rows = append(rows, QosRuleRow{
				ID:      item.DotName,
				Name:    item.Name,
				Enable:  statusLabel(item.Enable),
				Mode:    modeLabel(item.Mode),
				IPGroup: item.IPGroup,
				RateMax: item.RateMax,
				RateUp:  item.RateMaxMate,
				Time:    item.Time,
				IfPing:  item.IfPing,
				IfPong:  item.IfPong,
			})
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- id: %s\n  name: %s\n  enable: %s\n  mode: %s\n  ip_group: %s\n  rate_max: %s\n  rate_max_mate: %s\n  time: %s\n  if_ping: %s\n  if_pong: %s\n",
					r.ID, r.Name, r.Enable, r.Mode, r.IPGroup, r.RateMax, r.RateUp, r.Time, r.IfPing, r.IfPong)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("Qos规则为空")
				return nil
			}
			headers := []string{"ID", "NAME", "ENABLE", "MODE", "IP_GROUP", "RATE_DL", "RATE_UL", "TIME"}
			printTable(headers, rows, func(r QosRuleRow) []string {
				return []string{r.ID, r.Name, r.Enable, r.Mode, r.IPGroup, r.RateMax, r.RateUp, r.Time}
			})
			fmt.Printf("\n第 %d 页, 共 %d 条\n", page, total)
			return nil
		}
	},
}

var qosRuleAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加Qos规则",
	Long: `添加带宽控制规则。

带宽模式(mode): share(共享) | priv(独立)
IP类型(ip-type): src(源IP) | dest(目的IP)

示例:
  tplink qos rule add --name "限速规则" --if-ping LAN --if-pong WAN_ALL --ip-group test_ipgroup_1 --rate-max 1000 --rate-max-mate 500 --mode share --time Any`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		ifPing, _ := cmd.Flags().GetString("if-ping")
		ifPong, _ := cmd.Flags().GetString("if-pong")
		ipGroup, _ := cmd.Flags().GetString("ip-group")
		rateMax, _ := cmd.Flags().GetString("rate-max")
		rateMaxMate, _ := cmd.Flags().GetString("rate-max-mate")
		mode, _ := cmd.Flags().GetString("mode")
		time, _ := cmd.Flags().GetString("time")
		comment, _ := cmd.Flags().GetString("comment")
		enable, _ := cmd.Flags().GetString("enable")
		ipType, _ := cmd.Flags().GetString("ip-type")

		if name == "" || ifPing == "" || ifPong == "" || ipGroup == "" || rateMax == "" || rateMaxMate == "" {
			return fmt.Errorf("--name, --if-ping, --if-pong, --ip-group, --rate-max, --rate-max-mate 均为必填")
		}
		if mode == "" {
			mode = "share"
		}
		if time == "" {
			time = "Any"
		}
		if ipType == "" {
			ipType = "src"
		}
		if enable == "" {
			enable = "on"
		}

		para := map[string]interface{}{
			"name":           name,
			"if_ping":        ifPing,
			"if_pong":        ifPong,
			"ip_group":       ipGroup,
			"rate_max":       rateMax,
			"rate_max_mate":  rateMaxMate,
			"mode":           mode,
			"time":           time,
			"comment":        comment,
			"enable":         enable,
			"ip_type":        ipType,
			"position":       "",
		}

		id, err := client.AddQosRule(para)
		if err != nil {
			return err
		}
		fmt.Printf("Qos规则创建成功, ID: %s\n", id)
		return nil
	},
}

var qosRuleDelCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "删除Qos规则",
	Long: `根据ID删除Qos规则。ID 来自 list 输出。

示例:
  tplink qos rule del rule_1782974217`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		if err := client.DelQosRule(args[0]); err != nil {
			return err
		}
		fmt.Printf("删除成功, ID: %s\n", args[0])
		return nil
	},
}

// ============================================================
// acl 命令
// ============================================================

var aclCmd = &cobra.Command{
	Use:   "acl",
	Short: "访问控制(ACL)管理",
}

// --- acl rule ---

var aclRuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "ACL规则管理",
}

var aclRuleListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看ACL规则列表",
	Long: `查看访问控制规则列表，支持分页。

ID 列为 API 内部名称（如 rule_acl_inner_xxx），用于 del 操作。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		start := (page - 1) * pageSize
		end := page*pageSize - 1

		items, total, err := client.GetACLRules(start, end)
		if err != nil {
			return err
		}

		type ACLRuleRow struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Policy   string `json:"policy"`
			Service  string `json:"service"`
			Zone     string `json:"zone"`
			Src      string `json:"src"`
			Dest     string `json:"dest"`
			Time     string `json:"time"`
			Position int    `json:"position"`
		}

		rows := make([]ACLRuleRow, 0, len(items))
		for _, item := range items {
			rows = append(rows, ACLRuleRow{
				ID:       item.Name,
				Name:     item.Name,
				Policy:   policyLabel(item.Policy),
				Service:  item.Service,
				Zone:     item.Zone,
				Src:      item.Src,
				Dest:     item.Dest,
				Time:     item.Time,
				Position: item.Position,
			})
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- id: %s\n  name: %s\n  policy: %s\n  service: %s\n  zone: %s\n  src: %s\n  dest: %s\n  time: %s\n  position: %d\n",
					r.ID, r.Name, r.Policy, r.Service, r.Zone, r.Src, r.Dest, r.Time, r.Position)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("ACL规则为空")
				return nil
			}
			headers := []string{"ID", "NAME", "POLICY", "SRC", "DEST", "SERVICE", "ZONE", "TIME"}
			printTable(headers, rows, func(r ACLRuleRow) []string {
				return []string{r.ID, r.Name, r.Policy, r.Src, r.Dest, r.Service, r.Zone, r.Time}
			})
			fmt.Printf("\n第 %d 页, 共 %d 条\n", page, total)
			return nil
		}
	},
}

var aclRuleAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加ACL规则",
	Long: `添加访问控制规则（两步操作：先创建服务，再创建规则引用该服务）。

策略(policy): DROP(阻塞) | ACCEPT(允许)
协议(proto): tcp | udp | tcp-udp
端口范围格式: "起始-结束" (如 "2800-2811" 或 "1-65535")

示例:
  tplink acl rule add --name "阻止访问" --policy DROP --zone LAN --src test_ipgroup_1 --dest ISP_CHINA_TELECOM --proto tcp-udp --sport "1-65535" --dport "2800-2811"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		policy, _ := cmd.Flags().GetString("policy")
		zone, _ := cmd.Flags().GetString("zone")
		src, _ := cmd.Flags().GetString("src")
		dest, _ := cmd.Flags().GetString("dest")
		time, _ := cmd.Flags().GetString("time")
		proto, _ := cmd.Flags().GetString("proto")
		sport, _ := cmd.Flags().GetString("sport")
		dport, _ := cmd.Flags().GetString("dport")

		if name == "" || policy == "" || zone == "" || proto == "" || sport == "" || dport == "" {
			return fmt.Errorf("--name, --policy, --zone, --proto, --sport, --dport 均为必填")
		}
		if time == "" {
			time = "Any"
		}

		// 第一步：创建服务（内部名称用 _actl_xxx 格式）
		serviceName := "_actl_" + timeBasedSuffix()
		svcID, err := client.AddACLService(serviceName, proto, sport, dport)
		if err != nil {
			return err
		}
		fmt.Printf("ACL服务创建成功, ID: %s\n", svcID)

		// 第二步：创建规则
		para := map[string]interface{}{
			"name":     name,
			"policy":   policy,
			"service":  serviceName,
			"zone":     zone,
			"src":      src,
			"dest":     dest,
			"time":     time,
			"user":     "1",
			"position": "",
		}

		ruleID, err := client.AddACLRule(para)
		if err != nil {
			return fmt.Errorf("ACL服务已创建(%s), 但添加规则失败: %w", svcID, err)
		}
		fmt.Printf("ACL规则创建成功, ID: %s\n", ruleID)

		return nil
	},
}

var aclRuleDelCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "删除ACL规则或服务",
	Long: `根据ID删除ACL规则或关联的服务。ID 来自 list 输出。

删除规则(rule_acl_inner_xxx)时会自动执行两步操作:
  1. 从 access_ctl 表删除规则
  2. 从 service 表删除关联的服务

删除服务(service_xxx)时只执行一步操作。

示例:
  tplink acl rule del rule_acl_inner_1782974989
  tplink acl rule del service_17829749881`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		id := args[0]

		// 根据ID前缀判断
		if len(id) > 8 && id[:8] == "service_" {
			// 直接删除服务
			if err := client.DelACLService(id); err != nil {
				return err
			}
			fmt.Printf("ACL服务删除成功, ID: %s\n", id)
			return nil
		}

		// 删除规则: 需要先查询规则获取关联的 service 名称
		rules, _, err := client.GetACLRules(0, 499)
		if err != nil {
			return err
		}

		var serviceName string
		found := false
		for _, rule := range rules {
			if rule.DotName == id {
				serviceName = rule.Service
				found = true
				break
			}
		}

		// dry-run 模式下，GET 返回空列表，用模拟数据继续
		if !found && client.DryRun {
			serviceName = "service_dryrun"
			found = true
		}
		if !found {
			return fmt.Errorf("未找到ID为 %s 的ACL规则", id)
		}

		// 第一步: 删除 access_ctl 规则
		if err := client.DelACLRule(id); err != nil {
			return err
		}
		fmt.Printf("ACL规则删除成功, ID: %s\n", id)

		// 第二步: 删除关联的 service
		if serviceName != "" {
			if err := client.DelACLService(serviceName); err != nil {
				return fmt.Errorf("规则已删除, 但删除关联服务失败: %w", err)
			}
			fmt.Printf("关联服务删除成功, ID: %s\n", serviceName)
		}

		return nil
	},
}

// --- acl service ---

var aclServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "ACL服务列表",
}

var aclServiceListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看ACL服务列表",
	Long: `查看访问控制服务列表（含系统预置和自定义服务）。

ID 列为 API 内部名称（如 service_xxx），用于 del 操作。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		items, total, err := client.GetACLServices()
		if err != nil {
			return err
		}

		type ACLServiceRow struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Flag    string `json:"flag"`
			Proto   string `json:"proto"`
			SPort   string `json:"sport"`
			DPort   string `json:"dport"`
			Comment string `json:"comment"`
		}

		rows := make([]ACLServiceRow, 0, len(items))
		for _, item := range items {
			rows = append(rows, ACLServiceRow{
				ID:      item.Name,
				Name:    item.Name,
				Flag:    item.Flag,
				Proto:   item.Proto,
				SPort:   item.SPort,
				DPort:   item.DPort,
				Comment: item.Comment,
			})
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- id: %s\n  name: %s\n  flag: %s\n  proto: %s\n  sport: %s\n  dport: %s\n  comment: %s\n",
					r.ID, r.Name, r.Flag, r.Proto, r.SPort, r.DPort, r.Comment)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("ACL服务列表为空")
				return nil
			}
			headers := []string{"ID", "NAME", "FLAG", "PROTO", "SPORT", "DPORT", "COMMENT"}
			printTable(headers, rows, func(r ACLServiceRow) []string {
				return []string{r.ID, r.Name, flagLabel(r.Flag), r.Proto, r.SPort, r.DPort, r.Comment}
			})
			fmt.Printf("\n共 %d 条\n", total)
			return nil
		}
	},
}

// ============================================================
// 辅助函数
// ============================================================

// boolToOnOff 将布尔值标志转为 TP-Link 的 on/off 字符串
func boolToOnOff(v string) string {
	switch v {
	case "on", "1", "yes", "true":
		return "on"
	default:
		return "off"
	}
}

// flagLabel 将 flag 字段转为中文
func flagLabel(s string) string {
	switch s {
	case "system":
		return "系统"
	case "user":
		return "用户"
	case "inner":
		return "内部"
	default:
		return s
	}
}

// policyLabel 将策略转为中文
func policyLabel(s string) string {
	switch s {
	case "DROP":
		return "阻塞"
	case "ACCEPT":
		return "允许"
	default:
		return s
	}
}

// modeLabel 将Qos模式转为中文
func modeLabel(s string) string {
	switch s {
	case "share":
		return "共享"
	case "priv":
		return "独立"
	default:
		return s
	}
}

// weekdayLabel 将位掩码weekday转为可读字符串
func weekdayLabel(s string) string {
	names := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}

	var n int
	fmt.Sscanf(s, "%d", &n)
	if n <= 0 || n > 127 {
		return s
	}

	var parts []string
	for i := 0; i < 7; i++ {
		if n&(1<<i) != 0 {
			parts = append(parts, names[i])
		}
	}
	if len(parts) == 0 {
		return s
	}
	if len(parts) == 7 {
		return "每天"
	}

	result := ""
	for i, p := range parts {
		if i > 0 {
			result += ","
		}
		result += p
	}
	return result
}

// timeBasedSuffix 返回基于当前时间的后缀字符串（用于生成唯一内部名称）
func timeBasedSuffix() string {
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}

// ============================================================
// init 注册
// ============================================================

func init() {
	// --- ipgroup ---
	ipgroupListCmd.Flags().Int("page", 1, "页码")
	ipgroupListCmd.Flags().Int("page-size", 10, "每页条数")
	ipgroupAddCmd.Flags().StringP("name", "n", "", "IP组显示名称 (必填)")
	ipgroupAddCmd.Flags().String("scope-type", "", "范围类型: range|port (必填)")
	ipgroupAddCmd.Flags().String("scope", "", "IP范围, 如 192.168.0.20-192.168.0.25 (必填)")

	ipgroupCmd.AddCommand(ipgroupListCmd, ipgroupAddCmd, ipgroupDelCmd)
	rootCmd.AddCommand(ipgroupCmd)

	// --- timerange ---
	timerangeListCmd.Flags().Int("page", 1, "页码")
	timerangeListCmd.Flags().Int("page-size", 10, "每页条数")
	timerangeAddCmd.Flags().StringP("name", "n", "", "时间段名称 (必填)")
	timerangeAddCmd.Flags().String("mode", "manual", "时间模式: manual")
	timerangeAddCmd.Flags().String("weekday", "", "星期位掩码 (必填, 1=周日,2=周一,...,64=周六)")
	timerangeAddCmd.Flags().StringP("comment", "c", "", "备注说明")
	timerangeAddCmd.Flags().StringSliceP("time-section", "t", nil, "时间段, 格式 HHMM,HHMM (可多次指定, 必填)")

	timerangeCmd.AddCommand(timerangeListCmd, timerangeAddCmd, timerangeDelCmd)
	rootCmd.AddCommand(timerangeCmd)

	// --- qos ---
	qosConfigSetCmd.Flags().String("qos-enable", "", "启用Qos: on|off")
	qosConfigSetCmd.Flags().String("threshold-enable", "", "启用阈值控制: on|off")
	qosConfigSetCmd.Flags().String("qos-threshold", "", "带宽利用率阈值(%)")

	qosConfigCmd.AddCommand(qosConfigListCmd, qosConfigSetCmd)

	qosRuleListCmd.Flags().Int("page", 1, "页码")
	qosRuleListCmd.Flags().Int("page-size", 10, "每页条数")
	qosRuleAddCmd.Flags().StringP("name", "n", "", "规则名称 (必填)")
	qosRuleAddCmd.Flags().String("if-ping", "", "入口接口 (必填, 如 LAN)")
	qosRuleAddCmd.Flags().String("if-pong", "", "出口接口 (必填, 如 WAN_ALL)")
	qosRuleAddCmd.Flags().String("ip-group", "", "IP地址组名称 (必填)")
	qosRuleAddCmd.Flags().String("rate-max", "", "下行带宽(KB/s) (必填)")
	qosRuleAddCmd.Flags().String("rate-max-mate", "", "上行带宽(KB/s) (必填)")
	qosRuleAddCmd.Flags().String("mode", "share", "带宽模式: share|priv")
	qosRuleAddCmd.Flags().StringP("time", "t", "Any", "时间段名称")
	qosRuleAddCmd.Flags().StringP("comment", "c", "", "备注说明")
	qosRuleAddCmd.Flags().String("enable", "on", "启用: on|off")
	qosRuleAddCmd.Flags().String("ip-type", "src", "IP类型: src|dest")

	qosRuleCmd.AddCommand(qosRuleListCmd, qosRuleAddCmd, qosRuleDelCmd)
	qosCmd.AddCommand(qosConfigCmd, qosRuleCmd)
	rootCmd.AddCommand(qosCmd)

	// --- acl ---
	aclRuleListCmd.Flags().Int("page", 1, "页码")
	aclRuleListCmd.Flags().Int("page-size", 10, "每页条数")
	aclRuleAddCmd.Flags().StringP("name", "n", "", "规则名称 (必填)")
	aclRuleAddCmd.Flags().String("policy", "", "策略: DROP|ACCEPT (必填)")
	aclRuleAddCmd.Flags().String("zone", "", "区域 (必填, 如 LAN)")
	aclRuleAddCmd.Flags().StringP("src", "s", "", "源IP地址组名称")
	aclRuleAddCmd.Flags().StringP("dest", "d", "", "目的IP地址组名称")
	aclRuleAddCmd.Flags().StringP("time", "t", "Any", "时间段名称")
	aclRuleAddCmd.Flags().String("proto", "", "协议: tcp|udp|tcp-udp (必填)")
	aclRuleAddCmd.Flags().String("sport", "", "源端口范围, 如 1-65535 (必填)")
	aclRuleAddCmd.Flags().String("dport", "", "目的端口范围, 如 2800-2811 (必填)")

	aclRuleCmd.AddCommand(aclRuleListCmd, aclRuleAddCmd, aclRuleDelCmd)
	aclServiceCmd.AddCommand(aclServiceListCmd)
	aclCmd.AddCommand(aclRuleCmd, aclServiceCmd)
	rootCmd.AddCommand(aclCmd)
}
