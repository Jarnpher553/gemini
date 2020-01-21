package limit

import (
	"golang.org/x/time/rate"
)

// Limiter 访问频率限制类
type Limiter struct {
	*rate.Limiter
}

// New 构造函数
func New(limit rate.Limit, burst int) *Limiter {
	return &Limiter{Limiter: rate.NewLimiter(limit, burst)}
}
