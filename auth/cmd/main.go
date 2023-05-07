package main

import (
	"context"
	"log"

	"github.com/escalopa/fingo/pkg/global"
	"github.com/escalopa/fingo/pkg/pdb"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/pkg/validator"

	mypostgres "github.com/escalopa/fingo/auth/internal/adapters/db/postgres"
	"github.com/escalopa/fingo/auth/internal/adapters/db/redis"
	"github.com/escalopa/fingo/auth/internal/adapters/hasher"
	"github.com/escalopa/fingo/auth/internal/adapters/queue/rabbitmq"
	"github.com/escalopa/fingo/auth/internal/adapters/token"
	"github.com/escalopa/fingo/auth/internal/application"
)

func main() {
	appCtx, cancel := context.WithCancel(context.Background())
	go func() {
		<-global.CatchSignal()
		cancel()
	}()

	// Load cofigurations
	global.CheckError(global.LoadConfig(&cfg, "app", "./auth", "env"), "failed to load configurations")

	ph := hasher.NewBcryptHasher()
	v := validator.NewValidator()

	// Create a new token generator
	tg, err := token.NewPaseto(
		cfg.TokenSecret,
		cfg.AccessTokenDuration,
		cfg.RefreshTokenDuration,
	)
	global.CheckError(err, "failed to create token generator")
	log.Println("successfully parsed access token duration: ", cfg.AccessTokenDuration)
	log.Println("successfully parsed refresh token duration: ", cfg.RefreshTokenDuration)
	log.Println("successfully create token generator")

	// Create postgres conn
	pgConn, err := pdb.New(cfg.DatabaseUrl)
	global.CheckError(err, "failed to connect to postgres")
	log.Println("successfully connected to postgres")

	// Migrate database
	err = pdb.Migrate(pgConn, cfg.DatabaseMigrationPath)
	global.CheckError(err, "failed to migrate postgres db")
	log.Println("successfully migrated postgres db")

	// Create user repository
	ur, err := mypostgres.NewUserRepository(pgConn)
	global.CheckError(err, "failed to create user repository")
	log.Println("successfully created user repository")

	// Create session repository
	sr, err := mypostgres.NewSessionRepository(pgConn, mypostgres.WithSessionDuration(cfg.UserSessionDuration))
	global.CheckError(err, "failed to create session repository")
	log.Println("successfully parsed user session duration: ", cfg.UserSessionDuration)
	log.Println("successfully created session repository")

	// Connect to redis cache
	redisConn, err := redis.New(cfg.RedisUrl)
	global.CheckError(err, "failed to connect to redis cache")
	log.Println("successfully connected to redis cache")

	// Create token repository
	tr := redis.NewTokenRepository(redisConn, redis.WithTokenDuration(cfg.AccessTokenDuration))
	log.Println("successfully created token repository")

	// Connect to rabbitmq & Create a new message producer
	rbp, err := rabbitmq.NewProducer(cfg.RabbitmqUrl,
		rabbitmq.WithNewSignInSessionQueue(cfg.RabbitmqNewSigninSessionQueueName),
	)
	global.CheckError(err, "failed to connect to rabbitmq")
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

	// Start the server
	err = start(appCtx, uc)
	if err != nil {
		log.Println("failed to start auth grpc server")
	}
}
