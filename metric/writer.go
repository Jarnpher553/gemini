package metric

import (
	"github.com/rcrowley/go-metrics"
	"github.com/Jarnpher553/micro-core/log"
	"time"
)

// logWriter 日志输出类
type logWriter struct {
	*log.LogrusLogger
	freq time.Duration
}

// NewWriter 构造函数
func NewWriter(freq time.Duration) IWriter {
	return &logWriter{
		LogrusLogger: log.Logger,
		freq:         freq,
	}
}

// Write 实现IWriter接口
func (w *logWriter) Write(m *Metric) {
	metrics.Log(m.reg, w.freq, w.Mark("Metric"))
}
