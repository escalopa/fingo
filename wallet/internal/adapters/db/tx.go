package db

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/lordvidex/errs"
)

var (
	errorUniqueViolation = func(err error, msg string) error {
		return errs.B(err).Code(errs.AlreadyExists).Msg(msg).Err()
	}
	errorNotFound = func(err error, msg string) error {
		return errs.B(err).Code(errs.NotFound).Msg(msg).Err()
	}
	errorQuery = func(err error, msg string) error {
		return errs.B(err).Code(errs.Internal).Msg(msg).Err()
	}
	errorTxNotStarted = func(err error) error {
		return errs.B(err).Code(errs.Internal).Msg("transaction not started").Err()
	}
	errorTxNotCommitted = func(err, err2 error) error {
		return errs.B(err).Code(errs.Internal).Details(err2).Msg("transaction not committed").Err()
	}
	errorTxNotRolledBack = func(err, err2 error) error {
		return errs.B(err).Code(errs.Internal).Details(err2).Msg("transaction not rolled back").Err()
	}
	errorRollbackUnsupported = errs.B().Msg("rollback not supported for deposit & withdrawals transactions").Err()
)

func deferTx(tx *sql.Tx, err error) error {
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return errorTxNotRolledBack(err, err2)
		}
		return err
	}
	if err2 := tx.Commit(); err2 != nil {
		return errorTxNotCommitted(err, err2)
	}
	return nil
}

func IsUniqueViolationError(err error) bool {
	er, ok := err.(*pq.Error)
	return ok && er.Code == "23505"
}

func IsNotFoundError(err error) bool {
	if err == sql.ErrNoRows {
		return true
	}
	er, ok := err.(*pq.Error)
	if ok {
		return er.Code == "20000"
	}
	return false
}
