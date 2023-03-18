package application

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

func TestSendResetPasswordTokenCommandImpl_Execute(t *testing.T) {
	testCases := []struct {
		name        string
		params      SendResetPasswordTokenCommandParam
		expectError bool
	}{
		{
			name: "valid",
			params: SendResetPasswordTokenCommandParam{
				Name:  gofakeit.FirstName(),
				Email: gofakeit.Email(),
				Token: gofakeit.UUID(),
			},
			expectError: false,
		},
		{
			name: "invalid name",
			params: SendResetPasswordTokenCommandParam{
				Name:  "",
				Email: gofakeit.Email(),
				Token: gofakeit.UUID(),
			},
			expectError: true,
		},
		{
			name: "invalid email",
			params: SendResetPasswordTokenCommandParam{
				Name:  gofakeit.FirstName(),
				Email: "invalid",
				Token: gofakeit.UUID(),
			},
			expectError: true,
		},
		{
			name: "invalid token",
			params: SendResetPasswordTokenCommandParam{
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
			err := testUseCases.SendResetPasswordToken.Execute(context.Background(), SendResetPasswordTokenCommandParam{
				Name:  tc.params.Name,
				Email: tc.params.Email,
				Token: tc.params.Token,
			})
			if (err != nil) != tc.expectError {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
