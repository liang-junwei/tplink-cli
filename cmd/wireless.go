package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ljw/tplink-cli/internal/model"
	"github.com/spf13/cobra"
)

var wirelessCmd = &cobra.Command{
	Use:   "wireless",
	Short: "无线配置管理",
	Long:  `管理 TP-Link 路由器的无线网络配置，包括 WiFi 设置、访客网络、MAC 过滤、服务列表和客户端列表。`,
}

// ========== wireless config list ==========
var wirelessConfigListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看无线接口状态",
	Long: `查看 2.4G 和 5G 无线接口的当前配置状态。

可通过 --name 指定查看单个接口。支持的 name 值:
  wlan_host_2g   (2.4G)
  wlan_host_5g   (5G)

示例:
  tplink wireless config list
  tplink wireless config list --name wlan_host_2g`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")

		result, err := client.GetWirelessConfig()
		if err != nil {
			return err
		}

		type ConfigRow struct {
			Band    string `json:"band"`
			Enable  string `json:"enable"`
			SSID    string `json:"ssid"`
			Encrypt string `json:"encryption"`
			Channel string `json:"channel"`
			Mode    string `json:"mode"`
			Power   string `json:"power"`
		}

		rows := make([]ConfigRow, 0, 2)
		addRow := func(host *model.WlanHostConfig, band string) {
			if host == nil {
				return
			}
			rows = append(rows, ConfigRow{
				Band:    band,
				Enable:  statusLabel(host.Enable),
				SSID:    host.SSID,
				Encrypt: host.Encryption,
				Channel: host.Channel,
				Mode:    host.Mode,
				Power:   host.Power,
			})
		}

		switch name {
		case "wlan_host_2g":
			addRow(result.WlanHost2G, "2.4G")
		case "wlan_host_5g":
			addRow(result.WlanHost5G, "5G")
		case "":
			addRow(result.WlanHost2G, "2.4G")
			addRow(result.WlanHost5G, "5G")
		default:
			return fmt.Errorf("无效的 name: %s, 支持的值为 wlan_host_2g 或 wlan_host_5g", name)
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- band: %s\n  enable: %s\n  ssid: %s\n  encryption: %s\n  channel: %s\n  mode: %s\n  power: %s\n",
					r.Band, r.Enable, r.SSID, r.Encrypt, r.Channel, r.Mode, r.Power)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("没有找到无线配置")
				return nil
			}
			headers := []string{"BAND", "ENABLE", "SSID", "ENCRYPT", "CHANNEL", "MODE", "POWER"}
			printTable(headers, rows, func(r ConfigRow) []string {
				return []string{r.Band, r.Enable, r.SSID, r.Encrypt, r.Channel, r.Mode, r.Power}
			})
			return nil
		}
	},
}

