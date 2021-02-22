package git

import (
	"github.com/Scarlet-Fairy/cobold/pkg/clone"
	"github.com/opentracing/opentracing-go"
)

func MakeClone(tracer opentracing.Tracer) clone.Clone {
	var cloneInstance clone.Clone
	{
		cloneInstance = new()
		cloneInstance = clone.NewTraceDecorator(tracer, cloneInstance)
	}

	return cloneInstance
}
