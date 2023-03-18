package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/escalopa/fingo/auth/internal/adapters/db/redis"
	"github.com/escalopa/fingo/auth/internal/adapters/queue/rabbitmq"
	"github.com/escalopa/fingo/pkg/pkgError"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	mypostgres "github.com/escalopa/fingo/auth/internal/adapters/db/postgres"
	mygrpc "github.com/escalopa/fingo/auth/internal/adapters/grpc"
	"github.com/escalopa/fingo/auth/internal/adapters/hasher"
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
	pkgError.CheckError(err, "invalid access token duration")
	log.Println("successfully parsed access token duration: ", atd)
	rtd, err := time.ParseDuration(c.Get("AUTH_REFRESH_TOKEN_DURATION"))
	pkgError.CheckError(err, "invalid refresh token duration")
	log.Println("successfully parsed refresh token duration: ", rtd)
	tg, err := token.NewPaseto(c.Get("AUTH_TOKEN_SECRET"), atd, rtd)
	pkgError.CheckError(err, "failed to create token generator")
	log.Println("successfully create token generator")

	// Create postgres conn
	pgConn, err := mypostgres.New(c.Get("AUTH_DATABASE_URL"))
	pkgError.CheckError(err, "failed to connect to postgres")
	log.Println("successfully connected to postgres")

	// Migrate database
	err = mypostgres.Migrate(pgConn, c.Get("AUTH_DATABASE_MIGRATION_PATH"))
	pkgError.CheckError(err, "failed to migrate postgres db")
	log.Println("successfully migrated postgres db")

	// Create user repository
	ur, err := mypostgres.NewUserRepository(pgConn)
	pkgError.CheckError(err, "failed to create user repository")
	log.Println("successfully created user repository")

	// Create session repository
	std, err := time.ParseDuration(c.Get("AUTH_USER_SESSION_DURATION"))
	pkgError.CheckError(err, "invalid user session duration")
	log.Println("successfully parsed user session duration: ", std)
	sr, err := mypostgres.NewSessionRepository(pgConn, mypostgres.WithSessionDuration(std))
	pkgError.CheckError(err, "failed to create session repository")
	log.Println("successfully created session repository")

	// Connect to redis cache
	redisConn, err := redis.New(c.Get("AUTH_CACHE_URL"))
	pkgError.CheckError(err, "failed to connect to redis cache")
	log.Println("successfully connected to redis cache")

	// Create token repository
	tr := redis.NewTokenRepository(redisConn, redis.WithTokenDuration(atd))
	log.Println("successfully created token repository")

	// Connect to rabbitmq & Create a new message producer
	rbp, err := rabbitmq.NewProducer(c.Get("AUTH_RABBITMQ_URL"),
		rabbitmq.WithNewSignInSessionQueue(c.Get("AUTH_RABBITMQ_NEW_SIGNIN_SESSION_QUEUE_NAME")),
	)
	pkgError.CheckError(err, "failed to connect to rabbitmq")
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

	// Start the server
	pkgError.CheckError(start(c, uc), "failed to start gRPC server")
}

func start(c *goconfig.Config, uc *application.UseCases) error {
	// Create a new gRPC server
	handler := mygrpc.NewAuthHandler(uc)

	var opts []grpc.ServerOption

	// Enable TLS if required
	enableTLS := c.Get("AUTH_GRPC_TLS_ENABLE") == "true"
	log.Println("starting gRPC server with TLS:", enableTLS)
	if enableTLS {
		// Load TLS certificates
		creds, err := credentials.NewServerTLSFromFile(c.Get("AUTH_GRPC_TLS_CERT_FILE"), c.Get("AUTH_GRPC_TLS_KEY_FILE"))
		if err != nil {
			return errs.B(err).Msg("failed to load TLS certificates").Err()
		}
		opts = append(opts, grpc.Creds(creds))
		log.Println("loaded TLS certificates")
	}

	// Create a new gRPC interceptor
	enableTLS = c.Get("TOKEN_GRPC_TLS_ENABLE") == "true"
	log.Println("connecting to token gRPC server with TLS:", enableTLS)
	var tokenCreds credentials.TransportCredentials
	if enableTLS {
		// Load TLS certificates
		var err error
		tokenCreds, err = credentials.NewClientTLSFromFile(c.Get("AUTH_TOKEN_GRPC_TLS_USER_CERT_FILE"), "")
		if err != nil {
			return errs.B(err).Msg("failed to load TLS certificates").Err()
		}
		log.Println("loaded TLS certificates")
	} else {
		tokenCreds = insecure.NewCredentials()
	}
	interceptor, err := mygrpc.NewAuthInterceptor(c.Get("TOKEN_GRPC_URL"), tokenCreds)
	if err != nil {
		return errs.B(err).Msg("failed to create gRPC interceptor").Err()
	}
	opts = append(opts, grpc.UnaryInterceptor(interceptor.Unary()))

	// Create a new gRPC server
	grpcS := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(grpcS, handler)
	reflection.Register(grpcS)

	// Start the server
	port := c.Get("AUTH_GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return errs.B(err).Msg(fmt.Sprintf("Failed to listen on port %s", port)).Err()
	}
	log.Println("starting gRPC server on port", port)
	err = grpcS.Serve(lis)
	if err != nil {
		return errs.B(err).Msg(fmt.Sprintf("Failed to serve on port %s", port)).Err()
	}
	return nil
}
