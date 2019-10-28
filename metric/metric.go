package metric

import (
	"github.com/rcrowley/go-metrics"
)

// Metric 监控指标类
type Metric struct {
	reg         metrics.Registry
	ReqCount    metrics.Counter
	ReqDuration metrics.Timer
}

// IWriter 打印指标接口
type IWriter interface {
	Write(*Metric)
}

// New 构造函数
func New(writer IWriter) *Metric {
	metric := &Metric{reg: metrics.NewRegistry()}

	reqCount := metrics.NewCounter()
	reqDuration := metrics.NewCustomTimer(metrics.NewHistogram(metrics.NewUniformSample(255)), metrics.NewMeter())

	metric.ReqCount = reqCount
	metric.ReqDuration = reqDuration

	go writer.Write(metric)

	//metric.Start()

	return metric
}

// Start 开始打印
func (metric *Metric) Start() {
	metric.reg.GetOrRegister("reqCount", metric.ReqCount)
	metric.reg.GetOrRegister("reqDuration", metric.ReqDuration)
}

// Stop 取消打印
func (metric *Metric) Stop() {
	metric.reg.UnregisterAll()
}
