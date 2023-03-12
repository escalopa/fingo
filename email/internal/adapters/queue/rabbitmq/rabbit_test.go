package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConsumer_HandleSendVerificationsCode(t *testing.T) {
	// Start the consumer
	results := make(chan sendVerificationCodeMessage)
	go func() {
		err := testRabbitMQ.HandleSendVerificationsCode(func(ctx context.Context, email string, code string) error {
			results <- sendVerificationCodeMessage{
				Email: email,
				Code:  code,
			}
			return nil
		})
		require.NoError(t, err)
	}()
	// Create a channel
	ch, err := testRabbitMQ.q.Channel()
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
		msg  sendVerificationCodeMessage
	}{
		{
			name: "Send verification code",
			msg: sendVerificationCodeMessage{
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
	results := make(chan sendResetPasswordTokenMessage)
	go func() {
		err := testRabbitMQ.HandleSendResetPasswordToken(func(ctx context.Context, email string, token string) error {
			results <- sendResetPasswordTokenMessage{
				Email: email,
				Token: token,
			}
			return nil
		})
		require.NoError(t, err)
	}()
	// Create a channel
	ch, err := testRabbitMQ.q.Channel()
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
		msg  sendResetPasswordTokenMessage
	}{
		{
			name: "Send reset password token",
			msg: sendResetPasswordTokenMessage{
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
	results := make(chan sendNewSignInSessionMessage)
	go func() {
		err := testRabbitMQ.HandleSendNewSignInSession(func(ctx context.Context, email string, clientIP string, userAgent string) error {
			results <- sendNewSignInSessionMessage{
				Email:     email,
				ClientIP:  clientIP,
				UserAgent: userAgent,
			}
			return nil
		})
		require.NoError(t, err)
	}()
	// Create a channel
	ch, err := testRabbitMQ.q.Channel()
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
		msg  sendNewSignInSessionMessage
	}{
		{
			name: "Send new sign in session",
			msg: sendNewSignInSessionMessage{
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
