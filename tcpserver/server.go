package tcpserver

import (
	"fmt"
	"github.com/Jarnpher553/gemini/log"
	"github.com/Jarnpher553/gemini/util/random"
	"github.com/panjf2000/gnet"
	"go.uber.org/zap"
	"reflect"
	"time"
)

type LB = gnet.LoadBalancing

type Conn = gnet.Conn

type Action = gnet.Action

type Server = gnet.Server

type TcpServer struct {
	name    string
	logger  *log.ZapLogger
	addr    string
	opt     gnet.Options
	eh      EventService
	release []func() error
}

type Option func(*TcpServer)

func Addr(addr string) Option {
	return func(server *TcpServer) {
		server.addr = addr
	}
}

func Multicore(multicore bool) Option {
	return func(server *TcpServer) {
		server.opt.Multicore = multicore
	}
}

func Ticker(ticker bool) Option {
	return func(server *TcpServer) {
		server.opt.Ticker = ticker
	}
}

func Codec(codec ICodec) Option {
	return func(server *TcpServer) {
		server.opt.Codec = codec
	}
}

func TCPKeepAlive(tcpKeepAlive time.Duration) Option {
	return func(server *TcpServer) {
		server.opt.TCPKeepAlive = tcpKeepAlive
	}
}

func ReusePort(reusePort bool) Option {
	return func(server *TcpServer) {
		server.opt.ReusePort = reusePort
	}
}

func NumEventLoop(numEventLoop int) Option {
	return func(server *TcpServer) {
		server.opt.NumEventLoop = numEventLoop
	}
}

func LoadBalancing(lb LB) Option {
	return func(server *TcpServer) {
		server.opt.LB = gnet.LoadBalancing(lb)
	}
}

func LockOSThread(lockOSThread bool) Option {
	return func(server *TcpServer) {
		server.opt.LockOSThread = lockOSThread
	}
}

func New(opts ...Option) *TcpServer {
	name := "tcpserver-" + random.RandomString(6)

	s := &TcpServer{
		logger:  &log.ZapLogger{Logger: log.Logger.Mark("tcpserver").With(zap.String("name", name))},
		name:    name,
		opt:     gnet.Options{},
		release: make([]func() error, 0),
	}

	for _, option := range opts {
		option(s)
	}

	s.opt.Logger = s.logger.Sugar()

	return s
}

func (s *TcpServer) Serve(handler EventService, withPool bool) {
	service := &Service{
		EventServer: &gnet.EventServer{},
		logger:      s.logger,
	}
	if withPool {
		service.Pool = Pool()
		s.release = append(s.release, func() error {
			service.Pool.Release()
			return nil
		})
	}

	v := reflect.ValueOf(handler)
	switch v.Kind() {
	default:
		s.logger.Fatal("EventService is not a pointor")
	case reflect.Ptr:
		field := v.Elem().FieldByName("Service")
		if !field.IsValid() {
			s.logger.Fatal("EventService does not contains Service field")
		}
		if field.Kind() != reflect.Ptr {
			s.logger.Fatal(" EventService's Service field is not *tcpserver.Service")
		}
		v.Elem().FieldByName("Service").Set(reflect.ValueOf(service))
	}

	s.eh = handler
}

func (s *TcpServer) Run() {
	err := gnet.Serve(s.eh, fmt.Sprintf("tcp://%s", s.addr), gnet.WithOptions(s.opt))
	if err != nil {
		s.logger.Fatal(err.Error())
	}
	for _, r := range s.release {
		_ = r()
	}
}
