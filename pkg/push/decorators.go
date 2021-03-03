package push

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
	traceLog "github.com/opentracing/opentracing-go/log"
	"time"
)

type traceDecorator struct {
	jobID  string
	tracer opentracing.Tracer
	next   Push
}

func NewTraceDecorator(jobID string, tracer opentracing.Tracer, next Push) Push {
	return &traceDecorator{
		jobID:  jobID,
		tracer: tracer,
		next:   next,
	}
}

func (t *traceDecorator) Push(ctx context.Context, options Options) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "push")
	defer span.Finish()
	span.LogFields(
		traceLog.String("jobId", t.jobID),
		traceLog.String("name", options.Name),
		traceLog.String("tag", options.Tag))

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
