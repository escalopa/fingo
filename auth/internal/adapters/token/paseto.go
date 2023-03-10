package token

import (
	"github.com/google/uuid"
	"time"

	ac "github.com/escalopa/fingo/auth/internal/core"
	"github.com/lordvidex/errs"

	"github.com/o1egl/paseto"
)

const minSecretKeyLen = 32

type PasetoTokenizer struct {
	p   *paseto.V2
	sk  []byte
	atd time.Duration // access token's duration
	rtd time.Duration // refresh token's duration
}

// NewPaseto Creates a new instance of PasetoTokenizer
func NewPaseto(secretKey string, atd, rtd time.Duration) (*PasetoTokenizer, error) {
	// Check that secret key is exactly equal to `minSecretKeyLen`
	if len(secretKey) < minSecretKeyLen {
		return nil, errs.B().
			Code(errs.InvalidArgument).
			Msgf("secretKey len is less than the min value %d", minSecretKeyLen).
			Err()
	}
	// Return a new paseto tokenizer
	return &PasetoTokenizer{
		p:   paseto.NewV2(),
		sk:  []byte(secretKey),
		atd: atd,
		rtd: rtd,
	}, nil
}

// GenerateAccessToken Creates a new access token
func (pt *PasetoTokenizer) GenerateAccessToken(gtp GenerateTokenParam) (string, error) {
	return pt.generateToken(gtp.User, gtp.SessionID, pt.atd)
}

// GenerateRefreshToken Creates a new refresh token
func (pt *PasetoTokenizer) GenerateRefreshToken(gtp GenerateTokenParam) (string, error) {
	return pt.generateToken(gtp.User, gtp.SessionID, pt.rtd)
}

// generateToken Create a new token with user, sessionID, exp(Token life duration)
func (pt *PasetoTokenizer) generateToken(u ac.User, sID uuid.UUID, exp time.Duration) (string, error) {
	// Create userToken struct instance
	ut := ac.UserToken{
		User:      u,
		SessionID: sID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(exp),
	}
	// Encrypt userToken
	token, err := pt.p.Encrypt(pt.sk, ut, nil)
	if err != nil {
		return "", errs.B(err).Code(errs.Internal).Msg("failed to create token").Err()
	}
	return token, nil
}

// VerifyToken decrypts the token to get `UserToken` & verifies that token hasn't expired
func (pt *PasetoTokenizer) VerifyToken(token string) (ac.User, uuid.UUID, error) {
	// Decrypt token
	var ut ac.UserToken
	err := pt.p.Decrypt(token, pt.sk, &ut, nil)
	if err != nil {
		return ac.User{}, uuid.UUID{}, errs.B(err).Code(errs.InvalidArgument).
			Msg("failed to decrypt token, invalid token").Err()
	}
	// Check whether the token has expired
	if time.Now().After(ut.ExpiresAt) {
		return ac.User{}, uuid.UUID{}, errs.B().Code(errs.Unauthenticated).Msg("token expired").Err()
	}
	return ut.User, ut.SessionID, nil
}
