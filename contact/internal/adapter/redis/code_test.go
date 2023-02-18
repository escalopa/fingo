package redis

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/lordvidex/errs"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCodeRepository(t *testing.T) {
	exp := 3 * time.Second
	cr := NewCodeRepository(testRedis, WithCodeContext(testContext), WithExpiration(exp))

	testCases := []struct {
		name   string
		code   string
		userID string
	}{
		{
			name:   "Save and get code",
			code:   gofakeit.RandomString([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}),
			userID: gofakeit.UUID(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save code
			err := cr.Save(tc.code, tc.userID)
			require.NoError(t, err, fmt.Sprintf("error saving code: %v", err))
			// Get code
			userID, err := cr.Get(tc.code)
			require.NoError(t, err, fmt.Sprintf("error getting code: %v", err))
			require.NoError(t, err, fmt.Sprintf("expected userID %s, got %s", tc.userID, userID))
			// Wait for code to expire
			time.Sleep(exp)
			// Get code again (should return redis.Nil) since it has expired
			_, err = cr.Get(tc.code)
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
