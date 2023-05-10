package rabbitmq

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/stretchr/testify/require"
)

func TestProducer_NewProducer(t *testing.T) {
	t.Parallel()
	// Empty url for producer
	p, err := NewProducer("")
	require.Error(t, err)
	require.Nil(t, p)
	// Igonre setting `ssq` name
	p, err = NewProducer(rabbitmqUrl)
	require.Error(t, err)
	require.Nil(t, p)
	// Create the producer
	testProducer, err := NewProducer(rabbitmqUrl,
		WithNewSignInSessionQueue("new_sign_in_session_queue"),
	)
	require.NoError(t, err)
	require.NotNil(t, testProducer)
	defer func() {
		if err = testProducer.Close(); err != nil {
			require.NoError(t, err)
		}
	}()
}

func TestProducer_SendNewSignInSessionMessage(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		params      core.SendNewSignInSessionParams
		expectError bool
	}{
		{
			name: "valid",
			params: core.SendNewSignInSessionParams{
				Name:      gofakeit.FirstName(),
				Email:     gofakeit.Email(),
				ClientIP:  gofakeit.IPv4Address(),
				UserAgent: gofakeit.UserAgent(),
			},
			expectError: false,
		},
		{
			name: "invalid name",
			params: core.SendNewSignInSessionParams{
				Name:      "",
				Email:     gofakeit.Email(),
				ClientIP:  gofakeit.IPv4Address(),
				UserAgent: gofakeit.UserAgent(),
			},
			expectError: true,
		},
		{
			name: "invalid email",
			params: core.SendNewSignInSessionParams{
				Name:      gofakeit.FirstName(),
				Email:     "invalid",
				ClientIP:  gofakeit.IPv4Address(),
				UserAgent: gofakeit.UserAgent(),
			},
			expectError: true,
		},
		{
			name: "invalid client ip",
			params: core.SendNewSignInSessionParams{
				Name:      gofakeit.FirstName(),
				Email:     gofakeit.Email(),
				ClientIP:  "invalid",
				UserAgent: gofakeit.UserAgent(),
			},
			expectError: true,
		},
		{
			name: "invalid user agent",
			params: core.SendNewSignInSessionParams{
				Name:      gofakeit.FirstName(),
				Email:     gofakeit.Email(),
				ClientIP:  gofakeit.IPv4Address(),
				UserAgent: "",
			},
			expectError: true,
		},
	}
	// Create a new producer
	testProducer, err := NewProducer(rabbitmqUrl,
		WithNewSignInSessionQueue("new_sign_in_session_queue"),
	)
	defer func() {
		if err = testProducer.Close(); err != nil {
			require.NoError(t, err)
		}
	}()
	// Send messages to the queue
	for _, tc := range testCases {
		err := testProducer.SendNewSignInSessionMessage(context.Background(), tc.params)
		require.NoError(t, err)
	}
	messages, err := testProducer.msgChan.Consume(
		testProducer.ssqQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	require.NoError(t, err)
	// Read messages from the queue and compare them with the expected values
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := <-messages
			var receivedParams core.SendNewSignInSessionParams
			err = json.Unmarshal(msg.Body, &receivedParams)
			require.NoError(t, err)
			require.Equal(t, tc.params, receivedParams)
		})
	}
}

func TestProducer_Close(t *testing.T) {
	t.Parallel()
	// Create the producer
	testProducer, err := NewProducer(rabbitmqUrl,
		WithNewSignInSessionQueue("new_sign_in_session_queue"),
	)
	require.NoError(t, err)
	require.NotNil(t, testProducer)
	require.NoError(t, testProducer.Close())
	// Create a producer & close `msgChan` before `close`
	testProducer, err = NewProducer(rabbitmqUrl,
		WithNewSignInSessionQueue("new_sign_in_session_queue"),
	)
	require.NoError(t, err)
	require.NotNil(t, testProducer)
	require.NoError(t, testProducer.msgChan.Close())
	require.Error(t, testProducer.Close())
	// Create a producer & close `connection` before close
	testProducer, err = NewProducer(rabbitmqUrl,
		WithNewSignInSessionQueue("new_sign_in_session_queue"),
	)
	require.NoError(t, err)
	require.NotNil(t, testProducer)
	require.NoError(t, testProducer.q.Close())
	require.Error(t, testProducer.Close())

}
