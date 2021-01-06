package expire

import (
	"github.com/Jarnpher553/gemini/pkg/redis"
	"github.com/Jarnpher553/gemini/pkg/task"
	"testing"
	"time"
)

//func TestNewExpire(t *testing.T) {
//	rd := redis.New(redis.Addr("47.105.208.81:6379"), redis.DB(15), redis.PoolSize(10), redis.Pwd("Iseeyou2016isu1118"))
//	Bind(true, task.Redis(rd))
//	if exp.handles == nil {
//		t.Error("new expire error")
//	}
//}
//
//func TestExpire_Assign(t *testing.T) {
//	rd := redis.New(redis.Addr("47.105.208.81:6379"), redis.DB(15), redis.PoolSize(10), redis.Pwd("Iseeyou2016isu1118"))
//	Bind(true, task.Redis(rd))
//	Assign("talk_order", func(payload interface{}, opt *task.Options) {
//		fmt.Println("执行任务:", payload)
//	})
//}

func TestExpire_Run(t *testing.T) {
	rd := redis.New(redis.Addr("47.105.208.81:6379"), redis.DB(15), redis.PoolSize(10), redis.Pwd("Iseeyou2016isu1118"))
	Bind(true, task.Redis(rd))
	Assign("talk_order", func(payload interface{}, opt *task.Options) {
		t.Log("消费执行任务", payload)
	})
	Run()
	//go func() {
	//	//	rd.Set("talk_order_1", 10, time.Second*5)
	//	//	rd.Set("talk_order_2", 10, time.Second*5)
	//	//}()
	<-time.After(10 * time.Second)
	Stop()
	<-time.After(10 * time.Second)
}
