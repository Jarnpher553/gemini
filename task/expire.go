package task

import (
	"context"
	"github.com/Jarnpher553/micro-core/mongo"
	"github.com/Jarnpher553/micro-core/redis"
	"github.com/Jarnpher553/micro-core/repo"
	"strings"
	"sync"
)

//失效任务
type Expire struct {
	options *Options
	handles map[string]Handle
	m       *sync.Mutex
}

//处理程序
type Handle func(string, *Options)

//配置
type Option func(*Expire)

//redis配置
func Redis(rd *redis.RdClient) Option {
	return func(exp *Expire) {
		exp.options.redis = rd
	}
}

//repo配置
func Repo(rp *repo.Repository) Option {
	return func(exp *Expire) {
		exp.options.repo = rp
	}
}

//mongo配置
func Mongo(mg *mongo.MgoClient) Option {
	return func(exp *Expire) {
		exp.options.mgo = mg
	}
}

//绑定配置并运行
func NewExpire(options ...Option) *Expire {
	//失效任务实例
	var exp = &Expire{handles: make(map[string]Handle), m: &sync.Mutex{}, options: &Options{}}

	for _, op := range options {
		op(exp)
	}
	return exp
}

//分配失效任务
func (exp *Expire) Assign(name string, handle Handle) *Expire {
	exp.m.Lock()
	exp.handles[name] = handle
	exp.m.Unlock()
	return exp
}

//执行任务
func (exp *Expire) Run(stop context.Context) {
	pubSub := exp.options.redis.PSubscribe("__keyevent@*__:expired")
	ch := pubSub.Channel()

	go func() {
	f:
		for {
			select {
			case msg := <-ch:
				for k := range exp.handles {
					if strings.Contains(msg.Payload, k) {
						exp.handles[k](msg.Payload, exp.options)
						break
					}
				}
			case <-stop.Done():
				break f
			}
		}
	}()
}
