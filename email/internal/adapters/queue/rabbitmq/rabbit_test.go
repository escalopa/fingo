package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/escalopa/fingo/email/internal/core"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

func TestConsumer_HandleSendVerificationsCode(t *testing.T) {
	// Start the consumer
	results := make(chan core.SendVerificationCodeMessage)
	go func() {
		err := testConsumer.HandleSendVerificationsCode(func(ctx context.Context, params core.SendVerificationCodeMessage) error {
			results <- params
			return nil
		})
		require.NoError(t, err)
	}()
	// Create a channel
	ch, err := testConsumer.q.Channel()
	require.NoError(t, err)
	defer func() { require.NoError(t, ch.Close()) }()
	// Declare the queue
	queue, err := ch.QueueDeclare(
		"verification_code_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	require.NoError(t, err)
	// Publish a message to the queue
	testCases := []struct {
		name string
		msg  core.SendVerificationCodeMessage
	}{
		{
			name: "success",
			msg: core.SendVerificationCodeMessage{
				Email: gofakeit.Email(),
				Code:  gofakeit.UUID(),
			},
		},
	}
	// Process the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.msg)
			require.NoError(t, err)
			// Publish the message
			err = ch.PublishWithContext(context.Background(),
				"",
				queue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        b,
				},
			)
			require.NoError(t, err)
			// Wait for the message to be processed
			select {
			case result := <-results:
				require.Equal(t, tc.msg, result)
			case <-time.After(5 * time.Second):
				t.Fatal("timeout")
			}
		})
	}

}

func TestConsumer_HandleSendResetPasswordToken(t *testing.T) {
	// Start the consumer
	results := make(chan core.SendResetPasswordTokenMessage)
	go func() {
		err := testConsumer.HandleSendResetPasswordToken(func(ctx context.Context, params core.SendResetPasswordTokenMessage) error {
			results <- params
			return nil
		})
		require.NoError(t, err)
	}()
	// Create a channel
	ch, err := testConsumer.q.Channel()
	require.NoError(t, err)
	defer func() { require.NoError(t, ch.Close()) }()
	// Declare the queue
	queue, err := ch.QueueDeclare(
		"reset_password_token_queue",
		true,
		false,

		false,
		false,
		nil,
	)
	require.NoError(t, err)
	// Publish a message to the queue
	testCases := []struct {
		name string
		msg  core.SendResetPasswordTokenMessage
	}{
		{
			name: "success",
			msg: core.SendResetPasswordTokenMessage{
				Email: gofakeit.Email(),
				Token: gofakeit.UUID(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.msg)
			require.NoError(t, err)
			// Publish the message
			err = ch.PublishWithContext(context.Background(),
				"",
				queue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        b,
				},
			)
			require.NoError(t, err)
			// Wait for the message to be processed
			select {
			case result := <-results:
				require.Equal(t, tc.msg, result)
			case <-time.After(5 * time.Second):
				t.Fatal("timeout")
			}
		})
	}
}

func TestConsumer_HandleSendNewSignInSession(t *testing.T) {
	// Start the consumer
	results := make(chan core.SendNewSignInSessionMessage)
	go func() {
		err := testConsumer.HandleSendNewSignInSession(func(ctx context.Context, params core.SendNewSignInSessionMessage) error {
			results <- params
			return nil
		})
		require.NoError(t, err)
	}()
	// Create a channel
	ch, err := testConsumer.q.Channel()
	require.NoError(t, err)
	defer func() { require.NoError(t, ch.Close()) }()
	// Declare the queue
	queue, err := ch.QueueDeclare(
		"new_sign_in_session_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	require.NoError(t, err)
	// Publish a message to the queue
	testCases := []struct {
		name string
		msg  core.SendNewSignInSessionMessage
	}{
		{
			name: "success",
			msg: core.SendNewSignInSessionMessage{
				Email:     gofakeit.Email(),
				ClientIP:  gofakeit.IPv4Address(),
				UserAgent: gofakeit.UserAgent(),
			},
		},
	}
	// Process the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.msg)
			require.NoError(t, err)
			// Publish the message
			err = ch.PublishWithContext(context.Background(),
				"",
				queue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        b,
				},
			)
			require.NoError(t, err)
			// Wait for the message to be processed
			select {
			case result := <-results:
				require.Equal(t, tc.msg, result)
			case <-time.After(5 * time.Second):
				t.Fatal("timeout")
			}
		})
	}
}
