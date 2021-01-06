package httpclient

import (
	"context"
	"encoding/json"
	"github.com/Jarnpher553/gemini/pkg/tracing"
	"gopkg.in/resty.v1"
	"net/http"
)

// ReqClient http客户端
type ReqClient struct {
	*tracing.Tracer
	ServiceName string
}

// Option 配置函数
type Option func(*ReqClient)

// Tracer 链路跟踪配置
func Tracer(tracer *tracing.Tracer) Option {
	return func(client *ReqClient) {
		client.Tracer = tracer
	}
}

// Name 服务名称配置
func Name(name string) Option {
	return func(client *ReqClient) {
		client.ServiceName = name
	}
}

// New 构造函数
func New(options ...Option) *ReqClient {
	client := resty.GetClient()

	rc := &ReqClient{}

	for _, option := range options {
		option(rc)
	}

	client.Transport = &Transport{
		http.DefaultTransport,
		rc.ServiceName,
	}

	return rc
}

func R() *resty.Request {
	return resty.R()
}

// RGet get请求
func (c *ReqClient) RGet(url string, query map[string]string, ctx context.Context, v interface{}) error {
	request := resty.R()

	resp, err := request.SetContext(ctx).
		SetQueryParams(query).
		SetResult(v).
		Get(url)

	if resp.Header()["Content-Type"][0] == "text/plain" {
		_ = json.Unmarshal(resp.Body(), v)
	}

	if err != nil {
		return err
	} else {
		return nil
	}
}

// RPost post请求
func (c *ReqClient) RPost(url string, body interface{}, ctx context.Context, v interface{}) error {
	request := resty.R()

	_, err := request.SetContext(ctx).
		SetBody(body).
		SetResult(v).
		Post(url)

	if err != nil {
		return err
	} else {
		return nil
	}
}

// RPut put请求
func (c *ReqClient) RPut(url string, body interface{}, ctx context.Context, v interface{}) error {
	request := resty.R()

	_, err := request.SetContext(ctx).
		SetBody(body).
		SetResult(v).
		Put(url)

	if err != nil {
		return err
	} else {
		return nil
	}
}

// RDelete delete请求
func (c *ReqClient) RDelete(url string, query map[string]string, ctx context.Context, v interface{}) error {
	request := resty.R()

	request = request.SetContext(ctx).SetResult(v)
	if query != nil {
		request = request.SetQueryParams(query)
	}

	_, err := request.Delete(url)

	if err != nil {
		return err
	} else {
		return nil
	}
}
