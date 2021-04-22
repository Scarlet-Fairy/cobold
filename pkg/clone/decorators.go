package clone

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

func (t *traceDecorator) Clone(ctx context.Context, options Options) (err error) {
	ctx, span := t.tracer.Start(ctx, "clone")
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()
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

	defer func(start time.Time) {
		l.logger.Log(
			"options.Url", options.Url,
			"options.Path", options.Path,
			"err", err,
			"msg", "Clone completed",
			"took", time.Since(start),
		)
	}(time.Now())
	return l.next.Clone(ctx, options)
}

type metricDecorator struct {
	jobID     string
	histogram metrics.Histogram
	next      Clone
}

func NewMetricDecorator(jobID string, histogram metrics.Histogram, next Clone) Clone {
	return &metricDecorator{
		jobID:     jobID,
		histogram: histogram,
		next:      next,
	}
}
func (m metricDecorator) Clone(ctx context.Context, options Options) error {
	defer func(begin time.Time) {
		m.histogram.With("jobID", m.jobID).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.next.Clone(ctx, options)
}
