package core

import (
	"reflect"
	"testing"
)

func TestTokenCache_MarshalBinary(t *testing.T) {
	tests := []struct {
		name    string
		t       TokenCache
		want    []byte
		wantErr bool
	}{
		{
			name: "success",
			t: TokenCache{
				UserID:    "user-id",
				ClientIP:  "client-ip",
				UserAgent: "user-agent",
				Roles:     []string{"role1", "role2"},
			},
			want:    []byte(`{"UserID":"user-id","ClientIP":"client-ip","UserAgent":"user-agent","Roles":["role1","role2"]}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenCache.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenCache.MarshalBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenCache_UnmarshalBinary(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		t       TokenCache
		args    args
		wantErr bool
	}{
		{
			name: "success",
			t: TokenCache{
				UserID:    "user-id",
				ClientIP:  "client-ip",
				UserAgent: "user-agent",
				Roles:     []string{"role1", "role2"},
			},
			args: args{
				data: []byte(`{"UserID":"user-id","ClientIP":"client-ip","UserAgent":"user-agent","Roles":["role1","role2"]}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("TokenCache.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
