package mypostgres

import (
	"context"
	"database/sql"
	db "github.com/escalopa/gochat/auth/internal/adapters/db/postgres/sqlc"
	"github.com/escalopa/gochat/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
	"time"
)

type SessionRepository struct {
	db *sql.DB
	q  db.Querier
}

func NewSessionRepository(conn *sql.DB) (*UserRepository, error) {
	return &UserRepository{db: conn}, nil
}

type CreateSessionParams struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	RefreshToken string
	UserAgent    string
	ClientIp     string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type SetSessionIsBlockedParams struct {
	ID        uuid.UUID
	IsBlocked bool
}

func (sr *SessionRepository) CreateSession(ctx context.Context, arg CreateSessionParams) error {
	err := sr.q.CreateSession(ctx, db.CreateSessionParams{
		ID:           arg.ID,
		UserID:       arg.UserID,
		RefreshToken: arg.RefreshToken,
		UserAgent:    arg.UserAgent,
		ClientIp:     arg.ClientIp,
		ExpiresAt:    arg.ExpiresAt,
		CreatedAt:    arg.CreatedAt,
	})
	if err != nil {
		if isUniqueViolationError(err) { // Unique violation
			return errs.B(err).Msg("violation to unique keys in sessions table").Err()
		}
		return errs.B(err).Msg("failed to create new session").Err()
	}
	return nil
}

func (sr *SessionRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (core.Session, error) {
	session, err := sr.q.GetSessionByID(ctx, id)
	if err != nil {
		if isNotFoundError(err) {
			errs.B(err).Msgf("no user sessions found with the given id, id: %s", id.String())
		}
		return core.Session{}, errs.B(err).Msgf("failed to get user session with id, id: %s", id).Err()
	}
	return fromDbSessionToCore(session), nil
}

func (sr *SessionRepository) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]core.Session, error) {
	sessions, err := sr.q.GetUserSessions(ctx, userID)
	if err != nil {
		if isNotFoundError(err) {
			errs.B(err).Msgf("no user sessions found with the given userID, id: %s", userID.String())
		}
	}
	coreSessions := make([]core.Session, len(sessions))
	for i, v := range sessions {
		coreSessions[i] = fromDbSessionToCore(v)
	}
	return coreSessions, nil
}

func (sr *SessionRepository) GetUserDevices(ctx context.Context, userID uuid.UUID) ([]core.UserDevice, error) {
	devices, err := sr.q.GetUserDevices(ctx, userID)
	if err != nil {
		if isNotFoundError(err) {
			errs.B(err).Msgf("no user devices found with the given userID, userID: %s", userID.String())
		}
		return nil, errs.B(err).Msg("failed to get user devices").Err()
	}
	coreDevices := make([]core.UserDevice, len(devices))
	for i, v := range devices {
		coreDevices[i] = fromDbDeviceToCore(v)
	}
	return nil, nil
}

func (sr *SessionRepository) SetSessionIsBlocked(ctx context.Context, arg SetSessionIsBlockedParams) error {
	err := sr.q.SetSessionIsBlocked(ctx, db.SetSessionIsBlockedParams{
		ID:        arg.ID,
		IsBlocked: arg.IsBlocked,
	})
	if err != nil {
		if isNotFoundError(err) {
			errs.B(err).Msgf("no session found with the given id, id: %s", arg.ID.String())
		}
		return errs.B(err).Msg("failed to set session is_blocked value").Err()
	}
	return nil
}

func fromDbSessionToCore(session db.Session) core.Session {
	return core.Session{
		ID:           session.ID,
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		ClientIp:     session.ClientIp,
		IsBlocked:    session.IsBlocked,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}
}

func fromDbDeviceToCore(row db.GetUserDevicesRow) core.UserDevice {
	return core.UserDevice{
		UserAgent: row.UserAgent,
		ClientIP:  row.ClientIp,
	}
}
