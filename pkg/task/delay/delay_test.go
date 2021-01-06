package delay

import (
	"github.com/Jarnpher553/gemini/pkg/redis"
	"github.com/Jarnpher553/gemini/pkg/task"
	REDIS "github.com/go-redis/redis"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	rd := redis.New(redis.Addr("47.105.208.81:6379"), redis.DB(15), redis.PoolSize(10), redis.Pwd("Iseeyou2016isu1118"))
	Bind(true, task.Redis(rd))
	Assign("talk_order", func(payload interface{}, opt *task.Options) {
		t.Log("执行对话任务", payload)
	})
	Assign("sell_order", func(payload interface{}, opt *task.Options) {
		t.Log("执行购买任务", payload)
	})
	Run()
	go func() {
		rd.ZAdd("talk_order", REDIS.Z{Score: float64(time.Now().Add(1*time.Second).UnixNano() / 1e6), Member: "order_1"})
		rd.ZAdd("talk_order", REDIS.Z{Score: float64(time.Now().Add(2*time.Second).UnixNano() / 1e6), Member: "order_2"})
		rd.ZAdd("talk_order", REDIS.Z{Score: float64(time.Now().Add(3*time.Second).UnixNano() / 1e6), Member: "order_3"})
		rd.ZAdd("talk_order", REDIS.Z{Score: float64(time.Now().Add(4*time.Second).UnixNano() / 1e6), Member: "order_4"})
		rd.ZAdd("talk_order", REDIS.Z{Score: float64(time.Now().Add(5*time.Second).UnixNano() / 1e6), Member: "order_5"})
		rd.ZAdd("sell_order", REDIS.Z{Score: float64(time.Now().Add(1*time.Second).UnixNano() / 1e6), Member: "order_a"})
		rd.ZAdd("sell_order", REDIS.Z{Score: float64(time.Now().Add(2*time.Second).UnixNano() / 1e6), Member: "order_b"})
		rd.ZAdd("sell_order", REDIS.Z{Score: float64(time.Now().Add(3*time.Second).UnixNano() / 1e6), Member: "order_c"})
		rd.ZAdd("sell_order", REDIS.Z{Score: float64(time.Now().Add(4*time.Second).UnixNano() / 1e6), Member: "order_d"})
		rd.ZAdd("sell_order", REDIS.Z{Score: float64(time.Now().Add(5*time.Second).UnixNano() / 1e6), Member: "order_e"})
	}()
	<-time.After(10 * time.Second)
	Stop()

	<-time.After(10 * time.Second)
}

func TestRun2(t *testing.T) {
	rd := redis.New(redis.Addr("47.105.208.81:6379"), redis.DB(0), redis.PoolSize(10), redis.Pwd("Iseeyou2016isu1118"))
	a := rd.ZRank("bbb", "100").Err()
	t.Log(a)
	b := rd.ZRank("bbb", "345345345")
	if b.Err() != REDIS.Nil{
		t.Log(b.Err())
	}
}