// ========== wireless config set ==========
var wirelessConfigSetCmd = &cobra.Command{
	Use:   "set",
	Short: "修改无线配置",
	Long: `修改 2.4G 或 5G 无线接口配置。
需指定 --2g 或 --5g 选择目标频段，未指定的字段保持原值。

示例:
  tplink wireless config set --2g --ssid "MyWiFi" --key "mypassword"
  tplink wireless config set --5g --enable on --ssid "MyWiFi_5G"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		is2G, _ := cmd.Flags().GetBool("2g")
		is5G, _ := cmd.Flags().GetBool("5g")

		if !is2G && !is5G {
			return fmt.Errorf("必须指定 --2g 或 --5g")
		}
		if is2G && is5G {
			return fmt.Errorf("--2g 和 --5g 不能同时指定")
		}

		band := "wlan_host_5g"
		if is2G {
			band = "wlan_host_2g"
		}

		client, err := newAPIClient()
		if err != nil {
			return err
		}

		// 先获取当前配置（dry-run 模式下 GET 不真正执行，current 为 nil，用空结构体替代）
		result, err := client.GetWirelessConfig()
		if err != nil {
			return err
		}

		var current *model.WlanHostConfig
		if is2G {
			current = result.WlanHost2G
		} else {
			current = result.WlanHost5G
		}
		if current == nil {
			if client.DryRun {
				current = &model.WlanHostConfig{}
			} else {
				return fmt.Errorf("无法获取 %s 当前配置", band)
			}
		}

		// 合并字段：仅发送 SET API 支持的字段
		cfg := map[string]interface{}{}
		mergeFieldBool(cmd, "enable", cfg)
		mergeFieldStr(cmd, "ssid", current.SSID, cfg)
		mergeFieldStr(cmd, "ssid-code-type", current.SSIDCodeType, cfg)
		mergeFieldStr(cmd, "ssidbrd", current.SSIDBrd, cfg)
		mergeFieldStr(cmd, "isolate", current.Isolate, cfg)
		mergeFieldStr(cmd, "encryption", current.Encryption, cfg)
		mergeFieldStr(cmd, "auth", current.Auth, cfg)
		mergeFieldStr(cmd, "cipher", current.Cipher, cfg)
		mergeFieldStr(cmd, "key", current.Key, cfg)
		mergeFieldStr(cmd, "key-update-intv", current.KeyUpdateIntv, cfg)

		if err := client.SetWirelessConfig(band, cfg); err != nil {
			return err
		}

		fmt.Printf("无线配置更新成功, 频段: %s\n", band)
		return nil
	},
}

// ========== wireless guest list ==========
var wirelessGuestListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看访客网络配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		result, err := client.GetGuestNetwork()
		if err != nil {
			return err
		}

		g := result.Guest2G

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(g)
		case "yaml":
			fmt.Printf("enable: %s\nssid: %s\nencrypt: %s\nupload: %s\ndownload: %s\nradio_max_sta: %s\n",
				g.Enable, g.SSID, g.Encrypt, g.Upload, g.Download, g.RadioMaxSta)
			return nil
		default:
			type GuestRow struct {
				Enable      string
				SSID        string
				Encrypt     string
				MaxSta      string
				Upload      string
				Download    string
			}
			rows := []GuestRow{{
				Enable:   statusLabel(g.Enable),
				SSID:     g.SSID,
				Encrypt:  statusLabel(g.Encrypt),
				MaxSta:   g.RadioMaxSta,
				Upload:   g.Upload,
				Download: g.Download,
			}}
			headers := []string{"ENABLE", "SSID", "ENCRYPT", "MAX_STA", "UPLOAD", "DOWNLOAD"}
			printTable(headers, rows, func(r GuestRow) []string {
				return []string{r.Enable, r.SSID, r.Encrypt, r.MaxSta, r.Upload, r.Download}
			})
			return nil
		}
	},
}

// ========== wireless guest set ==========
var wirelessGuestSetCmd = &cobra.Command{
	Use:   "set",
	Short: "修改访客网络配置",
	Long: `修改访客网络配置，未指定的字段保持原值。

示例:
  tplink wireless guest set --enable on --ssid "GuestWiFi" --encrypt on --key "guest123"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		// 获取当前配置
		result, err := client.GetGuestNetwork()
		if err != nil {
			return err
		}
		current := result.Guest2G

		cfg := map[string]interface{}{}
		mergeFieldBool(cmd, "enable", cfg)
		mergeFieldStr(cmd, "ssid", current.SSID, cfg)
		mergeFieldStr(cmd, "encrypt", current.Encrypt, cfg)
		mergeFieldStr(cmd, "key", current.Key, cfg)
		mergeFieldStr(cmd, "upload", current.Upload, cfg)
		mergeFieldStr(cmd, "download", current.Download, cfg)
		mergeFieldStr(cmd, "radio-max-sta", current.RadioMaxSta, cfg)
		mergeFieldStr(cmd, "ssid-code-type", current.SSIDCodeType, cfg)

		if err := client.SetGuestNetwork(cfg); err != nil {
			return err
		}

		fmt.Println("访客网络配置更新成功")
		return nil
	},
}

// ========== wireless access config ==========
var wirelessAccessConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "MAC地址过滤配置",
}

var wirelessAccessConfigListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看MAC地址过滤配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		cfg, err := client.GetWlanAccessConfig()
		if err != nil {
			return err
		}

		modeStr := "白名单"
		if cfg.Mode == "2" {
			modeStr = "黑名单"
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(cfg)
		case "yaml":
			fmt.Printf("enable: %s\nmode: %s\nssid_list:\n", cfg.Enable, cfg.Mode)
			for _, s := range cfg.SSIDList {
				fmt.Printf("  - %s\n", decodeURL(s))
			}
			return nil
		default:
			fmt.Printf("启用: %s\n", statusLabel(cfg.Enable))
			fmt.Printf("模式: %s\n", modeStr)
			fmt.Println("生效 SSID:")
			for _, s := range cfg.SSIDList {
				fmt.Printf("  - %s\n", decodeURL(s))
			}
			return nil
		}
	},
}

