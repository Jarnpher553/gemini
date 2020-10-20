package scheduler

import (
	"github.com/Jarnpher553/gemini/mongo"
	"github.com/Jarnpher553/gemini/redis"
	"github.com/Jarnpher553/gemini/repo"
	"github.com/robfig/cron/v3"
	"time"
)

var ct = &CronTab{cron: cron.New(), conf: &Configuration{}}

type CronTab struct {
	cron *cron.Cron
	conf *Configuration
}

type Configuration struct {
	redis  *redis.RdClient
	repo   *repo.Repository
	mgo    *mongo.MgoClient
	Custom interface{}
}

type Conf func(*CronTab)

type Job func(*Configuration)

func Redis(rd *redis.RdClient) Conf {
	return func(ct *CronTab) {
		ct.conf.redis = rd
	}
}

func Repo(rp *repo.Repository) Conf {
	return func(ct *CronTab) {
		ct.conf.repo = rp
	}
}

func Mongo(mg *mongo.MgoClient) Conf {
	return func(ct *CronTab) {
		ct.conf.mgo = mg
	}
}

func Custom(custom interface{}) Conf {
	return func(ct *CronTab) {
		ct.conf.Custom = custom
	}
}

func Bind(conf ...Conf) {
	for _, v := range conf {
		v(ct)
	}
	ct.cron.Start()
}

func Assign(schedule interface{}, job Job) {
	switch t := schedule.(type) {
	case string:
		ct.cron.AddFunc(t, decorator(ct.conf, job))
	case cron.Schedule:
		ct.cron.Schedule(t, cron.FuncJob(decorator(ct.conf, job)))
	}
}

func Every(duration time.Duration) cron.Schedule {
	return cron.Every(duration)
}

func Stop() {
	ct.cron.Stop()
}
