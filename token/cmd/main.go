package main

import (
	"fmt"
	"log"
	"net"

	pkgvalidator "github.com/escalopa/fingo/pkg/validator"

	pkgerror "github.com/escalopa/fingo/pkg/error"
	grpctls "github.com/escalopa/fingo/pkg/tls"
	pkgtracer "github.com/escalopa/fingo/pkg/tracer"

	oteltracer "github.com/escalopa/fingo/token/internal/adapters/tracer"

	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/fingo/token/internal/adapters/cache"
	mygrpc "github.com/escalopa/fingo/token/internal/adapters/grpc"
	"github.com/escalopa/fingo/token/internal/application"
	"github.com/escalopa/goconfig"
	"github.com/lordvidex/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	c := goconfig.New()

	// Create validator
	v := pkgvalidator.NewValidator()
	log.Println("validator created")

	// Create redis client
	rc, err := cache.NewRedisClient(c.Get("TOKEN_REDIS_URL"))
	pkgerror.CheckError(err, "failed to create redis client")
	log.Println("redis client created")

	// Create token repository
	tr := cache.NewTokenRepositoryV1(rc)
	log.Println("token repository created")

	// Create use cases
	uc := application.NewUseCases(
		application.WithValidator(v),
		application.WithTokenRepository(tr),
	)

	// Create a new tracer
	t, err := pkgtracer.LoadTracer(
		c.Get("TOKEN_TRACING_ENABLE"),
		c.Get("TOKEN_TRACING_JAEGER_ENABLE"),
		c.Get("TOKEN_TRACING_JAEGER_AGENT_URL"),
		c.Get("TOKEN_TRACING_JAEGER_SERVICE_NAME"),
		c.Get("TOKEN_TRACING_JAEGER_ENVIRONMENT"),
	)
	pkgerror.CheckError(err, "failed to load tracer")
	oteltracer.SetTracer(t)
	log.Println("tracer created")

	// Load TLS certificates
	var opts []grpc.ServerOption
	err = loadTls(c, &opts)
	pkgerror.CheckError(err, "failed to load token tls certificates")

	// Start gRPC server
	pkgerror.CheckError(start(c, uc, opts), "failed to start gRPC server")
}

func start(c *goconfig.Config, uc *application.UseCases, opts []grpc.ServerOption) error {
	// Create a gRPC server object
	handler := mygrpc.NewTokenHandler(uc)
	grpcS := grpc.NewServer(opts...)
	pb.RegisterTokenServiceServer(grpcS, handler)
	reflection.Register(grpcS)

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

func loadTls(c *goconfig.Config, opts *[]grpc.ServerOption) error {
	// Enable TLS if required
	creds, err := grpctls.LoadServerTLS(
		c.Get("TOKEN_GRPC_TLS_ENABLE"),
		c.Get("TOKEN_GRPC_TLS_CERT_FILE"),
		c.Get("TOKEN_GRPC_TLS_KEY_FILE"),
	)
	if err != nil {
		return err
	}
	*opts = append(*opts, grpc.Creds(creds))
	return nil
}
