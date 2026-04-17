//go:build !debug

package helpers

import (
	"log"
	"net/http"
)

func WriteError(w http.ResponseWriter, err error, code int) {
	http.Error(w, http.StatusText(code), code)
	log.Println(err)
}
