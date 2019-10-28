package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

// DailyHook 每日日志钩子
type DailyHook struct {
	filePath string
	sync.Once
}

func NewDailyHook() *DailyHook {
	return &DailyHook{filePath: "./logs/"}
}

// Fire 钩子方法，根据日志时间和等级输出到不同的文件
func (h *DailyHook) Fire(entry *logrus.Entry) error {
	h.Do(func() {
		_ = os.MkdirAll(h.filePath, os.ModePerm)
	})

	level := entry.Level.String()

	fileName := fmt.Sprintf("%s_%s.log", entry.Time.Format("2006-01-02"), level)

	file := h.filePath + fileName

	fs, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	defer func() {
		_ = fs.Close()
	}()

	if err != nil {
		return err
	}

	b, err := entry.Logger.Formatter.Format(entry)

	if err != nil {
		return err
	}

	_, err = fs.Write(b)

	if err != nil {
		return err
	}

	return nil
}

// Levels 钩子打印所有级别日志
func (h *DailyHook) Levels() []logrus.Level {
	return logrus.AllLevels
}


// HourHook 每小时日志钩子
type HourHook struct {
	filePath string
	sync.Once
}

func NewHourHook() *DailyHook {
	return &DailyHook{filePath: "./logs/"}
}

// Fire 钩子方法，根据日志时间和等级输出到不同的文件
func (h *HourHook) Fire(entry *logrus.Entry) error {
	h.Do(func() {
		_ = os.MkdirAll(h.filePath, os.ModePerm)
	})

	level := entry.Level.String()

	fileName := fmt.Sprintf("%s_%s.log", entry.Time.Format("2006-01-02H15"), level)

	file := h.filePath + fileName

	fs, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	defer func() {
		_ = fs.Close()
	}()

	if err != nil {
		return err
	}

	b, err := entry.Logger.Formatter.Format(entry)

	if err != nil {
		return err
	}

	_, err = fs.Write(b)

	if err != nil {
		return err
	}

	return nil
}

// Levels 钩子打印所有级别日志
func (h *HourHook) Levels() []logrus.Level {
	return logrus.AllLevels
}