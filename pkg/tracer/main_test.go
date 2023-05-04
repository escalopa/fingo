package tracer

import (
	"context"
	"log"
	"testing"

	"github.com/escalopa/fingo/utils/testcontainer"
)

var (
	testJagerUrl string
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url, terminate, err := testcontainer.NewJagerContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	testJagerUrl = url
	m.Run()

	_ = terminate()
}
