package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ljw/tplink-cli/internal/api"
	"github.com/ljw/tplink-cli/internal/config"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api <method> <path>",
	Short: "发送原始 API 请求",
	Long: `发送原始 HTTP 请求到 TP-Link 路由器 API，用于应对未覆盖的 API 或难以命令化的操作。

方法: GET, POST, PUT, DELETE
路径: API 路径（不包含 host 和 port），例如 /users

示例:
  tplink api get /users -q dept=2
  tplink api post /users -d '{"account":"test","password":"Abc123456","realname":"测试"}'`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := args[0]
		path := args[1]

		queryStr, _ := cmd.Flags().GetString("query")
		dataStr, _ := cmd.Flags().GetString("data")

		// 构建完整路径（包含查询参数）
		fullPath := path
		if queryStr != "" {
			params, err := url.ParseQuery(queryStr)
			if err != nil {
				return fmt.Errorf("查询参数格式错误: %w", err)
			}
			fullPath = path + "?" + params.Encode()
		}

		// 加载配置获取当前 server
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		srvCfg, srvName, err := cfg.GetServer(serverName)
		if err != nil {
			return err
		}

		client := api.NewClient(srvCfg, srvName)
		client.DryRun = dryRun

		// 处理请求体
		var body []byte
		if dataStr != "" {
			// 验证 JSON 格式
			if !json.Valid([]byte(dataStr)) {
				return fmt.Errorf("请求体不是有效的 JSON")
			}
			// 压缩 JSON（去除空白）
			var compacted interface{}
			if err := json.Unmarshal([]byte(dataStr), &compacted); err != nil {
				return fmt.Errorf("请求体 JSON 解析失败: %w", err)
			}
			body, _ = json.Marshal(compacted)
		}

		respBody, err := client.DoRawRequest(method, fullPath, body)
		if err != nil {
			return err
		}

		// 格式化输出响应
		if respBody != nil && len(respBody) > 0 {
			var prettyJSON interface{}
			if json.Unmarshal(respBody, &prettyJSON) == nil {
				formatted, _ := json.MarshalIndent(prettyJSON, "", "  ")
				fmt.Println(string(formatted))
			} else {
				fmt.Println(string(respBody))
			}
		}

		return nil
	},
}

func init() {
	apiCmd.Flags().StringP("query", "q", "", "查询参数，格式: key=value 或 key=value&key=value")
	apiCmd.Flags().StringP("data", "d", "", "JSON 请求体")

	rootCmd.AddCommand(apiCmd)
}
