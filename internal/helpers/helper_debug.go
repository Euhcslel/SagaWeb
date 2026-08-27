// Package helpers содержит вспомогательные функции для обработки ошибок в http-хендлерах
//go:build debug

package helpers

import (
	"log"
	"net/http"
)

// WriteError выводит реальную ошибку на страницу и в логи
func WriteError(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
	log.Println(err)
}
