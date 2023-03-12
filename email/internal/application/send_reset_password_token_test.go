package application

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

func TestSendResetPasswordTokenCommandImpl_Execute(t *testing.T) {
	testCases := []struct {
		name        string
		param       SendResetPasswordTokenCommandParam
		expectError bool
	}{
		{
			name: "valid",
			param: SendResetPasswordTokenCommandParam{
				Name:  gofakeit.FirstName(),
				Email: gofakeit.Email(),
				Token: gofakeit.UUID(),
			},
			expectError: false,
		},
		{
			name: "invalid name",
			param: SendResetPasswordTokenCommandParam{
				Name:  "",
				Email: gofakeit.Email(),
				Token: gofakeit.UUID(),
			},
			expectError: true,
		},
		{
			name: "invalid email",
			param: SendResetPasswordTokenCommandParam{
				Name:  gofakeit.FirstName(),
				Email: "invalid",
				Token: gofakeit.UUID(),
			},
			expectError: true,
		},
		{
			name: "invalid token",
			param: SendResetPasswordTokenCommandParam{
				Name:  gofakeit.FirstName(),
				Email: gofakeit.Email(),
				Token: "",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// execute command
			err := testUseCases.SendResetPasswordToken.Execute(context.Background(), tc.param)
			if (err != nil) != tc.expectError {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
