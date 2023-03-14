package mypostgres

import (
	"context"
	"database/sql"
	db "github.com/escalopa/fingo/auth/internal/adapters/db/postgres/sqlc"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
	"time"
)

type SessionRepository struct {
	q   db.Querier
	std time.Duration // sessions time duration
}

// NewSessionRepository creates a new sessions repository with the given connection
func NewSessionRepository(conn *sql.DB, opts ...func(*SessionRepository)) (*SessionRepository, error) {
	sr := &SessionRepository{q: db.New(conn)}
	for _, opt := range opts {
		opt(sr)
	}
	if sr.std == 0 {
		return nil, errs.B().Msg("sessions time duration is not set").Err()
	}
	return sr, nil
}

// WithSessionDuration is a functional option to set the sessions time duration
func WithSessionDuration(d time.Duration) func(*SessionRepository) {
	return func(sr *SessionRepository) {
		sr.std = d
	}
}

// CreateSession creates a new sessions for a user
func (sr *SessionRepository) CreateSession(ctx context.Context, arg core.CreateSessionParams) error {
	err := sr.q.CreateSession(ctx, db.CreateSessionParams{
		ID:           arg.ID,
		UserID:       arg.UserID,
		RefreshToken: arg.RefreshToken,
		UserAgent:    arg.UserDevice.UserAgent,
		ClientIp:     arg.UserDevice.ClientIP,
		ExpiresAt:    time.Now().Add(sr.std),
	})
	if err != nil {
		if isUniqueViolationError(err) {
			return errs.B(err).Code(errs.AlreadyExists).Msg("violation to unique keys in sessions table").Err()
		}
		return errs.B(err).Code(errs.Internal).Msg("failed to create new sessions").Err()
	}
	return nil
}

// GetSessionByID returns a sessions by its id
func (sr *SessionRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (core.Session, error) {
	session, err := sr.q.GetSessionByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return core.Session{}, errs.B(err).Code(errs.NotFound).Msgf("no sessions found with the given id, id: %s", id).Err()
		}
		return core.Session{}, errs.B(err).Msgf("failed to get user sessions with id, id: %s", id).Err()
	}
	if session.ID == uuid.Nil {

	}
	return fromDbSessionToCore(session), nil
}

// GetUserSessions returns all sessions owned by a single user
func (sr *SessionRepository) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]core.Session, error) {
	sessions, err := sr.q.GetUserSessions(ctx, userID)
	if err != nil {
		return nil, errs.B(err).Code(errs.Internal).Msg("failed to get user sessions").Err()
	}
	// map from []db.Session to []core.Session
	coreSessions := make([]core.Session, len(sessions))
	for i, v := range sessions {
		coreSessions[i] = fromDbSessionToCore(v)
	}
	return coreSessions, nil
}

// UpdateSessionRefreshToken returns a sessions by its refresh token value
func (sr *SessionRepository) UpdateSessionRefreshToken(ctx context.Context, params core.UpdateSessionRefreshTokenParams) error {
	rows, err := sr.q.UpdateSessionRefreshToken(ctx, db.UpdateSessionRefreshTokenParams{
		ID:           params.ID,
		RefreshToken: params.RefreshToken,
		ExpiresAt:    time.Now().Add(sr.std),
	})
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("failed to set sessions refresh token").Err()
	}
	if rows == 0 {
		return errs.B().Code(errs.NotFound).Msgf("no sessions found with the given id, id: %s", params.ID.String()).Err()
	}
	return nil
}

// SetSessionIsBlocked sets the IsBlocked value of a sessions to true or false
func (sr *SessionRepository) SetSessionIsBlocked(ctx context.Context, arg core.SetSessionIsBlockedParams) error {
	// Update sessions value by setting IsBlocked value
	rows, err := sr.q.SetSessionIsBlocked(ctx, db.SetSessionIsBlockedParams{
		ID:        arg.ID,
		IsBlocked: arg.IsBlocked,
	})
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("failed to set sessions is_blocked value").Err()
	}
	if rows == 0 {
		return errs.B().Code(errs.NotFound).Msgf("no sessions found with the given id, id: %s", arg.ID.String()).Err()
	}
	return nil
}

// DeleteSessionByID deletes a sessions by its id
func (sr *SessionRepository) DeleteSessionByID(ctx context.Context, sessionID uuid.UUID) error {
	rows, err := sr.q.DeleteSessionByID(ctx, sessionID)
	if err != nil {
		return errs.B(err).Msg("failed to delete sessions with the given sID, sID: %s", sessionID.String()).Err()
	}
	if rows == 0 {
		return errs.B(err).Code(errs.NotFound).Msgf("no sessions found with the given id, id: %s", sessionID.String()).Err()
	}
	return nil
}

func fromDbSessionToCore(session db.Session) core.Session {
	return core.Session{
		ID:           session.ID,
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		UserDevice:   core.UserDevice{ClientIP: session.ClientIp, UserAgent: session.UserAgent},
		IsBlocked:    session.IsBlocked,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}
}
