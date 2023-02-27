package main

import (
	mypostgres "github.com/escalopa/gochat/auth/internal/adapters/db/postgres"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"time"

	mygrpc "github.com/escalopa/gochat/auth/internal/adapters/grpc"
	"github.com/escalopa/gochat/auth/internal/adapters/hasher"
	"github.com/escalopa/gochat/auth/internal/adapters/token"
	myvalidator "github.com/escalopa/gochat/auth/internal/adapters/validator"
	"github.com/escalopa/gochat/auth/internal/application"
	"github.com/escalopa/gochat/pb"
	"github.com/escalopa/goconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	c := goconfig.New()

	ph := hasher.NewBcryptHasher()
	v := myvalidator.NewValidator()

	// Create a new token generator
	atd, err := time.ParseDuration(c.Get("AUTH_ACCESS_TOKEN_DURATION"))
	if err != nil {
		log.Fatal(err, "Invalid access token duration")
	}
	rtd, err := time.ParseDuration(c.Get("AUTH_REFRESH_TOKEN_DURATION"))
	if err != nil {
		log.Fatal(err, "Invalid refresh token duration")
	}
	log.Println("Successfully parsed access token duration")
	tg, err := token.NewPaseto(c.Get("AUTH_TOKEN_SECRET"), atd, rtd)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully create token generator")

	// Create postgres conn
	pgConn, err := mypostgres.New(c.Get("AUTH_DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to postgres")

	// Migrate database
	err = mypostgres.Migrate(pgConn, c.Get("AUTH_DATABASE_MIGRATION_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully migrated postgres db")

	// Create user repository
	ur := mypostgres.NewUserRepository(pgConn)
	log.Println("Successfully created user repository")

	// Create session repository
	std, err := time.ParseDuration(c.Get("AUTH_USER_SESSION_DURATION"))
	if err != nil {
		log.Fatal(err, "Invalid user session duration")
	}
	log.Println("Successfully parsed user session duration")
	sr := mypostgres.NewSessionRepository(pgConn, std)
	log.Println("Successfully created session repository")

	// Connect to email service with gRPC
	conn, err := grpc.Dial(c.Get("AUTH_EMAIL_GRPC_URL"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	esc := pb.NewEmailServiceClient(conn)
	log.Println("Connected to email-service")

	// Create a new use case
	uc := application.NewUseCases(
		application.WithValidator(v),
		application.WithPasswordHasher(ph),
		application.WithTokenGenerator(tg),
		application.WithUserRepository(ur),
		application.WithSessionRepository(sr),
		application.WithEmailService(esc),
	)

	// Start the server
	StartGRPCServer(c, uc)
}

func StartGRPCServer(c *goconfig.Config, uc *application.UseCases) {
	// Create a new gRPC server
	grpcAH := mygrpc.NewAuthHandler(uc)
	grpcS := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcS, grpcAH)
	reflection.Register(grpcS)

	// Start the server
	port := c.Get("AUTH_GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	log.Println("Starting gRPC server on port", port)
	err = grpcS.Serve(lis)
	if err != nil {
		return
	}
}
