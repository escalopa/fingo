package global

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestCheckError(t *testing.T) {
	type args struct {
		err     error
		msg     string
		wantErr bool
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				err:     nil,
				msg:     "success",
				wantErr: false,
			},
		},
		{
			name: "success on error",
			args: args{
				err:     gofakeit.Error(),
				msg:     "fail",
				wantErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				x := recover()
				require.Equal(t, tt.args.wantErr, x != nil)
			}()
			CheckError(tt.args.err, tt.args.msg)
		})
	}
}
