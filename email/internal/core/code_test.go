package core

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestVerificationCode(t *testing.T) {
	vc := VerificationCode{
		Code:   string(rune(gofakeit.IntRange(10_000, 20_000))),
		SentAt: time.Now(),
	}
	b, err := vc.MarshalBinary()
	require.NoError(t, err)
	require.NotEmpty(t, b)
}
