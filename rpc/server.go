package rpc

import (
	"github.com/Jarnpher553/gemini/log"
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
	server *grpc.Server
}

func New(addr string) *GRpcServer {
	return &GRpcServer{
		addr:   "",
		logger: nil,
		server: grpc.NewServer(),
	}
}

func (s *GRpcServer) Run() {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logger.Fatal(err.Error())
	}

	if err := s.server.Serve(lis); err != nil {
		s.logger.Fatal(err.Error())
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	s.server.GracefulStop()
}
