package db_util

import (
	"errors"
	"github.com/jchavannes/jgo/jerr"
)

const (
	DuplicateEntryErrorMessage    = "Duplicate entry"
	InvalidConnectionErrorMessage = "invalid connection"
	NoRowsInResultSetErrorMessage = "sql: no rows in result set"
	DatabaseClosedErrorMessage    = "sql: database is closed"
	RecordNotFoundErrorMessage    = "record not found"
	LockTimeoutErrorMessage       = "Error 1205: Lock wait timeout exceeded; try restarting transaction"
	TableDoesntExistErrorMessage  = "Error 1146: Table '"
)

var (
	InvalidConnectionError = errors.New(InvalidConnectionErrorMessage)
	RecordNotFoundError    = errors.New(RecordNotFoundErrorMessage)
)

func IsDuplicateEntryError(e error) bool {
	return jerr.HasErrorPart(e, DuplicateEntryErrorMessage)
}

func IsRecordNotFoundError(e error) bool {
	return jerr.HasError(e, RecordNotFoundErrorMessage)
}

func IsConnectionError(e error) bool {
	return IsLockTimeoutError(e) || IsInvalidConnectionError(e) || IsDatabaseClosedError(e)
}

func IsLockTimeoutError(e error) bool {
	return jerr.HasError(e, LockTimeoutErrorMessage)
}

func IsInvalidConnectionError(e error) bool {
	return jerr.HasError(e, InvalidConnectionErrorMessage)
}

func IsNoRowsInResultSetError(e error) bool {
	return jerr.HasError(e, NoRowsInResultSetErrorMessage)
}

func IsDatabaseClosedError(e error) bool {
	return jerr.HasError(e, DatabaseClosedErrorMessage)
}

func IsTableDoesntExistError(e error) bool {
	return jerr.HasErrorPart(e, TableDoesntExistErrorMessage)
}
