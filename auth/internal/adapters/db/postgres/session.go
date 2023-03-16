package mypostgres

import (
	"context"
	"database/sql"
	"time"

	db "github.com/escalopa/fingo/auth/internal/adapters/db/postgres/sqlc"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
)

type SessionRepository struct {
	q   db.Querier
	std time.Duration // sessions time duration
}

// NewSessionRepository creates a new sessions repository with the given connection
func NewSessionRepository(conn *sql.DB, opts ...func(*SessionRepository)) (*SessionRepository, error) {
	if conn == nil {
		return nil, errs.B().Msg("passed connection is nil").Err()
	}
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
func (sr *SessionRepository) CreateSession(ctx context.Context, params core.CreateSessionParams) error {
	err := sr.q.CreateSession(ctx, db.CreateSessionParams{
		ID:           params.ID,
		UserID:       params.UserID,
		AccessToken:  params.AccessToken,
		RefreshToken: params.RefreshToken,
		UserAgent:    params.UserDevice.UserAgent,
		ClientIp:     params.UserDevice.ClientIP,
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
		return core.Session{}, errs.B(err).Code(errs.Internal).Msgf("failed to get user sessions with id, id: %s", id).Err()
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

// UpdateSessionTokens returns a sessions by its refresh token value
func (sr *SessionRepository) UpdateSessionTokens(ctx context.Context, params core.UpdateSessionTokenParams) error {
	rows, err := sr.q.UpdateSessionTokens(ctx, db.UpdateSessionTokensParams{
		ID:           params.ID,
		AccessToken:  params.AccessToken,
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

// DeleteSessionByID deletes a sessions by its id
func (sr *SessionRepository) DeleteSessionByID(ctx context.Context, sessionID uuid.UUID) error {
	rows, err := sr.q.DeleteSessionByID(ctx, sessionID)
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("failed to delete sessions with the given sID, sID: %s", sessionID.String()).Err()
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
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		UserDevice:   core.UserDevice{ClientIP: session.ClientIp, UserAgent: session.UserAgent},
		ExpiresAt:    session.ExpiresAt,
		UpdatedAt:    session.UpdatedAt,
	}
}
