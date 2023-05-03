package application

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/escalopa/fingo/auth/internal/mock"
	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
	"github.com/stretchr/testify/require"
)

func TestLogout_Execute(t *testing.T) {

	tests := []struct {
		name      string
		userID    string
		sessionID string
		stubs     func(
			userID string,
			sessionID string,
			v *mock.MockValidator,
			sr *mock.MockSessionRepository,
			tr *mock.MockTokenRepository,
		) context.Context
		check func(t *testing.T, err error)
	}{
		{
			name:      "success",
			userID:    gofakeit.UUID(),
			sessionID: gofakeit.UUID(),
			stubs: func(
				userID string,
				sessionID string,
				v *mock.MockValidator,
				sr *mock.MockSessionRepository,
				tr *mock.MockTokenRepository) context.Context {

				v.EXPECT().Validate(gomock.Any(), LogoutParams{SessionID: sessionID}).Return(nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{
					ID: uuid.MustParse(sessionID), UserID: uuid.MustParse(userID), AccessToken: "access_token",
				}, nil)
				tr.EXPECT().Delete(gomock.Any(), "access_token").Return(nil)
				sr.EXPECT().DeleteSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(nil)
				return contextutils.SetUserID(context.Background(), userID)
			},
			check: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:      "invalid session id",
			userID:    gofakeit.UUID(),
			sessionID: gofakeit.UUID(),
			stubs: func(
				userID string,
				sessionID string,
				v *mock.MockValidator,
				sr *mock.MockSessionRepository,
				tr *mock.MockTokenRepository) context.Context {

				v.EXPECT().Validate(gomock.Any(), LogoutParams{SessionID: sessionID}).Return(errs.B().Msg("some error").Err())
				return contextutils.SetUserID(context.Background(), userID)
			},
			check: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name:      "failed to parse user id from context",
			userID:    gofakeit.UUID(),
			sessionID: gofakeit.UUID(),
			stubs: func(
				userID string,
				sessionID string,
				v *mock.MockValidator,
				sr *mock.MockSessionRepository,
				tr *mock.MockTokenRepository) context.Context {

				v.EXPECT().Validate(gomock.Any(), LogoutParams{SessionID: sessionID}).Return(nil)
				return context.Background()
			},
			check: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name:      "session not found",
			userID:    gofakeit.UUID(),
			sessionID: gofakeit.UUID(),
			stubs: func(
				userID string,
				sessionID string,
				v *mock.MockValidator,
				sr *mock.MockSessionRepository,
				tr *mock.MockTokenRepository) context.Context {

				v.EXPECT().Validate(gomock.Any(), LogoutParams{SessionID: sessionID}).Return(nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{}, errs.B().Msg("not found").Err())
				return contextutils.SetUserID(context.Background(), userID)
			},
			check: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name:      "caller is not session owner",
			userID:    gofakeit.UUID(),
			sessionID: gofakeit.UUID(),
			stubs: func(
				userID string,
				sessionID string,
				v *mock.MockValidator,
				sr *mock.MockSessionRepository,
				tr *mock.MockTokenRepository) context.Context {

				v.EXPECT().Validate(gomock.Any(), LogoutParams{SessionID: sessionID}).Return(nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{
					ID: uuid.MustParse(sessionID), UserID: uuid.MustParse(userID),
				}, nil)
				return contextutils.SetUserID(context.Background(), gofakeit.UUID())
			},
			check: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name:      "failed to delete token with errs",
			userID:    gofakeit.UUID(),
			sessionID: gofakeit.UUID(),
			stubs: func(
				userID string,
				sessionID string,
				v *mock.MockValidator,
				sr *mock.MockSessionRepository,
				tr *mock.MockTokenRepository) context.Context {

				v.EXPECT().Validate(gomock.Any(), LogoutParams{SessionID: sessionID}).Return(nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{
					UserID: uuid.MustParse(userID), AccessToken: "access_token",
				}, nil)
				tr.EXPECT().Delete(gomock.Any(), "access_token").Return(errs.B().Msg("some error").Err())
				return contextutils.SetUserID(context.Background(), userID)
			},
			check: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name:      "failed to delete token with errs",
			userID:    gofakeit.UUID(),
			sessionID: gofakeit.UUID(),
			stubs: func(
				userID string,
				sessionID string,
				v *mock.MockValidator,
				sr *mock.MockSessionRepository,
				tr *mock.MockTokenRepository) context.Context {

				v.EXPECT().Validate(gomock.Any(), LogoutParams{SessionID: sessionID}).Return(nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{
					UserID: uuid.MustParse(userID), AccessToken: "access_token",
				}, nil)
				tr.EXPECT().Delete(gomock.Any(), "access_token").Return(fmt.Errorf("some error"))
				return contextutils.SetUserID(context.Background(), userID)
			},
			check: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name:      "failed to delete session",
			userID:    gofakeit.UUID(),
			sessionID: gofakeit.UUID(),
			stubs: func(
				userID string,
				sessionID string,
				v *mock.MockValidator,
				sr *mock.MockSessionRepository,
				tr *mock.MockTokenRepository) context.Context {

				v.EXPECT().Validate(gomock.Any(), LogoutParams{SessionID: sessionID}).Return(nil)
				sr.EXPECT().GetSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(core.Session{UserID: uuid.MustParse(userID)}, nil)
				tr.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				sr.EXPECT().DeleteSessionByID(gomock.Any(), uuid.MustParse(sessionID)).Return(errs.B().Msg("some error").Err())
				return contextutils.SetUserID(context.Background(), userID)
			},
			check: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			v := mock.NewMockValidator(ctrl)
			sr := mock.NewMockSessionRepository(ctrl)
			tr := mock.NewMockTokenRepository(ctrl)

			c := NewLogoutCommand(v, sr, tr)
			ctx := tt.stubs(tt.userID, tt.sessionID, v, sr, tr)
			err := c.Execute(ctx, LogoutParams{SessionID: tt.sessionID})
			tt.check(t, err)
		})
	}
}
