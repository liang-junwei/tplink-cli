package format

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/ljw/tplink-cli/internal/model"
	"gopkg.in/yaml.v3"
)

// OutputFormat 输出格式类型
type OutputFormat string

const (
	FormatTable OutputFormat = "table"
	FormatJSON  OutputFormat = "json"
	FormatYAML  OutputFormat = "yaml"
)

// Pagination 分页信息
type Pagination struct {
	Page     int `json:"page" yaml:"page"`
	PageSize int `json:"page_size" yaml:"page_size"`
	Total    int `json:"total" yaml:"total"`
}

// HasPagination 是否启用了分页
func (p *Pagination) HasPagination() bool {
	return p.Page > 0 && p.PageSize > 0
}

// TotalPages 总页数
func (p *Pagination) TotalPages() int {
	if p.PageSize <= 0 {
		return 0
	}
	return int(math.Ceil(float64(p.Total) / float64(p.PageSize)))
}

// RedirectRow 表格化展示的行结构
type RedirectRow struct {
	ID       string `json:"id" yaml:"id"`
	Name     string `json:"name" yaml:"name"`
	Proto    string `json:"proto" yaml:"proto"`
	Port     string `json:"port" yaml:"port"`
	DestIP   string `json:"dest_ip" yaml:"dest_ip"`
	DestPort string `json:"dest_port" yaml:"dest_port"`
	WAN      string `json:"wan" yaml:"wan"`
	Status   string `json:"status" yaml:"status"`
}

// FromRedirectRuleFull 将查询结果转换为表格行
func FromRedirectRuleFull(id string, r *model.RedirectRuleFull) RedirectRow {
	status := "off"
	if r.Enable == "on" {
		status = "on"
	}

	port := r.SrcDport
	if r.SrcDportStart != r.SrcDportEnd && r.SrcDportStart != "" {
		port = fmt.Sprintf("%s-%s", r.SrcDportStart, r.SrcDportEnd)
	}

	destPort := r.DestPort
	if r.DestPortStart != r.DestPortEnd && r.DestPortStart != "" {
		destPort = fmt.Sprintf("%s-%s", r.DestPortStart, r.DestPortEnd)
	}

	return RedirectRow{
		ID:       id,
		Name:     r.Name,
		Proto:    r.Proto,
		Port:     port,
		DestIP:   r.DestIP,
		DestPort: destPort,
		WAN:      r.If,
		Status:   status,
	}
}

// PrintRedirects 格式化输出端口映射规则列表
func PrintRedirects(items []model.RedirectItem, format OutputFormat, pagination *Pagination) error {
	rows := make([]RedirectRow, 0, len(items))
	for _, item := range items {
		for id, rule := range item {
			rows = append(rows, FromRedirectRuleFull(id, &rule))
		}
	}

	if err := printRedirectRows(rows, format); err != nil {
		return err
	}

	// 分页信息在数据之后输出
	if pagination != nil && pagination.HasPagination() {
		printPagination(pagination)
	}

	return nil
}

// PrintSingleRedirect 格式化输出单条规则
func PrintSingleRedirect(id string, rule *model.RedirectRuleFull, format OutputFormat) error {
	row := FromRedirectRuleFull(id, rule)
	return printRedirectRows([]RedirectRow{row}, format)
}

func printRedirectRows(rows []RedirectRow, format OutputFormat) error {
	switch format {
	case FormatJSON:
		return printJSON(rows)
	case FormatYAML:
		return printYAML(rows)
	case FormatTable:
		return printRedirectTable(rows)
	default:
		return printRedirectTable(rows)
	}
}

func printJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

func printYAML(v interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	return encoder.Encode(v)
}

func printPagination(p *Pagination) {
	fmt.Println()
	fmt.Println("---")
	fmt.Printf("Page %d/%d, Total %d rules\n", p.Page, p.TotalPages(), p.Total)
}

func printRedirectTable(rows []RedirectRow) error {
	if len(rows) == 0 {
		fmt.Println("没有找到端口映射规则")
		return nil
	}

	// 列定义
	headers := []string{"ID", "NAME", "PROTO", "PORT", "DEST IP", "DEST PORT", "WAN", "STATUS"}
	colWidths := make([]int, len(headers))

	// 计算每列宽度
	for i, h := range headers {
		colWidths[i] = len(h)
	}
	for _, row := range rows {
		vals := rowValues(row)
		for i, v := range vals {
			width := displayWidth(v)
			if width > colWidths[i] {
				colWidths[i] = width
			}
		}
	}

	// 打印表头
	printRow(headers, colWidths)

	// 打印分隔线
	seps := make([]string, len(headers))
	for i, w := range colWidths {
		seps[i] = strings.Repeat("-", w)
	}
	fmt.Println("  " + strings.Join(seps, "  "))

	// 打印数据行
	for _, row := range rows {
		vals := rowValues(row)
		printRow(vals, colWidths)
	}

	return nil
}

func rowValues(row RedirectRow) []string {
	return []string{row.ID, row.Name, row.Proto, row.Port, row.DestIP, row.DestPort, row.WAN, row.Status}
}

func printRow(cols []string, widths []int) {
	padded := make([]string, len(cols))
	for i, col := range cols {
		padded[i] = padRight(col, widths[i])
	}
	fmt.Println("  " + strings.Join(padded, "  "))
}

// padRight 右侧补空格到指定显示宽度
func padRight(s string, width int) string {
	dw := displayWidth(s)
	if dw >= width {
		return s
	}
	return s + strings.Repeat(" ", width-dw)
}

// displayWidth 计算字符串显示宽度（中文字符占2格）
func displayWidth(s string) int {
	width := 0
	for _, r := range s {
		if r > 0x7F {
			width += 2
		} else {
			width++
		}
	}
	return width
}
