package server

import (
	"fmt"
	"github.com/Jarnpher553/gemini/util/addr"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Jarnpher553/gemini/log"
	"github.com/Jarnpher553/gemini/router"
	"github.com/Jarnpher553/gemini/service"
)

// DefaultServer 默认服务器
type DefaultServer struct {
	*http.Server
	*service.Registry
	name    string
	runMode string
	env     string
	logger  *log.ZapLogger
}

type Option func(server *DefaultServer)

func Addr(addr string) Option {
	return func(server *DefaultServer) {
		server.Addr = addr
	}
}

func Router(r *router.Router) Option {
	return func(server *DefaultServer) {
		server.Handler = r
	}
}

func Registry(reg *service.Registry) Option {
	return func(server *DefaultServer) {
		server.Registry = reg
	}
}

func Name(name string) Option {
	return func(server *DefaultServer) {
		server.name = strings.ToLower(name)

		reg := regexp.MustCompile(`^[a-z]+(\.[a-z]+)*$`)
		if !reg.MatchString(server.name) {
			server.logger.Fatal(log.Message("wrong format of server name"))
		}
	}
}

func RunMode(mode string) Option {
	return func(server *DefaultServer) {
		server.runMode = mode
	}
}

func Env(env string) Option {
	return func(server *DefaultServer) {
		server.env = env
	}
}

// Default 构造函数
func Default(options ...Option) IBaseServer {
	server := &DefaultServer{
		Server: &http.Server{
			//ReadTimeout:    10 * time.Second,
			//WriteTimeout:   10 * time.Second,
			//MaxHeaderBytes: 1 << 20,
		},
		name:    "micro",
		logger:  log.Zap.Mark("server"),
		runMode: gin.ReleaseMode,
	}

	for _, op := range options {
		op(server)
	}

	r, ok := server.Handler.(*router.Router)
	if ok {
		server.printBanner()
		r.Startup(&router.Config{
			ServerName: server.name,
			RunMode:    server.runMode,
		})
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

// Run 实现IBaseServer接口
func (s *DefaultServer) Run() {
	defer s.logger.Sync()

	go func() {
		s.logger.Info(log.Message("start server"), []zapcore.Field{zap.String("name", s.name), zap.String("env", s.env), zap.String("addr", s.Server.Addr), zap.String("scheme", "http")}...)

		if err := s.ListenAndServe(); err != nil {
			s.logger.Fatal(log.Message(err))
		}
	}()

	if s.Registry != nil {
		<-time.After(1 * time.Second)
		_ = s.register()
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	if s.Registry != nil {
		_ = s.deregister()
	}
}

func (s *DefaultServer) register() []error {
	errs := make([]error, 0)
	errChan := make(chan error, len(s.Services))
	var wg sync.WaitGroup

	for _, node := range s.Services {
		wg.Add(1)

		address, _ := addr.Extract(s.Server.Addr)

		node.Address = address
		node.Port = strings.Split(s.Server.Addr, ":")[1]

		go s.Register(node, &wg, errChan)
	}

	wg.Wait()
	close(errChan)

	for {
		if err, ok := <-errChan; ok {
			errs = append(errs, err)
		} else {
			break
		}
	}
	return errs
}

func (s *DefaultServer) deregister() []error {
	errs := make([]error, 0)
	errChan := make(chan error, len(s.Services))
	var wg sync.WaitGroup

	for _, v := range s.Services {
		wg.Add(1)
		go s.Deregister(v, &wg, errChan)
	}

	wg.Wait()
	close(errChan)

	for {
		if err, ok := <-errChan; ok {
			errs = append(errs, err)
		} else {
			break
		}
	}
	return errs
}
