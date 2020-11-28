package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/Jarnpher553/gemini/breaker"
	"github.com/Jarnpher553/gemini/httpclient"
	"github.com/Jarnpher553/gemini/limit"
	"github.com/Jarnpher553/gemini/metric"
	"github.com/Jarnpher553/gemini/mongo"
	"github.com/Jarnpher553/gemini/redis"
	"github.com/Jarnpher553/gemini/repo"
	"github.com/Jarnpher553/gemini/tracing"
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
	Interceptor() *Interceptor
	SetInterceptor(*Interceptor)

	CustomContext(string) interface{}
	SetCustomContext(string, interface{})

	Area() string
}

type BaseService struct {
	repository  *repo.Repository
	redisClient *redis.RdClient
	mongoClient *mongo.MgoClient
	client      *httpclient.ReqClient
	node        *NodeInfo
	reg         *Registry

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
		s.client = httpclient.New(httpclient.Tracer(s.Interceptor().Tracer), httpclient.Name(node.ServerName+"."+node.Name))
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

func (s *BaseService) Area() string {
	return ""
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
// 	@method: get请求: get get:id get@action
// 			post请求: post post@action
// 			put请求: put:id
// 			delete请求: delete delete:id
// 	@args: 调用方法的入参
// 	@replay: 调用发放的返回
func (s *BaseService) Call(ctx context.Context, service string, method string, args interface{}, replay interface{}) error {
	if strings.Index(method, "post") != -1 {
		action := strings.TrimPrefix(method, "post")

		if strings.Index(action, "@") != -1 {
			return s.callPost(service, strings.Split(action, "@")[1], args, ctx, replay)
		} else {
			return s.callPost(service, "", args, ctx, replay)
		}

	} else if strings.Index(method, "get") != -1 {
		action := strings.TrimPrefix(method, "get")

		if strings.Index(action, ":") != -1 {
			return s.callGet(service, strings.Split(action, ":")[1], nil, ctx, replay)
		} else if strings.Index(action, "@") != -1 {
			var query map[string]string
			if args != nil {
				query = args.(map[string]string)
			}

			return s.callGet(service, strings.Split(action, "@")[1], query, ctx, replay)
		} else {
			var query map[string]string
			if args != nil {
				query = args.(map[string]string)
			}

			return s.callGet(service, "", query, ctx, replay)
		}
	} else if strings.Index(method, "put") != -1 {
		action := strings.TrimPrefix(method, "put")

		return s.callPut(service, strings.Split(action, ":")[1], nil, ctx, replay)
	} else if strings.Index(method, "delete") != -1 {
		action := strings.TrimPrefix(method, "put")

		if strings.Index(action, ":") != -1 {
			return s.callGet(service, strings.Split(action, ":")[1], nil, ctx, replay)
		} else {
			var query map[string]string
			if args != nil {
				query = args.(map[string]string)
			}
			return s.callDelete(service, "", query, ctx, replay)
		}
	}
	return nil
}

func (s *BaseService) callGet(serviceName string, paramOrAction string, query map[string]string, ctx context.Context, v interface{}) error {
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

func (s *BaseService) callPost(serviceName string, action string, body interface{}, ctx context.Context, v interface{}) error {
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

func (s *BaseService) callPut(serviceName string, param string, body interface{}, ctx context.Context, v interface{}) error {
	if param == "" {
		return errors.New("put request need param")
	}

	srv, err := s.reg.GetService(serviceName)
	if err != nil {
		return err
	}
	path := strings.Join(strings.Split(srv.Name, "."), "/")
	url := fmt.Sprintf("http://%s:%s/%s/%s", srv.Address, srv.Port, path, param)
	return s.Client().RPut(url, body, ctx, v)
}

func (s *BaseService) callDelete(serviceName string, param string, query map[string]string, ctx context.Context, v interface{}) error {
	if query != nil && param != "" {
		return errors.New("delete request not need param if query isn't empty")
	} else if query == nil && param == "" {
		return errors.New("delete request need one of param and query")
	}

	srv, err := s.reg.GetService(serviceName)
	if err != nil {
		return err
	}

	var url string
	path := strings.Join(strings.Split(srv.Name, "."), "/")
	if param == "" {
		url = fmt.Sprintf("http://%s:%s/%s", srv.Address, srv.Port, path)
	} else {
		url = fmt.Sprintf("http://%s:%s/%s/%s", srv.Address, srv.Port, path, param)
	}
	return s.Client().RDelete(url, query, ctx, v)
}
