package notify

import (
	"context"
	"github.com/go-kit/kit/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

func (t *traceDecorator) NotifyCompletion(ctx context.Context, options Options) error {
	ctx, span := t.tracer.Start(ctx, "notify")
	defer span.End()

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

func (l logDecorator) NotifyCompletion(ctx context.Context, options Options) error {
	l.logger.Log("reason", options.Reason, "error", options.Err)

	return l.next.NotifyCompletion(ctx, options)
}
