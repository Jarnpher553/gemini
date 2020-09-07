package tracing

import (
	"fmt"
	"github.com/Jarnpher553/gemini/log"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
)

// logrusReporter 日志报告者类，用于tracer输出
//type logrusReporter struct {
//	*log.LogrusLogger
//}
//
//// NewReporter 构造函数
//func NewLogrusReporter() reporter.Reporter {
//	return &logrusReporter{
//		LogrusLogger: log.Logrus,
//	}
//}
//
//// Send 实现Reporter接口
//func (r *logrusReporter) Send(s model.SpanModel) {
//	if b, err := s.MarshalJSON(); err == nil {
//		r.Mark("Tracer").Info(fmt.Sprintf("%s", string(b)))
//	}
//}

// Close 实现Reporter接口
//func (*logrusReporter) Close() error { return nil }

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
	if b, err := s.MarshalJSON(); err == nil {
		r.Info(fmt.Sprintf("%s", string(b)))
	}
}

// Close 实现Reporter接口
func (*ZapReporter) Close() error { return nil }
