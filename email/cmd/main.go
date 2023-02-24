package main

import (
	"github.com/escalopa/gochat/email/internal/adapter/codegen"
	"github.com/escalopa/gochat/email/internal/adapter/email/mycourier"
	mygrpc "github.com/escalopa/gochat/email/internal/adapter/mygrpc"
	"github.com/escalopa/gochat/email/internal/adapter/redis"
	"github.com/escalopa/gochat/email/internal/application"
	"github.com/escalopa/gochat/pb"
	"github.com/escalopa/goconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
	"time"
)

func main() {
	c := goconfig.New()

	// Create db connection
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	cache, err := redis.New(c.Get("EMAIL_CACHE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to cache")

	// Parse code expiration from config
	exp, err := time.ParseDuration(c.Get("EMAIL_USER_CODE_EXPIRATION"))
	if err != nil {
		log.Fatal(err, "Failed to parse code expiration")
	}
	log.Println("Using code-expiration:", exp)

	// Create a code repo
	cr := redis.NewCodeRepository(cache,
		redis.WithExpiration(exp),
	)
	// Close code repo on exit
	defer func(cr *redis.CodeRepository) {
		err := cr.Close()
		if err != nil {
			log.Println(err, "Failed to close code repo")
		}
	}(cr)
	log.Println("Connected to code-repo")

	// Create a courier sender
	cs, err := mycourier.New(c.Get("EMAIL_COURIER_TOKEN"),
		mycourier.WithExpiration(exp),
		mycourier.WithVerificationTemplate(c.Get("EMAIL_COURIER_VERIFICATION_TEMPLATE_ID")),
	)
	if err != nil {
		log.Fatal(err, "Failed to create courier sender")
	}
	log.Println("Connected to courier-sender")

	// Create a code generator
	codeLen, err := strconv.Atoi(c.Get("EMAIL_USER_CODE_LENGTH"))
	if err != nil {
		log.Fatal(err, "Failed to parse code length")
	}
	if codeLen < 1 {
		log.Fatal("Code length must be greater than 0")
	}
	log.Println("Using Code-length:", codeLen)
	cg := codegen.New(codeLen)

	// Create use cases
	mti, err := time.ParseDuration(c.Get("EMAIL_MIN_SEND_INTERVAL"))
	if err != nil {
		log.Fatal(err, "Failed to parse min send interval")
	}
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
	if err != nil {
		panic(err)
	}
	log.Println("Starting gRPC server on port", port)
	err = grpcS.Serve(lis)
	if err != nil {
		return
	}
}
