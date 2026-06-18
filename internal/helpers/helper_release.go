//go:build !debug

package helpers

import (
	"log"
	"net/http"
)

func WriteError(w http.ResponseWriter, err error, code int) {
	log.Println(err)
	msg := err.Error()
	if code >= 500 {
		msg = http.StatusText(code)
	}
	http.Error(w, msg, code)
}
