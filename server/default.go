package server

import (
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"gitee.com/jarnpher_rice/micro-core/log"
	"gitee.com/jarnpher_rice/micro-core/router"
	"gitee.com/jarnpher_rice/micro-core/service"
	"github.com/gin-gonic/gin"
)

// DefaultServer 默认服务器
type DefaultServer struct {
	*http.Server
	*service.Registry
	name    string
	runMode string
}

var entry = log.Logger.Mark("DefaultServer")

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
			entry.Fatalln("wrong format of server name")
		}
	}
}

func RunMode(mode string) Option {
	return func(server *DefaultServer) {
		server.runMode = mode
		gin.SetMode(server.runMode)
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
		name: "micro",
	}

	for _, op := range options {
		op(server)
	}

	r := server.Handler.(*router.Router)
	r.Engine = gin.Default()
	r.RootGroup(server.name)

	return server
}

// Run 实现IBaseServer接口
func (s *DefaultServer) Run() {
	r := s.Server.Handler.(*router.Router)
	for _, s := range r.Services {
		r.Register(s)
	}

	go func() {
		entry.Infof("server listening on %s...", s.Addr)
		if err := s.ListenAndServe(); err != nil {
			entry.Fatalln(err)
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

		node.Address = strings.Split(s.Server.Addr, ":")[0]
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
