package docker

import (
	"github.com/Scarlet-Fairy/cobold/pkg/push"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/go-kit/kit/log"
	"go.opentelemetry.io/otel/trace"
)

func MakePush(jobID string, dockerClient *docker.Client, logger log.Logger, tracer trace.Tracer) push.Push {
	var dockerBuild push.Push
	{
		dockerBuild = newPush(dockerClient)
		dockerBuild = push.NewLogDecorator(jobID, logger, dockerBuild)
		dockerBuild = push.NewTraceDecorator(jobID, tracer, dockerBuild)
	}

	return dockerBuild
}
