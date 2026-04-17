//go:build debug

package helpers

import (
	"log"
	"net/http"
)

func WriteError(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
	log.Println(err)
}
