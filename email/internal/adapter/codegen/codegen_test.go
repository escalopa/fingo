package codegen

import (
	"github.com/stretchr/testify/require"
	"testing"
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
			cg := New(tc.length)
			code, err := cg.GenerateCode()
			require.NoError(t, err)
			require.Len(t, code, tc.length)
			require.True(t, cg.VerifyCode(code))
		})
	}
}
