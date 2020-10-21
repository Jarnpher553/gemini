package queue

import (
	"fmt"
	"github.com/Jarnpher553/gemini/log"
	"github.com/Jarnpher553/gemini/mongo"
	"github.com/Jarnpher553/gemini/redis"
	"github.com/Jarnpher553/gemini/repo"
	"github.com/adjust/rmq/v3"
	cmap "github.com/orcaman/concurrent-map"
	"go.uber.org/zap"
	"time"
)

var conn = &RedisMessageConn{conf: &Configuration{}, openQueues: cmap.New()}
var logger = log.Logger.Mark("rmq")

type RedisMessageConn struct {
	conn       rmq.Connection
	name       string
	conf       *Configuration
	openQueues cmap.ConcurrentMap
}

type Configuration struct {
	redis  *redis.RdClient
	repo   *repo.Repository
	mgo    *mongo.MgoClient
	custom interface{}
}

func Redis(rd *redis.RdClient) Conf {
	return func(conn *RedisMessageConn) {
		conn.conf.redis = rd
	}
}

func Repo(rp *repo.Repository) Conf {
	return func(conn *RedisMessageConn) {
		conn.conf.repo = rp
	}
}

func Mongo(mg *mongo.MgoClient) Conf {
	return func(conn *RedisMessageConn) {
		conn.conf.mgo = mg
	}
}

func Custom(custom interface{}) Conf {
	return func(conn *RedisMessageConn) {
		conn.conf.custom = custom
	}
}

func Name(name string) Conf {
	return func(conn *RedisMessageConn) {
		conn.name = name
	}
}

func (c *Configuration) Redis() *redis.RdClient {
	return c.redis
}

func (c *Configuration) Repo() *repo.Repository {
	return c.repo
}

func (c *Configuration) Mongo() *mongo.MgoClient {
	return c.mgo
}

func (c *Configuration) Custom() interface{} {
	return c.custom
}

type Conf func(*RedisMessageConn)

func Bind(conf ...Conf) {
	conn.name = "rmq"

	for _, v := range conf {
		v(conn)
	}

	if conn.conf.redis == nil {
		logger.Fatal("has no redis client to initial")
	}

	var err error
	conn.conn, err = rmq.OpenConnectionWithRedisClient(conn.name, conn.conf.redis.Client, nil)
	if err != nil {
		logger.Fatal("can not open connection")
	}

	go startClean()
}

func queue(name string) (rmq.Queue, error) {
	val, ok := conn.openQueues.Get(name)
	if ok {
		return val.(rmq.Queue), nil
	} else {
		q, err := conn.conn.OpenQueue(name)
		if err != nil {
			return nil, err
		}
		conn.openQueues.Set(name, q)

		return q, nil
	}
}

func Publish(name string, payload interface{}) error {
	q, err := queue(name)
	if err != nil {
		return err
	}

	err = q.Publish(payload.(string))
	if err != nil {
		return err
	}
	return nil
}

type Func func(rmq.Delivery, *Configuration)
type FuncBatch func(rmq.Deliveries, *Configuration)
type BatchConsumerFunc func(deliveries rmq.Deliveries)

func (batchConsumerFunc BatchConsumerFunc) Consume(delivery rmq.Deliveries) {
	batchConsumerFunc(delivery)
}

func Assign(queueName string, prefetchLimit int64, duration time.Duration, f Func, pushQueueFunc ...Func) error {
	q, err := queue(queueName)
	if err != nil {
		return err
	}
	err = q.StartConsuming(prefetchLimit, duration)
	if err != nil {
		return err
	}

	if err := pushQueue(queueName, q, pushQueueFunc...); err != nil {
		return err
	}

	_, err = q.AddConsumerFunc(queueName+"-consumer", decorator(conn.conf, f))
	if err != nil {
		return err
	}

	return nil
}

func AssignBatch(queueName string, prefetchLimit int64, duration time.Duration, batchSize int64, timeout time.Duration, f FuncBatch, pushQueueFunc ...Func) error {
	q, err := queue(queueName)
	if err != nil {
		return err
	}
	err = q.StartConsuming(prefetchLimit, duration)
	if err != nil {
		return err
	}

	if err := pushQueue(queueName, q, pushQueueFunc...); err != nil {
		return err
	}

	_, err = q.AddBatchConsumer(queueName+"-consumer", batchSize, timeout, BatchConsumerFunc(decoratorBatch(conn.conf, f)))
	if err != nil {
		return err
	}

	return nil
}

func decorator(configuration *Configuration, f func(rmq.Delivery, *Configuration)) func(rmq.Delivery) {
	return func(delivery rmq.Delivery) {
		f(delivery, configuration)
	}
}

func decoratorBatch(configuration *Configuration, f func(rmq.Deliveries, *Configuration)) func(rmq.Deliveries) {
	return func(delivery rmq.Deliveries) {
		f(delivery, configuration)
	}
}

func startClean() {
	cleaner := rmq.NewCleaner(conn.conn)

	for range time.Tick(time.Minute) {
		returned, err := cleaner.Clean()
		if err != nil {
			logger.With(zap.String("err", err.Error())).Error("clean")
			continue
		}
		logger.With(zap.Int64("count", returned)).Error("clean")
	}
}

func StopConsuming(queueName string) error {
	q, err := queue(queueName)
	if err != nil {
		return err
	}
	<-q.StopConsuming()
	conn.openQueues.Remove(queueName)
	return nil
}

func StopAllConsuming() error {
	<-conn.conn.StopAllConsuming()
	for _, key := range conn.openQueues.Keys() {
		conn.openQueues.Remove(key)
	}
	return nil
}

func pushQueue(queueName string, q rmq.Queue, pushQueueFunc ...Func) error {
	var sq rmq.Queue
	sq = q
	for i, f := range pushQueueFunc {
		pq, err := queue(fmt.Sprintf("%s-%s-%d", queueName, "pushQ", i))
		if err != nil {
			return err
		}
		sq.SetPushQueue(pq)
		_, err = pq.AddConsumerFunc(fmt.Sprintf("%s-%s-%d", queueName, "pushQ", i), decorator(conn.conf, f))
		if err != nil {
			return err
		}
		sq = pq
	}
	sq.SetPushQueue(q)
	return nil
}

type QueueType int

const (
	Rejected QueueType = iota
	Ready
	Unacked
)

func Purge(queueName string, qt QueueType) error {
	q, err := queue(queueName)
	if err != nil {
		return err
	}
	if qt == Rejected {
		_, err := q.PurgeRejected()
		return err
	} else if qt == Ready {
		_, err := q.PurgeReady()
		return err
	} else {
		return nil
	}
}

func Return(queueName string, qt QueueType, max int64) error {
	q, err := queue(queueName)
	if err != nil {
		return err
	}
	if qt == Rejected {
		_, err := q.ReturnRejected(max)
		return err
	} else if qt == Unacked {
		_, err := q.ReturnUnacked(max)
		return err
	} else {
		return nil
	}
}

//func demo(delivery rmq.Delivery){
//	delivery.Ack()
//}
