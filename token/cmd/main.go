package main

import (
	"context"
	"log"

	"github.com/escalopa/fingo/pkg/global"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/pkg/validator"

	"github.com/escalopa/fingo/token/internal/adapters/cache"
	"github.com/escalopa/fingo/token/internal/application"
)

func main() {
	appCtx, cancel := context.WithCancel(context.Background())
	go func() {
		<-global.CatchSignal()
		cancel()
	}()

	// Load cofigurations
	global.CheckError(global.LoadConfig(&cfg, "app", "./token", "env"), "failed to load configurations")

	// Create validator
	v := validator.NewValidator()
	log.Println("validator created")

	// Create redis client
	rc, err := cache.NewRedisClient(cfg.RedisUrl)
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
		cfg.TracingEnable,
		cfg.TracingJaegerEnable,
		cfg.TracingJaegerAgentUrl,
		cfg.TracingJaegerServiceName,
		cfg.TracingJaegerEnvironment,
	)
	global.CheckError(err, "failed to load tracer")
	tracer.SetTracer(t)
	log.Println("tracer created")

	// Start gRPC server
	err = start(appCtx, uc)
	if err != nil {
		log.Println("Server start/shutdown failed: ", err)
	} else {
		log.Println("Server shutdown completed")
	}
}
