package docker

import (
	"bytes"
	"context"
	"github.com/Scarlet-Fairy/cobold/pkg/push"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
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
	ouputStream := bytes.NewBuffer(nil)

	pushOptions := docker.PushImageOptions{
		Registry:      options.Registry,
		Name:          options.Name,
		Tag:           options.Tag,
		Context:       ctx,
		OutputStream:  ouputStream,
		RawJSONStream: true,
	}
	if err := d.client.PushImage(pushOptions, docker.AuthConfiguration{
		Username: "docker",
	}); err != nil {
		return errors.Wrap(err, "could not push image")
	}

	return nil
}
