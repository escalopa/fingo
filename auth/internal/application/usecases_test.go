package application

import (
	"testing"

	"github.com/escalopa/fingo/auth/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewUseCases(t *testing.T) {
	type args struct {
		opts []func(*UseCases)
	}
	tests := []struct {
		name  string
		args  args
		stubs func(
			opts *[]func(*UseCases),
			v *mock.MockValidator,
			h *mock.MockPasswordHasher,
			tg *mock.MockTokenGenerator,
			ur *mock.MockUserRepository,
			sr *mock.MockSessionRepository,
			tr *mock.MockTokenRepository,
			mp *mock.MockMessageProducer,
		)
		check func(t *testing.T, uc *UseCases)
	}{
		{
			name: "success",
			args: args{
				opts: []func(*UseCases){},
			},
			stubs: func(
				opts *[]func(*UseCases),
				v *mock.MockValidator,
				h *mock.MockPasswordHasher,
				tg *mock.MockTokenGenerator,
				ur *mock.MockUserRepository,
				sr *mock.MockSessionRepository,
				tr *mock.MockTokenRepository,
				mp *mock.MockMessageProducer,
			) {
				*opts = append(*opts, []func(*UseCases){
					WithValidator(v),
					WithPasswordHasher(h),
					WithTokenGenerator(tg),
					WithUserRepository(ur),
					WithSessionRepository(sr),
					WithTokenRepository(tr),
					WithMessageProducer(mp),
				}...)
			},
			check: func(t *testing.T, uc *UseCases) {
				require.NotNil(t, uc)
				require.NotNil(t, uc.v)
				require.NotNil(t, uc.h)
				require.NotNil(t, uc.tg)
				require.NotNil(t, uc.ur)
				require.NotNil(t, uc.sr)
				require.NotNil(t, uc.tr)
				require.NotNil(t, uc.mp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			v := mock.NewMockValidator(ctrl)
			h := mock.NewMockPasswordHasher(ctrl)
			tg := mock.NewMockTokenGenerator(ctrl)
			ur := mock.NewMockUserRepository(ctrl)
			sr := mock.NewMockSessionRepository(ctrl)
			tr := mock.NewMockTokenRepository(ctrl)
			mp := mock.NewMockMessageProducer(ctrl)

			tt.stubs(&tt.args.opts, v, h, tg, ur, sr, tr, mp)
			tt.check(t, NewUseCases(tt.args.opts...))
		})
	}
}
