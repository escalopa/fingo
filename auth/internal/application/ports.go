package application

import (
	"context"

	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
)

// UserRepository is an interface for interacting with users in the database
type UserRepository interface {
	CreateUser(ctx context.Context, arg core.CreateUserParams) error
	GetUserByID(ctx context.Context, id uuid.UUID) (core.User, error)
	GetUserByEmail(ctx context.Context, email string) (core.User, error)
}

// SessionRepository is an interface for interacting with sessions in the database
type SessionRepository interface {
	CreateSession(ctx context.Context, arg core.CreateSessionParams) error
	GetSessionByID(ctx context.Context, id uuid.UUID) (core.Session, error)
	GetUserSessions(ctx context.Context, userID uuid.UUID) ([]core.Session, error)
	UpdateSessionTokens(ctx context.Context, params core.UpdateSessionTokenParams) error
	DeleteSessionByID(ctx context.Context, sessionID uuid.UUID) error
}

// TokenRepository is an interface for interacting with tokens in cache
type TokenRepository interface {
	Store(ctx context.Context, token string, params core.TokenPayload) error
	Delete(ctx context.Context, token string) error
}

// PasswordHasher is an interface for hashing and comparing passwords
type PasswordHasher interface {
	Hash(password string) (hashedPassword string, err error)
	Compare(password, hash string) (isSamePassword bool)
}

// TokenGenerator is an interface for generating and verifying tokens
type TokenGenerator interface {
	GenerateAccessToken(params core.GenerateTokenParam) (token string, err error)
	GenerateRefreshToken(params core.GenerateTokenParam) (token string, err error)
	DecryptToken(token string) (params core.TokenPayload, err error)
}

// MessageProducer is an interface for sending messages to a queue
type MessageProducer interface {
	SendNewSignInSessionMessage(ctx context.Context, params core.SendNewSignInSessionParams) error
}

// Validator is an interface for validating structs using tags
type Validator interface {
	Validate(s interface{}) (err error)
}
