package errors

import "errors"

var (
	ErrForbidden           = errors.New("Доступ запрещён")
	ErrUnauthorized        = errors.New("Необходима авторизация")
	ErrNotFound            = errors.New("Запись не найдена")
	ErrInvalidCredentials  = errors.New("Неверный логин или пароль")
	ErrInvalidOrderStatus  = errors.New("Недопустимый статус заказа")
	ErrInvalidGateType     = errors.New("Недопустимый тип ворот")
	ErrInvalidDocumentType = errors.New("Недопустимый тип документа")
	ErrInvalidTableType    = errors.New("Недопустимый тип таблицы")
	ErrBadRequest          = errors.New("Некорректный запрос")
	ErrNoDealerContract        = errors.New("У дилера отсутствует договор")
	ErrIncompleteCompanyDetails = errors.New("Необходимо заполнить реквизиты компании: ИНН, КПП, ОГРН и адрес")
)
