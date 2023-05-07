package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/fingo/pkg/global"
	"github.com/escalopa/fingo/pkg/tls"
	mygrpc "github.com/escalopa/fingo/token/internal/adapters/grpc"
	"github.com/escalopa/fingo/token/internal/application"
	"github.com/escalopa/goconfig"
	"github.com/lordvidex/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func start(appCtx context.Context, c *goconfig.Config, uc *application.UseCases) error {
	// Load TLS certificates
	var opts []grpc.ServerOption
	err := loadTls(c, &opts)
	if err != nil {
		log.Println("failed to load token tls certificates")
		return err
	}

	// Create a gRPC server object
	handler := mygrpc.NewTokenHandler(uc)
	server := grpc.NewServer(opts...)
	pb.RegisterTokenServiceServer(server, handler)
	reflection.Register(server)

	// Start the server
	port := c.Get("TOKEN_GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return errs.B(err).Msg(fmt.Sprintf("failed to listen on port %s", port)).Err()
	}

	// Terminate server on shutdown signals
	go global.Shutdown(appCtx, 10*time.Second, func() { server.GracefulStop() }, func() { server.Stop() })

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
		c.Get("TOKEN_GRPC_TLS_ENABLE") == "true",
		c.Get("TOKEN_GRPC_TLS_CERT_FILE"),
		c.Get("TOKEN_GRPC_TLS_KEY_FILE"),
	)
	if err != nil {
		return err
	}
	*opts = append(*opts, grpc.Creds(creds))
	return nil
}
