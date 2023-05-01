package token

import (
	"context"
	"time"

	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/escalopa/fingo/pkg/tracer"
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
func (pt *PasetoTokenizer) GenerateAccessToken(ctx context.Context, params core.GenerateTokenParam) (string, error) {
	ctx, span := tracer.Tracer().Start(ctx, "PasetoTokenizer.GenerateAccessToken")
	defer span.End()
	return pt.generateToken(ctx, params, pt.atd)
}

// GenerateRefreshToken Creates a new refresh token
func (pt *PasetoTokenizer) GenerateRefreshToken(ctx context.Context, params core.GenerateTokenParam) (string, error) {
	ctx, span := tracer.Tracer().Start(ctx, "PasetoTokenizer.GenerateRefreshToken")
	defer span.End()
	return pt.generateToken(ctx, params, pt.rtd)
}

// generateToken Create a new token with user, sessionID, exp(Token life duration)
func (pt *PasetoTokenizer) generateToken(ctx context.Context, params core.GenerateTokenParam, exp time.Duration) (string, error) {
	_, span := tracer.Tracer().Start(ctx, "PasetoTokenizer.generateToken")
	defer span.End()
	// Create userToken struct instance
	ut := core.TokenPayload{
		UserID:    params.UserID,
		SessionID: params.SessionID,
		ClientIP:  params.ClientIP,
		UserAgent: params.UserAgent,
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

// DecryptToken decrypts the token to get `TokenPayload` & verifies that token hasn't expired
func (pt *PasetoTokenizer) DecryptToken(ctx context.Context, token string) (core.TokenPayload, error) {
	_, span := tracer.Tracer().Start(ctx, "PasetoTokenizer.DecryptToken")
	defer span.End()
	// Decrypt token
	var payload core.TokenPayload
	err := pt.p.Decrypt(token, pt.sk, &payload, nil)
	if err != nil {
		return core.TokenPayload{}, errs.B(err).Code(errs.InvalidArgument).
			Msg("failed to decrypt token, invalid token").Err()
	}
	return payload, nil
}
