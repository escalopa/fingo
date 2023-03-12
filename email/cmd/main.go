package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
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
	checkError(err, "Failed to parse code expiration")
	log.Println("Using codes-expiration:", exp)

	// Create a courier sender
	cs, err := mycourier.New(c.Get("EMAIL_COURIER_TOKEN"),
		mycourier.WithExpiration(exp),
		mycourier.WithVerificationTemplate(c.Get("EMAIL_COURIER_VERIFICATION_TEMPLATE_ID")),
		mycourier.WithResetPasswordTemplate(c.Get("EMAIL_COURIER_RESET_PASSWORD_TEMPLATE_ID")),
		mycourier.WithNewSignInSessionTemplate(c.Get("EMAIL_COURIER_NEW_SIGNIN_SESSION_TEMPLATE_ID")),
	)
	checkError(err, "Failed to create courier sender")
	defer func() {
		log.Println("Closing courier-sender")
		_ = cs.Close()
	}()
	log.Println("Created courier-sender")

	// Create a rabbitmq consumer
	rbc, err := rabbitmq.NewConsumer(c.Get("EMAIL_RABBITMQ_URL"),
		rabbitmq.WithVerificationCodeQueue(c.Get("EMAIL_RABBITMQ_VERIFICATION_CODE_QUEUE_NAME")),
		rabbitmq.WithResetPasswordTokenQueue(c.Get("EMAIL_RABBITMQ_RESET_PASSWORD_TOKEN_QUEUE_NAME")),
		rabbitmq.WithNewSignInSessionQueue(c.Get("EMAIL_RABBITMQ_NEW_SIGNIN_SESSION_QUEUE_NAME")),
	)
	checkError(err, "Failed to create rabbitmq consumer")
	defer func() {
		log.Println("Closing rabbitmq consumer")
		_ = rbc.Close()
	}()

	// Parse send code min interval from config
	smi, err := time.ParseDuration(c.Get("EMAIL_SEND_CODE_MIN_INTERVAL"))
	checkError(err, "Failed to parse send code min interval")
	log.Println("Using send-min-interval:", smi)
	// Parse send reset password token min interval from config
	spi, err := time.ParseDuration(c.Get("EMAIL_SEND_RESET_PASSWORD_TOKEN_MIN_INTERVAL"))
	checkError(err, "Failed to parse send reset password token min interval")
	log.Println("Using send-min-interval:", spi)

	// Create a validator
	v := validator.NewValidator()
	log.Println("Created validator")

	// Create use cases
	uc := application.NewUseCases(
		application.WithValidator(v),
		application.WithEmailSender(cs),
		application.WithMinSendCodeInterval(smi),
		application.WithMinSendPasswordTokenInterval(spi),
	)

	// Create server
	s := server.NewServer(uc, rbc)
	log.Println("Created server")

	// Handle SIGINT and SIGTERM.
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt, syscall.SIGTERM:
			log.Println("Gracefully Stopping server")
			s.Stop()
		}
	}()

	// Start server
	log.Println("Starting server")
	if err = s.Start(); err != nil {
		log.Fatal(err, "Failed to start server")
	}
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatal(err, msg)
	}
}
