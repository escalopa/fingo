package contextutils

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func TestExecuteWithContextTimeout(t *testing.T) {
	type args struct {
		ctx     context.Context
		timeout time.Duration
		handler func() error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sucess",
			args: args{
				ctx:     context.Background(),
				timeout: 1 * time.Second,
				handler: func() error {
					time.Sleep(500 * time.Millisecond) // Do some work
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "timeout",
			args: args{
				ctx:     context.Background(),
				timeout: 1 * time.Second,
				handler: func() error {
					time.Sleep(2 * time.Second) // Do some work
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "handler error",
			args: args{
				ctx:     context.Background(),
				timeout: 1 * time.Second,
				handler: func() error {
					return gofakeit.Error()
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExecuteWithContextTimeout(tt.args.ctx, tt.args.timeout, tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("ExecuteWithContextTimeout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
