package docker

import (
	"bytes"
	"context"
	"github.com/Scarlet-Fairy/cobold/pkg/push"
	docker "github.com/fsouza/go-dockerclient"
)

type dockerPush struct {
	client *docker.Client
}

func newPush(dockerClient *docker.Client) push.Push {
	return &dockerPush{
		client: dockerClient,
	}
}

func (d *dockerPush) Push(ctx context.Context, options push.Options) error {
	outputStream := bytes.NewBuffer(nil)

	pushOptions := docker.PushImageOptions{
		Registry:      options.Registry,
		Name:          options.Name,
		Tag:           options.Tag,
		Context:       ctx,
		OutputStream:  outputStream,
		RawJSONStream: false,
	}

	if err := d.client.PushImage(pushOptions, docker.AuthConfiguration{
		Username: "docker",
	}); err != nil {
		return err
	}

	return nil
}
