package main

import (
	"context"
	"log"
	"time"

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
	"github.com/escalopa/goconfig"
)

func main() {
	appCtx, cancel := context.WithCancel(context.Background())
	go func() {
		<-global.CatchSignal()
		cancel()
	}()
	c := goconfig.New()

	ph := hasher.NewBcryptHasher()
	v := validator.NewValidator()

	// Create a new token generator
	atd, err := time.ParseDuration(c.Get("AUTH_ACCESS_TOKEN_DURATION"))
	global.CheckError(err, "invalid access token duration")
	log.Println("successfully parsed access token duration: ", atd)
	rtd, err := time.ParseDuration(c.Get("AUTH_REFRESH_TOKEN_DURATION"))
	global.CheckError(err, "invalid refresh token duration")
	log.Println("successfully parsed refresh token duration: ", rtd)
	tg, err := token.NewPaseto(c.Get("AUTH_TOKEN_SECRET"), atd, rtd)
	global.CheckError(err, "failed to create token generator")
	log.Println("successfully create token generator")

	// Create postgres conn
	pgConn, err := pdb.New(c.Get("AUTH_DATABASE_URL"))
	global.CheckError(err, "failed to connect to postgres")
	log.Println("successfully connected to postgres")

	// Migrate database
	err = pdb.Migrate(pgConn, c.Get("AUTH_DATABASE_MIGRATION_PATH"))
	global.CheckError(err, "failed to migrate postgres db")
	log.Println("successfully migrated postgres db")

	// Create user repository
	ur, err := mypostgres.NewUserRepository(pgConn)
	global.CheckError(err, "failed to create user repository")
	log.Println("successfully created user repository")

	// Create session repository
	std, err := time.ParseDuration(c.Get("AUTH_USER_SESSION_DURATION"))
	global.CheckError(err, "invalid user session duration")
	log.Println("successfully parsed user session duration: ", std)
	sr, err := mypostgres.NewSessionRepository(pgConn, mypostgres.WithSessionDuration(std))
	global.CheckError(err, "failed to create session repository")
	log.Println("successfully created session repository")

	// Connect to redis cache
	redisConn, err := redis.New(c.Get("AUTH_CACHE_URL"))
	global.CheckError(err, "failed to connect to redis cache")
	log.Println("successfully connected to redis cache")

	// Create token repository
	tr := redis.NewTokenRepository(redisConn, redis.WithTokenDuration(atd))
	log.Println("successfully created token repository")

	// Connect to rabbitmq & Create a new message producer
	rbp, err := rabbitmq.NewProducer(c.Get("AUTH_RABBITMQ_URL"),
		rabbitmq.WithNewSignInSessionQueue(c.Get("AUTH_RABBITMQ_NEW_SIGNIN_SESSION_QUEUE_NAME")),
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
		c.Get("AUTH_TRACING_ENABLE") == "true",
		c.Get("AUTH_TRACING_JAEGER_ENABLE") == "true",
		c.Get("AUTH_TRACING_JAEGER_AGENT_URL"),
		c.Get("AUTH_TRACING_JAEGER_SERVICE_NAME"),
		c.Get("AUTH_TRACING_JAEGER_ENVIRONMENT"),
	)
	global.CheckError(err, "failed to load tracer")
	tracer.SetTracer(t)

	// Start the server
	err = start(appCtx, c, uc)
	if err != nil {
		log.Println("failed to start auth grpc server")
	}
}
