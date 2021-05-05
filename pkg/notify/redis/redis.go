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

func (r *redisNotify) Init(_ context.Context) error {
	return nil
}

func (r *redisNotify) NotifyCompletion(ctx context.Context, options notify.Options) error {
	msg := notify.Message{
		Topic: options.Reason,
	}
	if options.Err != nil {
		msg.Error = options.Err.Error()
	}

	encodedMsg, err := notify.EncodeMessageToJson(msg)
	if err != nil {
		return err
	}

	if cmd := r.client.Publish(ctx, pubChannel(options.JobID), encodedMsg); cmd.Err() != nil {
		return errors.Wrap(cmd.Err(), "could not notify clone completed")
	}

	return nil
}
