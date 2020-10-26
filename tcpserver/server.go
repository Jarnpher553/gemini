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
	name   string
	logger *log.ZapLogger
	addr   string
	opt    gnet.Options
	eh     EventService
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
		logger: &log.ZapLogger{Logger: log.Logger.Mark("tcpserver").With(zap.String("name", name))},
		name:   name,
		opt:    gnet.Options{},
	}

	for _, option := range opts {
		option(s)
	}

	s.opt.Logger = s.logger.Sugar()

	return s
}

func (s *TcpServer) Serve(handler EventService) {
	v := reflect.ValueOf(handler)
	if v.Type().Kind() == reflect.Struct {
		v.FieldByName("logger").Set(reflect.ValueOf(s.logger))
	} else {
		v.Elem().FieldByName("logger").Set(reflect.ValueOf(s.logger))
	}

	s.eh = handler
}

func (s *TcpServer) Run() {
	err := gnet.Serve(s.eh, fmt.Sprintf("tcp://%s", s.addr), gnet.WithOptions(s.opt))
	if err != nil {
		s.logger.Fatal(err.Error())
	}
}