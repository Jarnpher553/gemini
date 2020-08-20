package mongo

import (
	"context"
	"fmt"
	"github.com/Jarnpher553/gemini/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MgoClient struct {
	*mongo.Client
	addr       string
	database   string
	collection string
}

type Option func(client *MgoClient)

func Addr(addr string) Option {
	return func(client *MgoClient) {
		client.addr = addr
	}
}

func Database(db string) Option {
	return func(client *MgoClient) {
		client.database = db
	}
}

func Collection(cl string) Option {
	return func(client *MgoClient) {
		client.collection = cl
	}
}

var entry = log.Zap.Mark("Mongo")

func New(opts ...Option) *MgoClient {
	mgo := &MgoClient{}

	for _, opt := range opts {
		opt(mgo)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", mgo.addr)))
	if err != nil {
		entry.Fatal(log.Message(err))
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		entry.Fatal(log.Message(err))
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		entry.Fatal(log.Message(err))
	}

	mgo.Client = client
	return mgo
}

func (c *MgoClient) DbCollection() *mongo.Collection {
	return c.Database(c.database).Collection(c.collection)
}
