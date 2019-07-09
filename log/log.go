package log

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"io"
	"runtime"
	"strings"

	//prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

// LogrusLogger 日志记录类
type LogrusLogger struct {
	*logrus.Logger
}

type LogrusEntry struct {
	*logrus.Entry
}

// 全局日志单例
var (
	Logger *LogrusLogger
)

// init 日志包初始化
func init() {
	Logger = &LogrusLogger{
		Logger: logrus.New(),
	}

	Logger.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		ShowFullLevel:   true,
		TrimMessages:    true,
	})

	Logger.SetReportCaller(true)

	//配置输出为标准输出
	Logger.SetOutput(os.Stdout)

	//配置钩子，根据日志时间和level打印到对应的文件
	Logger.AddHook(NewDailyHook())
}

// SetOutput 设置日志输出位置
func SetOutput(output io.Writer) {
	Logger.SetOutput(output)
}

// Mark 打标签，标识日志打印对象
func (l *LogrusLogger) Mark(key string) *LogrusEntry {
	return &LogrusEntry{Entry: l.Logger.WithField("", key)}
}

// Caller 标识日志打印方法
func (l *LogrusLogger) Caller(skip int) *LogrusEntry {
	p, _, _, _ := runtime.Caller(skip)
	caller := runtime.FuncForPC(p).Name()

	callerSplit := strings.Split(caller, ".")

	callerStruct := strings.Trim(callerSplit[len(callerSplit)-3], "()*")
	callerFunc := callerSplit[len(callerSplit)-2]

	return &LogrusEntry{Entry: l.Logger.WithField("a", "Service").WithField("b", callerStruct).WithField("c", callerFunc)}
}
