package main

import (
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"time"

	"github.com/escalopa/goconfig"
	"github.com/escalopa/gofly/auth/internal/adapters/cache/redis"
	mygrpc "github.com/escalopa/gofly/auth/internal/adapters/grpc"
	"github.com/escalopa/gofly/auth/internal/adapters/hasher"
	"github.com/escalopa/gofly/auth/internal/adapters/token"
	myValidator "github.com/escalopa/gofly/auth/internal/adapters/validator"
	"github.com/escalopa/gofly/auth/internal/application"
	"github.com/escalopa/gofly/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	c := goconfig.New()

	ph := hasher.NewBcryptHasher()
	v := myValidator.NewValidator()

	// Create a new token generator
	atd, err := time.ParseDuration(c.Get("TOKEN_ACCESS_DURATION"))
	if err != nil {
		log.Fatal(err, "Invalid access token duration")
	}
	tg, err := token.NewPaseto(c.Get("TOKEN_SECRET"), atd)
	if err != nil {
		log.Fatal(err)
	}

	// Create db connection and repository
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	cache, err := redis.New(c.Get("AUTH_CACHE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to cache")

	// Create a new user repository
	ur := redis.NewUserRepository(cache,
		redis.WithTimeout(5*time.Second),
	)
	log.Println("Connected to user-repository")

	// Connect to email service with gRPC
	conn, err := grpc.Dial(c.Get("EMAIL_GRPC_URL"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(1*time.Minute),
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
