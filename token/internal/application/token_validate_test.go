package application

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/escalopa/fingo/token/internal/core"
	"github.com/escalopa/fingo/token/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTokenValidateCommandImpl(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		ctx      context.Context
		uuid     uuid.UUID
		token    string
		stubs    func(uuid.UUID, string, *mock.MockValidator, *mock.MockTokenRepository)
		response func(t *testing.T, got uuid.UUID, exp uuid.UUID, err error)
	}{
		{
			name:  "success",
			uuid:  uuid.New(),
			token: gofakeit.RandomString([]string{"0123456789abcdef"}),
			ctx:   contextutils.ConvertContext(contextutils.SetForwardMetadata(context.Background(), "192.172.19.0", "Mozilla/5.0")),
			stubs: func(id uuid.UUID, token string, v *mock.MockValidator, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), TokenValidateParams{AccessToken: token}).Return(nil)
				tr.EXPECT().GetTokenPayload(gomock.Any(), token).Return(&core.TokenPayload{
					UserID:    id,
					ClientIP:  "192.172.19.0",
					UserAgent: "Mozilla/5.0",
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
			},
			response: func(t *testing.T, got uuid.UUID, exp uuid.UUID, err error) {
				require.NoError(t, err)
				require.Equal(t, got.String(), exp.String())
			},
		},
		{
			name:  "invalid token",
			uuid:  uuid.New(),
			ctx:   context.Background(),
			token: gofakeit.RandomString([]string{""}),
			stubs: func(id uuid.UUID, token string, v *mock.MockValidator, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), TokenValidateParams{AccessToken: token}).Return(gofakeit.Error())
			},
			response: func(t *testing.T, got uuid.UUID, exp uuid.UUID, err error) {
				require.Error(t, err)
				require.NotEqual(t, got.String(), exp.String())
			},
		},
		{
			name:  "token not found",
			uuid:  uuid.New(),
			ctx:   context.Background(),
			token: gofakeit.RandomString([]string{"0123456789abcdef"}),
			stubs: func(id uuid.UUID, token string, v *mock.MockValidator, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), TokenValidateParams{AccessToken: token}).Return(nil)
				tr.EXPECT().GetTokenPayload(gomock.Any(), token).Return(nil, gofakeit.Error())
			},
			response: func(t *testing.T, got uuid.UUID, exp uuid.UUID, err error) {
				require.Error(t, err)
				require.NotEqual(t, got.String(), exp.String())
			},
		},
		{
			name:  "token expired",
			uuid:  uuid.New(),
			token: gofakeit.RandomString([]string{"0123456789abcdef"}),
			ctx:   context.Background(),
			stubs: func(id uuid.UUID, token string, v *mock.MockValidator, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), TokenValidateParams{AccessToken: token}).Return(nil)
				tr.EXPECT().GetTokenPayload(gomock.Any(), token).Return(&core.TokenPayload{
					ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
				}, nil)
			},
			response: func(t *testing.T, got uuid.UUID, exp uuid.UUID, err error) {
				require.Error(t, err)
				require.NotEqual(t, got.String(), exp.String())
			},
		},
		{
			name:  "client ip mismatch",
			uuid:  uuid.New(),
			token: gofakeit.RandomString([]string{"0123456789abcdef"}),
			ctx:   contextutils.ConvertContext(contextutils.SetForwardMetadata(context.Background(), "192.172.19.0", "Mozilla/5.0")),
			stubs: func(id uuid.UUID, token string, v *mock.MockValidator, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), TokenValidateParams{AccessToken: token}).Return(nil)
				tr.EXPECT().GetTokenPayload(gomock.Any(), token).Return(&core.TokenPayload{
					ClientIP:  "192.172.19.1", // Mismatch
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
			},
			response: func(t *testing.T, got uuid.UUID, exp uuid.UUID, err error) {
				require.Error(t, err)
				require.NotEqual(t, got.String(), exp.String())
			},
		},
		{
			name:  "user agent mismatch",
			uuid:  uuid.New(),
			token: gofakeit.RandomString([]string{"0123456789abcdef"}),
			ctx:   contextutils.ConvertContext(contextutils.SetForwardMetadata(context.Background(), "192.172.19.0", "Mozilla/5.0")),
			stubs: func(id uuid.UUID, token string, v *mock.MockValidator, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), TokenValidateParams{AccessToken: token}).Return(nil)
				tr.EXPECT().GetTokenPayload(gomock.Any(), token).Return(&core.TokenPayload{
					ClientIP:  "192.172.19.0",
					UserAgent: "Mozilla/5.0.1", // Mismatch
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
			},
			response: func(t *testing.T, got uuid.UUID, exp uuid.UUID, err error) {
				require.Error(t, err)
				require.NotEqual(t, got.String(), exp.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			v := mock.NewMockValidator(ctrl)
			c := mock.NewMockTokenRepository(ctrl)

			// Create command
			tv := NewTokenValidateCommand(v, c)
			tt.stubs(tt.uuid, tt.token, v, c)
			id, err := tv.Execute(tt.ctx, TokenValidateParams{AccessToken: tt.token})
			tt.response(t, id, tt.uuid, err)
		})
	}
}
