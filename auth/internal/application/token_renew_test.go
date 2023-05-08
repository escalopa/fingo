package application

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/escalopa/fingo/auth/internal/mock"
	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRenewTokenCommand_Execute(t *testing.T) {
	type args struct {
		params RenewTokenParams
		ctx    func(userID string) context.Context
	}
	tests := []struct {
		name  string
		arg   args
		stubs func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository)
		check func(t *testing.T, response RenewTokenResponse, err error)
	}{
		{
			name: "success on token delete done in goroutine",
			arg: args{
				params: RenewTokenParams{
					RefreshToken: gofakeit.UUID(),
				},
				ctx: func(userID string) context.Context { return contextutils.SetUserID(context.Background(), userID) },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), arg.RefreshToken).Return(core.TokenPayload{
					UserID:    uuid.MustParse(userID),
					SessionID: uuid.MustParse(sessionID),
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{
					RefreshToken: arg.RefreshToken,
					AccessToken:  "old_token",
				}, nil)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("token1", nil)
				tg.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).Return("token2", nil)
				sr.EXPECT().UpdateSessionTokens(gomock.Any(), core.UpdateSessionTokenParams{
					ID:           uuid.MustParse(sessionID),
					AccessToken:  "token1",
					RefreshToken: "token2",
				}).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), "token1").Return(core.TokenPayload{}, nil)
				tr.EXPECT().Store(gomock.Any(), "token1", gomock.Any()).Return(nil)
				tr.EXPECT().Delete(gomock.Any(), "old_token").Return(nil)
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				time.Sleep(1 * time.Second) // Wait for tr.Delete to be called in the goroutine
				require.Equal(t, response, RenewTokenResponse{AccessToken: "token1", RefreshToken: "token2"})
				require.NoError(t, err)
			},
		},
		{
			name: "success on token delete failed in goroutine",
			arg: args{
				params: RenewTokenParams{
					RefreshToken: gofakeit.UUID(),
				},
				ctx: func(userID string) context.Context { return contextutils.SetUserID(context.Background(), userID) },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), arg.RefreshToken).Return(core.TokenPayload{
					UserID:    uuid.MustParse(userID),
					SessionID: uuid.MustParse(sessionID),
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{
					RefreshToken: arg.RefreshToken,
					AccessToken:  "old_token",
				}, nil)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("token1", nil)
				tg.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).Return("token2", nil)
				sr.EXPECT().UpdateSessionTokens(gomock.Any(), core.UpdateSessionTokenParams{
					ID:           uuid.MustParse(sessionID),
					AccessToken:  "token1",
					RefreshToken: "token2",
				}).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), "token1").Return(core.TokenPayload{}, nil)
				tr.EXPECT().Store(gomock.Any(), "token1", gomock.Any()).Return(nil)
				tr.EXPECT().Delete(gomock.Any(), "old_token").Return(gofakeit.Error())
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				time.Sleep(1 * time.Second) // Wait for tr.Delete to be called in the goroutine
				require.Equal(t, response, RenewTokenResponse{AccessToken: "token1", RefreshToken: "token2"})
				require.NoError(t, err)
			},
		},
		{
			name: "validation error",
			arg: args{
				params: RenewTokenParams{},
				ctx:    func(userID string) context.Context { return context.Background() },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(gofakeit.Error())
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "decrypt refresh token error",
			arg: args{
				params: RenewTokenParams{
					RefreshToken: gofakeit.UUID(),
				},
				ctx: func(userID string) context.Context { return context.Background() },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), arg.RefreshToken).Return(core.TokenPayload{}, gofakeit.Error())
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "get caller ID error",
			arg: args{
				params: RenewTokenParams{
					RefreshToken: gofakeit.UUID(),
				},
				ctx: func(userID string) context.Context { return context.Background() },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), arg.RefreshToken).Return(core.TokenPayload{}, nil)
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "caller is not the owner of the session",
			arg: args{
				params: RenewTokenParams{
					RefreshToken: gofakeit.UUID(),
				},
				ctx: func(userID string) context.Context { return contextutils.SetUserID(context.Background(), userID) },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), arg.RefreshToken).Return(core.TokenPayload{UserID: uuid.MustParse(gofakeit.UUID())}, nil)
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "referesh token expired",
			arg: args{
				params: RenewTokenParams{
					RefreshToken: gofakeit.UUID(),
				},
				ctx: func(userID string) context.Context { return contextutils.SetUserID(context.Background(), userID) },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), arg.RefreshToken).Return(core.TokenPayload{UserID: uuid.MustParse(userID), ExpiresAt: time.Now().Add(-1 * time.Second)}, nil)
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "session not found",
			arg: args{
				params: RenewTokenParams{
					RefreshToken: gofakeit.UUID(),
				},
				ctx: func(userID string) context.Context { return contextutils.SetUserID(context.Background(), userID) },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), arg.RefreshToken).Return(core.TokenPayload{
					UserID:    uuid.MustParse(userID),
					SessionID: uuid.MustParse(sessionID),
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{}, gofakeit.Error())
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "session refresh token mismatch",
			arg: args{
				params: RenewTokenParams{
					RefreshToken: gofakeit.UUID(),
				},
				ctx: func(userID string) context.Context { return contextutils.SetUserID(context.Background(), userID) },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), arg.RefreshToken).Return(core.TokenPayload{
					UserID:    uuid.MustParse(userID),
					SessionID: uuid.MustParse(sessionID),
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{RefreshToken: gofakeit.UUID()}, nil)
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "failed to create new access token",
			arg: args{
				params: RenewTokenParams{
					RefreshToken: gofakeit.UUID(),
				},
				ctx: func(userID string) context.Context { return contextutils.SetUserID(context.Background(), userID) },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), arg.RefreshToken).Return(core.TokenPayload{
					UserID:    uuid.MustParse(userID),
					SessionID: uuid.MustParse(sessionID),
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{RefreshToken: arg.RefreshToken}, nil)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("", gofakeit.Error())
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "failed to create new refresh token",
			arg: args{
				params: RenewTokenParams{
					RefreshToken: gofakeit.UUID(),
				},
				ctx: func(userID string) context.Context { return contextutils.SetUserID(context.Background(), userID) },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), arg.RefreshToken).Return(core.TokenPayload{
					UserID:    uuid.MustParse(userID),
					SessionID: uuid.MustParse(sessionID),
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{RefreshToken: arg.RefreshToken}, nil)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return(gofakeit.UUID(), nil)
				tg.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).Return("", gofakeit.Error())
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "update session tokens failed",
			arg: args{
				params: RenewTokenParams{
					RefreshToken: gofakeit.UUID(),
				},
				ctx: func(userID string) context.Context { return contextutils.SetUserID(context.Background(), userID) },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), arg.RefreshToken).Return(core.TokenPayload{
					UserID:    uuid.MustParse(userID),
					SessionID: uuid.MustParse(sessionID),
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{RefreshToken: arg.RefreshToken}, nil)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("token1", nil)
				tg.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).Return("token2", nil)
				sr.EXPECT().UpdateSessionTokens(gomock.Any(), core.UpdateSessionTokenParams{
					ID:           uuid.MustParse(sessionID),
					AccessToken:  "token1",
					RefreshToken: "token2",
				}).Return(gofakeit.Error())
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "dectypt new access token failed",
			arg: args{
				params: RenewTokenParams{
					RefreshToken: gofakeit.UUID(),
				},
				ctx: func(userID string) context.Context { return contextutils.SetUserID(context.Background(), userID) },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), arg.RefreshToken).Return(core.TokenPayload{
					UserID:    uuid.MustParse(userID),
					SessionID: uuid.MustParse(sessionID),
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{RefreshToken: arg.RefreshToken}, nil)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("token1", nil)
				tg.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).Return("token2", nil)
				sr.EXPECT().UpdateSessionTokens(gomock.Any(), core.UpdateSessionTokenParams{
					ID:           uuid.MustParse(sessionID),
					AccessToken:  "token1",
					RefreshToken: "token2",
				}).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), "token1").Return(core.TokenPayload{}, gofakeit.Error())
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "store token in cache failed",
			arg: args{
				params: RenewTokenParams{
					RefreshToken: gofakeit.UUID(),
				},
				ctx: func(userID string) context.Context { return contextutils.SetUserID(context.Background(), userID) },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), arg.RefreshToken).Return(core.TokenPayload{
					UserID:    uuid.MustParse(userID),
					SessionID: uuid.MustParse(sessionID),
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{RefreshToken: arg.RefreshToken}, nil)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("token1", nil)
				tg.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).Return("token2", nil)
				sr.EXPECT().UpdateSessionTokens(gomock.Any(), core.UpdateSessionTokenParams{
					ID:           uuid.MustParse(sessionID),
					AccessToken:  "token1",
					RefreshToken: "token2",
				}).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), "token1").Return(core.TokenPayload{}, nil)
				tr.EXPECT().Store(gomock.Any(), "token1", gomock.Any()).Return(gofakeit.Error())
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "store token in cache failed",
			arg: args{
				params: RenewTokenParams{
					RefreshToken: gofakeit.UUID(),
				},
				ctx: func(userID string) context.Context { return contextutils.SetUserID(context.Background(), userID) },
			},
			stubs: func(userID string, sessionID string, arg RenewTokenParams, v *mock.MockValidator, tg *mock.MockTokenGenerator, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), arg.RefreshToken).Return(core.TokenPayload{
					UserID:    uuid.MustParse(userID),
					SessionID: uuid.MustParse(sessionID),
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{RefreshToken: arg.RefreshToken}, nil)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("token1", nil)
				tg.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).Return("token2", nil)
				sr.EXPECT().UpdateSessionTokens(gomock.Any(), core.UpdateSessionTokenParams{
					ID:           uuid.MustParse(sessionID),
					AccessToken:  "token1",
					RefreshToken: "token2",
				}).Return(nil)
				tg.EXPECT().DecryptToken(gomock.Any(), "token1").Return(core.TokenPayload{}, nil)
				tr.EXPECT().Store(gomock.Any(), "token1", gomock.Any()).Return(gofakeit.Error())
			},
			check: func(t *testing.T, response RenewTokenResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			v := mock.NewMockValidator(ctrl)
			tg := mock.NewMockTokenGenerator(ctrl)
			sr := mock.NewMockSessionRepository(ctrl)
			tr := mock.NewMockTokenRepository(ctrl)

			c := NewRenewTokenCommand(v, tg, sr, tr)

			userID := gofakeit.UUID()
			sessionID := gofakeit.UUID()
			tt.stubs(userID, sessionID, tt.arg.params, v, tg, sr, tr)
			resp, err := c.Execute(tt.arg.ctx(userID), tt.arg.params)
			tt.check(t, resp, err)
		})
	}
}
