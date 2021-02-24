package docker

import (
	"bytes"
	"context"
	"github.com/Scarlet-Fairy/cobold/pkg/build"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"io"
)

type dockerBuild struct {
	client *docker.Client
}

func newBuild(endpoint string) (build.Build, error) {
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "could not create docker builder")
	}

	return &dockerBuild{
		client,
	}, nil
}

func (d dockerBuild) Build(_ context.Context, options build.Options) (io.Reader, error) {
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
