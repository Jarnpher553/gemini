package httpclient

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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
