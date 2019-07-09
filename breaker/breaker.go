package breaker

import "github.com/sony/gobreaker"

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	*gobreaker.CircuitBreaker
}

// New 构造函数
func New() *CircuitBreaker {
	return &CircuitBreaker{CircuitBreaker: gobreaker.NewCircuitBreaker(gobreaker.Settings{})}
}
