package tracing

import (
	"context"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/Janrpher553/micro-core/log"
)

// Tracer 跟踪类
type Tracer struct {
	*zipkin.Tracer
}

// New 构造函数
func New(reporter reporter.Reporter) *Tracer {
	t, err := zipkin.NewTracer(reporter, zipkin.WithSharedSpans(false))

	if err != nil {
		log.Logger.Mark("Tracer").Fatalln(err)
	}
	return &Tracer{
		Tracer: t,
	}
}

// SpanFromContext 从context中获取Span
func SpanFromContext(ctx context.Context) zipkin.Span {
	if s, ok := ctx.Value("span_key").(zipkin.Span); ok {
		return s
	}
	return nil
}

// NewContext 将Span写入context
func NewContext(ctx context.Context, s zipkin.Span) context.Context {
	return context.WithValue(ctx, "span_key", s)
}

// NewContextFromSpanContext 将SpanContext写入context
func NewContextFromSpanContext(ctx context.Context, spanContext *model.SpanContext) context.Context {
	return context.WithValue(ctx, "span_context_key", spanContext)
}

// SpanContextFromContext 从context中获取SpanContext
func SpanContextFromContext(ctx context.Context) *model.SpanContext {
	if s, ok := ctx.Value("span_context_key").(*model.SpanContext); ok {
		return s
	}
	return nil
}