var wirelessAccessConfigSetCmd = &cobra.Command{
	Use:   "set",
	Short: "设置MAC地址过滤",
	Long: `启用/禁用 MAC 地址过滤并设置模式。

mode: 1=白名单(仅允许列表内设备), 2=黑名单(禁止列表内设备)

示例:
  tplink wireless access config set --enable on --mode 1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		// 获取当前配置
		cfg, err := client.GetWlanAccessConfig()
		if err != nil {
			return err
		}

		merged := map[string]interface{}{}
		mergeFieldBool(cmd, "enable", merged)
		mergeFieldStr(cmd, "mode", cfg.Mode, merged)

		// 仅当用户指定了 ssid-list 时才包含
		if cmd.Flags().Changed("ssid-list") {
			ssidList, _ := cmd.Flags().GetStringSlice("ssid-list")
			merged["ssid_list"] = ssidList
		}

		if err := client.SetWlanAccessConfig(merged); err != nil {
			return err
		}

		fmt.Println("MAC 过滤配置更新成功")
		return nil
	},
}

// ========== wireless access white ==========
var wirelessAccessWhiteCmd = &cobra.Command{
	Use:   "white",
	Short: "MAC白名单管理",
}

	var wirelessAccessWhiteListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看白名单",
	Long: `查看白名单条目列表。ID 列为 API 标识名，用于 del 命令。
	
示例:
  tplink wireless access white list
  tplink wireless access white list -o json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		items, total, err := client.GetWlanAccessWhiteList()
		if err != nil {
			return err
		}

		type WhiteRow struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Mac  string `json:"mac"`
		}

		rows := make([]WhiteRow, 0, len(items))
		for _, item := range items {
			for key, data := range item {
				rows = append(rows, WhiteRow{
					ID:   key,
					Name: decodeURL(data.NameField),
					Mac:  data.Mac,
				})
			}
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
				fmt.Println("白名单为空")
				return nil
			}
			headers := []string{"ID", "NAME", "MAC"}
			printTable(headers, rows, func(r WhiteRow) []string {
				return []string{r.ID, r.Name, r.Mac}
			})
			if total > 0 {
				fmt.Printf("\n[总数: %d]\n", total)
			}
			return nil
		}
	},
}

var wirelessAccessWhiteAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加白名单",
	Long: `添加 MAC 地址到白名单。

示例:
  tplink wireless access white add -m F4-3B-D8-E1-C7-A6 -n "我的手机"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mac, _ := cmd.Flags().GetString("mac")
		name, _ := cmd.Flags().GetString("name")

		if mac == "" {
			return fmt.Errorf("--mac 为必选参数")
		}

		client, err := newAPIClient()
		if err != nil {
			return err
		}

		listName, err := client.AddWlanAccessWhite(mac, name)
		if err != nil {
			return err
		}

		fmt.Printf("白名单添加成功, ID: %s\n", listName)
		return nil
	},
}

	var wirelessAccessWhiteDelCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "删除白名单",
	Long: `根据 ID 删除白名单条目。ID 来自 list 输出中的 ID 列。
格式为 "white_list_NNNN"。

示例:
  tplink wireless access white del white_list_1782875167`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		if err := client.DelWlanAccessWhite(args[0]); err != nil {
			return err
		}

		fmt.Printf("白名单删除成功, ID: %s\n", args[0])
		return nil
	},
}

// ========== wireless access black ==========
var wirelessAccessBlackCmd = &cobra.Command{
	Use:   "black",
	Short: "MAC黑名单管理",
}

	var wirelessAccessBlackListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看黑名单",
	Long: `查看黑名单条目列表。ID 列为 API 标识名，用于 del 命令。
	
示例:
  tplink wireless access black list
  tplink wireless access black list -o json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		items, total, err := client.GetWlanAccessBlackList()
		if err != nil {
			return err
		}

		type BlackRow struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Mac  string `json:"mac"`
		}

		rows := make([]BlackRow, 0, len(items))
		for _, item := range items {
			for key, data := range item {
				rows = append(rows, BlackRow{
					ID:   key,
					Name: decodeURL(data.NameField),
					Mac:  data.Mac,
				})
			}
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
				fmt.Println("黑名单为空")
				return nil
			}
			headers := []string{"ID", "NAME", "MAC"}
			printTable(headers, rows, func(r BlackRow) []string {
				return []string{r.ID, r.Name, r.Mac}
			})
			if total > 0 {
				fmt.Printf("\n[总数: %d]\n", total)
			}
			return nil
		}
	},
}

var wirelessAccessBlackAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加黑名单",
	Long: `添加 MAC 地址到黑名单。

示例:
  tplink wireless access black add -m F4-3B-D8-E1-C7-A6 -n "未知设备"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mac, _ := cmd.Flags().GetString("mac")
		name, _ := cmd.Flags().GetString("name")

		if mac == "" {
			return fmt.Errorf("--mac 为必选参数")
		}

		client, err := newAPIClient()
		if err != nil {
			return err
		}

		listName, err := client.AddWlanAccessBlack(mac, name)
		if err != nil {
			return err
		}

		fmt.Printf("黑名单添加成功, ID: %s\n", listName)
		return nil
	},
}

	var wirelessAccessBlackDelCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "删除黑名单",
	Long: `根据 ID 删除黑名单条目。ID 来自 list 输出中的 ID 列。
格式为 "black_list_NNNN"。

示例:
  tplink wireless access black del black_list_1782875167`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		if err := client.DelWlanAccessBlack(args[0]); err != nil {
			return err
		}

		fmt.Printf("黑名单删除成功, ID: %s\n", args[0])
		return nil
	},
}

// ========== wireless service list ==========
var wirelessServiceListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看Wlan服务列表",
	Long:  `查看所有无线 Wlan 服务列表，包含 serv_id、radio_id、SSID 等信息。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		items, total, err := client.GetWlanServList()
		if err != nil {
			return err
		}

		type ServRow struct {
			ServID   string `json:"serv_id"`
			RadioID  string `json:"radio_id"`
			SSID     string `json:"ssid"`
			Enable   string `json:"enable"`
			NetType  string `json:"network_type"`
			Encrypt  string `json:"encryption"`
		}

		rows := make([]ServRow, 0, len(items))
		for _, item := range items {
			for _, data := range item {
				enableLabel := statusLabel(data.Enable)
				netType := data.NetworkType
				if netType == "1" {
					netType = "主网络"
				} else if netType == "2" {
					netType = "访客"
				}
				rows = append(rows, ServRow{
					ServID:  data.ServID,
					RadioID: data.RadioID,
					SSID:    data.SSID,
					Enable:  enableLabel,
					NetType: netType,
					Encrypt: data.Encryption,
				})
			}
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- serv_id: %s\n  radio_id: %s\n  ssid: %s\n  enable: %s\n  network_type: %s\n  encryption: %s\n",
					r.ServID, r.RadioID, r.SSID, r.Enable, r.NetType, r.Encrypt)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("没有找到Wlan服务")
				return nil
			}
			headers := []string{"SERV_ID", "RADIO_ID", "SSID", "ENABLE", "NETWORK", "ENCRYPT"}
			printTable(headers, rows, func(r ServRow) []string {
				return []string{r.ServID, r.RadioID, r.SSID, r.Enable, r.NetType, r.Encrypt}
			})
			if total > 0 {
				fmt.Printf("\n[总数: %d]\n", total)
			}
			return nil
		}
	},
}

// ========== wireless client list ==========
var wirelessClientListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看无线客户端列表",
	Long: `查看当前连接到 Wi-Fi 的客户端列表。

可通过 --radio-id 和 --serv-id 过滤（参考 service list 的输出）。

示例:
  tplink wireless client list
  tplink wireless client list --radio-id 1 --serv-id 20002`,
	RunE: func(cmd *cobra.Command, args []string) error {
		radioID, _ := cmd.Flags().GetString("radio-id")
		servID, _ := cmd.Flags().GetString("serv-id")

		client, err := newAPIClient()
		if err != nil {
			return err
		}

		items, total, err := client.GetWirelessClients(radioID, servID)
		if err != nil {
			return err
		}

		type ClientRow struct {
			Mac     string `json:"mac"`
			Name    string `json:"name"`
			IP      string `json:"ip"`
			SSID    string `json:"ssid"`
			RSSI    string `json:"rssi"`
			TxRate  string `json:"tx_rate"`
			RxRate  string `json:"rx_rate"`
			TxFlow  string `json:"tx_flow"`
			Status  string `json:"status"`
		}

		rows := make([]ClientRow, 0, len(items))
		for _, item := range items {
			for _, data := range item {
				statusStr := "离线"
				if data.Status == "1" {
					statusStr = "在线"
				}
				rows = append(rows, ClientRow{
					Mac:    data.Mac,
					Name:   data.Name,
					IP:     data.IP,
					SSID:   data.SSID,
					RSSI:   "-" + data.RSSI + "dBm",
					TxRate: data.TxRate,
					RxRate: data.RxRate,
					TxFlow: formatBytes(data.TxFlow),
					Status: statusStr,
				})
			}
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- mac: %s\n  name: %s\n  ip: %s\n  ssid: %s\n  rssi: %s\n  tx_rate: %s\n  rx_rate: %s\n  tx_flow: %s\n  status: %s\n",
					r.Mac, r.Name, r.IP, r.SSID, r.RSSI, r.TxRate, r.RxRate, r.TxFlow, r.Status)
			}
			if total > 0 {
				fmt.Printf("\n# total: %d\n", total)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("没有找到无线客户端")
				return nil
			}
			headers := []string{"MAC", "NAME", "IP", "SSID", "RSSI", "TX_RATE", "RX_RATE", "TX_FLOW", "STATUS"}
			printTable(headers, rows, func(r ClientRow) []string {
				return []string{r.Mac, r.Name, r.IP, r.SSID, r.RSSI, r.TxRate, r.RxRate, r.TxFlow, r.Status}
			})
			if total > 0 {
				fmt.Printf("\n[总数: %d]\n", total)
			}
			return nil
		}
	},
}

