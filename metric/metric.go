package metric

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"sync"
	"time"
)

// Metric 监控指标类
type Metric struct {
	name        string
	reg         metrics.Registry
	ReqCount    metrics.Counter
	ReqDuration metrics.Timer
	printer     IPrinter
	freq        time.Duration
	once        *sync.Once
}

type Config struct {
	ServiceName string
	Printer     IPrinter
	Freq        time.Duration
}

// IWriter 打印指标接口
type IPrinter interface {
	Print(m map[string]string)
}

// New 构造函数
func New(conf *Config) *Metric {
	metric := &Metric{reg: metrics.NewRegistry(), name: conf.ServiceName}

	reqCount := metrics.NewCounter()
	reqDuration := metrics.NewCustomTimer(metrics.NewHistogram(metrics.NewUniformSample(255)), metrics.NewMeter())

	metric.ReqCount = reqCount
	metric.ReqDuration = reqDuration
	metric.printer = conf.Printer
	metric.freq = conf.Freq
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
	metric.register()

	metric.once.Do(func() {
		go metric.log(metric.reg, metric.freq, time.Millisecond, metric.printer)
	})
}

// Stop 取消打印
func (metric *Metric) Stop() {
	metric.unregister()
}

func (metric *Metric) log(registry metrics.Registry, freq time.Duration, scale time.Duration, logger IPrinter) {
	du := float64(scale)
	duSuffix := scale.String()[1:]

	for range time.Tick(freq) {
		registry.Each(func(name string, i interface{}) {
			switch m := i.(type) {
			case metrics.Counter:
				logger.Print(map[string]string{"name": metric.name, "counter": name, "count": fmt.Sprintf("%d", m.Count())})
			case metrics.Timer:
				t := m.Snapshot()
				ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				logger.Print(map[string]string{
					"name":        metric.name,
					"timer":       name,
					"count":       fmt.Sprintf("%d", t.Count()),
					"min":         fmt.Sprintf("%.2f%s", float64(t.Min())/du, duSuffix),
					"max":         fmt.Sprintf("%.2f%s", float64(t.Max())/du, duSuffix),
					"mean":        fmt.Sprintf("%.2f%s", t.Mean()/du, duSuffix),
					"stddev":      fmt.Sprintf("%.2f%s", t.StdDev()/du, duSuffix),
					"median":      fmt.Sprintf("%.2f%s", ps[0]/du, duSuffix),
					"75%%":        fmt.Sprintf("%.2f%s", ps[1]/du, duSuffix),
					"95%%":        fmt.Sprintf("%.2f%s", ps[2]/du, duSuffix),
					"99%%":        fmt.Sprintf("%.2f%s", ps[3]/du, duSuffix),
					"99.9%%":      fmt.Sprintf("%.2f%s", ps[4]/du, duSuffix),
					"1-min rate":  fmt.Sprintf("%.2f", t.Rate1()),
					"5-min rate":  fmt.Sprintf("%.2f", t.Rate5()),
					"15-min rate": fmt.Sprintf("%.2f", t.Rate15()),
					"mean rate":   fmt.Sprintf("%.2f", t.RateMean()),
				})
			}
		})
	}
}
