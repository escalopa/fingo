package main

import (
	"log"
	"time"

	"github.com/escalopa/fingo/email/internal/adapters/server"
	"github.com/escalopa/fingo/email/internal/adapters/validator"

	"github.com/escalopa/fingo/email/internal/adapters/email/mycourier"
	"github.com/escalopa/fingo/email/internal/adapters/queue/rabbitmq"
	"github.com/escalopa/fingo/email/internal/application"
	"github.com/escalopa/goconfig"
)

func main() {
	// Create a new config instance
	c := goconfig.New()
	// Parse code expiration from config
	exp, err := time.ParseDuration(c.Get("EMAIL_CODES_EXPIRATION"))
	checkError(err, "failed to parse code expiration")
	log.Println("using codes-expiration:", exp)

	// Create a courier sender
	cs, err := mycourier.New(c.Get("EMAIL_COURIER_TOKEN"),
		mycourier.WithExpiration(exp),
		mycourier.WithVerificationTemplate(c.Get("EMAIL_COURIER_VERIFICATION_TEMPLATE_ID")),
		mycourier.WithResetPasswordTemplate(c.Get("EMAIL_COURIER_RESET_PASSWORD_TEMPLATE_ID")),
		mycourier.WithNewSignInSessionTemplate(c.Get("EMAIL_COURIER_NEW_SIGNIN_SESSION_TEMPLATE_ID")),
	)
	checkError(err, "failed to create courier sender")
	defer func() {
		log.Println("closing courier-sender")
		_ = cs.Close()
	}()
	log.Println("created courier-sender")

	// Create a rabbitmq consumer
	rbc, err := rabbitmq.NewConsumer(c.Get("EMAIL_RABBITMQ_URL"),
		rabbitmq.WithVerificationCodeQueue(c.Get("EMAIL_RABBITMQ_VERIFICATION_CODE_QUEUE_NAME")),
		rabbitmq.WithResetPasswordTokenQueue(c.Get("EMAIL_RABBITMQ_RESET_PASSWORD_TOKEN_QUEUE_NAME")),
		rabbitmq.WithNewSignInSessionQueue(c.Get("EMAIL_RABBITMQ_NEW_SIGNIN_SESSION_QUEUE_NAME")),
	)
	checkError(err, "failed to create rabbitmq consumer")
	defer func() {
		log.Println("closing rabbitmq consumer")
		_ = rbc.Close()
	}()
	log.Println("connected to rabbitmq")

	// Parse send code min interval from config
	smi, err := time.ParseDuration(c.Get("EMAIL_SEND_CODE_MIN_INTERVAL"))
	checkError(err, "failed to parse send code min interval")
	log.Println("using send-min-interval:", smi)
	// Parse send reset password token min interval from config
	spi, err := time.ParseDuration(c.Get("EMAIL_SEND_RESET_PASSWORD_TOKEN_MIN_INTERVAL"))
	checkError(err, "failed to parse send reset password token min interval")
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

	// Create server
	s := server.NewServer(uc, rbc)
	log.Println("created server")

	// Start server
	log.Println("starting server")
	checkError(s.Start(), "failed to start server")
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
