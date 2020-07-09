package log

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
	"strings"

	//prefixed "github.com/x-cray/logrus-prefixed-formatter"
	//"os"
	"regexp"
)

// LogrusLogger 日志记录类
type LogrusLogger struct {
	*logrus.Logger
	fire bool
}

type LogrusEntry struct {
	*logrus.Entry
}

// 全局日志单例
var (
	Logrus *LogrusLogger
)

//init 日志包初始化
func LogrusN() {
	Logrus = &LogrusLogger{
		Logger: logrus.New(),
	}

	Logrus.SetFormatter(&nested.Formatter{
		FieldsOrder: []string{
			"source",
			"module",
			"method",
		},
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		ShowFullLevel:   true,
		TrimMessages:    true,
		NoColors:        true,
	})

	//Logrus.SetFormatter(&logrus.JSONFormatter{})

	Logrus.SetReportCaller(true)

	//配置输出为标准输出
	Logrus.SetOutput(os.Stdout)

	//配置钩子，根据日志时间和level打印到对应的文件
	//Logrus.AddHook(NewHourHook())
}

// SetOutput 设置日志输出位置
func SetOutput(output io.Writer) {
	Logrus.SetOutput(output)
}

// Mark 打标签，标识日志打印对象
func (l *LogrusLogger) Mark(key string) *LogrusEntry {
	return &LogrusEntry{Entry: l.Logger.WithField("source", key)}
}

// Caller 标识日志打印方法
func (l *LogrusLogger) Caller(skip int) *LogrusEntry {
	p, _, _, _ := runtime.Caller(skip)
	caller := runtime.FuncForPC(p).Name()

	callerSplit := strings.Split(caller, ".")

	callers := make(map[string]interface{})

	reg1 := regexp.MustCompile(`^func[0-9]$`)
	reg2 := regexp.MustCompile(`^[0-9]$`)
	for i := range callerSplit {
		if reg1.MatchString(callerSplit[i]) || reg2.MatchString(callerSplit[i]) {
			continue
		}
		if strings.Contains(callerSplit[i], "/") {
			callers["source"] = strings.Title(strings.Split(callerSplit[i], "/")[1])
		} else if strings.Contains(callerSplit[i], "*") {
			callers["module"] = strings.Trim(callerSplit[i], "()*")
		} else {
			callers["method"] = callerSplit[i]
		}
	}

	return &LogrusEntry{Entry: l.Logger.WithFields(callers)}
}