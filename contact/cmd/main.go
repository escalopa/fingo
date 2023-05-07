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
	"github.com/escalopa/goconfig"
)

func main() {
	// Create application context
	appCtx, cancel := context.WithCancel(context.Background())
	go func() {
		<-global.CatchSignal()
		cancel()
	}()

	// Create a new config instance
	c := goconfig.New()
	// Parse code expiration from config
	exp, err := time.ParseDuration(c.Get("CONTACT_CODES_EXPIRATION"))
	global.CheckError(err, "failed to parse code expiration")
	log.Println("using codes-expiration:", exp)

	// Create a courier sender
	cs, err := mycourier.New(c.Get("CONTACT_COURIER_TOKEN"),
		mycourier.WithExpiration(exp),
		mycourier.WithVerificationTemplate(c.Get("CONTACT_COURIER_VERIFICATION_TEMPLATE_ID")),
		mycourier.WithResetPasswordTemplate(c.Get("CONTACT_COURIER_RESET_PASSWORD_TEMPLATE_ID")),
		mycourier.WithNewSignInSessionTemplate(c.Get("CONTACT_COURIER_NEW_SIGNIN_SESSION_TEMPLATE_ID")),
	)
	global.CheckError(err, "failed to create courier sender")
	defer func() {
		log.Println("closing courier-sender")
		_ = cs.Close()
	}()
	log.Println("created courier-sender")

	// Create a rabbitmq consumer
	rbc, err := rabbitmq.NewConsumer(c.Get("CONTACT_RABBITMQ_URL"),
		rabbitmq.WithVerificationCodeQueue(c.Get("CONTACT_RABBITMQ_VERIFICATION_CODE_QUEUE_NAME")),
		rabbitmq.WithResetPasswordTokenQueue(c.Get("CONTACT_RABBITMQ_RESET_PASSWORD_TOKEN_QUEUE_NAME")),
		rabbitmq.WithNewSignInSessionQueue(c.Get("CONTACT_RABBITMQ_NEW_SIGNIN_SESSION_QUEUE_NAME")),
	)
	global.CheckError(err, "failed to create rabbitmq consumer")
	defer func() {
		log.Println("closing rabbitmq consumer")
		_ = rbc.Close()
	}()
	log.Println("connected to rabbitmq")

	// Parse send code min interval from config
	smi, err := time.ParseDuration(c.Get("CONTACT_SEND_CODE_MIN_INTERVAL"))
	global.CheckError(err, "failed to parse send code min interval")
	log.Println("using send-min-interval:", smi)
	// Parse send reset password token min interval from config
	spi, err := time.ParseDuration(c.Get("CONTACT_SEND_RESET_PASSWORD_TOKEN_MIN_INTERVAL"))
	global.CheckError(err, "failed to parse send reset password token min interval")
	log.Println("using send-min-interval:", spi)

	// Create a validator
	v := validator.NewValidator()
	log.Println("created validator")

	// Create use cases
	uc := application.NewUseCases(
		application.WithValidator(v),
		application.WithEmailSender(cs),
		application.WithMinSendCodeInterval(smi),
		application.WithMinSendPasswordTokenInterval(spi),
	)

	// Create a new tracer
	t, err := tracer.LoadTracer(
		c.Get("CONTACT_TRACING_ENABLE") == "true",
		c.Get("CONTACT_TRACING_JAEGER_ENABLE") == "true",
		c.Get("CONTACT_TRACING_JAEGER_AGENT_URL"),
		c.Get("CONTACT_TRACING_JAEGER_SERVICE_NAME"),
		c.Get("CONTACT_TRACING_JAEGER_ENVIRONMENT"),
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
