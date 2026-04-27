package errors

import "errors"

var (
	ErrForbidden          = errors.New("forbidden")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrNotFound           = errors.New("not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidOrderStatus = errors.New("invalid order status")
	ErrInvalidGateType    = errors.New("invalid gate type")
)
