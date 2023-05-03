package application

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/auth/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSignupCommand_Execute(t *testing.T) {
	tests := []struct {
		name  string
		arg   SignupParams
		stubs func(args SignupParams, v *mock.MockValidator, ph *mock.MockPasswordHasher, ur *mock.MockUserRepository)
		check func(t *testing.T, err error)
	}{
		{
			name: "success",
			arg: SignupParams{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Username:  gofakeit.Username(),
				Email:     gofakeit.Email(),
				Password:  gofakeit.Password(true, true, true, false, false, 8),
			},
			stubs: func(args SignupParams, v *mock.MockValidator, ph *mock.MockPasswordHasher, ur *mock.MockUserRepository) {
				v.EXPECT().Validate(gomock.Any(), args).Return(nil)
				ph.EXPECT().Hash(gomock.Any(), args.Password).Return("hashedPassword", nil)
				ur.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)
			},
			check: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "validation error",
			arg: SignupParams{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Username:  gofakeit.Username(),
				Email:     gofakeit.Date().String(), // invalid email
				Password:  gofakeit.Password(true, true, true, false, false, 8),
			},
			stubs: func(args SignupParams, v *mock.MockValidator, ph *mock.MockPasswordHasher, ur *mock.MockUserRepository) {
				v.EXPECT().Validate(gomock.Any(), args).Return(gofakeit.Error())
			},
			check: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "hash error",
			arg: SignupParams{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Username:  gofakeit.Username(),
				Email:     gofakeit.Email(),
				Password:  gofakeit.Password(true, true, true, false, false, 7),
			},
			stubs: func(args SignupParams, v *mock.MockValidator, ph *mock.MockPasswordHasher, ur *mock.MockUserRepository) {
				v.EXPECT().Validate(gomock.Any(), args).Return(nil)
				ph.EXPECT().Hash(gomock.Any(), args.Password).Return("", gofakeit.Error())
			},
			check: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "failed to save user",
			arg: SignupParams{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Username:  gofakeit.Username(),
				Email:     gofakeit.Email(),
				Password:  gofakeit.Password(true, true, true, false, false, 8),
			},
			stubs: func(args SignupParams, v *mock.MockValidator, ph *mock.MockPasswordHasher, ur *mock.MockUserRepository) {
				v.EXPECT().Validate(gomock.Any(), args).Return(nil)
				ph.EXPECT().Hash(gomock.Any(), args.Password).Return("hashedPassword", nil)
				ur.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(gofakeit.Error())
			},
			check: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			v := mock.NewMockValidator(ctrl)
			ur := mock.NewMockUserRepository(ctrl)
			ph := mock.NewMockPasswordHasher(ctrl)

			c := NewSignupCommand(v, ph, ur)
			tt.stubs(tt.arg, v, ph, ur)
			err := c.Execute(context.Background(), tt.arg)
			tt.check(t, err)
		})
	}
}
