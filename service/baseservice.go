package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Jarnpher553/micro-core/breaker"
	"github.com/Jarnpher553/micro-core/httpclient"
	"github.com/Jarnpher553/micro-core/limit"
	"github.com/Jarnpher553/micro-core/metric"
	"github.com/Jarnpher553/micro-core/mongo"
	"github.com/Jarnpher553/micro-core/redis"
	"github.com/Jarnpher553/micro-core/repo"
	"github.com/Jarnpher553/micro-core/tracing"
	"github.com/satori/go.uuid"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type IBaseService interface {
	Get(*Handler) HandlerFunc
	GetList(*Handler) HandlerFunc
	Post(*Handler) HandlerFunc
	Delete(*Handler) HandlerFunc
	DeleteBatch(*Handler) HandlerFunc
	Put(*Handler) HandlerFunc
	Head(*Handler) HandlerFunc
	Patch(*Handler) HandlerFunc
	Options(*Handler) HandlerFunc

	Use(*Handler)

	Repo() *repo.Repository
	SetRepo(*repo.Repository)
	Redis() *redis.RdClient
	SetRedis(*redis.RdClient)
	Mongo() *mongo.MgoClient
	SetMongo(*mongo.MgoClient)
	Client() *httpclient.ReqClient
	SetClient(*httpclient.ReqClient)
	Node() *NodeInfo
	SetNode(*NodeInfo)
	Reg() *Registry
	SetReg(*Registry)
	Option() *Options
	SetOption(*Options)
}

type BaseService struct {
	repository  *repo.Repository
	redisClient *redis.RdClient
	mongoClient *mongo.MgoClient
	client      *httpclient.ReqClient
	node        *NodeInfo
	reg         *Registry

	option *Options
}

type Options struct {
	Metric  *metric.Metric
	Tracer  *tracing.Tracer
	Limiter *limit.Limiter
	Cb      *breaker.CircuitBreaker
}

type NodeInfo struct {
	Id         string
	ServerName string
	Name       string
	Port       string
	Address    string
}

type Option func(service IBaseService)

func Repository(repository *repo.Repository) Option {
	return func(service IBaseService) {
		service.SetRepo(repository)
	}
}

func RedisClient(redisClient *redis.RdClient) Option {
	return func(service IBaseService) {
		service.SetRedis(redisClient)
	}
}

func MongoClient(client *mongo.MgoClient) Option {
	return func(service IBaseService) {
		service.SetMongo(client)
	}
}

func Tracer(tracer *tracing.Tracer) Option {
	return func(service IBaseService) {
		service.Option().Tracer = tracer
	}
}

func Metric(m *metric.Metric) Option {
	return func(service IBaseService) {
		service.Option().Metric = m
	}
}

func Limiter(limiter *limit.Limiter) Option {
	return func(service IBaseService) {
		service.Option().Limiter = limiter
	}
}

func Cb(circuitBreaker *breaker.CircuitBreaker) Option {
	return func(service IBaseService) {
		service.Option().Cb = circuitBreaker
	}
}

func NewService(service IBaseService, option ...Option) IBaseService {
	v := reflect.ValueOf(service)
	t := reflect.TypeOf(service)

	name := strings.TrimSuffix(t.Elem().Name(), "Service")
	name = strings.ToLower(name[:1]) + name[1:]
	bs := &BaseService{
		node: &NodeInfo{
			Id:   uuid.NewV4().String(),
			Name: name,
		},
		option: &Options{},
	}

	for _, op := range option {
		op(bs)
	}

	if bs.option.Tracer == nil {
		bs.option.Tracer = tracing.New(tracing.NewReporter())
	}

	if bs.option.Limiter == nil {
		bs.option.Limiter = limit.New(time.Second*1, 100)

	}

	if bs.option.Metric == nil {
		bs.option.Metric = metric.New(metric.NewWriter(1 * time.Minute))

	}

	if bs.option.Cb == nil {
		bs.option.Cb = breaker.New()
	}

	v.Elem().FieldByName("BaseService").Set(reflect.ValueOf(bs))

	return service
}

func (s *BaseService) Use(handler *Handler) {
	return
}

func (s *BaseService) SetRedis(redisClient *redis.RdClient) {
	s.redisClient = redisClient
}

func (s *BaseService) Redis() *redis.RdClient {
	return s.redisClient
}

func (s *BaseService) SetMongo(client *mongo.MgoClient) {
	s.mongoClient = client
}

func (s *BaseService) Mongo() *mongo.MgoClient {
	return s.mongoClient
}

func (s *BaseService) SetClient(client *httpclient.ReqClient) {
	s.client = client
}

func (s *BaseService) Client() *httpclient.ReqClient {
	if s.client == nil {
		node := s.Node()
		s.client = httpclient.New(httpclient.Tracer(s.Option().Tracer), httpclient.Name(node.ServerName+"."+node.Name))
	}
	return s.client
}

