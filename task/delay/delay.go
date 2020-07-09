package delay

import (
	"context"
	"github.com/Jarnpher553/micro-core/log"
	"github.com/Jarnpher553/micro-core/task"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

//失效任务实例
var delay = &Delay{handles: make(map[string]task.Handle), m: &sync.Mutex{}, options: &task.Options{}, logger: log.Zap.Mark("TaskDelay")}

type Delay struct {
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
		op(delay.options)
	}
	delay.stop = ctx
	delay.cancel = cancel
	if autoRun {
		Run()
	}
}

//分配失效任务
func Assign(name string, handle task.Handle) {
	delay.m.Lock()
	delay.handles[name] = handle
	delay.m.Unlock()
}

//执行任务
func Run() {
	for key := range delay.handles {
		go func(k string) {
		For:
			for {
				select {
				case <-delay.stop.Done():
					delay.logger.Info(log.Message(k, "stopped"))
					break For
				default:
					<-time.After(100 * time.Millisecond)
					now := time.Now().UnixNano() / 1e6
					zset := delay.options.Redis.ZRangeByScoreWithScores(k, redis.ZRangeBy{"-inf", "+inf", 0, 1}).Val()
					if len(zset) == 0 {
						continue
					}
					score := zset[0].Score
					if float64(now) >= score {
						go delay.handles[k](zset[0].Member, delay.options)
						delay.options.Redis.ZRem(k, zset[0].Member)
					}
				}
			}
		}(key)
	}
}

func Stop() {
	delay.cancel()
}

func Join(taskName string, duration time.Duration, value string) {
	delay.options.Redis.ZAdd(taskName, redis.Z{Score: float64(time.Now().Add(duration).UnixNano() / 1e6), Member: value})
}

func Timestamp(taskName string, value string) float64 {
	return delay.options.Redis.ZScore(taskName, value).Val()
}

func Exist(taskName string, value string) bool {
	err := delay.options.Redis.ZRank(taskName, value).Err()
	if err == redis.Nil {
		return false
	}
	return true
}
