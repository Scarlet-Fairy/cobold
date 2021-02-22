package clone

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type tracing struct {
	tracer opentracing.Tracer
	next   Clone
}

func NewTraceDecorator(tracer opentracing.Tracer, next Clone) Clone {
	return &tracing{
		tracer: tracer,
		next:   next,
	}
}

func (t *tracing) Clone(ctx context.Context, url string, path string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "clone")
	defer span.Finish()
	span.LogFields(
		log.String("url", url),
		log.String("path", path))

	return t.next.Clone(ctx, url, path)
}
