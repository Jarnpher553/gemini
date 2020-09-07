package router

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Jarnpher553/gemini/log"
	"github.com/Jarnpher553/gemini/service"
	_ "github.com/Jarnpher553/gemini/validator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router 路由类
type Router struct {
	sync.Once
	sync.Mutex
	*gin.Engine
	services []service.IBaseService
	static   string
	template string
}

var zapLogger = log.Zap.Mark("gin")

type Option func(router *Router)

// New 构造函数
//		permission 角色鉴权中间件
func New(opts ...Option) *Router {
	r := &Router{}

	for _, opt := range opts {
		opt(r)
	}
	return r
}

func HTMLGlod(pattern string) Option {
	return func(router *Router) {
		router.template = pattern
	}
}

func StaticFs(path string) Option {
	return func(router *Router) {
		router.static = path
	}
}

type Config struct {
	ServerName string
	RunMode    string
}

func (r *Router) Startup(config *Config) {
	gin.SetMode(config.RunMode)

	r.rootGroup(config.ServerName)
	r.register()
	r.printRoutes()
}

func (r *Router) printRoutes() {
	for _, route := range r.Engine.Routes() {
		zapLogger.Info("add route", zap.String("method", route.Method), zap.String("path", route.Path))
	}
}

func ginEngine() *gin.Engine {
	engine := gin.New()
	engine.Use(recoverMiddleware(500))
	return engine
}

func (r *Router) rootGroup(group string) {
	r.Engine = ginEngine()

	for i := range r.services {
		r.services[i].Node().ServerName = group
	}

	groupList := strings.Split(group, ".")

	var routerGroup = &(r.RouterGroup)
	for _, v := range groupList {
		routerGroup = routerGroup.Group(v)
	}

	r.RouterGroup = *(routerGroup)

	//挂载跨域
	r.cors()

	//注册静态文件路径
	if r.static != "" {
		r.registerStatic(r.static)
	}

	if r.template != "" {
		r.LoadHTMLGlob(r.template)
	}
}

func (r *Router) registerStatic(path string) {
	_ = os.MkdirAll(path, os.ModePerm)
	r.Static("/static", path)
}

func (r *Router) cors() {
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "access-token"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		ExposeHeaders:    []string{"Content-Disposition"},
		MaxAge:           12 * time.Hour,
	}))
}

func (r *Router) Assign(service ...service.IBaseService) {
	r.Lock()
	defer r.Unlock()
	r.services = append(r.services, service...)
}

func (r *Router) register() {
	for _, s := range r.services {
		r.doRegister(s)
	}
}

// Register 自定义注册
func (r *Router) doRegister(srv service.IBaseService) {
	serviceType := reflect.TypeOf(srv)

	serviceVal := reflect.ValueOf(srv)

	node := srv.Node()
	name := node.ServerName + "." + node.Name

	//获取服务全局中间件
	var handler service.Handler
	srv.Use(&handler)
	var middleware []gin.HandlerFunc
	for _, h := range handler.GinMiddleware {
		middleware = append(middleware, h)
	}
	for _, m := range handler.Middleware {
		middleware = append(middleware, service.Wrapper(m(srv)))
	}
	//构建服务对应路由
	group := r.Group(fmt.Sprintf("%s", node.Name))

	//服务注册中间件
	group.Use(service.Wrapper(service.ReserveLimiterMiddleware(srv.Interceptor().Limiter)(srv)))
	group.Use(service.Wrapper(service.BreakerMiddleware(srv.Interceptor().Cb)(srv)))
	group.Use(service.Wrapper(service.MetricMiddleware(srv.Interceptor().Metric)(srv)))
	//group.Use(service.Wrapper(service.ExtractHttpMiddleware()(srv)))
	group.Use(service.Wrapper(service.TracerMiddleware(srv.Interceptor().Tracer, name)(srv)))

	//注册自定义中间件
	group.Use(middleware...)

	//获取服务类型所有方法
	numMethod := serviceType.NumMethod()
	for i := 0; i < numMethod; i++ {
		var handler service.Handler
		method := serviceType.Method(i)
		methodName := method.Name
		_func := method.Func

		//入参不满足
		if _func.Type().NumIn() != 2 || _func.Type().In(1) != reflect.TypeOf(&service.Handler{}) {
			continue
		}
		//出参不满足
		if _func.Type().NumOut() != 1 || _func.Type().Out(0) != reflect.TypeOf(service.HandlerFunc(func(ctx *service.Ctx) {})) {
			continue
		}

		re := regexp.MustCompile(`(?i:(post|get|delete|put|head|patch|options|)(.*))`)
		matches := re.FindAllStringSubmatch(methodName, -1)
		if matches == nil {
			continue
		}
		ret := _func.Call([]reflect.Value{serviceVal, reflect.ValueOf(&handler)})

		var httpMethod string
		if handler.HttpMethod == "" {
			hm := strings.ToTitle(matches[0][1])
			if hm != "" {
				httpMethod = strings.ToTitle(hm)
			}
		} else {
			httpMethod = handler.HttpMethod
		}
		if httpMethod == "" {
			httpMethod = "GET"
		}
		var relativePath string
		if handler.RelativePath == "" {
			path := matches[0][2]
			if path != "" {
				relativePath = strings.ToLower(path[0:1]) + path[1:]
			}
		} else {
			relativePath = handler.RelativePath
		}
		var middleware []gin.HandlerFunc
		for _, h := range handler.GinMiddleware {
			middleware = append(middleware, h)
		}
		for _, m := range handler.Middleware {
			middleware = append(middleware, service.Wrapper(m(srv)))
		}
		middleware = append(middleware, service.Wrapper(ret[0].Interface().(service.HandlerFunc)))
		group.Handle(httpMethod, relativePath, middleware...)
	}
}
