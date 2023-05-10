package testcontainer

import (
	"context"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

func TestNewRabbitMQContainer(t *testing.T) {
	url, terminate, err := NewRabbitMQContainer(context.Background())
	require.NoError(t, err)
	require.NotNil(t, url)

	defer func() { require.NoError(t, terminate()) }()

	conn, err := amqp.Dial(url)
	require.NoError(t, err)
	require.NotNil(t, conn)
}
