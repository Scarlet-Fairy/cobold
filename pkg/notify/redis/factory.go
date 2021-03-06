package redis

import (
	"github.com/Scarlet-Fairy/cobold/pkg/notify"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/trace"
)

func MakeNotify(redisClient *redis.Client, duration metrics.Histogram, logger log.Logger, tracer trace.Tracer) notify.Notify {
	redisClient.AddHook(redisotel.TracingHook{})

	var notifyInstance notify.Notify
	{
		notifyInstance = newNotify(redisClient)
		notifyInstance = notify.NewMetricDecorator(duration, notifyInstance)
		notifyInstance = notify.NewLogDecorator(logger, notifyInstance)
		notifyInstance = notify.NewTraceDecorator(tracer, notifyInstance)
	}

	return notifyInstance
}
