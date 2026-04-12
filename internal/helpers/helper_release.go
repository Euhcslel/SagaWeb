//go:build !debug

package helpers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var LogFile *os.File
var mutex sync.Mutex

func WriteError(w http.ResponseWriter, err error, code int) {
	http.Error(w, http.StatusText(code), code)
	log.Println(err)

	mutex.Lock()
	defer mutex.Unlock()

	if LogFile == nil {
		return
	}

	line := fmt.Sprintf("%s http %d: %v\n",
		time.Now().Format(time.RFC3339), code, err)

	if _, err = LogFile.WriteString(line); err != nil {
		log.Println(err)
	}
}
