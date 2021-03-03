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

func newBuild(client *docker.Client) build.Build {
	return &dockerBuild{
		client,
	}
}

func (d dockerBuild) Build(ctx context.Context, options build.Options) (io.Reader, error) {
	inputStream, err := createTarball(options.Directory)
	if err != nil {
		return nil, errors.Wrap(err, "could not build image")
	}

	outputStream := bytes.NewBuffer(nil)

	buildOptions := docker.BuildImageOptions{
		Dockerfile:   "Dockerfile",
		Context:      ctx,
		Name:         options.Name,
		InputStream:  inputStream,
		OutputStream: outputStream,
	}

	if err := d.client.BuildImage(buildOptions); err != nil {
		return nil, errors.Wrap(err, "could not build image")
	}

	serialized, err := serializeOutputStreamBuffer(outputStream)
	if err != nil {
		return nil, errors.Wrap(err, "could not serialized buffer")
	}

	return bytes.NewBuffer(serialized), nil
}
