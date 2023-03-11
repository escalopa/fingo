package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/escalopa/fingo/email/internal/adapter/codegen"
	"github.com/escalopa/fingo/email/internal/adapter/email/mycourier"
	mygrpc "github.com/escalopa/fingo/email/internal/adapter/grpc"
	"github.com/escalopa/fingo/email/internal/adapter/redis"
	"github.com/escalopa/fingo/email/internal/application"
	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/goconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Return code
	returnCode := 0
	defer func() { os.Exit(returnCode) }()

	// Create a new config instance
	c := goconfig.New()

	// Create redis client
	cache, err := redis.New(c.Get("EMAIL_CACHE_URL"))
	checkError(err, "Failed to connect to cache")
	log.Println("Connected to cache")

	// Parse code expiration from config
	exp, err := time.ParseDuration(c.Get("EMAIL_USER_CODE_EXPIRATION"))
	checkError(err, "Failed to parse code expiration")
	log.Println("Using code-expiration:", exp)

	// Create a code repo
	cr := redis.NewCodeRepository(cache,
		redis.WithExpiration(exp),
	)
	// Close code repo on exit
	defer func() {
		err := cr.Close()
		checkError(err, "Failed to close code repo")
	}()
	log.Println("Connected to code-repo")

	// Create a courier sender
	cs, err := mycourier.New(c.Get("EMAIL_COURIER_TOKEN"),
		mycourier.WithExpiration(exp),
		mycourier.WithVerificationTemplate(c.Get("EMAIL_COURIER_VERIFICATION_TEMPLATE_ID")),
	)
	checkError(err, "Failed to create courier sender")
	log.Println("Connected to courier-sender")

	// Create a code generator
	codeLen, err := strconv.Atoi(c.Get("EMAIL_USER_CODE_LENGTH"))
	checkError(err, "Failed to parse code length")
	if codeLen < 1 {
		log.Println("Code length must be greater than 0")
	}
	log.Println("Using Code-length:", codeLen)
	cg := codegen.New(codeLen)

	// Create use cases
	mti, err := time.ParseDuration(c.Get("EMAIL_MIN_SEND_INTERVAL"))
	checkError(err, "Failed to parse min send interval")
	log.Println("Using min-send-interval:", mti)

	uc := application.NewUseCases(
		application.WithCodeRepository(cr),
		application.WithCodeGenerator(cg),
		application.WithEmailSender(cs),
		application.WithMinTimeInterval(mti),
	)

	StartGRPCServer(c, uc)
}

func StartGRPCServer(c *goconfig.Config, uc *application.UseCases) {
	// Create a new gRPC server
	grpcH := mygrpc.New(uc)
	grpcS := grpc.NewServer()
	pb.RegisterEmailServiceServer(grpcS, grpcH)
	reflection.Register(grpcS)

	// Start the server
	port := c.Get("EMAIL_GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	checkError(err, fmt.Sprintf("Failed to listen on port %s", port))

	log.Println("Starting gRPC server on port", port)
	err = grpcS.Serve(lis)
	checkError(err, "Failed to start gRPC server")
}

func checkError(err error, msg string) {
	if err != nil {
		log.Println(err, msg)
		os.Exit(1)
	}
}
