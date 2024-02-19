package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var httpTransport = http.DefaultTransport

func main() {
	fmt.Println("\t>>> SimpleGoProxy >>>")
	fmt.Println()

	// "custom" server that handles differently simple http requests and https
	server := http.Server{
		Addr: ":1337",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
                log.Println("Handleing Tunneling to: ", r.URL)
                handleTunneling(w, r)
			} else {
                log.Println("Handling simple request to: ", r.URL)
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

	// Copy headers and their values
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Copy status code and body
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
    // Check if we can take over the connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported!", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(dest io.WriteCloser, source io.ReadCloser) {
	defer dest.Close()
	defer source.Close()

	io.Copy(dest, source)
}
