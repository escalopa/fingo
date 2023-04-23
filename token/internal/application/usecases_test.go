package application

import (
	"testing"

	"github.com/escalopa/fingo/token/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := mock.NewMockValidator(ctrl)
	tr := mock.NewMockTokenRepository(ctrl)

	tests := []struct {
		name  string
		opts  []func(*UseCases)
		check func(*testing.T, *UseCases)
	}{
		{
			name: "success nil",
			opts: []func(*UseCases){
				WithTokenRepository(nil),
				WithValidator(nil),
			},
			check: func(t *testing.T, uc *UseCases) {
				require.NotNil(t, uc)
				require.Nil(t, uc.v)
				require.Nil(t, uc.tr)
			},
		},
		{
			name: "success mock",
			opts: []func(*UseCases){
				WithTokenRepository(tr),
				WithValidator(v),
			},
			check: func(t *testing.T, uc *UseCases) {
				require.NotNil(t, uc)
				require.Equal(t, uc.v, v)
				require.Equal(t, uc.tr, tr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUseCases(tt.opts...)
			tt.check(t, uc)
		})
	}
}
