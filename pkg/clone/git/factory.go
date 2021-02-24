package git

import (
	"github.com/Scarlet-Fairy/cobold/pkg/clone"
	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
)

func MakeClone(jobID string, logger log.Logger, tracer opentracing.Tracer) clone.Clone {
	var cloneInstance clone.Clone
	{
		cloneInstance = new()
		cloneInstance = clone.NewTraceDecorator(jobID, tracer, cloneInstance)
		cloneInstance = clone.NewLogDecorator(jobID, logger, cloneInstance)
	}

	return cloneInstance
}