func (s *BaseService) Node() *NodeInfo {
	return s.node
}

func (s *BaseService) SetNode(node *NodeInfo) {
	s.node = node
}

func (s *BaseService) Reg() *Registry {
	return s.reg
}

func (s *BaseService) SetReg(reg *Registry) {
	s.reg = reg
}

func (s *BaseService) Repo() *repo.Repository {
	return s.repository
}

func (s *BaseService) SetRepo(repository *repo.Repository) {
	s.repository = repository
}

func (s *BaseService) Option() *Options {
	return s.option
}

func (s *BaseService) SetOption(op *Options) {
	s.option = op
}

func (s *BaseService) Get(handler *Handler) HandlerFunc {
	return func(ctx *Ctx) {
		ctx.String(http.StatusNotFound, "404 page not found")
	}
}

func (s *BaseService) GetList(handler *Handler) HandlerFunc {
	return func(ctx *Ctx) {
		ctx.String(http.StatusNotFound, "404 page not found")
	}
}

func (s *BaseService) Post(handler *Handler) HandlerFunc {
	return func(ctx *Ctx) {
		ctx.String(http.StatusNotFound, "404 page not found")
	}
}

func (s *BaseService) Delete(handler *Handler) HandlerFunc {
	return func(ctx *Ctx) {
		ctx.String(http.StatusNotFound, "404 page not found")
	}
}

func (s *BaseService) DeleteBatch(handler *Handler) HandlerFunc {
	return func(ctx *Ctx) {
		ctx.String(http.StatusNotFound, "404 page not found")
	}
}

func (s *BaseService) Put(handler *Handler) HandlerFunc {
	return func(ctx *Ctx) {
		ctx.String(http.StatusNotFound, "404 page not found")
	}
}

func (s *BaseService) Head(handler *Handler) HandlerFunc {
	return func(ctx *Ctx) {
		ctx.String(http.StatusNotFound, "404 page not found")
	}
}

func (s *BaseService) Patch(handler *Handler) HandlerFunc {
	return func(ctx *Ctx) {
		ctx.String(http.StatusNotFound, "404 page not found")
	}
}

func (s *BaseService) Options(handler *Handler) HandlerFunc {
	return func(ctx *Ctx) {
		ctx.String(http.StatusNotFound, "404 page not found")
	}
}

func (s *BaseService) DoGet(serviceName string, paramOrAction string, query map[string]string, ctx context.Context, v interface{}) error {
	srv, err := s.reg.GetService(serviceName)
	if err != nil {
		return err
	}
	var url string
	path := strings.Join(strings.Split(srv.Name, "."), "/")
	if paramOrAction == "" {
		url = fmt.Sprintf("http://%s:%s/%s", srv.Address, srv.Port, path)
	} else {
		url = fmt.Sprintf("http://%s:%s/%s/%s", srv.Address, srv.Port, path, paramOrAction)
	}
	return s.Client().RGet(url, query, ctx, v)
}

func (s *BaseService) DoPost(serviceName string, action string, body interface{}, ctx context.Context, v interface{}) error {
	srv, err := s.reg.GetService(serviceName)
	if err != nil {
		return err
	}
	var url string
	path := strings.Join(strings.Split(srv.Name, "."), "/")
	if action == "" {
		url = fmt.Sprintf("http://%s:%s/%s", srv.Address, srv.Port, path)
	} else {
		url = fmt.Sprintf("http://%s:%s/%s/%s", srv.Address, srv.Port, path, action)
	}

	return s.Client().RPost(url, body, ctx, v)
}

func (s *BaseService) DoPut(serviceName string, param int, body interface{}, ctx context.Context, v interface{}) error {
	if param == 0 {
		return errors.New("put request need param")
	}

	srv, err := s.reg.GetService(serviceName)
	if err != nil {
		return err
	}
	path := strings.Join(strings.Split(srv.Name, "."), "/")
	url := fmt.Sprintf("http://%s:%s/%s/%d", srv.Address, srv.Port, path, param)
	return s.Client().RPut(url, body, ctx, v)
}

func (s *BaseService) DoDelete(serviceName string, param int, query map[string]string, ctx context.Context, v interface{}) error {
	if query != nil && param != 0 {
		return errors.New("delete request not need param if query isn't empty")
	} else if query == nil && param == 0 {
		return errors.New("delete request need one of param and query")
	}

	srv, err := s.reg.GetService(serviceName)
	if err != nil {
		return err
	}

	var url string
	path := strings.Join(strings.Split(srv.Name, "."), "/")
	if param == 0 {
		url = fmt.Sprintf("http://%s:%s/%s", srv.Address, srv.Port, path)
	} else {
		url = fmt.Sprintf("http://%s:%s/%s/%d", srv.Address, srv.Port, path, param)
	}
	return s.Client().RDelete(url, query, ctx, v)
}
