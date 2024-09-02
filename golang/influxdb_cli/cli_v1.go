package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chzyer/readline"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

type InfluxDBInteractiveClient struct {
	url      string
	username string
	password string
	org      string
	history  []string
}

func NewInfluxDBInteractiveClient(url, username, password, org string) *InfluxDBInteractiveClient {
	return &InfluxDBInteractiveClient{
		url:      url,
		username: username,
		password: password,
		org:      org,
	}
}

func (c *InfluxDBInteractiveClient) addToHistory(query string) {
	if len(c.history) >= 20 {
		c.history = c.history[1:]
	}
	c.history = append(c.history, query)
}

func (c *InfluxDBInteractiveClient) parseCustomQuery(customQuery string) string {
	parts := strings.Split(customQuery, "|")
	baseQuery := strings.TrimSpace(parts[0])

	var conditions, orderBy, groupBy, limit string
	var timeRange string

	limit = " LIMIT 20"
	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		switch {
		case strings.HasPrefix(part, "id ="):
			conditions += fmt.Sprintf(" AND _item_id = %s", strings.TrimPrefix(part, "id ="))
		case strings.HasPrefix(part, "ob "):
			orderBy = fmt.Sprintf(" ORDER BY %s", strings.TrimPrefix(part, "ob "))
		case strings.HasPrefix(part, "gb "):
			groupBy = fmt.Sprintf(" GROUP BY %s", strings.TrimPrefix(part, "gb "))
		case strings.HasPrefix(part, "limit "):
			limit = fmt.Sprintf(" LIMIT %s", strings.TrimPrefix(part, "limit "))
		case strings.HasPrefix(part, "time "):
			timeRange = parseTimeRange(part)
		default:
			conditions += fmt.Sprintf(" AND %s", part)
		}
	}

	query := fmt.Sprintf("%s WHERE 1=1%s%s%s%s%s", baseQuery, conditions, timeRange, groupBy, orderBy, limit)
	return query
}

func parseTimeRange(timeStr string) string {
	timeStr = strings.TrimPrefix(timeStr, "time ")

	// 匹配 "5m", "1h" 等格式
	if match, _ := regexp.MatchString(`^\d+[mhs]$`, timeStr); match {
		return fmt.Sprintf(" AND time > now() - %s", timeStr)
	}

	// 匹配 "in [2020-08-28T09:35 , 2020-08-08T10:25]" 格式
	inRangeRegex := regexp.MustCompile(`in \[(.+?)\s*,\s*(.+?)\]`)
	if matches := inRangeRegex.FindStringSubmatch(timeStr); len(matches) == 3 {
		startTime := parseTime(matches[1])
		endTime := parseTime(matches[2])
		return fmt.Sprintf(" AND time >= '%s' AND time <= '%s'", startTime, endTime)
	}

	// 如果都不匹配，返回空字符串
	return ""
}

func parseTime(timeStr string) string {
	// 尝试解析时间字符串
	t, err := time.Parse("2006-01-02T15:04", timeStr)
	if err != nil {
		// 如果解析失败，返回原始字符串
		return timeStr
	}
	t = t.Add(-8 * time.Hour)
	// 返回格式化的时间字符串
	return t.Format(time.RFC3339)
}

func (c *InfluxDBInteractiveClient) ExecuteQuery(query string) {
	query = c.parseCustomQuery(query)
	fmt.Println("real query : ", query)
	url := fmt.Sprintf("%s/query?db=%s&q=%s", c.url, "fio", url.QueryEscape(query))

	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(c.username, c.password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("查询执行出错: %s\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("查询执行失败，状态码: %d, 响应: %s\n", resp.StatusCode, string(body))
		return
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	var headers []string
	var data [][]string

	if results, ok := result["results"].([]interface{}); ok && len(results) > 0 {
		if firstResult, ok := results[0].(map[string]interface{}); ok {
			if series, ok := firstResult["series"].([]interface{}); ok && len(series) > 0 {
				if firstSeries, ok := series[0].(map[string]interface{}); ok {
					if columns, ok := firstSeries["columns"].([]interface{}); ok {
						for _, col := range columns {
							headers = append(headers, fmt.Sprintf("%v", col))
						}
						table.SetHeader(headers)

						if values, ok := firstSeries["values"].([]interface{}); ok {
							for _, row := range values {
								if rowData, ok := row.([]interface{}); ok {
									var dataRow []string
									for _, val := range rowData {
										dataRow = append(dataRow, fmt.Sprintf("%v", val))
									}
									data = append(data, dataRow)
								}
							}
						}
					}
				}
			}
		}
	}

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	table.AppendBulk(data)
	table.Render()
	fmt.Println(buf.String())
}

func (c *InfluxDBInteractiveClient) InteractiveMode() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "InfluxDB > ",
		HistoryFile:     "/tmp/influxdb_history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	fmt.Println("欢迎使用InfluxDB交互式客户端。输入'exit'退出。")
	fmt.Println("使用上下箭头键浏览历史记录。")

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}
		line = strings.TrimSpace(line)
		if line == "exit" {
			break
		}
		if line == "" {
			continue
		}
		c.ExecuteQuery(line)
		c.addToHistory(line)
		rl.SaveHistory(line)
	}
}

func main() {
	var url string
	// 让用户输出ip地址
	fmt.Print("请输入InfluxDB地址: ")
	fmt.Scanln(&url)
	url = "http://" + url + ":9086"
	username := "fio"
	password := "Fio#1234"
	org := "fio"

	client := NewInfluxDBInteractiveClient(url, username, password, org)
	client.InteractiveMode()
}
