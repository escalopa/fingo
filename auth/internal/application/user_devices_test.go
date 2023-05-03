package application

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/escalopa/fingo/auth/internal/mock"
	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetUserDevices_Execute(t *testing.T) {

	tests := []struct {
		name   string
		userID string
		stubs  func(userID string, ur *mock.MockValidator, sr *mock.MockSessionRepository) context.Context
		check  func(t *testing.T, sessions []core.Session, err error)
	}{
		{
			name:   "success",
			userID: gofakeit.UUID(),
			stubs: func(userID string, v *mock.MockValidator, sr *mock.MockSessionRepository) context.Context {
				v.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
				sr.EXPECT().GetUserSessions(gomock.Any(), uuid.MustParse(userID)).Return([]core.Session{}, nil)
				return contextutils.SetUserID(context.Background(), userID)
			},
			check: func(t *testing.T, sessions []core.Session, err error) {
				require.NoError(t, err)
				require.NotNil(t, sessions)
			},
		},
		{
			name:   "validation error",
			userID: gofakeit.UUID(),
			stubs: func(userID string, v *mock.MockValidator, sr *mock.MockSessionRepository) context.Context {
				v.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(gofakeit.Error())
				return context.Background()
			},
			check: func(t *testing.T, sessions []core.Session, err error) {
				require.Error(t, err)
				require.Nil(t, sessions)
			},
		},
		{
			name: "user id not set in context",
			stubs: func(userID string, v *mock.MockValidator, sr *mock.MockSessionRepository) context.Context {
				v.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
				return context.Background()
			},
			check: func(t *testing.T, sessions []core.Session, err error) {
				require.Error(t, err)
				require.Nil(t, sessions)
			},
		},
		{
			name:   "get user sessions error",
			userID: gofakeit.UUID(),
			stubs: func(userID string, v *mock.MockValidator, sr *mock.MockSessionRepository) context.Context {
				v.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
				sr.EXPECT().GetUserSessions(gomock.Any(), uuid.MustParse(userID)).Return(nil, gofakeit.Error())
				return contextutils.SetUserID(context.Background(), userID)
			},
			check: func(t *testing.T, sessions []core.Session, err error) {
				require.Error(t, err)
				require.Nil(t, sessions)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			v := mock.NewMockValidator(ctrl)
			sr := mock.NewMockSessionRepository(ctrl)

			cmd := NewGetUserDevicesCommand(v, sr)

			ctx := tt.stubs(tt.userID, v, sr)
			sessions, err := cmd.Execute(ctx, GetUserDevicesParams{})
			tt.check(t, sessions, err)
		})
	}
}
