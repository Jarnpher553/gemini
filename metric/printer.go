package metric

import (
	"github.com/Jarnpher553/gemini/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logPrinter 日志输出类
type logPrinter struct {
	logger *log.ZapLogger
}

// NewWriter 构造函数
func NewPrinter() IPrinter {
	return &logPrinter{
		logger: log.Zap.Mark("metric"),
	}
}

func (lw *logPrinter) Print(m map[string]string) {
	var fields []zapcore.Field
	for key := range m {
		fields = append(fields, zap.String(key, m[key]))
	}
	lw.logger.Info("monitor", fields...)
}
