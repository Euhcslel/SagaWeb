package errors

import "errors"

var (
	ErrForbidden           = errors.New("forbidden")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrNotFound            = errors.New("not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidOrderStatus  = errors.New("invalid order status")
	ErrInvalidGateType     = errors.New("invalid gate type")
	ErrInvalidDocumentType = errors.New("invalid document type")
	ErrInvalidTableType    = errors.New("invalid table type")
	ErrBadRequest          = errors.New("bad user request")
	ErrNoDealerContract    = errors.New("у дилера отсутствует договор")
)
