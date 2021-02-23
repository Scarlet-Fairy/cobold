package docker

import (
	"bytes"
	"github.com/Scarlet-Fairy/cobold/pkg/build"
	"github.com/Scarlet-Fairy/cobold/pkg/log"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"io"
)

type dockerBuild struct {
	client *docker.Client
}

func new(endpoint string) build.Build {
	client, err := docker.NewClient(endpoint)
	if err != nil {
		log.Logger.WithField("endpoint", endpoint).WithField("error", err).Fatalf("Could not create client")
	}

	return &dockerBuild{
		client,
	}
}

func (d dockerBuild) Build(options build.BuildOptions) (io.Reader, error) {
	outputStream := bytes.NewBuffer(nil)

	if err := d.client.BuildImage(docker.BuildImageOptions{
		ContextDir:   options.Directory,
		Name:         options.Tag,
		OutputStream: outputStream,
	}); err != nil {
		return nil, errors.Wrap(err, "could not build image")
	}

	return outputStream, nil
}
