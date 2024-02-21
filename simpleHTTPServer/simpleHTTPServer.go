package simpleHTTPServer

import (
	"io"
	"log"
	"net/http"
)

func StartBackendServer() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, "Hello there !")
    })

    err := http.ListenAndServe(":6969", nil)
    if err != nil {
        log.Fatal("Server error: ", err)
    }
}
