package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/lordvidex/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	mygrpc "github.com/escalopa/fingo/auth/internal/adapters/grpc"
	"github.com/escalopa/fingo/auth/internal/application"
	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/fingo/pkg/global"
	"github.com/escalopa/fingo/pkg/tls"
	"github.com/escalopa/goconfig"
)

func start(appCtx context.Context, c *goconfig.Config, uc *application.UseCases) error {
	var opts []grpc.ServerOption
	// Load TLS certificates
	err := loadTls(c, &opts)
	global.CheckError(err, "failed to load auth TLS certificates")

	// Load auth interceptor
	global.CheckError(loadInterceptor(c, &opts), "failed to load auth interceptor")

	// Create a new gRPC server
	handler := mygrpc.NewAuthHandler(uc)
	server := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(server, handler)
	reflection.Register(server)

	// Terminate server on shutdown signals
	go global.Shutdown(appCtx, 10*time.Second, func() { server.GracefulStop() }, func() { server.Stop() })

	// Start the server
	port := c.Get("AUTH_GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return errs.B(err).Msg(fmt.Sprintf("failed to listen on port %s", port)).Err()
	}
	log.Println("starting gRPC server on port", port)
	err = server.Serve(lis)
	if err != nil {
		if err != grpc.ErrServerStopped {
			return errs.B(err).Msg(fmt.Sprintf("failed to serve on port %s", port)).Err()
		}
	}
	return nil
}

func loadTls(c *goconfig.Config, opts *[]grpc.ServerOption) error {
	// Enable TLS if required
	creds, err := tls.LoadServerTLS(
		c.Get("AUTH_GRPC_TLS_ENABLE") == "true",
		c.Get("AUTH_GRPC_TLS_CERT_FILE"),
		c.Get("AUTH_GRPC_TLS_KEY_FILE"),
	)
	if err != nil {
		return err
	}
	*opts = append(*opts, grpc.Creds(creds))
	return nil
}

func loadInterceptor(c *goconfig.Config, opts *[]grpc.ServerOption) error {
	creds, err := tls.LoadClientTLS(
		c.Get("TOKEN_GRPC_TLS_ENABLE") == "true",
		c.Get("AUTH_TOKEN_GRPC_TLS_USER_CERT_FILE"),
	)
	if err != nil {
		return err
	}
	interceptor, err := mygrpc.NewAuthInterceptor(c.Get("TOKEN_GRPC_URL"), creds)
	if err != nil {
		return errs.B(err).Msg("failed to create token gRPC interceptor").Err()
	}
	*opts = append(*opts, grpc.UnaryInterceptor(interceptor.Unary()))
	log.Println("created gRPC token interceptor")
	return nil
}
