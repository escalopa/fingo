package application

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

func TestSendNewSingInSessionCommandImpl_Execute(t *testing.T) {
	testCases := []struct {
		name        string
		param       SendNewSingInSessionCommandParam
		expectError bool
	}{
		{
			name: "valid",
			param: SendNewSingInSessionCommandParam{
				Name:      gofakeit.FirstName(),
				Email:     gofakeit.Email(),
				ClientIP:  gofakeit.IPv4Address(),
				UserAgent: gofakeit.UserAgent(),
			},
			expectError: false,
		},
		{
			name: "invalid name",
			param: SendNewSingInSessionCommandParam{
				Name:      "",
				Email:     gofakeit.Email(),
				ClientIP:  gofakeit.IPv4Address(),
				UserAgent: gofakeit.UserAgent(),
			},
			expectError: true,
		},
		{
			name: "invalid email",
			param: SendNewSingInSessionCommandParam{
				Name:      gofakeit.FirstName(),
				Email:     "invalid",
				ClientIP:  gofakeit.IPv4Address(),
				UserAgent: gofakeit.UserAgent(),
			},
			expectError: true,
		},
		{
			name: "invalid client ip",
			param: SendNewSingInSessionCommandParam{
				Name:      gofakeit.FirstName(),
				Email:     gofakeit.Email(),
				ClientIP:  "",
				UserAgent: gofakeit.UserAgent(),
			},
			expectError: true,
		},
		{
			name: "invalid user agent",
			param: SendNewSingInSessionCommandParam{
				Name:      gofakeit.FirstName(),
				Email:     gofakeit.Email(),
				ClientIP:  gofakeit.IPv4Address(),
				UserAgent: "",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// execute command
			err := testUseCases.SendNewSignInSession.Execute(context.Background(), tc.param)
			if (err != nil) != tc.expectError {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
