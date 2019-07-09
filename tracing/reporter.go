package tracing

import (
	"fmt"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/Janrpher553/micro-core/log"
)

// logrusReporter 日志报告者类，用于tracer输出
type logrusReporter struct {
	*log.LogrusLogger
}

// NewReporter 构造函数
func NewReporter() reporter.Reporter {
	return &logrusReporter{
		LogrusLogger: log.Logger,
	}
}

// Send 实现Reporter接口
func (r *logrusReporter) Send(s model.SpanModel) {
	if b, err := s.MarshalJSON(); err == nil {
		r.Mark( "Tracer").Info(fmt.Sprintf("%s", string(b)))
	}
}

// Close 实现Reporter接口
func (*logrusReporter) Close() error { return nil }