// mergeFieldStr 从 cobra flag 合并字符串字段，仅当用户显式指定时才写入
// 删除 else 分支：不应将未修改的字段也发送给 SET API，
// 否则会导致请求体包含全部字段，TP-Link API 可能拒绝或产生意外行为。
func mergeFieldStr(cmd *cobra.Command, flagName, defaultVal string, out map[string]interface{}) {
	f := cmd.Flags().Lookup(flagName)
	if f != nil && f.Changed {
		out[flagName] = f.Value.String()
	}
}

// mergeFieldBool 合并布尔类字段，将用户友好的 on/off 转换为 API 期望的 "1"/"0"
// 支持: on/off, 1/0, yes/no, true/false
func mergeFieldBool(cmd *cobra.Command, flagName string, out map[string]interface{}) {
	f := cmd.Flags().Lookup(flagName)
	if f != nil && f.Changed {
		val := f.Value.String()
		switch strings.ToLower(val) {
		case "on", "1", "yes", "true":
			val = "1"
		case "off", "0", "no", "false":
			val = "0"
		}
		out[flagName] = val
	}
}

// mergeFieldBoolMapped 同 mergeFieldBool，但使用不同的 map key（flagName 为 CLI flag，apiField 为 JSON key）
func mergeFieldBoolMapped(cmd *cobra.Command, flagName, apiField string, out map[string]interface{}) {
	f := cmd.Flags().Lookup(flagName)
	if f != nil && f.Changed {
		val := f.Value.String()
		switch strings.ToLower(val) {
		case "on", "1", "yes", "true":
			val = "1"
		case "off", "0", "no", "false":
			val = "0"
		}
		out[apiField] = val
	}
}

// mergeFieldStrMapped 同 mergeFieldStr，但使用不同的 map key（flagName 为 CLI flag，apiField 为 JSON key）
func mergeFieldStrMapped(cmd *cobra.Command, flagName, apiField, defaultVal string, out map[string]interface{}) {
	f := cmd.Flags().Lookup(flagName)
	if f != nil && f.Changed {
		out[apiField] = f.Value.String()
	}
}

