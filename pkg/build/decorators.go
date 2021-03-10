package build

import (
	"context"
	"github.com/go-kit/kit/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"io"
	"time"
)

type traceDecorator struct {
	jobID  string
	tracer trace.Tracer
	next   Build
}

func NewTraceDecorator(jobID string, tracer trace.Tracer, next Build) Build {
	return &traceDecorator{
		jobID:  jobID,
		tracer: tracer,
		next:   next,
	}
}

func (b *traceDecorator) Build(ctx context.Context, options Options) (io.Reader, error) {
	ctx, span := b.tracer.Start(ctx, "build")
	defer span.End()
	span.SetAttributes(
		attribute.String("jobID", b.jobID),
		attribute.String("directory", options.Directory),
		attribute.String("tag", options.Name))

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
