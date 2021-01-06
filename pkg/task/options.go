package task

import (
	"github.com/Jarnpher553/gemini/pkg/mongo"
	"github.com/Jarnpher553/gemini/pkg/redis"
	"github.com/Jarnpher553/gemini/pkg/repo"
)

type Options struct {
	Redis *redis.RdClient
	Repo  *repo.Repository
	Mgo   *mongo.MgoClient
}

//配置
type Option func(*Options)

//redis配置
func Redis(rd *redis.RdClient) Option {
	return func(opt *Options) {
		opt.Redis = rd
	}
}

//repo配置
func Repo(rp *repo.Repository) Option {
	return func(opt *Options) {
		opt.Repo = rp
	}
}

//mongo配置
func Mongo(mg *mongo.MgoClient) Option {
	return func(opt *Options) {
		opt.Mgo = mg
	}
}

//处理程序
type Handle func(interface{}, *Options)