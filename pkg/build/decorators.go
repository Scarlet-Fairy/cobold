package build

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

func (b *traceDecorator) Build(ctx context.Context, options Options) (reader io.Reader, err error) {
	ctx, span := b.tracer.Start(ctx, "build")
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	span.SetAttributes(
		attribute.String("jobID", b.jobID),
		attribute.String("directory", options.Directory),
		attribute.String("name", options.Name),
	)

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
	defer func(begin time.Time) {
		l.logger.Log(
			"options.Name", options.Name,
			"options.Directory", options.Directory,
			"err", err,
			"msg", "Build completed",
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.next.Build(ctx, options)
}

type metricDecorator struct {
	jobID     string
	histogram metrics.Histogram
	next      Build
}

func NewMetricDecorator(jobID string, histogram metrics.Histogram, next Build) Build {
	return &metricDecorator{
		jobID:     jobID,
		histogram: histogram,
		next:      next,
	}
}
func (m metricDecorator) Build(ctx context.Context, options Options) (io.Reader, error) {
	defer func(begin time.Time) {
		m.histogram.With("jobID", m.jobID).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.next.Build(ctx, options)
}
