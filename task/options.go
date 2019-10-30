package task

import (
	"github.com/Jarnpher553/micro-core/mongo"
	"github.com/Jarnpher553/micro-core/redis"
	"github.com/Jarnpher553/micro-core/repo"
)

type Options struct {
	redis *redis.RdClient
	repo  *repo.Repository
	mgo   *mongo.MgoClient
}
