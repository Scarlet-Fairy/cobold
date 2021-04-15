package git

import (
	"github.com/Scarlet-Fairy/cobold/pkg/clone"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"go.opentelemetry.io/otel/trace"
)

func MakeClone(jobID string, duration metrics.Histogram, logger log.Logger, tracer trace.Tracer) clone.Clone {
	var cloneInstance clone.Clone
	{
		cloneInstance = newClone()
		cloneInstance = clone.NewMetricDecorator(jobID, duration, cloneInstance)
		cloneInstance = clone.NewLogDecorator(jobID, logger, cloneInstance)
		cloneInstance = clone.NewTraceDecorator(jobID, tracer, cloneInstance)
	}

	return cloneInstance
}
