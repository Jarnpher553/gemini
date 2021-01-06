package service

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/Jarnpher553/gemini/pkg/breaker"
	"github.com/Jarnpher553/gemini/pkg/httpclient"
	"github.com/Jarnpher553/gemini/pkg/limit"
	"github.com/Jarnpher553/gemini/pkg/metric"
	"github.com/Jarnpher553/gemini/pkg/mongo"
	"github.com/Jarnpher553/gemini/pkg/redis"
	"github.com/Jarnpher553/gemini/pkg/repo"
	"github.com/Jarnpher553/gemini/pkg/tracing"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/time/rate"
	"runtime"
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

	Repo() repo.Repository
	SetRepo(repo.Repository)
	Redis() *redis.RdClient
	SetRedis(*redis.RdClient)
	Mongo() *mongo.MgoClient
	SetMongo(*mongo.MgoClient)
	Client() *httpclient.ReqClient
	SetClient(*httpclient.ReqClient)
	Node() *NodeInfo
	SetNode(*NodeInfo)
	Interceptor() *Interceptor
	SetInterceptor(*Interceptor)

	CustomContext(string) interface{}
	SetCustomContext(string, interface{})
}

type BaseService struct {
	repository  repo.Repository
	redisClient *redis.RdClient
	mongoClient *mongo.MgoClient
	client      *httpclient.ReqClient
	node        *NodeInfo

	interceptor *Interceptor

	customContext map[string]interface{}
}

type Interceptor struct {
	Metric  *metric.Metric
	Tracer  *tracing.Tracer
	Limiter *limit.Limiter
	Cb      *breaker.CircuitBreaker
}

type NodeInfo struct {
	Id       string
	RootName string
	Name     string
	Port     string
	Address  string
}

type Option func(service IBaseService)

func Repository(repository repo.Repository) Option {
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
		service.Interceptor().Tracer = tracer
	}
}

func Metric(m *metric.Metric) Option {
	return func(service IBaseService) {
		service.Interceptor().Metric = m
	}
}

func Limiter(limiter *limit.Limiter) Option {
	return func(service IBaseService) {
		service.Interceptor().Limiter = limiter
	}
}

func Cb(circuitBreaker *breaker.CircuitBreaker) Option {
	return func(service IBaseService) {
		service.Interceptor().Cb = circuitBreaker
	}
}

func CustomContext(key string, value interface{}) Option {
	return func(service IBaseService) {
		service.SetCustomContext(key, value)
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
		interceptor: &Interceptor{},
	}

	for _, op := range option {
		op(bs)
	}

	if bs.interceptor.Tracer == nil {
		bs.interceptor.Tracer = tracing.New(tracing.NewZapReporter())
	}

	if bs.interceptor.Limiter == nil {
		bs.interceptor.Limiter = limit.New(rate.Limit(200*runtime.NumCPU()), 200*runtime.NumCPU())
	}

	if bs.interceptor.Metric == nil {
		bs.interceptor.Metric = metric.New(&metric.Config{ServiceName: name, Printer: metric.NewPrinter(), Freq: 1 * time.Minute})
	}

	if bs.interceptor.Cb == nil {
		bs.interceptor.Cb = breaker.New()
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
		s.client = httpclient.New(httpclient.Tracer(s.Interceptor().Tracer), httpclient.Name(node.RootName+"."+node.Name))
	}
	return s.client
}

func (s *BaseService) Node() *NodeInfo {
	return s.node
}

func (s *BaseService) SetNode(node *NodeInfo) {
	s.node = node
}

func (s *BaseService) Repo() repo.Repository {
	return s.repository
}

func (s *BaseService) SetRepo(repository repo.Repository) {
	s.repository = repository
}

func (s *BaseService) Interceptor() *Interceptor {
	return s.interceptor
}

func (s *BaseService) SetInterceptor(op *Interceptor) {
	s.interceptor = op
}

func (s *BaseService) CustomContext(key string) interface{} {
	return s.customContext[key]
}

func (s *BaseService) SetCustomContext(key string, value interface{}) {
	if s.customContext == nil {
		s.customContext = make(map[string]interface{})
	}
	s.customContext[key] = value
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

// Call 调用其它服务方法
// params:
// 	@method: 请求方法 get post put delete
//	@url: 请求地址
// 	@args: 调用方法的入参
// 	@replay: 调用发放的返回
func (s *BaseService) Call(ctx context.Context, method string, url string, args interface{}, replay interface{}) error {
	m := strings.ToLower(method)
	if strings.Index(m, "post") != -1 {
		return s.callPost(url, args, ctx, replay)
	} else if strings.Index(m, "get") != -1 {
		return s.callGet(url, args.(map[string]string), ctx, replay)
	} else if strings.Index(m, "put") != -1 {
		return s.callPut(url, args, ctx, replay)
	} else if strings.Index(m, "delete") != -1 {
		return s.callDelete(url, args.(map[string]string), ctx, replay)
	}
	return nil
}

func (s *BaseService) callGet(url string, query map[string]string, ctx context.Context, v interface{}) error {
	return s.Client().RGet(url, query, ctx, v)
}

func (s *BaseService) callPost(url string, body interface{}, ctx context.Context, v interface{}) error {
	return s.Client().RPost(url, body, ctx, v)
}

func (s *BaseService) callPut(url string, body interface{}, ctx context.Context, v interface{}) error {
	return s.Client().RPut(url, body, ctx, v)
}

func (s *BaseService) callDelete(url string, query map[string]string, ctx context.Context, v interface{}) error {
	return s.Client().RDelete(url, query, ctx, v)
}
