package service

import (
	"errors"
	"fmt"
	"github.com/Jarnpher553/gemini/pkg/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Jarnpher553/gemini/pkg/breaker"
	"github.com/Jarnpher553/gemini/pkg/erro"
	"github.com/Jarnpher553/gemini/pkg/limit"
	"github.com/Jarnpher553/gemini/pkg/metric"
	"github.com/Jarnpher553/gemini/pkg/tracing"
)

// Middleware 中间件
type Middleware func(IBaseService) HandlerFunc

// MetricMiddleware 指标监控中间件
func MetricMiddleware(m *metric.Metric) Middleware {
	return func(srv IBaseService) HandlerFunc {
		name := srv.Node().RootName + "." + srv.Node().AreaName + "." + srv.Node().Name
		m.SetName(name)
		return func(context *Ctx) {
			defer func(begin time.Time) {
				m.ReqCount.Inc(1)
				m.ReqDuration.UpdateSince(begin)
			}(time.Now())
			context.Next()
		}
	}
}

// TracerMiddleware 服务跟踪中间件
func TracerMiddleware(t *tracing.Tracer) Middleware {
	return func(srv IBaseService) HandlerFunc {
		return func(context *Ctx) {
			sc, _ := t.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(context.Request.Header))

			span := t.StartSpan(srv.Node().RootName+"."+srv.Node().AreaName+"."+srv.Node().Name,
				opentracing.ChildOf(sc),
				ext.SpanKindRPCServer,
				opentracing.StartTime(time.Now()),
				opentracing.Tag{Key: string(ext.HTTPUrl), Value: context.Request.URL.Path},
				opentracing.Tag{Key: string(ext.HTTPMethod), Value: context.Request.Method},
			)

			defer func() {
				code := context.Writer.Status()
				span.SetTag(string(ext.HTTPStatusCode), strconv.Itoa(code))
				defer span.Finish()
			}()

			ctx := opentracing.ContextWithSpan(context.Request.Context(), span)

			rNew := context.Request.WithContext(ctx)
			context.Request = rNew
			context.Next()
		}
	}
}

// BreakerMiddleware 断路器中间件
func BreakerMiddleware(cb *breaker.CircuitBreaker) Middleware {
	return func(srv IBaseService) HandlerFunc {
		return func(ctx *Ctx) {
			_, err := cb.Execute(func() (i interface{}, e error) {
				defer func() {
					if err := recover(); err != nil {
						e = fmt.Errorf("%v", err)
						_, f, l, _ := runtime.Caller(2)
						fSlice := strings.Split(f, "/")
						f = strings.Join(fSlice[len(fSlice)-2:], "/")
						log.Logger.Error("break",
							zap.String("service", srv.Node().RootName+"."+srv.Node().AreaName+"."+srv.Node().Name),
							zap.String("target", fmt.Sprintf("%s:%d", f, l)),
							zap.Error(e))
					}
				}()
				ctx.Next()
				return nil, nil
			})

			if err != nil {
				switch cb.State() {
				case gobreaker.StateClosed:
					ctx.Failure(erro.ErrDefault, err)
				case gobreaker.StateOpen:
					ctx.Failure(erro.ErrBreaker, err)
				case gobreaker.StateHalfOpen:
					ctx.Failure(erro.ErrMaxRequest, err)
				}

				ctx.Abort()
				return
			}
		}
	}
}

// RateLimiterMiddleware 频率限制中间件
func RateLimiterMiddleware(limiter *limit.Limiter) Middleware {
	return func(srv IBaseService) HandlerFunc {
		return func(ctx *Ctx) {
			if !limiter.Allow() {
				ctx.Failure(erro.ErrRateLimiter, errors.New("rate limit exceeded"))
				ctx.Abort()
				return
			}
		}
	}
}

// DelayLimiterMiddleware 频率延迟中间件
func DelayLimiterMiddleware(limiter *limit.Limiter) Middleware {
	return func(srv IBaseService) HandlerFunc {
		return func(ctx *Ctx) {
			if err := limiter.Wait(ctx.Request.Context()); err != nil {
				ctx.Failure(erro.ErrDelayLimiter, err)
				ctx.Abort()
				return
			}
		}
	}
}

func ReserveLimiterMiddleware(limiter *limit.Limiter) Middleware {
	return func(srv IBaseService) HandlerFunc {
		return func(ctx *Ctx) {
			r := limiter.Reserve()
			if !r.OK() {
				ctx.Failure(erro.ErrReserveLimiter, errors.New("lim.burst must to be > 0"))
				ctx.Abort()
				return
			}
			<-time.After(r.Delay())
		}
	}
}
