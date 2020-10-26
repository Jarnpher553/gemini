package tcpserver

import (
	"github.com/Jarnpher553/gemini/log"
	"github.com/panjf2000/gnet"
	"time"
)

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

func (s *Service) Logger() *log.ZapLogger {
	return s.logger
}
