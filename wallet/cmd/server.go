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

	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/fingo/pkg/global"
	"github.com/escalopa/fingo/pkg/interceptors"
	"github.com/escalopa/fingo/pkg/tls"
	mygrpc "github.com/escalopa/fingo/wallet/internal/adapters/grpc"
	"github.com/escalopa/fingo/wallet/internal/application"
)

func start(appCtx context.Context, uc *application.UseCases) error {
	opts := []grpc.ServerOption{grpc.ChainUnaryInterceptor(
		interceptors.TracingUnaryInterceptor(),
		interceptors.LoggingUnaryInterceptor(),
	)} // Load TLS certificates
	err := loadTls(&opts)
	global.CheckError(err, "failed to load wallet TLS certificates")

	// Load auth interceptor
	err = loadInterceptor(&opts)
	global.CheckError(err, "failed to load auth interceptor")

	// Create a gRPC server object
	handler := mygrpc.NewWalletHandler(uc)
	server := grpc.NewServer(opts...)
	pb.RegisterWalletServiceServer(server, handler)
	reflection.Register(server)

	// Terminate server on shutdown signals
	go global.Shutdown(appCtx, 10*time.Second, func() { server.GracefulStop() }, func() { server.Stop() })

	// Start the server
	port := cfg.GrpcPort
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return errs.B(err).Msg(fmt.Sprintf("failed to listen on port %s", port)).Err()
	}
	log.Println("starting gRPC server on port", port)
	err = server.Serve(lis)
	if err != nil {
		return errs.B(err).Msg(fmt.Sprintf("failed to serve on port %s", port)).Err()
	}
	return nil
}

func loadTls(opts *[]grpc.ServerOption) error {
	// Enable TLS if required
	creds, err := tls.LoadServerTLS(
		cfg.GrpcTlsEnable,
		cfg.GrpcTlsCertFile,
		cfg.GrpcTlsKeyFile,
	)
	if err != nil {
		return err
	}
	*opts = append(*opts, grpc.Creds(creds))
	return nil
}

func loadInterceptor(opts *[]grpc.ServerOption) error {
	creds, err := tls.LoadClientTLS(
		cfg.TokenGrpcTlsEnable,
		cfg.TokenGrpcTlsUserCertFile,
	)
	if err != nil {
		return err
	}
	interceptor, err := mygrpc.NewAuthInterceptor(cfg.TokenGrpcUrl, creds)
	if err != nil {
		return errs.B(err).Msg("failed to create token gRPC interceptor").Err()
	}
	*opts = append(*opts, grpc.UnaryInterceptor(interceptor.Unary()))
	log.Println("created gRPC token interceptor")
	return nil
}
