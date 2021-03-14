package notify

import (
	"context"
	"github.com/go-kit/kit/log"
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
	l.logger.Log("msg", "Start Notify")
	defer func(begin time.Time) {
		if err == nil {
			l.logger.Log("took", time.Since(begin), "msg", "End Notify")
		}
	}(time.Now())

	return l.next.NotifyCompletion(ctx, options)
}
