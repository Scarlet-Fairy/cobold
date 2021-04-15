package docker

import (
	"github.com/Scarlet-Fairy/cobold/pkg/build"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"go.opentelemetry.io/otel/trace"
)

func MakeBuild(jobID string, dockerClient *docker.Client, duration metrics.Histogram, logger log.Logger, tracer trace.Tracer) build.Build {
	var buildInstance build.Build
	{
		buildInstance = newBuild(dockerClient)
		buildInstance = build.NewMetricDecorator(jobID, duration, buildInstance)
		buildInstance = build.NewTraceDecorator(jobID, tracer, buildInstance)
		buildInstance = build.NewLogDecorator(jobID, logger, buildInstance)
	}

	return buildInstance
}
