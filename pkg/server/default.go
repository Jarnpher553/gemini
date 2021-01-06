package server

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jarnpher553/gemini/pkg/log"
	"github.com/Jarnpher553/gemini/pkg/router"
)

// DefaultServer 默认服务器
type DefaultServer struct {
	*http.Server
	env     string
	logger  *log.ZapLogger
	startup func() error
	release func() error
}

type Option func(server *DefaultServer)

func Addr(addr string) Option {
	return func(server *DefaultServer) {
		server.Addr = addr
	}
}

func Env(env string) Option {
	return func(server *DefaultServer) {
		server.env = env
	}
}

func Startup(startup func() error) Option {
	return func(server *DefaultServer) {
		server.startup = startup
	}
}

func Release(release func() error) Option {
	return func(server *DefaultServer) {
		server.release = release
	}
}

func Router(r *router.Router) Option {
	return func(server *DefaultServer) {
		server.Handler = r
	}
}

func Route(route func() *router.Router) Option {
	return func(server *DefaultServer) {
		server.Handler = route()
	}
}

// Default 构造函数
func Default(options ...Option) IBaseServer {
	server := &DefaultServer{
		Server: &http.Server{
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		logger: log.Zap.Mark("server"),
	}

	for _, op := range options {
		op(server)
	}

	server.printBanner()

	if server.Handler == nil {
		server.logger.Fatal("the router of server hasn't been initialized")
	}

	server.printRoutes()

	if server.startup != nil {
		if err := server.startup(); err != nil {
			server.logger.Fatal(err.Error())
		}
	}

	return server
}

func (s *DefaultServer) printBanner() {
	const banner = `
      _____     
    /  ___  \    ________    _________    __    _____    __
   | |____|  |  |  ____  |  |  _   _  |  |__|  |  _  |  |__|
    \_____   |  | |____| |  | | | | | |   __   | | | |   __
    _____/   |  |  ______|  |_| |_| |_|  |  |  |_| |_|  |  |
   \ ______ /   | |_____                 |__|           |__|
                |________\

    Welcome to gemini, starting application ...
`
	fmt.Println(fmt.Sprintf("\x1b[32m%s\x1b[0m", banner))
}

func (s *DefaultServer) printRoutes() {
	routes := s.Handler.(*router.Router).Routes()
	for _, route := range routes {
		s.logger.Info("add route", zap.String("method", route.Method), zap.String("path", route.Path))
	}
}

// Run 实现IBaseServer接口
func (s *DefaultServer) Run() {
	defer s.logger.Sync()

	go func() {
		s.logger.Info(log.Message("start server"), []zapcore.Field{zap.String("env", s.env), zap.String("addr", s.Server.Addr), zap.String("scheme", "http")}...)

		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal(log.Message(err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	if s.release != nil {
		if err := s.release(); err != nil {
			s.logger.Fatal(err.Error())
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		s.logger.With(zap.String("err", err.Error())).Fatal("server forced to shutdown")
	}
	s.logger.Info("server exiting")
}
