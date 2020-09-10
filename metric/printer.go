package metric

import (
	"fmt"
	"github.com/Jarnpher553/gemini/log"
	"time"
)

// logPrinter 日志输出类
type logPrinter struct {
	logger *log.ZapLogger
	freq   time.Duration
}

// NewWriter 构造函数
func NewPrinter() IPrinter {
	return &logPrinter{
		logger: log.Zap.Mark("metric"),
	}
}

func (lw *logPrinter) Printf(format string, v ...interface{}) {
	lw.logger.Info(fmt.Sprintf(format, v))
}
