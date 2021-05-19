package amqp

import (
	"context"
	"github.com/Scarlet-Fairy/cobold/pkg/notify"
	"github.com/streadway/amqp"
)

const (
	Exchanger = "build_image"
)

type rabbitNotify struct {
	channel *amqp.Channel
	jobId   string
}

func newNotify(ch *amqp.Channel, jobId string) notify.Notify {
	return &rabbitNotify{
		channel: ch,
		jobId:   jobId,
	}
}

func (r *rabbitNotify) Init(_ context.Context) error {
	if err := r.channel.ExchangeDeclare(
		Exchanger,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := r.declareQueues(r.jobId); err != nil {
		return err
	}

	return nil
}

func (r *rabbitNotify) declareQueues(id string) error {
	_, err := r.channel.QueueDeclare(
		id,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	if err := r.channel.QueueBind(
		id,
		id,
		Exchanger,
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (r *rabbitNotify) NotifyCompletion(_ context.Context, options notify.Options) error {
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

	if err := r.channel.Publish(
		Exchanger,
		options.JobID,
		false,
		false,
		amqp.Publishing{
			Body: encodedMsg,
		},
	); err != nil {
		return err
	}

	return nil
}
