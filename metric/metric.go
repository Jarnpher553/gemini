package metric

import (
	"github.com/rcrowley/go-metrics"
	"sync"
	"time"
)

// Metric 监控指标类
type Metric struct {
	reg         metrics.Registry
	ReqCount    metrics.Counter
	ReqDuration metrics.Timer
	printer     IPrinter
	freq        time.Duration
	once        *sync.Once
}

// IWriter 打印指标接口
type IPrinter interface {
	Printf(format string, v ...interface{})
}

// New 构造函数
func New(printer IPrinter, freq time.Duration) *Metric {
	metric := &Metric{reg: metrics.NewRegistry()}

	reqCount := metrics.NewCounter()
	reqDuration := metrics.NewCustomTimer(metrics.NewHistogram(metrics.NewUniformSample(255)), metrics.NewMeter())

	metric.ReqCount = reqCount
	metric.ReqDuration = reqDuration
	metric.printer = printer
	metric.freq = freq
	metric.once = &sync.Once{}

	return metric
}

func (metric *Metric) register() {
	metric.reg.GetOrRegister("reqCount", metric.ReqCount)
	metric.reg.GetOrRegister("reqDuration", metric.ReqDuration)
}

func (metric *Metric) unregister() {
	metric.reg.UnregisterAll()
}

// Start 开始打印
func (metric *Metric) Start() {
	metric.once.Do(func() {
		go metrics.Log(metric.reg, metric.freq, metric.printer)
	})
	metric.register()
}

// Stop 取消打印
func (metric *Metric) Stop() {
	metric.unregister()
}
