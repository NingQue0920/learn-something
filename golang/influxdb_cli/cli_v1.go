package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type InfluxDBInteractiveClient struct {
	url      string
	username string
	password string
	org      string
}

func NewInfluxDBInteractiveClient(url, username, password, org string) *InfluxDBInteractiveClient {
	return &InfluxDBInteractiveClient{
		url:      url,
		username: username,
		password: password,
		org:      org,
	}
}

func (c *InfluxDBInteractiveClient) ExecuteQuery(query string) {
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

	table := tablewriter.NewWriter(os.Stdout)
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
}

func (c *InfluxDBInteractiveClient) InteractiveMode() {
	fmt.Println("欢迎使用InfluxDB交互式客户端。输入'exit'退出。")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("请输入SQL查询 > ")
		if !scanner.Scan() {
			break
		}
		query := scanner.Text()
		if len(query) < 10 {
			continue
		}
		if strings.ToLower(query) == "exit" {
			break
		}
		c.ExecuteQuery(query)
	}
}

func main() {
	url := "http://10.20.28.235:9086"
	username := "fio"
	password := "Fio#1234"
	org := "fio"

	client := NewInfluxDBInteractiveClient(url, username, password, org)
	client.InteractiveMode()
}
