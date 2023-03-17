package main

import (
	"fmt"
	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/fingo/pkg/pkgError"
	"github.com/escalopa/fingo/token/internal/adapters/cache"
	mygrpc "github.com/escalopa/fingo/token/internal/adapters/grpc"
	"github.com/escalopa/fingo/token/internal/adapters/validator"
	"github.com/escalopa/fingo/token/internal/application"
	"github.com/escalopa/goconfig"
	"github.com/lordvidex/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	c := goconfig.New()

	// Create validator
	v := validator.NewValidator()
	log.Println("validator created")

	// Create redis client
	rc, err := cache.NewRedisClient(c.Get("TOKEN_REDIS_URL"))
	pkgError.CheckError(err, "failed to create redis client")
	log.Println("redis client created")

	// Create token repository
	tr, err := cache.NewTokenRepositoryV1(rc)
	pkgError.CheckError(err, "failed to create token repository")
	log.Println("token repository created")

	// Create use cases
	uc := application.NewUseCases(
		application.WithValidator(v),
		application.WithTokenRepository(tr),
	)

	pkgError.CheckError(start(c, uc), "failed to start gRPC server")
}

func start(c *goconfig.Config, uc *application.UseCases) error {
	// Create a new gRPC server with TLS enabled
	handler := mygrpc.NewTokenHandler(uc)

	var opts []grpc.ServerOption

	// Enable TLS if required
	enableTLS := c.Get("TOKEN_GRPC_TLS") == "true"
	log.Println("starting gRPC server with TLS:", enableTLS)
	if enableTLS {
		// Load TLS certificates
		creds, err := credentials.NewServerTLSFromFile(c.Get("TOKEN_GRPC_TLS_CERT"), c.Get("TOKEN_GRPC_TLS_KEY"))
		if err != nil {
			return errs.B(err).Msg("failed to load TLS certificates").Err()
		}
		opts = append(opts, grpc.Creds(creds))
		log.Println("loaded TLS certificates")
	}

	// Create a gRPC server object
	grpcS := grpc.NewServer(opts...)
	pb.RegisterTokenServiceServer(grpcS, handler)
	reflection.Register(grpcS)
	log.Println("gRPC server instance created")

	// Start the server
	port := c.Get("TOKEN_GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return errs.B(err).Msg(fmt.Sprintf("failed to listen on port %s", port)).Err()
	}
	log.Println("starting gRPC server on port", port)
	err = grpcS.Serve(lis)
	if err != nil {
		return errs.B(err).Msg(fmt.Sprintf("failed to serve on port %s", port)).Err()
	}
	return nil
}
