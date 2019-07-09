package scheduler

import (
	"github.com/Janrpher553/micro-core/mongo"
	"github.com/Janrpher553/micro-core/redis"
	"github.com/Janrpher553/micro-core/repo"
	"github.com/roylee0704/gron"
	"time"
)

var sch = &Scheduler{cron: gron.New(), Options: &Options{}}

type Scheduler struct {
	cron *gron.Cron
	*Options
}

type Options struct {
	redis *redis.RdClient
	repo  *repo.Repository
	mgo   *mongo.MgoClient
}

func (o *Options) Redis() *redis.RdClient {
	return o.redis
}

func (o *Options) Repo() *repo.Repository {
	return o.repo
}

func (o *Options) Mongo() *mongo.MgoClient {
	return o.mgo
}

type Option func(*Scheduler)

type Job func(*Options)

func Redis(rd *redis.RdClient) Option {
	return func(sch *Scheduler) {
		sch.redis = rd
	}
}

func Repo(rp *repo.Repository) Option {
	return func(sch *Scheduler) {
		sch.repo = rp
	}
}

func Mongo(mg *mongo.MgoClient) Option {
	return func(sch *Scheduler) {
		sch.mgo = mg
	}
}

func Bind(ops ...Option) {
	for _, v := range ops {
		v(sch)
	}

	sch.cron.Start()
}

func Every(duration time.Duration) gron.AtSchedule {
	return gron.Every(duration)
}

func Assign(schedule gron.Schedule, job Job) {
	sch.cron.AddFunc(schedule, wrapper(sch, job))
}
