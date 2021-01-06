package tracing

import (
	"github.com/Jarnpher553/gemini/pkg/log"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapReporter struct {
	*log.ZapLogger
}

// NewReporter 构造函数
func NewZapReporter() reporter.Reporter {
	return &ZapReporter{
		ZapLogger: log.Zap.Mark("tracer"),
	}
}

// Send 实现Reporter接口
func (r *ZapReporter) Send(s model.SpanModel) {
	fields := make([]zapcore.Field, 0)
	fields = append(fields, zap.String("cost", s.Duration.String()))
	fields = append(fields, zap.String("trace_id", s.TraceID.String()))
	fields = append(fields, zap.String("id", s.ID.String()))
	fields = append(fields, zap.String("name", s.Name))
	fields = append(fields, zap.String("kind", string(s.Kind)))
	fields = append(fields, zap.String("http.method", s.Tags["http.method"]))
	fields = append(fields, zap.String("http.status_code", s.Tags["http.status_code"]))
	fields = append(fields, zap.String("http.url", s.Tags["http.url"]))

	r.Info("tracing", fields...)
}

// Close 实现Reporter接口
func (*ZapReporter) Close() error { return nil }
