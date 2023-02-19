package main

import (
	"context"
	"github.com/escalopa/goconfig"
	"github.com/escalopa/gofly/contact/internal/adapter/email/mycourier"
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

	// Parse code expiration from config
	exp, err := time.ParseDuration(c.Get("CODE_EXPIRATION"))
	if err != nil {
		log.Fatal(err, "Failed to parse code expiration")
	}

	// Create a code repo
	cr := redis.NewCodeRepository(cache,
		redis.WithCodeContext(ctx),
		redis.WithExpiration(exp),
	)
	// Close code repo on exit
	defer func(cr *redis.CodeRepository) {
		err := cr.Close()
		if err != nil {
			log.Println(err, "Failed to close code repo")
		}
	}(cr)

	// Create an email sender
	//from := c.Get("EMAIL_FROM")
	//host := c.Get("EMAIL_HOST")
	//port, err := strconv.Atoi(c.Get("EMAIL_PORT"))
	//if err != nil {
	//	log.Fatal(err, "Failed to parse email port")
	//}
	//es, err := mysmtp.New(
	//	smtp.WithExpiration(exp),
	//	smtp.WithHost(host),
	//	smtp.WithPort(port),
	//	smtp.WithFrom(from),
	//)
	//if err != nil {
	//	log.Fatal(err, "Failed to create email sender")
	//}
	// Close email sender on exit
	//defer func(es *smtp.Sender) {
	//	err := es.Close()
	//	if err != nil {
	//		log.Println(err, "Failed to close email sender")
	//	}
	//}(es)

	// Create a courier sender
	cs, err := mycourier.New(c.Get("COURIER_TOKEN"),
		mycourier.WithExpiration(exp),
		mycourier.WithVerificationTemplate(c.Get("COURIER_VERIFICATION_TEMPLATE_ID")),
	)
	if err != nil {
		log.Fatal(err, "Failed to create courier sender")
	}

	_, _ = cs, cr
}
