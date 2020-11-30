package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Jarnpher553/gemini/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sony/gobreaker"
	"strconv"
	"strings"
	"time"

	"github.com/Jarnpher553/gemini/breaker"
	"github.com/Jarnpher553/gemini/erro"
	"github.com/Jarnpher553/gemini/jwt"
	"github.com/Jarnpher553/gemini/limit"
	"github.com/Jarnpher553/gemini/metric"
	"github.com/Jarnpher553/gemini/tracing"
)

// Middleware 中间件
type Middleware func(IBaseService) HandlerFunc

func ExtractHttpMiddleware() Middleware {
	return func(srv IBaseService) HandlerFunc {
		return func(context *Ctx) {
			scPtr := ExtractHttp(context.Request)()

			ctx := tracing.NewContextFromSpanContext(context.Request.Context(), scPtr)
			context.Request = context.Request.WithContext(ctx)
		}
	}
}

// MetricMiddleware 指标监控中间件
func MetricMiddleware(m *metric.Metric) Middleware {
	return func(srv IBaseService) HandlerFunc {
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

			span := t.StartSpan(srv.Node().ServerName+"."+srv.Node().Name,
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
						action := strings.Split(ctx.Request.URL.Path, "/")
						e = fmt.Errorf("%v service %s action %s", err, srv.Node().ServerName, action[3]+"."+action[4])
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

func AuthMiddleware() Middleware {
	return func(baseService IBaseService) HandlerFunc {
		return func(ctx *Ctx) {

			token := ctx.GetHeader("Authorization")
			claims, err := jwt.Parse(token)
			if err != nil {
				rdClient := baseService.Redis()

				if rdClient == nil {
					ctx.Failure(erro.ErrAuthor, err)
					ctx.Abort()
					return
				} else {

					uid := rdClient.Get(token)

					if uid == "" {
						ctx.Failure(erro.ErrAuthor, err)
						ctx.Abort()
						return
					} else {
						var cc context.Context
						if err := uuid.IsGUID(uid); err != nil {
							uidInt, _ := strconv.Atoi(uid)
							cc = context.WithValue(ctx.Request.Context(), "auth_user_id", uidInt)
						} else {
							cc = context.WithValue(ctx.Request.Context(), "auth_user_id", uuid.GUID(uid))
						}
						ctx.Request = ctx.Request.WithContext(cc)
					}
				}
			} else {
				var cc context.Context
				if claims.UserIdInt != 0 {
					cc = context.WithValue(ctx.Request.Context(), "auth_user_id", claims.UserIdInt)
				} else {
					cc = context.WithValue(ctx.Request.Context(), "auth_user_id", claims.UserIdUUID)
				}
				ctx.Request = ctx.Request.WithContext(cc)
			}
		}
	}
}
