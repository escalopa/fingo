package testcontainer

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/require"
)

func TestNewRedisContainer(t *testing.T) {
	// Create redis container
	url, terminate, err := NewRedisContainer(context.Background())
	require.NoError(t, err)
	require.NotNil(t, url)

	defer func() { require.NoError(t, terminate()) }()

	// Create redis client
	URL, err := redis.ParseURL(url)
	require.NoError(t, err)
	c := redis.NewClient(URL)
	require.NoError(t, err)

	// Ping redis client
	r := c.Conn().Ping(context.Background())
	require.NoError(t, r.Err())
}
