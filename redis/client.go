package redis

import (
	"github.com/Jarnpher553/gemini/log"
	"github.com/go-redis/redis/v7"
	"time"
)

// RdClient redis客户端类
type RdClient struct {
	*redis.Client
	logger *log.ZapLogger
}

// Option 配置项方法
type Option func(*redis.Options)

// Addr 地址配置项
func Addr(addr string) Option {
	return func(option *redis.Options) {
		option.Addr = addr
	}
}

// DB 数据库配置项
func DB(db int) Option {
	return func(option *redis.Options) {
		option.DB = db
	}
}

// Pwd 密码配置项
func Pwd(pwd string) Option {
	return func(option *redis.Options) {
		option.Password = pwd
	}
}

// PoolSize 池大小配置项
func PoolSize(size int) Option {
	return func(option *redis.Options) {
		option.PoolSize = size
	}
}

// New 构造函数
func New(options ...Option) *RdClient {
	option := &redis.Options{}

	for _, op := range options {
		op(option)
	}

	client := &RdClient{
		redis.NewClient(option),
		log.Zap.Mark("redis"),
	}

	err := client.Ping().Err()
	if err != nil {
		client.logger.Fatal(log.Message("redis connected error:", err.Error()))
	}
	return client
}

// 以下是redis操作

func (r *RdClient) IncrStr(key string) string {
	return r.Client.Incr(key).String()
}

func (r *RdClient) IncrInt(key string) int64 {
	return r.Client.Incr(key).Val()
}

func (r *RdClient) Get(key string) string {
	return r.Client.Get(key).Val()
}

func (r *RdClient) Del(key string) bool {
	return r.Client.Del(key).Val() == 1
}

func (r *RdClient) Set(key string, val interface{}, expiration time.Duration) bool {
	return r.Client.Set(key, val, expiration).Val() == "OK"
}

func (r *RdClient) Publish(channel string, msg interface{}) bool {
	return r.Client.Publish(channel, msg).Val() == 1
}

func (r *RdClient) Exists(key string) bool {
	return r.Client.Exists(key).Val() == 1
}

func (r *RdClient) Expire(key string, expiration time.Duration) bool {
	return r.Client.Expire(key, expiration).Val()
}

func (r *RdClient) SetNX(key string, val interface{}, expiration time.Duration) bool {
	return r.Client.SetNX(key, val, expiration).Val()
}

func (r *RdClient) SexEX(key string, val interface{}, expiration time.Duration) bool {
	return r.Client.SetXX(key, val, expiration).Val()
}

func (r *RdClient) HGet(key, field string) string {
	return r.Client.HGet(key, field).Val()
}

func (r *RdClient) HSet(key, field string, val interface{}) bool {
	return r.Client.HSet(key, field, val).Val() != 0
}

func (r *RdClient) HSetNX(key, field string, val interface{}) bool {
	return r.Client.HSetNX(key, field, val).Val()
}

func (r *RdClient) LPush(key string, val interface{}) bool {
	return r.Client.LPush(key, val).Val() > 0
}

func (r *RdClient) LRange(key string, start int64, stop int64) []string {
	return r.Client.LRange(key, start, stop).Val()
}

func (r *RdClient) Keys(pattern string) []string {
	return r.Client.Keys(pattern).Val()
}

func (r *RdClient) BRPop(timeout time.Duration, keys ...string) []string {
	return r.Client.BRPop(timeout, keys...).Val()
}

func (r *RdClient) TTL(key string) float64 {
	return r.Client.TTL(key).Val().Seconds()
}
