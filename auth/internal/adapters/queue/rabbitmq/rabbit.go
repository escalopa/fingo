package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/lordvidex/errs"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	q        *amqp.Connection
	ssq      string // newSignInSessionQueueName
	ssqQueue amqp.Queue
	msgChan  *amqp.Channel
}

func NewProducer(url string, opts ...func(*Producer)) (*Producer, error) {
	// Connect to rabbitmq server
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, errs.B(err).Code(errs.InvalidArgument).Msg("failed to connect to rabbitmq").Err()
	}
	// Create producer
	p := &Producer{q: conn}
	for _, opt := range opts {
		opt(p)
	}
	// Check if queue name is set for sending new login sessions message
	if p.ssq == "" {
		return nil, errs.B(err).Code(errs.InvalidArgument).
			Msg("failed to parse queue name for sending new login sessions message").Err()
	}
	// Create user channel
	p.msgChan, err = p.q.Channel()
	if err != nil {
		return nil, errs.B(err).Code(errs.InvalidArgument).Msg("failed to open a channel").Err()
	}
	p.ssqQueue, err = p.msgChan.QueueDeclare(
		p.ssq, // name
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errs.B(err).Code(errs.Internal).Msg("failed to declare a queue for producer").Err()
	}
	return p, nil
}

// WithNewSignInSessionQueue sets the queue name for sending new login sessions message
func WithNewSignInSessionQueue(name string) func(*Producer) {
	return func(r *Producer) {
		r.ssq = name
	}
}

// SendNewSignInSessionMessage sends a message to the queue to send a new login session email
func (r *Producer) SendNewSignInSessionMessage(ctx context.Context, params core.SendNewSignInSessionParams) error {
	// Marshal message
	b, err := json.Marshal(params)
	if err != nil {
		return errs.B(err).Code(errs.InvalidArgument).Msg("failed to marshal message").Err()
	}
	// Publish message to queue
	err = r.msgChan.PublishWithContext(ctx,
		"",
		r.ssq,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})
	if err != nil {
		return errs.B(err).Code(errs.InvalidArgument).Msg("failed to publish message").Err()
	}
	return nil
}

// Close closes the connection to the queue
func (r *Producer) Close() error {
	err := r.msgChan.Close()
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("failed to close message channel").Err()
	}
	err = r.q.Close()
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("failed to close connection with rabbitmq").Err()
	}
	return nil
}
