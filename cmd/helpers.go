package cmd

import (
	"fmt"
	"strings"
)

// statusLabel 将 enable 字段值转为中文标签
// 兼容两种约定：on/off（redirect、brv6mode、dhcp等）和 1/0（wireless各子命令）
func statusLabel(s string) string {
	switch s {
	case "on", "1", "yes", "true":
		return "启用"
	case "off", "0", "no", "false", "":
		return "禁用"
	default:
		return "禁用"
	}
}

// stateLabel 将连接状态转为中文
func stateLabel(s string) string {
	if s == "connected" {
		return "已连接"
	}
	return "未连接"
}

// linkLabel 将 link_status 转为中文
func linkLabel(s string) string {
	if s == "up" {
		return "已连接"
	}
	return "未连接"
}

// protoLabel 将协议转为中文
func protoLabel(s string) string {
	switch s {
	case "static":
		return "静态IP"
	case "dhcp":
		return "DHCP"
	case "pppoe":
		return "PPPoE"
	default:
		return s
	}
}

// formatBytes 将字节数字符串转为可读格式
func formatBytes(s string) string {
	var n float64
	fmt.Sscanf(s, "%f", &n)
	if n >= 1<<40 {
		return fmt.Sprintf("%.1fT", n/(1<<40))
	}
	if n >= 1<<30 {
		return fmt.Sprintf("%.1fG", n/(1<<30))
	}
	if n >= 1<<20 {
		return fmt.Sprintf("%.1fM", n/(1<<20))
	}
	if n >= 1<<10 {
		return fmt.Sprintf("%.1fK", n/(1<<10))
	}
	return s
}

// decodeURL 简单解码 URL 编码的字符串
func decodeURL(s string) string {
	s = strings.ReplaceAll(s, "%2b", "+")
	s = strings.ReplaceAll(s, "%5e", "^")
	s = strings.ReplaceAll(s, "%20", " ")
	return s
}

// formatExpires 格式化 DHCP 过期时间
func formatExpires(s string) string {
	return s + "s"
}

// printTable 通用表格输出（泛型版本）
func printTable[T any](headers []string, rows []T, getVals func(T) []string) {
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}
	for _, row := range rows {
		vals := getVals(row)
		for i, v := range vals {
			w := 0
			for _, r := range v {
				if r > 0x7F {
					w += 2
				} else {
					w++
				}
			}
			if w > colWidths[i] {
				colWidths[i] = w
			}
		}
	}

	rowCols := make([]string, len(headers))
	for i, h := range headers {
		rowCols[i] = padRight3(h, colWidths[i])
	}
	fmt.Println("  " + strings.Join(rowCols, "  "))

	seps := make([]string, len(headers))
	for i, w := range colWidths {
		seps[i] = strings.Repeat("-", w)
	}
	fmt.Println("  " + strings.Join(seps, "  "))

	for _, row := range rows {
		vals := getVals(row)
		for i, v := range vals {
			rowCols[i] = padRight3(v, colWidths[i])
		}
		fmt.Println("  " + strings.Join(rowCols, "  "))
	}
}

func padRight3(s string, width int) string {
	w := 0
	for _, r := range s {
		if r > 0x7F {
			w += 2
		} else {
			w++
		}
	}
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}
