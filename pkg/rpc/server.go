package rpc

import (
	"github.com/Jarnpher553/gemini/pkg/log"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var logger = log.Logger.Mark("grpc")

type GRpcServer struct {
	addr   string
	logger *log.ZapLogger
	grpc   *grpc.Server
}

func New(addr string) *GRpcServer {
	return &GRpcServer{
		addr:   addr,
		logger: nil,
		grpc:   grpc.NewServer(),
	}
}

func (s *GRpcServer) RegisterService(f func(*grpc.Server)) {
	f(s.grpc)
}

func (s *GRpcServer) Run() {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logger.Fatal(err.Error())
	}

	if err := s.grpc.Serve(lis); err != nil {
		s.logger.Fatal(err.Error())
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	s.logger.Info("aaaa")
	s.grpc.GracefulStop()
}
