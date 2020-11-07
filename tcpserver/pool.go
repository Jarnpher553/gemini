package tcpserver

import "github.com/panjf2000/gnet/pool/goroutine"

type GoroutinePool = goroutine.Pool

func Pool() *GoroutinePool {
	return goroutine.Default()
}
