//go:build !debug

package helpers

import (
	"log"
	"net/http"
)

// WriteError в release выводит ошибки на страницу,
// Но скрывает точную причину ошибок на сервере.
func WriteError(w http.ResponseWriter, err error, code int) {
	log.Println(err)
	msg := err.Error()
	if code >= 500 {
		msg = http.StatusText(code)
	}
	http.Error(w, msg, code)
}
