package breaker

import (
	"github.com/Jarnpher553/gemini/log"
	"github.com/sony/gobreaker"
	"time"
)

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	*gobreaker.CircuitBreaker
}

var l = log.Zap.Mark("breaker")
var state = map[gobreaker.State]string{
	0: "close",
	1: "half-open",
	2: "open",
}

// New 构造函数
func New() *CircuitBreaker {
	return &CircuitBreaker{CircuitBreaker: gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "micro-breaker",
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			l.Info(log.Messagef("%s change from %d to %d", name, state[from], state[to]))
		},
	})}
}
