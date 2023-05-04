package contextutils

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestGetAccessToken(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		expErr bool
	}{
		{
			name: "success",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				authorizationHeader, authorizationType+gofakeit.UUID(),
			)),
			expErr: false,
		},
		{
			name: "missing auth header",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				"not auth header", authorizationType+gofakeit.UUID(),
			)),
			expErr: true,
		},
		{
			name: "missing auth type",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				authorizationHeader, "not bearer"+gofakeit.UUID(),
			)),
			expErr: true,
		},
		{
			name: "missing token",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				authorizationHeader, authorizationType+" ",
			)),
			expErr: true,
		},
		{
			name:   "missing metadata",
			ctx:    context.Background(),
			expErr: true,
		},
	}

	for _, tt := range tests {
		token, err := GetAccessToken(tt.ctx)
		require.Equal(t, err != nil, tt.expErr)
		if !tt.expErr {
			require.NotEmpty(t, token)
		}
	}
}
