package validator

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidator(t *testing.T) {
	v := NewValidator()

	testCases := []struct {
		name  string
		tv    TestValidatorStruct
		isErr bool
	}{
		{
			name: "valid",
			tv: TestValidatorStruct{
				Name:     gofakeit.Name(),
				Email:    gofakeit.Email(),
				Password: gofakeit.Password(true, true, true, true, false, 8),
			},
			isErr: false,
		}, {
			name: "empty name",
			tv: TestValidatorStruct{
				Name:     "",
				Email:    gofakeit.Email(),
				Password: gofakeit.Password(true, true, true, true, false, 8),
			},
			isErr: true,
		}, {
			name: "invalid email 1",
			tv: TestValidatorStruct{
				Name:     gofakeit.Name(),
				Email:    "ahmad@gmail",
				Password: gofakeit.Password(true, true, true, true, false, 8),
			},
			isErr: true,
		}, {
			name: "invalid email 2",
			tv: TestValidatorStruct{
				Name:     gofakeit.Name(),
				Email:    "ahmad@.com",
				Password: gofakeit.Password(true, true, true, true, false, 8),
			},
			isErr: true,
		}, {
			name: "invalid email 3",
			tv: TestValidatorStruct{
				Name:     gofakeit.Name(),
				Email:    "@gmail.com",
				Password: gofakeit.Password(true, true, true, true, false, 8),
			},
			isErr: true,
		}, {
			name: "invalid password",
			tv: TestValidatorStruct{
				Name:     gofakeit.Name(),
				Email:    gofakeit.Email(),
				Password: gofakeit.Password(true, true, true, true, false, 7),
			},
			isErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := v.Validate(tc.tv)
			require.True(t, (err != nil) == tc.isErr)
		})
	}
}
