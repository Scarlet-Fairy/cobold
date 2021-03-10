package push

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
	next   Push
}

func NewTraceDecorator(jobID string, tracer trace.Tracer, next Push) Push {
	return &traceDecorator{
		jobID:  jobID,
		tracer: tracer,
		next:   next,
	}
}

func (t *traceDecorator) Push(ctx context.Context, options Options) error {
	ctx, span := t.tracer.Start(ctx, "push")
	defer span.End()
	span.SetAttributes(
		attribute.String("jobId", t.jobID),
		attribute.String("name", options.Name),
		attribute.String("tag", options.Tag),
		attribute.String("registry", options.Registry))

	return t.next.Push(ctx, options)
}

type logDecorator struct {
	jobID  string
	logger log.Logger
	next   Push
}

func NewLogDecorator(jobID string, logger log.Logger, next Push) Push {
	return &logDecorator{
		jobID:  jobID,
		logger: logger,
		next:   next,
	}
}

func (l logDecorator) Push(ctx context.Context, options Options) (err error) {
	l.logger.Log("msg", "Start Push")
	defer func(begin time.Time) {
		if err == nil {
			l.logger.Log("took", time.Since(begin), "msg", "End Push")
		}
	}(time.Now())

	return l.next.Push(ctx, options)
}
