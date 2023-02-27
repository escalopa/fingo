package token

import (
	"github.com/escalopa/gochat/auth/internal/core"
	"github.com/google/uuid"
)

type GenerateTokenParam struct {
	User      core.User
	SessionID uuid.UUID
}
