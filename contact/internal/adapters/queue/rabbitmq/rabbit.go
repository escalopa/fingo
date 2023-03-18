package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/escalopa/fingo/contact/internal/core"

	"github.com/lordvidex/errs"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	q *amqp.Connection

	vcq string // verificationCodeQueueName
	rsq string // resetPasswordTokenQueueName
	ssq string // newSignInSessionQueueName
}

func NewConsumer(url string, opts ...func(*Consumer)) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, errs.B(err).Code(errs.InvalidArgument).Msg("failed to connect to rabbitmq").Err()
	}

	r := &Consumer{q: conn}
	for _, opt := range opts {
		opt(r)
	}

	if r.vcq == "" {
		return nil, errs.B(err).Code(errs.InvalidArgument).
			Msg("RabbitMQ Consumer: sendVerificationCodeQueueName is not set").Err()
	}
	if r.rsq == "" {
		return nil, errs.B(err).Code(errs.InvalidArgument).
			Msg("RabbitMQ Consumer: sendResetPasswordTokenQueueName is not set").Err()
	}
	if r.ssq == "" {
		return nil, errs.B(err).Code(errs.InvalidArgument).
			Msg("RabbitMQ Consumer: sendNewSignInSessionQueueName is not set").Err()
	}

	return r, nil
}

func WithVerificationCodeQueue(name string) func(*Consumer) {
	return func(r *Consumer) {
		r.vcq = name
	}
}

func WithResetPasswordTokenQueue(name string) func(*Consumer) {
	return func(r *Consumer) {
		r.rsq = name
	}
}

func WithNewSignInSessionQueue(name string) func(*Consumer) {
	return func(r *Consumer) {
		r.ssq = name
	}
}

func (r *Consumer) HandleSendVerificationsCode(handler func(ctx context.Context, params core.SendVerificationCodeMessage) error) error {
	messages, err := r.setupQueue(r.vcq)
	if err != nil {
		return errs.B(err).Code(errs.InvalidArgument).Msg("failed to setup queue on send verification code").Err()
	}
	for d := range messages {
		go func(d amqp.Delivery) {
			var m core.SendVerificationCodeMessage
			r.handleMessage(d, &m, func(ctx context.Context) error {
				return handler(ctx, m)
			})
		}(d)
	}
	return nil
}

func (r *Consumer) HandleSendResetPasswordToken(handler func(ctx context.Context, params core.SendResetPasswordTokenMessage) error) error {
	messages, err := r.setupQueue(r.rsq)
	if err != nil {
		return errs.B(err).Code(errs.InvalidArgument).Msg("failed to setup queue on send verification code").Err()
	}
	for d := range messages {
		go func(d amqp.Delivery) {
			var m core.SendResetPasswordTokenMessage
			r.handleMessage(d, &m, func(ctx context.Context) error {
				return handler(ctx, m)
			})
		}(d)
	}
	return nil
}

func (r *Consumer) HandleSendNewSignInSession(handler func(ctx context.Context, params core.SendNewSignInSessionMessage) error) error {
	messages, err := r.setupQueue(r.ssq)
	if err != nil {
		return errs.B(err).Code(errs.InvalidArgument).Msg("failed to setup queue on send verification code").Err()
	}
	for d := range messages {
		go func(d amqp.Delivery) {
			var m core.SendNewSignInSessionMessage
			r.handleMessage(d, &m, func(ctx context.Context) error {
				return handler(ctx, m)
			})
		}(d)
	}
	return nil
}

func (r *Consumer) handleMessage(msg amqp.Delivery, body interface{}, handle func(context.Context) error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		if err := msg.Ack(false); err != nil {
			log.Println("failed to ack message: ", err, msg.MessageId)
		}
		cancel()
	}()
	// Read message from queue
	err := json.Unmarshal(msg.Body, &body)
	if err != nil {
		log.Println("failed to unmarshal message: ", err)
		if err = msg.Nack(false, true); err != nil {
			log.Println("failed to nack message: ", err, msg.MessageId)
		}
		return
	}
	// Handle message
	err = handle(ctx)
	if err != nil {
		log.Println("failed to handle message: ", err)
		return
	}
	if err = msg.Ack(false); err != nil {
		log.Println("failed to ack message: ", err, msg.MessageId)
	}
}

func (r *Consumer) setupQueue(queueName string) (<-chan amqp.Delivery, error) {
	ch, err := r.q.Channel()
	if err != nil {
		return nil, errs.B(err).Code(errs.InvalidArgument).Msg("failed to open a channel").Err()
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, errs.B(err).Code(errs.InvalidArgument).
			Msg("failed to declare a queue").Err()
	}

	err = ch.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, errs.B(err).Code(errs.InvalidArgument).
			Msg("failed to set QoS").Err()
	}

	messages, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, errs.B(err).Code(errs.InvalidArgument).
			Msg("failed to register a consumer").Err()
	}

	return messages, nil
}

func (r *Consumer) Close() error {
	err := r.q.Close()
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("failed to close connection").Err()
	}
	return nil
}
