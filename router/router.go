package router

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Jarnpher553/micro-core/log"
	"github.com/Jarnpher553/micro-core/service"
	"github.com/Jarnpher553/micro-core/uuid"
	_ "github.com/Jarnpher553/micro-core/validator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router 路由类
type Router struct {
	sync.Once
	sync.Mutex
	*gin.Engine
	Services []service.IBaseService
	static   string
	template string
}

// 初始化 初始化gin输出位置
func init() {
	gin.DefaultWriter = log.Logger.Mark("Gin").Writer()
	gin.DefaultErrorWriter = log.Logger.Mark("Gin").Writer()
}

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

func (r *Router) RootGroup(group string) {
	for i := range r.Services {
		r.Services[i].Node().ServerName = group
	}

	groupList := strings.Split(group, ".")

	var routerGroup = &(r.RouterGroup)
	for _, v := range groupList {
		routerGroup = routerGroup.Group(v)
	}

	r.RouterGroup = *(routerGroup)

	//挂载跨域
	r.Cors()

	//注册静态文件路径
	if r.static != "" {
		r.RegisterStatic(r.static)
	}

	if r.template != "" {
		r.LoadHTMLGlob(r.template)
	}
}

func (r *Router) RegisterStatic(path string) {
	_ = os.MkdirAll(path, os.ModePerm)
	r.Static("/static", path)
}

func (r *Router) Cors() {
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "access-token"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		ExposeHeaders:    []string{"Content-Disposition"},
		MaxAge:           12 * time.Hour,
	}))
}

func (r *Router) InjectSlice(services ...service.IBaseService) {
	for _, v := range services {
		r.Inject(v)
	}
}

func (r *Router) Inject(service service.IBaseService) {
	r.Lock()
	defer r.Unlock()
	r.Services = append(r.Services, service)
}

// Register 自定义注册
func (r *Router) Register(srv service.IBaseService) {
	serviceType := reflect.TypeOf(srv)

	serviceVal := reflect.ValueOf(srv)

	node := srv.Node()
	name := node.ServerName + "." + node.Name

	//构建服务对应路由
	group := r.Group(fmt.Sprintf("%s", node.Name))

	//服务注册中间件
	group.Use(service.Wrapper(service.RateLimiterMiddleware(srv.Option().Limiter)(srv)))
	group.Use(service.Wrapper(service.BreakerMiddleware(srv.Option().Cb)(srv)))
	group.Use(service.Wrapper(service.MetricMiddleware(srv.Option().Metric)(srv)))
	group.Use(service.Wrapper(service.ExtractHttpMiddleware()(srv)))
	group.Use(service.Wrapper(service.TracerMiddleware(srv.Option().Tracer, name)(srv)))

	//服务注册路由
	group.Handle("POST", "", service.Wrapper(func(context *service.Ctx) {
		var handler service.Handler
		srv.Use(&handler)

		post := srv.Post(&handler)

		for _, m := range handler.Middleware {
			m(srv)(context)

			if context.IsAborted() {
				return
			}
		}

		post(context)
	}))

	group.Handle("POST", "/:action", service.Wrapper(func(context *service.Ctx) {
		action := strings.Title(context.Param("action"))
		method, exist := serviceType.MethodByName("Post" + action)
		if !exist {
			context.String(http.StatusNotFound, "404 page not found")
			return
		}

		var handler service.Handler
		srv.Use(&handler)
		ret := method.Func.Call([]reflect.Value{serviceVal, reflect.ValueOf(&handler)})

		for _, m := range handler.Middleware {
			m(srv)(context)

			if context.IsAborted() {
				return
			}
		}

		ret[0].Interface().(service.HandlerFunc)(context)
	}))

	group.Handle("DELETE", "", service.Wrapper(func(context *service.Ctx) {
		var handler service.Handler
		srv.Use(&handler)

		deleteBatch := srv.DeleteBatch(&handler)

		for _, m := range handler.Middleware {
			m(srv)(context)

			if context.IsAborted() {
				return
			}
		}

		deleteBatch(context)
	}))

	group.Handle("DELETE", "/:id", service.Wrapper(func(context *service.Ctx) {
		idStr := context.Param("id")

		_, err1 := strconv.Atoi(idStr)
		err2 := uuid.IsGUID(idStr)
		if err1 != nil && err2 != nil {
			context.String(http.StatusNotFound, "404 page not found")
			return
		} else {
			context.Params = gin.Params{
				gin.Param{
					Key:   "id",
					Value: idStr,
				},
			}

			var handler service.Handler
			srv.Use(&handler)

			del := srv.Delete(&handler)

			for _, m := range handler.Middleware {
				m(srv)(context)

				if context.IsAborted() {
					return
				}
			}

			del(context)
		}
	}))

	group.Handle("PUT", "/:id", service.Wrapper(func(context *service.Ctx) {
		idStr := context.Param("id")

		_, err1 := strconv.Atoi(idStr)
		err2 := uuid.IsGUID(idStr)
		if err1 != nil && err2 != nil {
			context.String(http.StatusNotFound, "404 page not found")
			return
		} else {
			context.Params = gin.Params{
				gin.Param{
					Key:   "id",
					Value: idStr,
				},
			}

			var handler service.Handler
			srv.Use(&handler)

			put := srv.Put(&handler)

			for _, m := range handler.Middleware {
				m(srv)(context)

				if context.IsAborted() {
					return
				}
			}

			put(context)
		}
	}))

	group.Handle("GET", "", service.Wrapper(func(context *service.Ctx) {
		var handler service.Handler
		srv.Use(&handler)

		getList := srv.GetList(&handler)

		for _, m := range handler.Middleware {
			m(srv)(context)

			if context.IsAborted() {
				return
			}
		}

		getList(context)
	}))

	group.Handle("GET", "/:actionOrID", service.Wrapper(func(context *service.Ctx) {
		actionOrID := context.Param("actionOrID")

		_, err1 := strconv.Atoi(actionOrID)
		err2 := uuid.IsGUID(actionOrID)
		if err1 != nil && err2 != nil {
			action := strings.Title(actionOrID)
			method, exist := serviceType.MethodByName("Get" + action)
			if !exist {
				context.String(http.StatusNotFound, "404 page not found")
				return
			}

			var handler service.Handler
			srv.Use(&handler)

			ret := method.Func.Call([]reflect.Value{serviceVal, reflect.ValueOf(&handler)})

			for _, m := range handler.Middleware {
				m(srv)(context)

				if context.IsAborted() {
					return
				}
			}

			ret[0].Interface().(service.HandlerFunc)(context)
		} else {
			context.Params = gin.Params{
				gin.Param{
					Key:   "id",
					Value: actionOrID,
				},
			}

			var handler service.Handler
			srv.Use(&handler)

			get := srv.Get(&handler)

			for _, m := range handler.Middleware {
				m(srv)(context)

				if context.IsAborted() {
					return
				}
			}

			get(context)
		}
	}))
}
