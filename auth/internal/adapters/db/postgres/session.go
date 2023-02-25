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
	q   db.Querier
	std time.Duration // session time duration
}

func NewSessionRepository(conn *sql.DB, std time.Duration) *SessionRepository {
	return &SessionRepository{q: db.New(conn), std: std}
}

func (sr *SessionRepository) CreateSession(ctx context.Context, arg core.CreateSessionParams) error {
	err := sr.q.CreateSession(ctx, db.CreateSessionParams{
		ID:           arg.ID,
		UserID:       arg.UserID,
		RefreshToken: arg.RefreshToken,
		UserAgent:    sql.NullString{String: arg.UserAgent},
		ClientIp:     sql.NullString{String: arg.ClientIp},
		ExpiresAt:    time.Now().Add(sr.std),
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
	// Get a specific session by its id
	session, err := sr.q.GetSessionByID(ctx, id)
	// Handle error
	if err != nil {
		if isNotFoundError(err) {
			errs.B(err).Msgf("no user sessions found with the given id, id: %s", id.String())
		}
		return core.Session{}, errs.B(err).Msgf("failed to get user session with id, id: %s", id).Err()
	}
	return fromDbSessionToCore(session), nil
}

func (sr *SessionRepository) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]core.Session, error) {
	// Get all user sessions owned by a single user
	sessions, err := sr.q.GetUserSessions(ctx, userID)
	// Handle error
	if err != nil {
		if isNotFoundError(err) {
			errs.B(err).Msgf("no user sessions found with the given userID, id: %s", userID.String())
		}
	}
	// map from []db.Session to []core.Session
	coreSessions := make([]core.Session, len(sessions))
	for i, v := range sessions {
		coreSessions[i] = fromDbSessionToCore(v)
	}
	return coreSessions, nil
}

func (sr *SessionRepository) GetUserDevices(ctx context.Context, userID uuid.UUID) ([]core.UserDevice, error) {
	// Get all user devices (Client-IP, User-Agent)
	devices, err := sr.q.GetUserDevices(ctx, userID)
	// Handel error
	if err != nil {
		if isNotFoundError(err) {
			errs.B(err).Msgf("no user devices found with the given userID, userID: %s", userID.String())
		}
		return nil, errs.B(err).Msg("failed to get user devices").Err()
	}
	// map from []db.UserDevice to []core.UserDevice
	coreDevices := make([]core.UserDevice, len(devices))
	for i, v := range devices {
		coreDevices[i] = fromDbDeviceToCore(v)
	}
	return nil, nil
}

func (sr *SessionRepository) SetSessionIsBlocked(ctx context.Context, arg core.SetSessionIsBlockedParams) error {
	// Update session value by setting IsBlocked value
	err := sr.q.SetSessionIsBlocked(ctx, db.SetSessionIsBlockedParams{
		ID:        arg.ID,
		IsBlocked: arg.IsBlocked,
	})
	// Handle error
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
		UserAgent:    session.UserAgent.String,
		ClientIp:     session.ClientIp.String,
		IsBlocked:    session.IsBlocked,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}
}

func fromDbDeviceToCore(row db.GetUserDevicesRow) core.UserDevice {
	return core.UserDevice{
		UserAgent: row.UserAgent.String,
		ClientIP:  row.ClientIp.String,
	}
}
