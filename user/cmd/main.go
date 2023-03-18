package main

import (
	"log"
	"strconv"
	"time"

	"github.com/escalopa/fingo/pkg/pkgError"
	"github.com/escalopa/fingo/user/internal/adapters/codegen"
	"github.com/escalopa/fingo/user/internal/adapters/redis"
	"github.com/escalopa/goconfig"
)

func main() {
	// Create a new config instance
	c := goconfig.New()

	// Parse code expiration from config
	exp, err := time.ParseDuration(c.Get("EMAIL_USER_CODE_EXPIRATION"))
	pkgError.CheckError(err, "Failed to parse code expiration")
	log.Println("Using code-expiration:", exp)

	// Create redis client
	cache, err := redis.New(c.Get("EMAIL_CACHE_URL"))
	pkgError.CheckError(err, "Failed to connect to cache")
	log.Println("Connected to cache")

	// Create a code repo
	cr := redis.NewCodeRepository(cache,
		redis.WithExpiration(exp),
	)
	defer func() { _ = cr.Close() }()
	log.Println("Connected to code-repo")

	// Create a code generator
	codeLen, err := strconv.Atoi(c.Get("EMAIL_USER_CODE_LENGTH"))
	pkgError.CheckError(err, "Failed to parse code length")
	cg, err := codegen.New(codeLen)
	pkgError.CheckError(err, "Failed to create code-generator")
	log.Println("Using Code-length:", codeLen)

	_ = cg
}
