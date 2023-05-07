package main

import (
	"context"
	"log"

	"github.com/escalopa/fingo/pkg/global"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/pkg/validator"

	"github.com/escalopa/fingo/token/internal/adapters/cache"
	"github.com/escalopa/fingo/token/internal/application"
	"github.com/escalopa/goconfig"
)

func main() {
	appCtx, cancel := context.WithCancel(context.Background())
	go func() {
		<-global.CatchSignal()
		cancel()
	}()
	c := goconfig.New()

	// Create validator
	v := validator.NewValidator()
	log.Println("validator created")

	// Create redis client
	rc, err := cache.NewRedisClient(c.Get("TOKEN_REDIS_URL"))
	global.CheckError(err, "failed to create redis client")
	log.Println("redis client created")

	// Create token repository
	tr := cache.NewTokenRepositoryV1(rc)
	log.Println("token repository created")

	// Create use cases
	uc := application.NewUseCases(
		application.WithValidator(v),
		application.WithTokenRepository(tr),
	)

	// Create a new tracer
	t, err := tracer.LoadTracer(
		c.Get("TOKEN_TRACING_ENABLE") == "true",
		c.Get("TOKEN_TRACING_JAEGER_ENABLE") == "true",
		c.Get("TOKEN_TRACING_JAEGER_AGENT_URL"),
		c.Get("TOKEN_TRACING_JAEGER_SERVICE_NAME"),
		c.Get("TOKEN_TRACING_JAEGER_ENVIRONMENT"),
	)
	global.CheckError(err, "failed to load tracer")
	tracer.SetTracer(t)
	log.Println("tracer created")

	// Start gRPC server
	err = start(appCtx, c, uc)
	if err != nil {
		log.Println("Server start/shutdown failed: ", err)
	} else {
		log.Println("Server shutdown completed")
	}
}
