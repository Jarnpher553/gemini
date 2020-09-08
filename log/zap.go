package log

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"regexp"
	"runtime"
	"strings"
)

type ZapLogger struct {
	*zap.Logger
}

var Zap *ZapLogger

func init() {
	new()
}

func new() {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          "json",
		DisableStacktrace: true,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := config.Build()
	Zap = &ZapLogger{Logger: logger}
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
	return &ZapLogger{l.Logger.With(zap.String("mod", strings.ToLower(key)))}
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
			fields = append(fields, zap.String("mod", strings.ToLower(strings.Split(callerSplit[i], "/")[1])))
		} else if strings.Contains(callerSplit[i], "*") {
			fields = append(fields, zap.String("struct", strings.ToLower(strings.Trim(callerSplit[i], "()*"))))
		} else {
			fields = append(fields, zap.String("func", strings.ToLower(callerSplit[i])))
		}
	}

	return &ZapLogger{l.Logger.With(fields...)}
}

func Message(messages ...interface{}) string {
	buf := bytes.Buffer{}
	for i := range messages {
		if i != 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprint(messages[i]))
	}

	return buf.String()
}

func Messagef(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}
