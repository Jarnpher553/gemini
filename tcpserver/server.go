package tcpserver

import (
	"fmt"
	"github.com/Jarnpher553/gemini/log"
	"github.com/Jarnpher553/gemini/util/random"
	"github.com/panjf2000/gnet"
	"go.uber.org/zap"
	"time"
)

type LB = gnet.LoadBalancing

type Conn = gnet.Conn

type Action = gnet.Action

type Server = gnet.Server

type EventService interface {
	OnInitComplete(server Server) (action Action)

	OnShutdown(server Server)

	OnOpened(c Conn) (out []byte, action Action)

	OnClosed(c Conn, err error) (action Action)

	PreWrite()

	React(frame []byte, c Conn) (out []byte, action Action)

	Tick() (delay time.Duration, action Action)
}

type Service struct {
	*gnet.EventServer
	logger *log.ZapLogger
}

func NewService(logger *log.ZapLogger) *Service {
	return &Service{
		EventServer: &gnet.EventServer{},
	}
}

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
	handler.(*Service).logger = s.logger
	s.eh = handler
}

func (s *TcpServer) Run() {
	err := gnet.Serve(s.eh, fmt.Sprintf("tcp://%s", s.addr), gnet.WithOptions(s.opt))
	if err != nil {
		s.logger.Fatal(err.Error())
	}
}
