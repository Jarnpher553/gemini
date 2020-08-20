package expire

import (
	"context"
	"github.com/Jarnpher553/gemini/log"
	"github.com/Jarnpher553/gemini/task"
	"strings"
	"sync"
)

//失效任务实例
var exp = &Expire{handles: make(map[string]task.Handle), m: &sync.Mutex{}, options: &task.Options{}, logger: log.Zap.Mark("TaskExpire")}

//失效任务
type Expire struct {
	options *task.Options
	handles map[string]task.Handle
	m       *sync.Mutex
	stop    context.Context
	cancel  context.CancelFunc
	logger  *log.ZapLogger
}

//绑定配置并运行
func Bind(autoRun bool, options ...task.Option) {
	ctx, cancel := context.WithCancel(context.Background())
	for _, op := range options {
		op(exp.options)
	}
	exp.stop = ctx
	exp.cancel = cancel
	if autoRun {
		Run()
	}
}

//分配失效任务
func Assign(name string, handle task.Handle) *Expire {
	exp.m.Lock()
	exp.handles[name] = handle
	exp.m.Unlock()
	return exp
}

//执行任务
func Run() {
	pubSub := exp.options.Redis.PSubscribe("__keyevent@*__:expired")
	ch := pubSub.Channel()

	go func() {
	For:
		for {
			select {
			case msg := <-ch:
				for k := range exp.handles {
					if strings.Contains(msg.Payload, k) {
						exp.handles[k](msg.Payload, exp.options)
						break
					}
				}
			case <-exp.stop.Done():
				for k := range exp.handles {
					exp.logger.Info(log.Message(k, "stopped"))
				}
				break For
			default:

			}

		}
	}()
}

func Stop() {
	exp.cancel()
}
