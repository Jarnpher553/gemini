package metric

import (
	"fmt"
	"github.com/Jarnpher553/gemini/log"
	"github.com/rcrowley/go-metrics"
	"time"
)

// logWriter 日志输出类
type logWriter struct {
	logger *log.ZapLogger
	freq   time.Duration
}

// NewWriter 构造函数
func NewWriter(freq time.Duration) IWriter {
	return &logWriter{
		logger: log.Zap.Mark("metric"),
		freq:   freq,
	}
}

func (lw *logWriter) Printf(format string, v ...interface{}) {
	lw.logger.Info(fmt.Sprintf(format, v))
}

// Write 实现IWriter接口
func (lw *logWriter) Write(m *Metric) {
	metrics.Log(m.reg, lw.freq, lw)
}
