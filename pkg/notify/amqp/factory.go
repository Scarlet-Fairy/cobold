package amqp

import (
	"github.com/Scarlet-Fairy/cobold/pkg/notify"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
)

func MakeNotify(rabbitMqChannel *amqp.Channel, jobId string, duration metrics.Histogram, logger log.Logger, tracer trace.Tracer) notify.Notify {

	var notifyInstance notify.Notify
	{
		notifyInstance = newNotify(rabbitMqChannel, jobId)
		notifyInstance = notify.NewMetricDecorator(duration, notifyInstance)
		notifyInstance = notify.NewLogDecorator(logger, notifyInstance)
		notifyInstance = notify.NewTraceDecorator(tracer, notifyInstance)
	}

	return notifyInstance
}
