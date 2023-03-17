package application

import (
	"context"
	"github.com/escalopa/fingo/token/internal/core"
)

type TokenRepository interface {
	GetTokenPayload(ctx context.Context, accessToken string) (*core.TokenPayload, error)
}

type Validator interface {
	Validate(params interface{}) error
}
