package main

import (
	"fmt"
	"log"
	"net"
	"time"

	oteltracer "github.com/escalopa/fingo/auth/internal/adapters/tracer"
	"github.com/escalopa/fingo/pkg/grpctls"
	"github.com/escalopa/fingo/pkg/pkgerror"
	"github.com/escalopa/fingo/pkg/pkgtracer"

	mypostgres "github.com/escalopa/fingo/auth/internal/adapters/db/postgres"
	"github.com/escalopa/fingo/auth/internal/adapters/db/redis"
	mygrpc "github.com/escalopa/fingo/auth/internal/adapters/grpc"
	"github.com/escalopa/fingo/auth/internal/adapters/hasher"
	"github.com/escalopa/fingo/auth/internal/adapters/queue/rabbitmq"
	"github.com/escalopa/fingo/auth/internal/adapters/token"
	myvalidator "github.com/escalopa/fingo/auth/internal/adapters/validator"
	"github.com/escalopa/fingo/auth/internal/application"
	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/goconfig"
	"github.com/lordvidex/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	c := goconfig.New()

	ph := hasher.NewBcryptHasher()
	v := myvalidator.NewValidator()

	// Create a new token generator
	atd, err := time.ParseDuration(c.Get("AUTH_ACCESS_TOKEN_DURATION"))
	pkgerror.CheckError(err, "invalid access token duration")
	log.Println("successfully parsed access token duration: ", atd)
	rtd, err := time.ParseDuration(c.Get("AUTH_REFRESH_TOKEN_DURATION"))
	pkgerror.CheckError(err, "invalid refresh token duration")
	log.Println("successfully parsed refresh token duration: ", rtd)
	tg, err := token.NewPaseto(c.Get("AUTH_TOKEN_SECRET"), atd, rtd)
	pkgerror.CheckError(err, "failed to create token generator")
	log.Println("successfully create token generator")

	// Create postgres conn
	pgConn, err := mypostgres.New(c.Get("AUTH_DATABASE_URL"))
	pkgerror.CheckError(err, "failed to connect to postgres")
	log.Println("successfully connected to postgres")

	// Migrate database
	err = mypostgres.Migrate(pgConn, c.Get("AUTH_DATABASE_MIGRATION_PATH"))
	pkgerror.CheckError(err, "failed to migrate postgres db")
	log.Println("successfully migrated postgres db")

	// Create user repository
	ur, err := mypostgres.NewUserRepository(pgConn)
	pkgerror.CheckError(err, "failed to create user repository")
	log.Println("successfully created user repository")

	// Create session repository
	std, err := time.ParseDuration(c.Get("AUTH_USER_SESSION_DURATION"))
	pkgerror.CheckError(err, "invalid user session duration")
	log.Println("successfully parsed user session duration: ", std)
	sr, err := mypostgres.NewSessionRepository(pgConn, mypostgres.WithSessionDuration(std))
	pkgerror.CheckError(err, "failed to create session repository")
	log.Println("successfully created session repository")

	// Connect to redis cache
	redisConn, err := redis.New(c.Get("AUTH_CACHE_URL"))
	pkgerror.CheckError(err, "failed to connect to redis cache")
	log.Println("successfully connected to redis cache")

	// Create token repository
	tr := redis.NewTokenRepository(redisConn, redis.WithTokenDuration(atd))
	log.Println("successfully created token repository")

	// Connect to rabbitmq & Create a new message producer
	rbp, err := rabbitmq.NewProducer(c.Get("AUTH_RABBITMQ_URL"),
		rabbitmq.WithNewSignInSessionQueue(c.Get("AUTH_RABBITMQ_NEW_SIGNIN_SESSION_QUEUE_NAME")),
	)
	pkgerror.CheckError(err, "failed to connect to rabbitmq")
	log.Println("successfully connected to rabbitmq")

	// Create a new use case
	uc := application.NewUseCases(
		application.WithValidator(v),
		application.WithPasswordHasher(ph),
		application.WithTokenGenerator(tg),
		application.WithUserRepository(ur),
		application.WithSessionRepository(sr),
		application.WithTokenRepository(tr),
		application.WithMessageProducer(rbp),
	)

	var opts []grpc.ServerOption
	// Load TLS certificates
	pkgerror.CheckError(loadTls(c, &opts), "failed to load auth TLS certificates")

	// Load auth interceptor
	pkgerror.CheckError(loadInterceptor(c, &opts), "failed to load auth interceptor")

	// Create a new tracer
	t, err := pkgtracer.LoadTracer(
		c.Get("AUTH_TRACING_ENABLE"),
		c.Get("AUTH_TRACING_JAEGER_ENABLE"),
		c.Get("AUTH_TRACING_JAEGER_AGENT_URL"),
		c.Get("AUTH_TRACING_JAEGER_SERVICE_NAME"),
		c.Get("AUTH_TRACING_JAEGER_ENVIRONMENT"),
	)
	pkgerror.CheckError(err, "failed to load tracer")
	oteltracer.SetTracer(t)

	// Start the server
	pkgerror.CheckError(start(c, uc, opts), "failed to start gRPC server")
}

func start(c *goconfig.Config, uc *application.UseCases, opts []grpc.ServerOption) error {
	// Create a new gRPC server
	handler := mygrpc.NewAuthHandler(uc)
	server := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(server, handler)
	reflection.Register(server)

	// Start the server
	port := c.Get("AUTH_GRPC_PORT")
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

func loadTls(c *goconfig.Config, opts *[]grpc.ServerOption) error {
	// Enable TLS if required
	creds, err := grpctls.LoadServerTLS(
		c.Get("AUTH_GRPC_TLS_ENABLE"),
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
	creds, err := grpctls.LoadClientTLS(
		c.Get("TOKEN_GRPC_TLS_ENABLE"),
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
