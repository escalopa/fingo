package main

import (
	"fmt"
	"github.com/escalopa/fingo/auth/internal/adapters/db/redis"
	"github.com/escalopa/fingo/auth/internal/adapters/queue/rabbitmq"
	"log"
	"net"
	"time"

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
	checkError(err, "invalid access token duration")
	log.Println("successfully parsed access token duration: ", atd)
	rtd, err := time.ParseDuration(c.Get("AUTH_REFRESH_TOKEN_DURATION"))
	checkError(err, "invalid refresh token duration")
	log.Println("successfully parsed refresh token duration: ", rtd)
	tg, err := token.NewPaseto(c.Get("AUTH_TOKEN_SECRET"), atd, rtd)
	checkError(err, "failed to create token generator")
	log.Println(fmt.Sprintf("successfully create token generator"))

	// Create postgres conn
	pgConn, err := mypostgres.New(c.Get("AUTH_DATABASE_URL"))
	checkError(err, "failed to connect to postgres")
	log.Println("successfully connected to postgres")

	// Migrate database
	err = mypostgres.Migrate(pgConn, c.Get("AUTH_DATABASE_MIGRATION_PATH"))
	checkError(err, "failed to migrate postgres db")
	log.Println("successfully migrated postgres db")

	// Create user repository
	ur := mypostgres.NewUserRepository(pgConn)
	log.Println("successfully created user repository")

	// Create session repository
	std, err := time.ParseDuration(c.Get("AUTH_USER_SESSION_DURATION"))
	checkError(err, "invalid user session duration")
	log.Println("successfully parsed user session duration: ", std)
	sr, err := mypostgres.NewSessionRepository(pgConn, mypostgres.WithSessionDuration(std))
	checkError(err, "failed to create session repository")
	log.Println("successfully created session repository")

	// Create role repository
	rr := mypostgres.NewRolesRepository(pgConn)
	log.Println("successfully created role repository")

	// Connect to redis cache
	redisConn, err := redis.New(c.Get("AUTH_CACHE_URL"))
	checkError(err, "failed to connect to redis cache")
	log.Println("successfully connected to redis cache")

	// Create token repository
	tr := redis.NewTokenRepository(redisConn, redis.WithTokenDuration(atd))
	log.Println("successfully created token repository")

	// Connect to rabbitmq & Create a new message producer
	rbp, err := rabbitmq.NewProducer(c.Get("AUTH_RABBITMQ_URL"),
		rabbitmq.WithNewSignInSessionQueue(c.Get("AUTH_RABBITMQ_NEW_SIGNIN_SESSION_QUEUE")),
	)
	checkError(err, "failed to connect to rabbitmq")
	log.Println("successfully connected to rabbitmq")

	// Create a new use case
	uc := application.NewUseCases(
		application.WithValidator(v),
		application.WithPasswordHasher(ph),
		application.WithTokenGenerator(tg),
		application.WithUserRepository(ur),
		application.WithSessionRepository(sr),
		application.WithTokenRepository(tr),
		application.WithRoleRepository(rr),
		application.WithMessageProducer(rbp),
	)

	// Start the server
	checkError(startGRPCServer(c, uc), "Failed to start gRPC server")
}

func startGRPCServer(c *goconfig.Config, uc *application.UseCases) error {
	// Create a new gRPC server with TLS enabled
	grpcAH := mygrpc.NewAuthHandler(uc)
	grpcS := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcS, grpcAH)
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

func checkError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
