package service

import (
	"github.com/openzipkin/zipkin-go/model"
	"net/http"
	"strconv"
)

// ExtractHttp 解包http头为SpanContext
func ExtractHttp(r *http.Request) func() *model.SpanContext {
	return func() *model.SpanContext {
		var (
			traceIDHeader      = r.Header.Get("jar-traceid")
			spanIDHeader       = r.Header.Get("jar-spanid")
			parentSpanIDHeader = r.Header.Get("jar-parentid")
			sampledHeader      = r.Header.Get("jar-sampled")
			flagsHeader        = r.Header.Get("jar-flags")
		)

		var sc = &model.SpanContext{}

		traceID, err := model.TraceIDFromHex(traceIDHeader)
		if err == nil {
			sc.TraceID = traceID
		}

		spanID, err := strconv.ParseUint(spanIDHeader, 16, 64)
		if err == nil {
			sc.ID = model.ID(spanID)
		}

		parentID, err := strconv.ParseUint(parentSpanIDHeader, 16, 64)
		if err == nil {
			pID := model.ID(parentID)
			sc.ParentID = &pID
		}

		if sampledHeader == "0" {
			sampled := false
			sc.Sampled = &sampled
		} else if sampledHeader == "1" {
			sampled := true
			sc.Sampled = &sampled
		}

		if flagsHeader == "1" {
			sc.Debug = true
		}

		return sc
	}
}
