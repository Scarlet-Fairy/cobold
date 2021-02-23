package docker

import "github.com/Scarlet-Fairy/cobold/pkg/build"

func MakeDockerBuild(endpoint string) build.Build {
	var docker build.Build
	{
		docker = new(endpoint)
	}

	return docker
}
