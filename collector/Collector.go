package main

import (
	"fmt"
	"go_remote_control/collector/data"
	"go_remote_control/collector/parser"
	"go_remote_control/collector/request"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

var client = &http.Client{}

type Collector struct {
	Request *request.Request //网页访问请求
	Parser  *parser.Parser   //解析器
	Data    *data.Data       //数据
	Info    string           //状态信息
}

func (c *Collector) Collect() {
	//请求构造
	c.Request.GenerateRequest()
	//请求发起
	resp, err := client.Do(c.Request.Req)
	if err != nil {
		c.Info = err.Error()
		return
	}
	//数据读取，延迟关闭body
	defer resp.Body.Close()
	temp, _ := io.ReadAll(resp.Body)
	c.Data.RawHtml = string(temp)
	//数据解析
	c.Parser.GetLink(c.Data)
	c.Info = "success"
}

func main() {
	baseUrl := "http://www.iyanghua.com/huahui/"
	//baseUrl := "http://127.0.0.1:20001/flowerKnowledge/getAllFlowerKnowledgeData"
	start := time.Now()
	wg := sync.WaitGroup{}
	var dataList []*data.Data
	var urlList []string
	var collectorList []Collector
	for i := 0; i < 1; i++ {
		urlList = append(urlList, fmt.Sprintf("%s%c.html", baseUrl, byte(65+i)))
	}
	for _, v := range urlList {
		wg.Add(1)
		c := Collector{
			Request: &request.Request{
				Rule: &request.Rule{
					Method: "GET",
					Url:    v,
				},
			},
			Data: &data.Data{},
		}
		collectorList = append(collectorList, c)
		go func() {
			c.Collect()
			wg.Done()
		}()
	}
	wg.Wait()
	for _, v := range collectorList {
		//fmt.Println(v.Data.RawHtml)
		fmt.Println(strings.ReplaceAll(v.Data.RawHtml, "\n", ""))
		//fmt.Println(v.Data.ParsedData)
	}
	for _, v := range dataList {
		fmt.Println(v.RawHtml)
		fmt.Println(v.ParsedData)
	}
	//时间统计
	end := time.Since(start)
	fmt.Println(end)
}
