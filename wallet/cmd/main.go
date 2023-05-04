package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/escalopa/fingo/wallet/internal/adapters/db"
	"github.com/escalopa/fingo/wallet/internal/adapters/locker"

	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/fingo/pkg/pdb"
	"github.com/escalopa/fingo/pkg/perror"
	"github.com/escalopa/fingo/pkg/tls"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/pkg/validator"
	mygrpc "github.com/escalopa/fingo/wallet/internal/adapters/grpc"
	"github.com/escalopa/fingo/wallet/internal/adapters/numgen"
	"github.com/escalopa/fingo/wallet/internal/application"
	"github.com/lordvidex/errs"

	"github.com/escalopa/goconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	c := goconfig.New()
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create validator
	v := validator.NewValidator()
	log.Println("validator created")

	// Create database connection
	conn, err := pdb.New(c.Get("WALLET_DATABASE_URL"))
	perror.CheckError(err, "failed to create database connection")
	log.Print("database connection created")

	// Migrate database
	perror.CheckError(pdb.Migrate(conn, c.Get("WALLET_DATABASE_MIGRATION_PATH")), "failed to migrate database")
	log.Print("database migrated")

	ur := db.NewUserRepository(conn)
	cr := db.NewCardRepository(conn)
	ar := db.NewAccountRepository(conn)
	tr := db.NewTransactionRepository(conn)

	// Create a new number generator
	cardNumberLength, err := strconv.Atoi(c.Get("WALLET_CARD_NUMBER_LENGTH"))
	perror.CheckError(err, "failed to convert WALLET_CARD_NUMBER_LENGTH to int")
	cng := numgen.NewNumGen(cardNumberLength)
	log.Println("Card number generator created with card number length:", cardNumberLength)

	// Create a new tracer
	t, err := tracer.LoadTracer(
		c.Get("WALLET_TRACING_ENABLE"),
		c.Get("WALLET_TRACING_JAEGER_ENABLE"),
		c.Get("WALLET_TRACING_JAEGER_AGENT_URL"),
		c.Get("WALLET_TRACING_JAEGER_SERVICE_NAME"),
		c.Get("WALLET_TRACING_JAEGER_ENVIRONMENT"),
	)
	perror.CheckError(err, "failed to load tracer")
	tracer.SetTracer(t)
	log.Println("tracer created")

	// Create an ids locker
	cleanupDuration, err := time.ParseDuration(c.Get("WALLET_LOCKER_CLEANUP_DURATION"))
	perror.CheckError(err, "failed to convert WALLET_LOCKER_CLEANUP_DURATION to time.Duration")
	l := locker.NewLocker(appCtx, cleanupDuration)

	// Create use cases
	uc := application.NewUseCases(
		application.WithValidator(v),
		application.WithLocker(l),
		application.WithUserRepository(ur),
		application.WithCardRepository(cr),
		application.WithAccountRepository(ar),
		application.WithTransactionRepository(tr),
		application.WithCardNumberGenerator(cng),
	)

	var opts []grpc.ServerOption
	// Load TLS certificates
	err = loadTls(c, &opts)
	perror.CheckError(err, "failed to load wallet TLS certificates")

	// Load auth interceptor
	err = loadInterceptor(c, &opts)
	perror.CheckError(err, "failed to load auth interceptor")

	// Start gRPC server
	perror.CheckError(start(c, uc, opts), "failed to start gRPC server")
}

func start(c *goconfig.Config, uc *application.UseCases, opts []grpc.ServerOption) error {
	// Create a gRPC server object
	handler := mygrpc.NewWalletHandler(uc)
	grpcS := grpc.NewServer(opts...)
	pb.RegisterWalletServiceServer(grpcS, handler)
	reflection.Register(grpcS)

	// Start the server
	port := c.Get("WALLET_GRPC_PORT")
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
	creds, err := tls.LoadServerTLS(
		c.Get("WALLET_GRPC_TLS_ENABLE") == "true",
		c.Get("WALLET_GRPC_TLS_CERT_FILE"),
		c.Get("WALLET_GRPC_TLS_KEY_FILE"),
	)
	if err != nil {
		return err
	}
	*opts = append(*opts, grpc.Creds(creds))
	return nil
}

func loadInterceptor(c *goconfig.Config, opts *[]grpc.ServerOption) error {
	creds, err := tls.LoadClientTLS(
		c.Get("TOKEN_GRPC_TLS_ENABLE") == "true",
		c.Get("WALLET_TOKEN_GRPC_TLS_USER_CERT_FILE"),
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
