package log

import (
	"fmt"
	"go.uber.org/zap"
	"regexp"
	"runtime"
	"strings"
)

type ZapLogger struct {
	*zap.Logger
}

var Zap *ZapLogger

func init() {
	Producation()
}

func Producation() {
	logger, _ := zap.NewProduction()
	Zap = &ZapLogger{
		logger,
	}
}

func Development() {
	logger, _ := zap.NewDevelopment()
	Zap = &ZapLogger{
		logger,
	}
}

func (l *ZapLogger) Mark(key string) *ZapLogger {
	return &ZapLogger{l.Logger.With(zap.String("source", key))}
}

func (l *ZapLogger) Caller(skip int) *ZapLogger {
	p, _, _, _ := runtime.Caller(skip)
	caller := runtime.FuncForPC(p).Name()

	callerSplit := strings.Split(caller, ".")

	reg1 := regexp.MustCompile(`^func[0-9]$`)
	reg2 := regexp.MustCompile(`^[0-9]$`)

	var fields []zap.Field
	for i := range callerSplit {
		if reg1.MatchString(callerSplit[i]) || reg2.MatchString(callerSplit[i]) {
			continue
		}
		if strings.Contains(callerSplit[i], "/") {
			fields = append(fields, zap.String("source", strings.Title(strings.Split(callerSplit[i], "/")[1])))
		} else if strings.Contains(callerSplit[i], "*") {
			fields = append(fields, zap.String("module", strings.Trim(callerSplit[i], "()*")))
		} else {
			fields = append(fields, zap.String("method", callerSplit[i]))
		}
	}

	return &ZapLogger{l.Logger.With(fields...)}
}

func Message(messages ...interface{}) string {
	return fmt.Sprint(messages)
}

func Messagef(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}
