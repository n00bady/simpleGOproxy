package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

var httpTransport = http.DefaultTransport

func main() {
	fmt.Println(">>> SimpleGoProxy >>>")

	server := http.Server{
		Addr: ":1337",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				log.Println("Tunneling not implemented yet...")
			} else {
				handleSimpleHttp(w, r)
			}
		}),
	}

	log.Println("Starting proxy server on port 1337")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting proxy server: ", err)
	}
}

func handleSimpleHttp(w http.ResponseWriter, r *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, "Error sending proxy request: ", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
