package hasher

import (
	"context"

	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/lordvidex/errs"
	"golang.org/x/crypto/bcrypt"
)

const minPasswordLen = 8

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (c *BcryptHasher) Hash(ctx context.Context, password string) (string, error) {
	_, span := tracer.Tracer().Start(ctx, "BcryptHasher.Hash")
	defer span.End()
	if password == "" {
		return "", errs.B().Code(errs.InvalidArgument).Msgf("password length is less than %d", minPasswordLen).Err()
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (c *BcryptHasher) Compare(ctx context.Context, hash, password string) bool {
	_, span := tracer.Tracer().Start(ctx, "BcryptHasher.Compare")
	defer span.End()
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
