package git

import (
	"github.com/Scarlet-Fairy/cobold/pkg/clone"
	"github.com/go-kit/kit/log"
	"go.opentelemetry.io/otel/trace"
)

func MakeClone(jobID string, logger log.Logger, tracer trace.Tracer) clone.Clone {
	var cloneInstance clone.Clone
	{
		cloneInstance = newClone()
		cloneInstance = clone.NewLogDecorator(jobID, logger, cloneInstance)
		cloneInstance = clone.NewTraceDecorator(jobID, tracer, cloneInstance)
	}

	return cloneInstance
}
