package main

import (
	"context"
	"log"

	"github.com/escalopa/fingo/pkg/global"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/pkg/validator"

	"github.com/escalopa/fingo/wallet/internal/adapters/db"
	"github.com/escalopa/fingo/wallet/internal/adapters/locker"
	"github.com/escalopa/fingo/wallet/internal/adapters/numgen"
	"github.com/escalopa/fingo/wallet/internal/application"
)

func main() {
	appCtx, cancel := context.WithCancel(context.Background())
	go func() {
		<-global.CatchSignal()
		cancel()
	}()

	// Load cofigurations
	global.CheckError(global.LoadConfig(&cfg, "app", "./wallet", "env"), "failed to load configurations")

	// Create validator
	v := validator.NewValidator()
	log.Println("validator created")

	// Create database connection
	conn, err := db.New(cfg.DatabaseUrl)
	global.CheckError(err, "failed to create database connection")
	log.Print("database connection created")

	// Migrate database
	global.CheckError(db.Migrate(conn, cfg.DatabaseMigrationPath), "failed to migrate database")
	log.Print("database migrated")

	ur := db.NewUserRepository(conn)
	cr := db.NewCardRepository(conn)
	ar := db.NewAccountRepository(conn)
	tr := db.NewTransactionRepository(conn)

	// Create a new number generator
	cng := numgen.NewNumGen(cfg.CardNumberLength)
	log.Println("Card number generator created with card number length:", cfg.CardNumberLength)

	// Create a new tracer
	t, err := tracer.LoadTracer(
		cfg.TracingEnable,
		cfg.TracingJaegerEnable,
		cfg.TracingJaegerAgentUrl,
		cfg.TracingJaegerServiceName,
		cfg.TracingJaegerEnvironment,
	)
	global.CheckError(err, "failed to load tracer")
	tracer.SetTracer(t)
	log.Println("tracer created")

	// Create an ids locker
	l := locker.NewLocker(appCtx, cfg.LockerCleanupDuration)

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
	global.CheckError(start(appCtx, uc), "failed to start gRPC server")
}
