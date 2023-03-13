package application

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

func TestSendVerificationCodeCommandImpl_Execute(t *testing.T) {

	testCases := []struct {
		name        string
		params      SendVerificationCodeCommandParam
		expectError bool
	}{
		{
			name: "valid",
			params: SendVerificationCodeCommandParam{
				Name:  gofakeit.FirstName(),
				Email: gofakeit.Email(),
				Code:  gofakeit.UUID(),
			},
			expectError: false,
		},
		{
			name: "invalid name",
			params: SendVerificationCodeCommandParam{
				Name:  "",
				Email: gofakeit.Email(),
				Code:  gofakeit.UUID(),
			},
			expectError: true,
		},
		{
			name: "invalid email",
			params: SendVerificationCodeCommandParam{
				Name:  gofakeit.FirstName(),
				Email: "invalid",
				Code:  gofakeit.UUID(),
			},
			expectError: true,
		},
		{
			name: "invalid code",
			params: SendVerificationCodeCommandParam{
				Name:  gofakeit.FirstName(),
				Email: gofakeit.Email(),
				Code:  "",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// execute command
			err := testUseCases.SendVerificationCode.Execute(context.Background(), SendVerificationCodeCommandParam{
				Name:  tc.params.Name,
				Email: tc.params.Email,
				Code:  tc.params.Code,
			})
			if (err != nil) != tc.expectError {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
