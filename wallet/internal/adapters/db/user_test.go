package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_CreateUser(t *testing.T) {
	id := uuid.New()

	type args struct {
		uuid uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{uuid: id},
			wantErr: false,
		},
		{
			name:    "duplicate",
			args:    args{uuid: id},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := NewUserRepository(conn)
			err := ur.CreateUser(context.Background(), tt.args.uuid)
			require.Truef(t, (err != nil) == tt.wantErr, "UserRepository.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	id := uuid.New()

	type args struct {
		store uuid.UUID
		get   uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{store: id, get: id},
			wantErr: false,
		},
		{
			name:    "not found",
			args:    args{store: id, get: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := NewUserRepository(conn)
			ur.CreateUser(context.Background(), tt.args.store)
			got, err := ur.GetUser(context.Background(), tt.args.get)
			require.Truef(t, (err != nil) == tt.wantErr, "UserRepository.GetUser() error = %v, wantErr %v", err, tt.wantErr)
			if err == nil {
				require.NotEqual(t, int64(0), got)
			}
		})
	}
}
