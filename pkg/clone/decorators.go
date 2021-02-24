package clone

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
	next   Clone
}

func NewTraceDecorator(jobID string, tracer opentracing.Tracer, next Clone) Clone {
	return &traceDecorator{
		jobID:  jobID,
		tracer: tracer,
		next:   next,
	}
}

func (t *traceDecorator) Clone(ctx context.Context, options Options) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "clone")
	defer span.Finish()
	span.LogFields(
		traceLog.String("jobId", t.jobID),
		traceLog.String("url", options.Url),
		traceLog.String("path", options.Path))

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
	l.logger.Log("jobID", l.jobID, "msg", "Start clone")
	defer func(begin time.Time) {
		if err == nil {
			l.logger.Log("jobID", l.jobID, "took", time.Since(begin), "msg", "End Clone")
		}
	}(time.Now())

	return l.next.Clone(ctx, options)
}
