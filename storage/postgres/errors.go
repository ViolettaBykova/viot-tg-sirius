package postgres

import "errors"

var (
	ErrTxNotFound = errors.New("tx not found in context")
)
