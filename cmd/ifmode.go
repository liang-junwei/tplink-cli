package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ifmodeCmd = &cobra.Command{
	Use:   "ifmode",
	Short: "接口模式管理",
	Long:  `查看和设置 TP-Link 路由器的 WAN 接口模式（单WAN口/双WAN口等）。`,
}

// ifmode list
var ifmodeListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看接口模式",
	Long:  `查看当前 WAN 接口模式配置。wan_mode: 1=单WAN口, 2=双WAN口, 3=三WAN口, 4=四WAN口`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		data, err := client.GetIfMode()
		if err != nil {
			return err
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(data)
		case "yaml":
			fmt.Printf("wan_mode: %s\nsingle_wan: %s\n", data.WanMode, data.SingleWan)
			return nil
		default:
			desc := map[string]string{"1": "单WAN口", "2": "双WAN口", "3": "三WAN口", "4": "四WAN口"}
			modeDesc := desc[data.WanMode]
			if modeDesc == "" {
				modeDesc = "未知"
			}
			fmt.Printf("WAN模式: %s (%s)\n", data.WanMode, modeDesc)
			fmt.Printf("单WAN标识: %s\n", data.SingleWan)
			return nil
		}
	},
}

// ifmode set
var ifmodeSetCmd = &cobra.Command{
	Use:   "set",
	Short: "设置接口模式",
	Long:  `设置 WAN 接口模式。wan_mode: 1=单WAN口, 2=双WAN口, 3=三WAN口, 4=四WAN口`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wanMode, _ := cmd.Flags().GetString("wan-mode")
		if wanMode == "" {
			return fmt.Errorf("--wan-mode 为必选参数 (1/2/3/4)")
		}

		client, err := newAPIClient()
		if err != nil {
			return err
		}
		if err := client.SetIfMode(wanMode); err != nil {
			return err
		}

		fmt.Printf("接口模式设置成功: wan_mode=%s\n", wanMode)
		return nil
	},
}

func init() {
	ifmodeSetCmd.Flags().StringP("wan-mode", "m", "", "WAN口数量: 1|2|3|4 (必选)")

	ifmodeCmd.AddCommand(ifmodeListCmd)
	ifmodeCmd.AddCommand(ifmodeSetCmd)
	rootCmd.AddCommand(ifmodeCmd)
}
