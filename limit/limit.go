package limit

import (
	"golang.org/x/time/rate"
	"time"
)

// Limiter 访问频率限制类
type Limiter struct {
	*rate.Limiter
}

// New 构造函数
func New(duration time.Duration, burst int) *Limiter {
	return &Limiter{Limiter: rate.NewLimiter(rate.Every(duration), burst)}
}
