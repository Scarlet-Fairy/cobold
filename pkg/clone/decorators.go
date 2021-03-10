package clone

import (
	"context"
	"github.com/go-kit/kit/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type traceDecorator struct {
	jobID  string
	tracer trace.Tracer
	next   Clone
}

func NewTraceDecorator(jobID string, tracer trace.Tracer, next Clone) Clone {
	return &traceDecorator{
		jobID:  jobID,
		tracer: tracer,
		next:   next,
	}
}

func (t *traceDecorator) Clone(ctx context.Context, options Options) error {
	ctx, span := t.tracer.Start(ctx, "clone")
	defer span.End()
	span.SetAttributes(
		attribute.String("jobId", t.jobID),
		attribute.String("url", options.Url),
		attribute.String("path", options.Path),
	)

	return t.next.Clone(ctx, options)
}

type logDecorator struct {
	jobID  string
	logger log.Logger
	next   Clone
}

func NewLogDecorator(jobID string, logger log.Logger, next Clone) Clone {
	return &logDecorator{
		jobID:  jobID,
		logger: logger,
		next:   next,
	}
}

func (l logDecorator) Clone(ctx context.Context, options Options) (err error) {
	l.logger.Log("msg", "Start clone")
	defer func(begin time.Time) {
		if err == nil {
			l.logger.Log("took", time.Since(begin), "msg", "End Clone")
		}
	}(time.Now())

	return l.next.Clone(ctx, options)
}
