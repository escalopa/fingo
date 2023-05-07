package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/escalopa/fingo/wallet/internal/adapters/db"
	"github.com/escalopa/fingo/wallet/internal/adapters/locker"

	"github.com/escalopa/fingo/pkg/global"
	"github.com/escalopa/fingo/pkg/pdb"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/pkg/validator"
	"github.com/escalopa/fingo/wallet/internal/adapters/numgen"
	"github.com/escalopa/fingo/wallet/internal/application"

	"github.com/escalopa/goconfig"
)

func main() {
	c := goconfig.New()
	appCtx, cancel := context.WithCancel(context.Background())
	go func() {
		<-global.CatchSignal()
		cancel()
	}()
	// Create validator
	v := validator.NewValidator()
	log.Println("validator created")

	// Create database connection
	conn, err := pdb.New(c.Get("WALLET_DATABASE_URL"))
	global.CheckError(err, "failed to create database connection")
	log.Print("database connection created")

	// Migrate database
	global.CheckError(pdb.Migrate(conn, c.Get("WALLET_DATABASE_MIGRATION_PATH")), "failed to migrate database")
	log.Print("database migrated")

	ur := db.NewUserRepository(conn)
	cr := db.NewCardRepository(conn)
	ar := db.NewAccountRepository(conn)
	tr := db.NewTransactionRepository(conn)

	// Create a new number generator
	cardNumberLength, err := strconv.Atoi(c.Get("WALLET_CARD_NUMBER_LENGTH"))
	global.CheckError(err, "failed to convert WALLET_CARD_NUMBER_LENGTH to int")
	cng := numgen.NewNumGen(cardNumberLength)
	log.Println("Card number generator created with card number length:", cardNumberLength)

	// Create a new tracer
	t, err := tracer.LoadTracer(
		c.Get("WALLET_TRACING_ENABLE") == "true",
		c.Get("WALLET_TRACING_JAEGER_ENABLE") == "true",
		c.Get("WALLET_TRACING_JAEGER_AGENT_URL"),
		c.Get("WALLET_TRACING_JAEGER_SERVICE_NAME"),
		c.Get("WALLET_TRACING_JAEGER_ENVIRONMENT"),
	)
	global.CheckError(err, "failed to load tracer")
	tracer.SetTracer(t)
	log.Println("tracer created")

	// Create an ids locker
	cleanupDuration, err := time.ParseDuration(c.Get("WALLET_LOCKER_CLEANUP_DURATION"))
	global.CheckError(err, "failed to convert WALLET_LOCKER_CLEANUP_DURATION to time.Duration")
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

	// Start gRPC server
	global.CheckError(start(appCtx, c, uc), "failed to start gRPC server")
}
