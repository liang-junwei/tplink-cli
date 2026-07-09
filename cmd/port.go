package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var portCmd = &cobra.Command{
	Use:   "port",
	Short: "有线接口状态",
	Long:  `查看 TP-Link 路由器所有有线接口的状态信息。`,
}

var portListCmd = &cobra.Command{
	Use:   "list",
	Short: "查看有线接口状态",
	Long:  `列出所有有线接口的状态，包括连接状态、速率、收发流量等。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		items, err := client.GetPorts()
		if err != nil {
			return err
		}

		// 转换为行数据
		type PortRow struct {
			ID     string `json:"id"`
			State  string `json:"state"`
			Speed  string `json:"speed"`
			Duplex string `json:"duplex"`
			TxAll  string `json:"tx_all"`
			RxAll  string `json:"rx_all"`
		}

		rows := make([]PortRow, 0, len(items))
		for _, item := range items {
			for _, data := range item {
				rows = append(rows, PortRow{
					ID:     data.PortID,
					State:  stateLabel(data.PortState),
					Speed:  data.LinkSpeed,
					Duplex: data.LinkDuplex,
					TxAll:  formatBytes(data.TxAll),
					RxAll:  formatBytes(data.RxAll),
				})
			}
		}

		switch output {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(rows)
		case "yaml":
			for _, r := range rows {
				fmt.Printf("- id: %s\n  state: %s\n  speed: %s\n  duplex: %s\n  tx: %s\n  rx: %s\n",
					r.ID, r.State, r.Speed, r.Duplex, r.TxAll, r.RxAll)
			}
			return nil
		default:
			if len(rows) == 0 {
				fmt.Println("没有找到接口信息")
				return nil
			}

			headers := []string{"PORT", "STATE", "SPEED", "DUPLEX", "TX", "RX"}
			printTable(headers, rows, func(r PortRow) []string {
				return []string{r.ID, r.State, r.Speed, r.Duplex, r.TxAll, r.RxAll}
			})
			return nil
		}
	},
}

func init() {
	portCmd.AddCommand(portListCmd)
	rootCmd.AddCommand(portCmd)
}
