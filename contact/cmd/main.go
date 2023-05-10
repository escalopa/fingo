package main

import (
	"context"
	"log"
	"time"

	"github.com/escalopa/fingo/pkg/global"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/pkg/validator"

	"github.com/escalopa/fingo/contact/internal/adapters/server"

	"github.com/escalopa/fingo/contact/internal/adapters/email/mycourier"
	"github.com/escalopa/fingo/contact/internal/adapters/queue/rabbitmq"
	"github.com/escalopa/fingo/contact/internal/application"
)

func main() {
	// Create application context
	appCtx, cancel := context.WithCancel(context.Background())
	go func() {
		<-global.CatchSignal()
		cancel()
	}()

	// Load cofigurations
	global.CheckError(global.LoadConfig(&cfg, "app", "./contact", "env"), "failed to load configurations")

	// Create a validator
	v := validator.NewValidator()
	log.Println("created validator")

	// Parse code expiration from config
	log.Println("using codes-expiration:", cfg.CodesExpiration)

	// Create a courier sender
	cs, err := mycourier.New(
		cfg.CourierToken,
		mycourier.WithExpiration(cfg.CodesExpiration),
		mycourier.WithVerificationTemplate(cfg.CourierVerificationTemplateID),
		mycourier.WithResetPasswordTemplate(cfg.CourierResetPasswordTemplateID),
		mycourier.WithNewSignInSessionTemplate(cfg.CourierNewSigninSessionTemplateID),
	)
	global.CheckError(err, "failed to create courier sender")
	defer func() {
		if err := cs.Close(); err != nil {
			log.Println("failed to close courier sender", err)
		} else {
			log.Println("closing courier-sender")
		}
	}()
	log.Println("created courier-sender")

	// Create a rabbitmq consumer
	rbc, err := rabbitmq.NewConsumer(cfg.RabbitmqUrl,
		rabbitmq.WithVerificationCodeQueue(cfg.RabbitmqVerificationCodeQueueName),
		rabbitmq.WithResetPasswordTokenQueue(cfg.RabbitmqResetPasswordTokenQueueName),
		rabbitmq.WithNewSignInSessionQueue(cfg.RabbitmqNewSigninSessionQueueName),
	)
	global.CheckError(err, "failed to create rabbitmq consumer")
	log.Println("connected to rabbitmq")

	log.Println("using send-min-interval:", cfg.SendCodeMinInterval)
	log.Println("using send-min-interval:", cfg.SendResetPasswordTokenMinInterval)

	// Create use cases
	uc := application.NewUseCases(
		application.WithValidator(v),
		application.WithEmailSender(cs),
		application.WithMinSendCodeInterval(cfg.SendCodeMinInterval),
		application.WithMinSendPasswordTokenInterval(cfg.SendResetPasswordTokenMinInterval),
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

	// Create server
	s := server.NewServer(uc, rbc)
	log.Println("created server")

	// Terminate server on shutdown signals
	go global.Shutdown(appCtx, 10*time.Second, func() { s.Stop() }, func() {})

	// Start server
	log.Println("starting server")
	global.CheckError(s.Start(), "failed to start server")
}
