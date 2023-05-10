package mygrpc

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/escalopa/fingo/token/internal/application"
	"github.com/escalopa/fingo/token/internal/core"
	"github.com/escalopa/fingo/token/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTokenHandlerValidateToken(t *testing.T) {
	// Create mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	v := mock.NewMockValidator(ctrl)
	tr := mock.NewMockTokenRepository(ctrl)

	// Create grpc server
	conn := setup(t, v, tr)
	defer conn.Close()

	// Create grpc client
	client := pb.NewTokenServiceClient(conn)

	test := []struct {
		name  string
		id    uuid.UUID
		ctx   context.Context
		req   *pb.ValidateTokenRequest
		stubs func(uuid.UUID, string, *mock.MockValidator, *mock.MockTokenRepository)
		check func(t *testing.T, got *pb.ValidateTokenResponse, err error)
	}{
		{
			name: "valid token",
			id:   uuid.New(),
			req:  &pb.ValidateTokenRequest{AccessToken: gofakeit.RandomString([]string{"0123456789abcdef"})},
			ctx:  contextutils.ConvertContext(contextutils.SetForwardMetadata(context.Background(), "192.172.19.0", "Mozilla/5.0")),
			stubs: func(id uuid.UUID, token string, v *mock.MockValidator, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), application.TokenValidateParams{AccessToken: token}).Return(nil)
				tr.EXPECT().GetTokenPayload(gomock.Any(), token).Return(&core.TokenPayload{
					UserID:    id,
					ClientIP:  "192.172.19.0",
					UserAgent: "Mozilla/5.0",
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
			},
			check: func(t *testing.T, got *pb.ValidateTokenResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, got.GetUserId(), got.GetUserId())
			},
		},
		{
			name: "invalid token",
			id:   uuid.New(),
			req:  &pb.ValidateTokenRequest{AccessToken: gofakeit.RandomString([]string{"0123456789abcdef"})},
			ctx:  contextutils.ConvertContext(contextutils.SetForwardMetadata(context.Background(), "", "")),
			stubs: func(id uuid.UUID, token string, v *mock.MockValidator, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), application.TokenValidateParams{AccessToken: token}).Return(gofakeit.Error())
			},
			check: func(t *testing.T, got *pb.ValidateTokenResponse, err error) {
				require.Error(t, err)
				require.Nil(t, got)
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.stubs(tt.id, tt.req.AccessToken, v, tr)
			resp, err := client.ValidateToken(tt.ctx, tt.req)
			tt.check(t, resp, err)
		})
	}
}
