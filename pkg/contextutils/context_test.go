package contextutils

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name   string
		ctx    func(userID string) context.Context
		userID string
		expErr bool
	}{
		{
			name: "success",
			ctx: func(userID string) context.Context {
				return metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
					contextKeyUserID.String(), userID,
				))
			},
			userID: gofakeit.UUID(),
		},
		{
			name: "failed id not found",
			ctx: func(userID string) context.Context {
				return metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
					"not_user_id_header", userID,
				))
			},
			userID: gofakeit.UUID(),
			expErr: true,
		},
		{
			name: "failed to parse userID",
			ctx: func(userID string) context.Context {
				return metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
					contextKeyUserID.String(), userID[:1], // Wrong uuid
				))
			},
			userID: gofakeit.UUID(),
			expErr: true,
		},
		{
			name: "missing metadata",
			ctx: func(userID string) context.Context {
				return context.Background()
			},
			userID: gofakeit.UUID(),
			expErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := GetUserID(tt.ctx(tt.userID))
			require.Equal(t, err != nil, tt.expErr)
			if !tt.expErr {
				require.Equal(t, tt.userID, userID.String())
			}
		})
	}
}

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

func TestConvertContext(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "success context",
			ctx:  context.Background(),
		},
		{
			name: "success md context incoming context",
			ctx:  metadata.NewIncomingContext(context.Background(), nil),
		},
		{
			name: "success md context outgoing context",
			ctx:  metadata.NewOutgoingContext(context.Background(), nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ConvertContext(tt.ctx)
			require.NotNil(t, ctx)
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
