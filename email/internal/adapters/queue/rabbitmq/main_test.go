package rabbitmq

import (
	"context"
	"github.com/escalopa/fingo/utils/testcontainer"
	"log"
	"testing"
)

var testRabbitMQ *Consumer

func TestMain(m *testing.M) {
	// Start the RabbitMQ container
	var err error
	url, terminate, err := testcontainer.NewRabbitMQContainer(context.Background())
	if err != nil {
		log.Fatal(err, "failed to start rabbitmq container")
	}
	log.Println(url)
	// Stop the container
	defer func() {
		if err = terminate(); err != nil {
			log.Fatal(err, "failed to stop rabbitmq container")
		}
	}()
	// Create the consumer
	testRabbitMQ, err = NewConsumer(url,
		WithVerificationCodeQueue("verification_code_queue"),
		WithResetPasswordTokenQueue("reset_password_token_queue"),
		WithNewSignInSessionQueue("new_sign_in_session_queue"),
	)
	if err != nil {
		log.Fatal(err, "failed to create rabbitmq consumer")
	}
	// Run the tests
	m.Run()
}
