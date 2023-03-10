package application

import (
	"context"
	"github.com/escalopa/fingo/auth/internal/adapters/token"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg core.CreateUserParams) error
	GetUserByID(ctx context.Context, id uuid.UUID) (core.User, error)
	GetUserByEmail(ctx context.Context, email string) (core.User, error)
	GetUserByUsername(ctx context.Context, username string) (core.User, error)
	SetUserIsVerified(ctx context.Context, arg core.SetUserIsVerifiedParams) error
	ChangeUserEmail(ctx context.Context, arg core.ChangeUserEmailParams) error
	ChangePassword(ctx context.Context, arg core.ChangePasswordParams) error
	DeleteUserByID(ctx context.Context, id uuid.UUID) error
}

type SessionRepository interface {
	CreateSession(ctx context.Context, arg core.CreateSessionParams) error
	GetSessionByID(ctx context.Context, id uuid.UUID) (core.Session, error)
	GetUserSessions(ctx context.Context, userID uuid.UUID) ([]core.Session, error)
	GetUserDevices(ctx context.Context, userID uuid.UUID) ([]core.UserDevice, error)
	SetSessionIsBlocked(ctx context.Context, arg core.SetSessionIsBlockedParams) error
	DeleteSessionByID(ctx context.Context, sessionID uuid.UUID) error
}

type PasswordHasher interface {
	Hash(password string) (hashedPassword string, err error)
	Compare(password, hash string) (isSamePassword bool)
}

type TokenGenerator interface {
	GenerateAccessToken(param token.GenerateTokenParam) (token string, err error)
	GenerateRefreshToken(param token.GenerateTokenParam) (token string, err error)
	VerifyToken(token string) (user core.User, sessionID uuid.UUID, err error)
}

type Validator interface {
	Validate(s interface{}) (err error)
}
