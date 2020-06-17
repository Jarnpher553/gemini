package httpclient

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/openzipkin/zipkin-go/model"
	"net/http"
	"time"
)

// Transport http客户端自定义传输类
type Transport struct {
	http.RoundTripper
	ServiceName string
}

// RoundTrip 实现传输类必要方法
func (tran *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx := r.Context()

	span, _ := opentracing.StartSpanFromContext(ctx, tran.ServiceName,
		ext.SpanKindRPCClient,
		opentracing.StartTime(time.Now()),
		opentracing.Tag{Key: string(ext.HTTPUrl), Value: r.URL.Path},
		opentracing.Tag{Key: string(ext.HTTPMethod), Value: r.Method})
	defer span.Finish()

	_ = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

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

		if sc.ParentID != nil {
			r.Header.Set("jar-parentid", sc.ParentID.String())
		} else {
			r.Header.Set("jar-parentid", "0")
		}
	}
}
