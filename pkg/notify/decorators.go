package notify

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
	tracer trace.Tracer
	next   Notify
}

func NewTraceDecorator(tracer trace.Tracer, next Notify) Notify {
	return &traceDecorator{
		tracer: tracer,
		next:   next,
	}
}

func (t *traceDecorator) NotifyCompletion(ctx context.Context, options Options) (err error) {
	ctx, span := t.tracer.Start(ctx, "notify")
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	var errorAttribute attribute.KeyValue
	{
		if options.Err != nil {
			errorAttribute = attribute.String("error", options.Err.Error())
		}
	}
	span.SetAttributes(
		attribute.String("jobId", options.JobID),
		attribute.String("reason", options.Reason),
		errorAttribute,
	)

	return t.next.NotifyCompletion(ctx, options)
}

type logDecorator struct {
	logger log.Logger
	next   Notify
}

func NewLogDecorator(logger log.Logger, next Notify) Notify {
	return &logDecorator{
		logger: logger,
		next:   next,
	}
}

func (l logDecorator) NotifyCompletion(ctx context.Context, options Options) (err error) {
	defer func(begin time.Time) {
		if err == nil {
			l.logger.Log(
				"took", time.Since(begin),
				"options.err", options.Err,
				"options.JobId", options.JobID,
				"options.Reason", options.Reason,
				"err", err,
				"msg", "Notification completed",
			)
		}
	}(time.Now())

	return l.next.NotifyCompletion(ctx, options)
}

type metricDecorator struct {
	histogram metrics.Histogram
	next      Notify
}

func NewMetricDecorator(histogram metrics.Histogram, next Notify) Notify {
	return &metricDecorator{
		histogram: histogram,
		next:      next,
	}
}
func (m metricDecorator) NotifyCompletion(ctx context.Context, options Options) error {
	defer func(begin time.Time) {
		m.histogram.With("jobID", options.JobID).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.next.NotifyCompletion(ctx, options)
}
