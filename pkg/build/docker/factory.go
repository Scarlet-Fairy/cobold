package docker

import (
	"github.com/Scarlet-Fairy/cobold/pkg/build"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/opentracing/opentracing-go"
	"os"
)

func MakeBuild(jobID, dockerEnpoint string, logger log.Logger, tracer opentracing.Tracer) build.Build {
	var buildInstance build.Build
	var err error
	{
		buildInstance, err = newBuild(dockerEnpoint)
		if err != nil {
			level.Error(logger).Log("docker-endpoint", dockerEnpoint, "err", err.Error())
			os.Exit(1)
		}

		buildInstance = build.NewTraceDecorator(jobID, tracer, buildInstance)
		buildInstance = build.NewLogDecorator(jobID, logger, buildInstance)
	}

	return buildInstance
}
