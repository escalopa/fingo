package hasher

import (
	"github.com/lordvidex/errs"
	"golang.org/x/crypto/bcrypt"
)

const minPasswordLen = 8

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (c *BcryptHasher) Hash(password string) (string, error) {
	if password == "" {
		return "", errs.B().Code(errs.InvalidArgument).Msgf("password length is less than %d", minPasswordLen).Err()
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (c *BcryptHasher) Compare(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
