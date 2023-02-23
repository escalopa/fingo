package token

import (
	"time"

	ac "github.com/escalopa/gochat/auth/internal/core"
	"github.com/lordvidex/errs"

	"github.com/o1egl/paseto"
)

const minSecretKeyLen = 32

type PasetoTokenizer struct {
	p   *paseto.V2
	sk  []byte
	ate time.Duration // access token expiration time
}

func NewPaseto(secretKey string, accessTokenDuration time.Duration) (*PasetoTokenizer, error) {
	if len(secretKey) < minSecretKeyLen {
		return nil, errs.B().
			Code(errs.InvalidArgument).
			Msgf("secretKet len is less than the min value %d", minSecretKeyLen).
			Err()
	}

	return &PasetoTokenizer{
		p:   paseto.NewV2(),
		sk:  []byte(secretKey),
		ate: accessTokenDuration,
	}, nil
}

func (pt *PasetoTokenizer) GenerateToken(u ac.User) (string, error) {
	ut := ac.UserToken{
		User:      u,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(pt.ate),
	}
	token, err := pt.p.Encrypt(pt.sk, ut, nil)
	if err != nil {
		return "", err
	}
	return token, err
}

func (pt *PasetoTokenizer) VerifyToken(token string) (ac.User, error) {
	var ut ac.UserToken
	err := pt.p.Decrypt(token, pt.sk, &ut, nil)
	if err != nil {
		return ac.User{}, err
	}
	if time.Now().After(ut.ExpiresAt) {
		return ac.User{}, errs.B().Code(errs.Unauthenticated).Msg("token expired").Err()
	}
	return ut.User, nil
}
