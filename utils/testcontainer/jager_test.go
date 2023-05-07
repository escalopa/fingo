package testcontainer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewJagerContainer(t *testing.T) {
	// Create redis container
	url, terminate, err := NewJagerContainer(context.Background())
	require.NoError(t, err)
	require.NotNil(t, url)
	require.NoError(t, terminate())
}
