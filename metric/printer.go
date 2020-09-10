package metric

import (
	"fmt"
	"github.com/Jarnpher553/gemini/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"sync"
	"time"
)

// logPrinter 日志输出类
type logPrinter struct {
	logger *log.ZapLogger

	fields []zapcore.Field
	once   *sync.Once
}

// NewWriter 构造函数
func NewPrinter() IPrinter {
	return &logPrinter{
		logger: log.Zap.Mark("metric"),
		once:   &sync.Once{},
	}
}

func (lw *logPrinter) Printf(format string, v ...interface{}) {
	lw.once.Do(func() {
		go lw.timing()
	})
	if strings.Index(format, ":") != -1 {
		lw.fields = append(lw.fields, zap.Any(strings.Split(format, " ")[0], v[0]))
	} else {
		formatSlice := strings.Split(format, ":")
		key := strings.TrimSpace(formatSlice[0])
		value := fmt.Sprintf(strings.TrimSpace(formatSlice[1]), v...)
		lw.fields = append(lw.fields, zap.Any(key, value))
	}

}

func (lw *logPrinter) timing() {
	for range time.Tick(30 * time.Second) {
		log.Logger.Info("monitor", lw.fields...)
		lw.fields = lw.fields[:0]
	}
}
