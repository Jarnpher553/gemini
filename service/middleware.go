package service

import (
	"context"
	"fmt"
	"github.com/Jarnpher553/micro-core/uuid"
	"strconv"
	"strings"
	"time"

	"github.com/Jarnpher553/micro-core/breaker"
	"github.com/Jarnpher553/micro-core/erro"
	"github.com/Jarnpher553/micro-core/jwt"
	"github.com/Jarnpher553/micro-core/limit"
	"github.com/Jarnpher553/micro-core/log"
	"github.com/Jarnpher553/micro-core/metric"
	"github.com/Jarnpher553/micro-core/tracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
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
func TracerMiddleware(t *tracing.Tracer, name string) Middleware {
	return func(srv IBaseService) HandlerFunc {
		return func(context *Ctx) {
			var sc model.SpanContext

			if parentSc := tracing.SpanContextFromContext(context.Request.Context()); parentSc != nil {
				sc = *parentSc
			}

			ep, _ := zipkin.NewEndpoint("", context.ClientIP())
			sp := t.StartSpan(name, zipkin.Parent(sc), zipkin.Kind(model.Server), zipkin.StartTime(time.Now()), zipkin.RemoteEndpoint(ep))

			zipkin.TagHTTPMethod.Set(sp, context.Request.Method)
			zipkin.TagHTTPPath.Set(sp, context.Request.URL.Path)

			defer func() {
				code := context.Writer.Status()
				zipkin.TagHTTPStatusCode.Set(sp, strconv.Itoa(code))
				sp.Finish()
			}()

			ctx := tracing.NewContext(context.Request.Context(), sp)

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
						//var buf [2 << 10]byte
						//stack := string(buf[:runtime.Stack(buf[:], true)])
						action := strings.Split(ctx.Request.URL.Path, "/")
						e = fmt.Errorf("%v service %s action %s", err, srv.Node().ServerName, action[2]+"."+action[3])
					}
				}()
				ctx.Next()
				return nil, nil
			})

			if err != nil {
				log.Logger.Mark("Breaker").Errorln(erro.ErrMsg[erro.ErrBreaker], err)

				ctx.Response(erro.ErrBreaker, nil)
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
				log.Logger.Mark("Limiter").Errorln(erro.ErrMsg[erro.ErrRateLimiter], "rate limit exceeded")

				ctx.Response(erro.ErrRateLimiter, nil)
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
				log.Logger.Mark("Limiter").Errorln(erro.ErrMsg[erro.ErrDelayLimiter], err)

				ctx.Response(erro.ErrDelayLimiter, nil)
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
				log.Logger.Mark("Limiter").Errorln(erro.ErrMsg[erro.ErrReserveLimiter], "Did you remember to set lim.burst to be > 0 ?")

				ctx.Response(erro.ErrReserveLimiter, nil)
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
					log.Logger.Mark("Author").Errorln(erro.ErrAuthor, erro.ErrMsg[erro.ErrAuthor], err)
					ctx.Response(erro.ErrAuthor, nil)
					ctx.Abort()
					return
				} else {

					uid := rdClient.Get(token)

					if uid == "" {
						log.Logger.Mark("Author").Errorln(erro.ErrAuthor, erro.ErrMsg[erro.ErrAuthor], err)
						ctx.Response(erro.ErrAuthor, nil)
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
