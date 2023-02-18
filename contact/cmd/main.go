package main

import (
	"context"
	"github.com/escalopa/goconfig"
	"github.com/escalopa/gofly/contact/internal/adapter/redis"
	"log"
	"time"
)

func main() {
	c := goconfig.New()

	// Create db connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cache, err := redis.New(c.Get("CACHE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to cache")

	// Create a code repo
	cr := redis.NewCodeRepository(cache, redis.WithCodeContext(ctx), redis.WithExpiration(30*time.Minute))
	_ = cr
}
