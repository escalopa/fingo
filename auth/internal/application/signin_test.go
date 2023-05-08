package application

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/escalopa/fingo/auth/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSigninCommand_Execute(t *testing.T) {
	type args struct {
		params SigninParams
	}
	tests := []struct {
		name  string
		arg   args
		stubs func(arg SigninParams, v *mock.MockValidator, h *mock.MockPasswordHasher, tg *mock.MockTokenGenerator, ur *mock.MockUserRepository, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository, mp *mock.MockMessageProducer)
		check func(t *testing.T, response SigninResponse, err error)
	}{
		{
			name: "success on send new sign in session message done",
			arg: args{
				params: SigninParams{
					Email:     gofakeit.Email(),
					Password:  gofakeit.Password(true, true, true, true, false, 32),
					ClientIP:  gofakeit.IPv4Address(),
					UserAgent: gofakeit.UserAgent(),
				},
			},
			stubs: func(arg SigninParams, v *mock.MockValidator, h *mock.MockPasswordHasher, tg *mock.MockTokenGenerator, ur *mock.MockUserRepository, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository, mp *mock.MockMessageProducer) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				ur.EXPECT().GetUserByEmail(gomock.Any(), arg.Email).Return(core.User{
					FirstName:      "fingo_user",
					Email:          arg.Email,
					HashedPassword: gofakeit.NewCrypto().Password(true, true, true, true, false, 32),
				}, nil)
				h.EXPECT().Compare(gomock.Any(), gomock.Any(), arg.Password).Return(true)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("access_token", nil)
				tg.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).Return("refresh_token", nil)
				tg.EXPECT().DecryptToken(gomock.Any(), gomock.Any()).Return(core.TokenPayload{}, nil)
				tr.EXPECT().Store(gomock.Any(), "access_token", gomock.Any()).Return(nil)
				sr.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(nil)
				mp.EXPECT().SendNewSignInSessionMessage(gomock.Any(), core.SendNewSignInSessionParams{
					Name:      "fingo_user",
					Email:     arg.Email,
					ClientIP:  arg.ClientIP,
					UserAgent: arg.UserAgent,
				}).Return(nil)
			},
			check: func(t *testing.T, response SigninResponse, err error) {
				time.Sleep(1 * time.Second)
				require.Equal(t, response, SigninResponse{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
				})
				require.NoError(t, err)
			},
		},
		{
			name: "success on send new sign in session message error",
			arg: args{
				params: SigninParams{
					Email:     gofakeit.Email(),
					Password:  gofakeit.Password(true, true, true, true, false, 32),
					ClientIP:  gofakeit.IPv4Address(),
					UserAgent: gofakeit.UserAgent(),
				},
			},
			stubs: func(arg SigninParams, v *mock.MockValidator, h *mock.MockPasswordHasher, tg *mock.MockTokenGenerator, ur *mock.MockUserRepository, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository, mp *mock.MockMessageProducer) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				ur.EXPECT().GetUserByEmail(gomock.Any(), arg.Email).Return(core.User{
					FirstName:      "fingo_user",
					Email:          arg.Email,
					HashedPassword: gofakeit.NewCrypto().Password(true, true, true, true, false, 32),
				}, nil)
				h.EXPECT().Compare(gomock.Any(), gomock.Any(), arg.Password).Return(true)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("access_token", nil)
				tg.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).Return("refresh_token", nil)
				tg.EXPECT().DecryptToken(gomock.Any(), gomock.Any()).Return(core.TokenPayload{}, nil)
				tr.EXPECT().Store(gomock.Any(), "access_token", gomock.Any()).Return(nil)
				sr.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(nil)
				mp.EXPECT().SendNewSignInSessionMessage(gomock.Any(), core.SendNewSignInSessionParams{
					Name:      "fingo_user",
					Email:     arg.Email,
					ClientIP:  arg.ClientIP,
					UserAgent: arg.UserAgent,
				}).Return(gofakeit.Error())
			},
			check: func(t *testing.T, response SigninResponse, err error) {
				time.Sleep(1 * time.Second)
				require.Equal(t, response, SigninResponse{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
				})
				require.NoError(t, err)
			},
		},
		{
			name: "validation error",
			arg: args{
				params: SigninParams{},
			},
			stubs: func(arg SigninParams, v *mock.MockValidator, h *mock.MockPasswordHasher, tg *mock.MockTokenGenerator, ur *mock.MockUserRepository, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository, mp *mock.MockMessageProducer) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(gofakeit.Error())
			},
			check: func(t *testing.T, response SigninResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "get user by email error",
			arg: args{
				params: SigninParams{
					Email: gofakeit.Email(),
				},
			},
			stubs: func(arg SigninParams, v *mock.MockValidator, h *mock.MockPasswordHasher, tg *mock.MockTokenGenerator, ur *mock.MockUserRepository, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository, mp *mock.MockMessageProducer) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				ur.EXPECT().GetUserByEmail(gomock.Any(), arg.Email).Return(core.User{}, gofakeit.Error())
			},
			check: func(t *testing.T, response SigninResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "compare password error",
			arg: args{
				params: SigninParams{
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, false, 32),
				},
			},
			stubs: func(arg SigninParams, v *mock.MockValidator, h *mock.MockPasswordHasher, tg *mock.MockTokenGenerator, ur *mock.MockUserRepository, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository, mp *mock.MockMessageProducer) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				ur.EXPECT().GetUserByEmail(gomock.Any(), arg.Email).Return(core.User{
					HashedPassword: gofakeit.NewCrypto().Password(true, true, true, true, false, 32),
				}, nil)
				h.EXPECT().Compare(gomock.Any(), gomock.Any(), arg.Password).Return(false)
			},
			check: func(t *testing.T, response SigninResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "create access token error",
			arg: args{
				params: SigninParams{
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, false, 32),
				},
			},
			stubs: func(arg SigninParams, v *mock.MockValidator, h *mock.MockPasswordHasher, tg *mock.MockTokenGenerator, ur *mock.MockUserRepository, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository, mp *mock.MockMessageProducer) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				ur.EXPECT().GetUserByEmail(gomock.Any(), arg.Email).Return(core.User{
					HashedPassword: gofakeit.NewCrypto().Password(true, true, true, true, false, 32),
				}, nil)
				h.EXPECT().Compare(gomock.Any(), gomock.Any(), arg.Password).Return(true)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("", gofakeit.Error())
			},
			check: func(t *testing.T, response SigninResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "create refresh token error",
			arg: args{
				params: SigninParams{
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, false, 32),
				},
			},
			stubs: func(arg SigninParams, v *mock.MockValidator, h *mock.MockPasswordHasher, tg *mock.MockTokenGenerator, ur *mock.MockUserRepository, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository, mp *mock.MockMessageProducer) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				ur.EXPECT().GetUserByEmail(gomock.Any(), arg.Email).Return(core.User{
					HashedPassword: gofakeit.NewCrypto().Password(true, true, true, true, false, 32),
				}, nil)
				h.EXPECT().Compare(gomock.Any(), gomock.Any(), arg.Password).Return(true)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("access_token", nil)
				tg.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).Return("", gofakeit.Error())
			},
			check: func(t *testing.T, response SigninResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "decrypt token error",
			arg: args{
				params: SigninParams{
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, false, 32),
				},
			},
			stubs: func(arg SigninParams, v *mock.MockValidator, h *mock.MockPasswordHasher, tg *mock.MockTokenGenerator, ur *mock.MockUserRepository, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository, mp *mock.MockMessageProducer) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				ur.EXPECT().GetUserByEmail(gomock.Any(), arg.Email).Return(core.User{
					HashedPassword: gofakeit.NewCrypto().Password(true, true, true, true, false, 32),
				}, nil)
				h.EXPECT().Compare(gomock.Any(), gomock.Any(), arg.Password).Return(true)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("access_token", nil)
				tg.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).Return("refresh_token", nil)
				tg.EXPECT().DecryptToken(gomock.Any(), gomock.Any()).Return(core.TokenPayload{}, gofakeit.Error())
			},
			check: func(t *testing.T, response SigninResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "store token in cache error",
			arg: args{
				params: SigninParams{
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, false, 32),
				},
			},
			stubs: func(arg SigninParams, v *mock.MockValidator, h *mock.MockPasswordHasher, tg *mock.MockTokenGenerator, ur *mock.MockUserRepository, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository, mp *mock.MockMessageProducer) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				ur.EXPECT().GetUserByEmail(gomock.Any(), arg.Email).Return(core.User{
					HashedPassword: gofakeit.NewCrypto().Password(true, true, true, true, false, 32),
				}, nil)
				h.EXPECT().Compare(gomock.Any(), gomock.Any(), arg.Password).Return(true)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("access_token", nil)
				tg.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).Return("refresh_token", nil)
				tg.EXPECT().DecryptToken(gomock.Any(), gomock.Any()).Return(core.TokenPayload{}, nil)
				tr.EXPECT().Store(gomock.Any(), "access_token", gomock.Any()).Return(gofakeit.Error())
			},
			check: func(t *testing.T, response SigninResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "create session in databse error",
			arg: args{
				params: SigninParams{
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, false, 32),
				},
			},
			stubs: func(arg SigninParams, v *mock.MockValidator, h *mock.MockPasswordHasher, tg *mock.MockTokenGenerator, ur *mock.MockUserRepository, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository, mp *mock.MockMessageProducer) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				ur.EXPECT().GetUserByEmail(gomock.Any(), arg.Email).Return(core.User{
					HashedPassword: gofakeit.NewCrypto().Password(true, true, true, true, false, 32),
				}, nil)
				h.EXPECT().Compare(gomock.Any(), gomock.Any(), arg.Password).Return(true)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("access_token", nil)
				tg.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).Return("refresh_token", nil)
				tg.EXPECT().DecryptToken(gomock.Any(), gomock.Any()).Return(core.TokenPayload{}, nil)
				tr.EXPECT().Store(gomock.Any(), "access_token", gomock.Any()).Return(nil)
				sr.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(gofakeit.Error())
			},
			check: func(t *testing.T, response SigninResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
		{
			name: "create session in databse error",
			arg: args{
				params: SigninParams{
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, false, 32),
				},
			},
			stubs: func(arg SigninParams, v *mock.MockValidator, h *mock.MockPasswordHasher, tg *mock.MockTokenGenerator, ur *mock.MockUserRepository, sr *mock.MockSessionRepository, tr *mock.MockTokenRepository, mp *mock.MockMessageProducer) {
				v.EXPECT().Validate(gomock.Any(), arg).Return(nil)
				ur.EXPECT().GetUserByEmail(gomock.Any(), arg.Email).Return(core.User{
					HashedPassword: gofakeit.NewCrypto().Password(true, true, true, true, false, 32),
				}, nil)
				h.EXPECT().Compare(gomock.Any(), gomock.Any(), arg.Password).Return(true)
				tg.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("access_token", nil)
				tg.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).Return("refresh_token", nil)
				tg.EXPECT().DecryptToken(gomock.Any(), gomock.Any()).Return(core.TokenPayload{}, nil)
				tr.EXPECT().Store(gomock.Any(), "access_token", gomock.Any()).Return(nil)
				sr.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(gofakeit.Error())
			},
			check: func(t *testing.T, response SigninResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			v := mock.NewMockValidator(ctrl)
			h := mock.NewMockPasswordHasher(ctrl)
			ur := mock.NewMockUserRepository(ctrl)
			sr := mock.NewMockSessionRepository(ctrl)
			tr := mock.NewMockTokenRepository(ctrl)
			tg := mock.NewMockTokenGenerator(ctrl)
			mp := mock.NewMockMessageProducer(ctrl)

			c := NewSigninCommand(v, h, tg, ur, sr, tr, mp)

			tt.stubs(tt.arg.params, v, h, tg, ur, sr, tr, mp)
			resp, err := c.Execute(context.Background(), tt.arg.params)
			tt.check(t, resp, err)
		})
	}
}
