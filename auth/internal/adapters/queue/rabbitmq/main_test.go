package rabbitmq

import (
	"context"
	"log"
	"testing"

	"github.com/escalopa/fingo/utils/testcontainer"
)

var rabbitmqUrl string

func TestMain(m *testing.M) {
	// Start the RabbitMQ container
	var err error
	url, terminate, err := testcontainer.NewRabbitMQContainer(context.Background())
	if err != nil {
		log.Fatal(err, "failed to start rabbitmq container")
	}
	// Stop the container
	defer func() {
		if err = terminate(); err != nil {
			log.Fatal(err, "failed to stop rabbitmq container")
		}
	}()
	if err != nil {
		log.Fatal(err, "failed to create rabbitmq consumer")
	}
	rabbitmqUrl = url
	// Run the tests
	m.Run()
}
