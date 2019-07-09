package httpclient

import (
	"github.com/openzipkin/zipkin-go/model"
	"net/http"
	"gitee.com/jarnpher_rice/micro-core/tracing"
)

// Transport http客户端自定义传输类
type Transport struct {
	http.RoundTripper
}

// RoundTrip 实现传输类必要方法
func (tran *Transport) RoundTrip(r *http.Request) (*http.Response, error) {

	ctx := r.Context()

	var sc model.SpanContext
	if span := tracing.SpanFromContext(ctx); span != nil {
		sc = span.Context()
	}

	InjectHttp(r)(sc)

	return tran.RoundTripper.RoundTrip(r)
}

// InjectHttp 将context写入http请求头
func InjectHttp(r *http.Request) func(model.SpanContext) {
	return func(sc model.SpanContext) {
		if sc.Debug {
			r.Header.Set("jar-flags", "1")
		}

		if *sc.Sampled {
			r.Header.Set("jar-sampled", "1")
		} else {
			r.Header.Set("jar-sampled", "0")
		}

		r.Header.Set("jar-traceid", sc.TraceID.String())
		r.Header.Set("jar-spanid", sc.ID.String())
		r.Header.Set("jar-parentid", sc.ParentID.String())
	}
}
