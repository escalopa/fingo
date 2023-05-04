package numgen

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNumGen_GenCardNumber(t *testing.T) {
	type fields struct {
		l int
	}
	type args struct {
		ctx func() context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		check  func(t *testing.T, card string, err error)
	}{
		{
			name: "conetxt timeout",
			args: args{ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				return ctx
			}},
			check: func(t *testing.T, s string, err error) {
				require.Error(t, err)
				require.Len(t, s, 0)
			},
		},
		{
			name:   "16",
			fields: fields{l: 16},
			args:   args{ctx: func() context.Context { return context.Background() }},
			check: func(t *testing.T, s string, err error) {
				require.NoError(t, err)
				require.Len(t, s, 16)
			},
		},
		{
			name:   "32",
			fields: fields{l: 32},
			args:   args{ctx: func() context.Context { return context.Background() }},
			check: func(t *testing.T, s string, err error) {
				require.NoError(t, err)
				require.Len(t, s, 32)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNumGen(tt.fields.l)
			got, err := n.GenCardNumber(tt.args.ctx())
			tt.check(t, got, err)
		})
	}
}
