package main

import (
	"fmt"
	"go_remote_control/collector/data"
	"go_remote_control/collector/parser"
	"go_remote_control/collector/request"
	"io"
	"net/http"
	"sync"
	"time"
)

type Scheduler struct {
}

var client = &http.Client{}

type Collector struct {
	Request *request.Request
	Data    *data.Data
	Parser  *parser.Parser
	Info    string
}

func (c *Collector) Collect() {
	//获取请求网址
	resp, err := client.Do(c.Request.Req)
	if err != nil {
		c.Info = err.Error()
		return
	}
	//爬取的内容在Body，结束时应该关闭
	defer resp.Body.Close()
	temp, _ := io.ReadAll(resp.Body)
	c.Data = &data.Data{
		RawHtml: string(temp),
	}

	//fmt.Println(c.RawHtml)
	c.Info = "success"
	c.Parser.GetLink(c.Data)
}

//func (c *Collector) Parse() {
//	temp := strings.ReplaceAll(c.data.RawHtml, "\n", "")
//	for _, regex := range c.Regex {
//		re := regexp.MustCompile(regex)
//		c.data.ParsedData = re.FindAllString(temp, -1)
//	}
//}

func main() {
	baseUrl := "http://www.iyanghua.com/huahui/"
	//baseUrl := "http://127.0.0.1:20001/flowerKnowledge/getAllFlowerKnowledgeData"
	start := time.Now()
	//c := http.Client{}
	wg := sync.WaitGroup{}
	for i := 0; i < 26; i++ {
		wg.Add(1)
		c := Collector{
			Request: &request.Request{
				Rule: &request.Rule{
					Method: "GET",
					Url:    baseUrl,
				},
			},
		}
		go func() {
			c.Collect()
			wg.Done()
		}()
		wg.Wait()
		//fmt.Println(collector.data)
	}
	//var reqList []*http.Request
	//for i := 0; i < 26; i++ {
	//	req, _ := http.GenerateRequest("GET", fmt.Sprintf("%s%c.html", baseUrl, byte(65+i)), nil)
	//	reqList = append(reqList, req)
	//}
	//wg := sync.WaitGroup{}
	//for _, req := range reqList {
	//	wg.Add(1)
	//	go func() {
	//		resp, err := client.Do(req)
	//		if err != nil {
	//			return
	//		}
	//		html, _ := io.ReadAll(resp.Body)
	//		fmt.Println(string(html))
	//		wg.Done()
	//	}()
	//}
	//c := &Collector{
	//	req: RequestRule{
	//		Url: baseUrl,
	//	},
	//}
	//wg.Wait()
	//var urlList []string
	////Do something
	//for i := 0; i < 26; i++ {
	//	collector := &Collector{
	//		req: RequestRule{
	//			Client: &http.Client{},
	//			Url:    fmt.Sprintf("%s%c.html", baseUrl, byte(65+i)),
	//			Headers: map[string]string{
	//				"1": "2",
	//				"3": "4",
	//			},
	//		},
	//		data: Data{},
	//	}
	//	collector.Collect()
	//	reLink := regexp.MustCompile(`src="(.*?)"`)
	//	urls := reLink.FindAllString(collector.data.RawHtml, -1)
	//	urlList = append(urlList, urls...)
	//	fmt.Println(collector.Info)
	//}
	//
	//fmt.Println(urlList)
	//collector := Collector{
	//	Client: http.Client{},
	//	Url:    baseUrl,
	//}
	//collector.Collect()
	//fmt.Println(collector.RawHtml)

	//时间统计
	end := time.Since(start)
	fmt.Println(end)
}
