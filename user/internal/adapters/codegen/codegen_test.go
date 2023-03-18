package codegen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCodeGen(t *testing.T) {
	testCases := []struct {
		name   string
		length int
	}{
		{"test1", 10},
		{"test2", 20},
		{"test3", 30},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cg, err := New(tc.length)
			require.NoError(t, err)
			code, err := cg.GenerateCode()
			require.NoError(t, err)
			require.Len(t, code, tc.length)
			require.Truef(t, cg.VerifyCode(code), "Failed to verify code")
		})
	}
}
