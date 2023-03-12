package application

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

func TestSendVerificationCodeCommandImpl_Execute(t *testing.T) {

	testCases := []struct {
		name        string
		param       SendVerificationCodeCommandParam
		expectError bool
	}{
		{
			name: "valid",
			param: SendVerificationCodeCommandParam{
				Name:  gofakeit.FirstName(),
				Email: gofakeit.Email(),
				Code:  gofakeit.UUID(),
			},
			expectError: false,
		},
		{
			name: "invalid name",
			param: SendVerificationCodeCommandParam{
				Name:  "",
				Email: gofakeit.Email(),
				Code:  gofakeit.UUID(),
			},
			expectError: true,
		},
		{
			name: "invalid email",
			param: SendVerificationCodeCommandParam{
				Name:  gofakeit.FirstName(),
				Email: "invalid",
				Code:  gofakeit.UUID(),
			},
			expectError: true,
		},
		{
			name: "invalid code",
			param: SendVerificationCodeCommandParam{
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
			err := testUseCases.SendVerificationCode.Execute(context.Background(), tc.param)
			if (err != nil) != tc.expectError {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
