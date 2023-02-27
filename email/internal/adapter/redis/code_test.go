package redis

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/gochat/email/internal/core"
	"github.com/lordvidex/errs"
	"github.com/stretchr/testify/require"
)

func TestCodeRepository(t *testing.T) {
	exp := 3 * time.Second
	cr := NewCodeRepository(redisClient, WithExpiration(exp))

	testCases := []struct {
		name  string
		code  string
		email string
	}{
		{
			name:  "Save and get code",
			code:  gofakeit.RandomString([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}),
			email: gofakeit.Email(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save code
			vc := core.VerificationCode{
				Code:   tc.code,
				SentAt: time.Now(),
			}
			err := cr.Save(testContext, tc.email, vc)
			require.NoError(t, err, "error saving code from cache")
			// Get code
			vcGet, err := cr.Get(testContext, tc.email)
			require.NoError(t, err, "error getting code from cache")
			require.Equalf(t, vc.Code, vcGet.Code,
				"stored value doesn't match one set expected: %+v, got: %+v", vc.Code, vcGet.Code)
			// Wait for code to expire
			time.Sleep(exp)
			// Get code again (should return redis.Nil) since it has expired
			_, err = cr.Get(testContext, tc.email)
			if err != nil {
				if er, ok := err.(*errs.Error); ok {
					require.Equal(t, errs.NotFound, er.Code)
				} else {
					t.Fatalf("unsupported error type: expected %T, got %T", &errs.Error{}, er)
				}
			}
		})
	}
}
