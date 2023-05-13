package request

import "net/http"

type Rule struct {
	Headers map[string]string //请求头
	Method  string            //请求方法
	Url     string            //爬取url
	Query   map[string]string //参数
	//Form    map[string]string //表单
}

type Request struct {
	*Rule
	Req *http.Request
}

func (n *Request) GenerateRequest() {
	req, _ := http.NewRequest(n.Method, n.Url, nil)
	for k, v := range n.Headers {
		req.Header.Set(k, v)
	}
	query := req.URL.Query()
	for k, v := range n.Query {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()
	n.Req = req
}
