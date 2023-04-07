package contextutils

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestSetUserID(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		userID string
		expErr bool
	}{
		{
			name:   "success context",
			ctx:    context.Background(),
			userID: gofakeit.UUID(),
			expErr: false,
		},
		{
			name:   "success md context",
			ctx:    metadata.NewIncomingContext(context.Background(), nil),
			userID: gofakeit.UUID(),
			expErr: false,
		},
		{
			name:   "failed to parse userID",
			ctx:    metadata.NewIncomingContext(context.Background(), nil),
			userID: gofakeit.UUID()[:1], // Wrong uuid
			expErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := SetUserID(tt.ctx, tt.userID)
			userID, err := GetUserID(ctx)
			require.Equal(t, err != nil, tt.expErr)
			if !tt.expErr {
				require.Equal(t, tt.userID, userID.String())
			}
		})
	}
}

func TestGetMetadata(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		clientIP  string
		userAgent string
	}{
		{
			name:      "success",
			ctx:       context.Background(),
			clientIP:  gofakeit.IPv4Address(),
			userAgent: gofakeit.UserAgent(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(tt.ctx, metadata.Pairs(
				contextKeyUserAgent.String(), tt.userAgent,
				contextKeyClientIP.String(), tt.clientIP,
			))
			clientIP, userAgent := GetMetadata(ctx)
			require.Equal(t, tt.clientIP, clientIP)
			require.Equal(t, tt.userAgent, userAgent)
		})
	}
}

func TestSetForwardMetadata(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		clientIP  string
		userAgent string
	}{
		{
			name:      "success context",
			ctx:       context.Background(),
			clientIP:  gofakeit.IPv4Address(),
			userAgent: gofakeit.UserAgent(),
		},
		{
			name:      "success md context",
			ctx:       metadata.NewIncomingContext(context.Background(), nil),
			clientIP:  gofakeit.IPv4Address(),
			userAgent: gofakeit.UserAgent(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create out going context
			ctx := SetForwardMetadata(tt.ctx, tt.clientIP, tt.userAgent)
			// Convert outgoing to ingoing
			md, ok := metadata.FromOutgoingContext(ctx)
			require.True(t, ok)
			ctx = metadata.NewIncomingContext(ctx, md)
			// Send ctx as incoming context
			clientIP, userAgent := GetForwardMetadata(ctx)
			require.Equal(t, tt.clientIP, clientIP)
			require.Equal(t, tt.userAgent, userAgent)
		})
	}
}

func TestFormatClientIP(t *testing.T) {
	tests := []struct {
		name       string
		sentIP     string
		expectedIP string
	}{
		{
			name:       "success ipv4",
			sentIP:     "172.19.0.1:45066",
			expectedIP: "172.19.0.1",
		},
		{
			name:       "success ipv6",
			sentIP:     "[123:123:123:123]:45066",
			expectedIP: "123:123:123:123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIP := formatClientIP(tt.sentIP)
			require.Equal(t, gotIP, tt.expectedIP)
		})
	}
}
