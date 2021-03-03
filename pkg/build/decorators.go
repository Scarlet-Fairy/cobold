package build

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
	traceLog "github.com/opentracing/opentracing-go/log"
	"io"
	"time"
)

type traceDecorator struct {
	jobID  string
	tracer opentracing.Tracer
	next   Build
}

func NewTraceDecorator(jobID string, tracer opentracing.Tracer, next Build) Build {
	return &traceDecorator{
		jobID:  jobID,
		tracer: tracer,
		next:   next,
	}
}

func (b *traceDecorator) Build(ctx context.Context, options Options) (io.Reader, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, b.tracer, "build")
	defer span.Finish()
	span.LogFields(
		traceLog.String("jobID", b.jobID),
		traceLog.String("directory", options.Directory),
		traceLog.String("tag", options.Name))

	return b.next.Build(ctx, options)
}

type logDecorator struct {
	jobID  string
	logger log.Logger
	next   Build
}

func NewLogDecorator(jobID string, logger log.Logger, next Build) Build {
	return &logDecorator{
		jobID:  jobID,
		logger: logger,
		next:   next,
	}
}

func (l logDecorator) Build(ctx context.Context, options Options) (reader io.Reader, err error) {
	l.logger.Log("msg", "Start build")
	defer func(begin time.Time) {
		if err == nil {
			l.logger.Log("took", time.Since(begin), "msg", "End build")
		}
	}(time.Now())

	return l.next.Build(ctx, options)
}
