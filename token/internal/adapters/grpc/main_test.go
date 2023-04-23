package mygrpc

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/fingo/token/internal/application"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func setup(t *testing.T, v application.Validator, tr application.TokenRepository) *grpc.ClientConn {
	uc := application.NewUseCases(
		application.WithTokenRepository(tr),
		application.WithValidator(v),
	)

	h := NewTokenHandler(uc)
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterTokenServiceServer(s, h)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	require.NoError(t, err)
	return conn
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
