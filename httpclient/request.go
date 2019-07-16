package httpclient

import (
	"context"
	"github.com/Jarnpher553/micro-core/tracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"gopkg.in/resty.v1"
	"net/http"
	"time"
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
	client.Transport = &Transport{
		http.DefaultTransport,
	}

	rc := &ReqClient{}

	for _, option := range options {
		option(rc)
	}

	return rc
}

// RGet get请求
func (c *ReqClient) RGet(url string, query map[string]string, ctx context.Context, v interface{}) error {
	request := resty.R()

	if c.Tracer != nil {
		sp := c.getSpan(ctx, url, "GET")
		defer sp.Finish()

		ctx = tracing.NewContext(request.Context(), sp)
	}

	_, err := request.SetContext(ctx).
		SetQueryParams(query).
		SetResult(v).
		Get(url)

	if err != nil {
		return err
	} else {
		return nil
	}
}

// RPost post请求
func (c *ReqClient) RPost(url string, body interface{}, ctx context.Context, v interface{}) error {
	request := resty.R()

	if c.Tracer != nil {
		sp := c.getSpan(ctx, url, "POST")
		defer sp.Finish()

		ctx = tracing.NewContext(request.Context(), sp)
	}

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

	if c.Tracer != nil {
		sp := c.getSpan(ctx, url, "PUT")
		defer sp.Finish()

		ctx = tracing.NewContext(request.Context(), sp)
	}

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

	if c.Tracer != nil {
		sp := c.getSpan(ctx, url, "DELETE")
		defer sp.Finish()

		ctx = tracing.NewContext(request.Context(), sp)
	}

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

// getSpan 获取上下文跟踪对象
func (c *ReqClient) getSpan(ctx context.Context, url string, method string) zipkin.Span {
	var sc model.SpanContext
	if parentSpan := tracing.SpanFromContext(ctx); parentSpan != nil {
		sc = parentSpan.Context()
	}

	sp := c.Tracer.StartSpan(c.ServiceName, zipkin.Kind(model.Client), zipkin.Parent(sc), zipkin.StartTime(time.Now()))
	zipkin.TagHTTPUrl.Set(sp, url)
	zipkin.TagHTTPMethod.Set(sp, method)

	return sp
}
