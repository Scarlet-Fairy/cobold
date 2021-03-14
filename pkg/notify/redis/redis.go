package redis

import (
	"context"
	"github.com/Scarlet-Fairy/cobold/pkg/notify"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type redisNotify struct {
	client *redis.Client
}

func newNotify(client *redis.Client) notify.Notify {
	return &redisNotify{
		client: client,
	}

}

func (r *redisNotify) NotifyCompletion(ctx context.Context, options notify.Options) error {
	msg, err := encodeMessageToJson(message{
		Error: options.Err,
	})
	if err != nil {
		return err
	}

	if cmd := r.client.Publish(ctx, pubChannel(options.JobID, options.Reason), msg); cmd.Err() != nil {
		return errors.Wrap(cmd.Err(), "could not notify clone completed")
	}

	return nil
}
