package docker

import (
	"github.com/Scarlet-Fairy/cobold/pkg/build"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/go-kit/kit/log"
	"go.opentelemetry.io/otel/trace"
)

func MakeBuild(jobID string, dockerClient *docker.Client, logger log.Logger, tracer trace.Tracer) build.Build {
	var buildInstance build.Build
	{
		buildInstance = newBuild(dockerClient)
		buildInstance = build.NewTraceDecorator(jobID, tracer, buildInstance)
		buildInstance = build.NewLogDecorator(jobID, logger, buildInstance)
	}

	return buildInstance
}
