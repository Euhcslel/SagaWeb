package errors

import "errors"

var (
	ErrForbidden = errors.New("forbidden")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrNotFound           = errors.New("not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidGateType    = errors.New("invalid gate type")
)
