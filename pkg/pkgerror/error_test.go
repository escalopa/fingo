package pkgerror

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{
			name: "success on nil error",
			err:  nil,
			msg:  "success",
		},
		{
			name: "success on non-nil error",
			err:  errors.New("error"),
			msg:  "fail",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					msg, ok := r.(string)
					require.True(t, ok)
					require.Equal(t, fmt.Sprintf("%s: %s", tt.msg, tt.err), msg)
				}
			}()
			CheckError(tt.err, tt.msg)
		})
	}
}
