package cache

import (
	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/errs"
)

// NewRedisClient creates a new redis client from the given url
func NewRedisClient(url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, errs.B(err).Code(errs.InvalidArgument).Msg("failed to parse cache url").Err()
	}
	return redis.NewClient(opts), nil
}