func init() {
	// wireless config 子命令
	var wirelessConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "无线接口配置",
	}
	wirelessConfigListCmd.Flags().StringP("name", "n", "", "接口名称: wlan_host_2g | wlan_host_5g")
	wirelessConfigCmd.AddCommand(wirelessConfigListCmd)
	wirelessConfigCmd.AddCommand(wirelessConfigSetCmd)

	wirelessConfigSetCmd.Flags().Bool("2g", false, "修改 2.4G 频段")
	wirelessConfigSetCmd.Flags().Bool("5g", false, "修改 5G 频段")
	wirelessConfigSetCmd.Flags().String("enable", "", "启用: on|off|1|0")
	wirelessConfigSetCmd.Flags().String("ssid", "", "WiFi 名称")
	wirelessConfigSetCmd.Flags().String("ssid-code-type", "", "SSID编码类型")
	wirelessConfigSetCmd.Flags().String("ssidbrd", "", "SSID广播: 0|1")
	wirelessConfigSetCmd.Flags().String("isolate", "", "AP隔离: 0|1")
	wirelessConfigSetCmd.Flags().String("encryption", "", "加密方式")
	wirelessConfigSetCmd.Flags().String("auth", "", "认证方式")
	wirelessConfigSetCmd.Flags().String("cipher", "", "加密算法")
	wirelessConfigSetCmd.Flags().String("key", "", "WiFi 密码")
	wirelessConfigSetCmd.Flags().String("key-update-intv", "", "密钥更新间隔(s)")

	// wireless guest 子命令
	var wirelessGuestCmd = &cobra.Command{
		Use:   "guest",
		Short: "访客网络管理",
	}
	wirelessGuestCmd.AddCommand(wirelessGuestListCmd)
	wirelessGuestCmd.AddCommand(wirelessGuestSetCmd)

	wirelessGuestSetCmd.Flags().String("enable", "", "启用状态: on|off|1|0")
	wirelessGuestSetCmd.Flags().String("ssid", "", "访客WiFi名称")
	wirelessGuestSetCmd.Flags().String("encrypt", "", "加密: 0|1")
	wirelessGuestSetCmd.Flags().String("key", "", "访客WiFi密码")
	wirelessGuestSetCmd.Flags().String("upload", "", "上行限速(KB/s)")
	wirelessGuestSetCmd.Flags().String("download", "", "下行限速(KB/s)")
	wirelessGuestSetCmd.Flags().String("radio-max-sta", "", "最大连接数")
	wirelessGuestSetCmd.Flags().String("ssid-code-type", "", "SSID编码类型")

	// wireless access 子命令
	var wirelessAccessCmd = &cobra.Command{
		Use:   "access",
		Short: "MAC地址过滤管理",
	}

	// access config
	wirelessAccessConfigCmd.AddCommand(wirelessAccessConfigListCmd)
	wirelessAccessConfigCmd.AddCommand(wirelessAccessConfigSetCmd)
	wirelessAccessConfigSetCmd.Flags().String("enable", "", "启用: on|off|1|0")
	wirelessAccessConfigSetCmd.Flags().String("mode", "", "模式: 1=白名单, 2=黑名单")
	wirelessAccessConfigSetCmd.Flags().StringSlice("ssid-list", nil, "生效SSID列表")

	// access white
	wirelessAccessWhiteCmd.AddCommand(wirelessAccessWhiteListCmd)
	wirelessAccessWhiteCmd.AddCommand(wirelessAccessWhiteAddCmd)
	wirelessAccessWhiteCmd.AddCommand(wirelessAccessWhiteDelCmd)
	wirelessAccessWhiteAddCmd.Flags().StringP("mac", "m", "", "MAC 地址 (必选)")
	wirelessAccessWhiteAddCmd.Flags().StringP("name", "n", "", "设备名称")

	// access black
	wirelessAccessBlackCmd.AddCommand(wirelessAccessBlackListCmd)
	wirelessAccessBlackCmd.AddCommand(wirelessAccessBlackAddCmd)
	wirelessAccessBlackCmd.AddCommand(wirelessAccessBlackDelCmd)
	wirelessAccessBlackAddCmd.Flags().StringP("mac", "m", "", "MAC 地址 (必选)")
	wirelessAccessBlackAddCmd.Flags().StringP("name", "n", "", "设备名称")

	wirelessAccessCmd.AddCommand(wirelessAccessConfigCmd)
	wirelessAccessCmd.AddCommand(wirelessAccessWhiteCmd)
	wirelessAccessCmd.AddCommand(wirelessAccessBlackCmd)

	// wireless service 子命令
	var wirelessServiceCmd = &cobra.Command{
		Use:   "service",
		Short: "Wlan服务管理",
	}
	wirelessServiceCmd.AddCommand(wirelessServiceListCmd)

	// wireless client 子命令
	var wirelessClientCmd = &cobra.Command{
		Use:   "client",
		Short: "无线客户端管理",
	}
	wirelessClientListCmd.Flags().String("radio-id", "", "Radio ID 过滤")
	wirelessClientListCmd.Flags().String("serv-id", "", "Service ID 过滤")
	wirelessClientCmd.AddCommand(wirelessClientListCmd)

	wirelessCmd.AddCommand(wirelessConfigCmd)
	wirelessCmd.AddCommand(wirelessGuestCmd)
	wirelessCmd.AddCommand(wirelessAccessCmd)
	wirelessCmd.AddCommand(wirelessServiceCmd)
	wirelessCmd.AddCommand(wirelessClientCmd)

	rootCmd.AddCommand(wirelessCmd)
}

// avoid unused import
var _ = strconv.Itoa
