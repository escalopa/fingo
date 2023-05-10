package testcontainer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPostgresContainer(t *testing.T) {
	conn, terminate, err := NewPostgresContainer(context.Background())
	require.NoError(t, err)
	require.NotNil(t, conn)

	defer func() { require.NoError(t, terminate()) }()

	require.NoError(t, conn.Ping())
}
