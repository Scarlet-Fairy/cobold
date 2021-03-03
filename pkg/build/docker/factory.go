package docker

import (
	"github.com/Scarlet-Fairy/cobold/pkg/build"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
)

func MakeBuild(jobID string, dockerClient *docker.Client, logger log.Logger, tracer opentracing.Tracer) build.Build {
	var buildInstance build.Build
	{
		buildInstance = newBuild(dockerClient)
		buildInstance = build.NewTraceDecorator(jobID, tracer, buildInstance)
		buildInstance = build.NewLogDecorator(jobID, logger, buildInstance)
	}

	return buildInstance
}
