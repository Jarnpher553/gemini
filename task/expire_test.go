package task

import (
	"context"
	"fmt"
	"github.com/Jarnpher553/micro-core/redis"
	"testing"
)

func TestNewExpire(t *testing.T) {
	rd := redis.New(redis.Addr("47.105.208.81:6379"), redis.DB(15), redis.PoolSize(10), redis.Pwd("47.105.208.81"))
	exp := NewExpire(Redis(rd))
	if exp.handles == nil {
		t.Error("new expire error")
	}
}

func TestExpire_Assign(t *testing.T) {
	rd := redis.New(redis.Addr("47.105.208.81:6379"), redis.DB(15), redis.PoolSize(10), redis.Pwd("47.105.208.81"))
	exp := NewExpire(Redis(rd))
	exp.Assign("talk_order", func(payload string, opt *Options) {
		fmt.Println("执行任务:", payload)
	})
}

func TestExpire_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rd := redis.New(redis.Addr("47.105.208.81:6379"), redis.DB(15), redis.PoolSize(10), redis.Pwd("47.105.208.81"))
	exp := NewExpire(Redis(rd))
	exp.Assign("talk_order", func(payload string, opt *Options) {
		//fmt.Println("执行任务:", payload)
		t.Log("aa")
	}).Run(ctx)
	cancel()
}
