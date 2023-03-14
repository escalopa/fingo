package rabbitmq

import (
	"context"
	"log"
	"testing"

	"github.com/escalopa/fingo/utils/testcontainer"
)

var rabbitmqURL string

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
	if err != nil {
		log.Fatal(err, "failed to create rabbitmq consumer")
	}
	rabbitmqURL = url
	// Run the tests
	m.Run()
}
