package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func StartReverseProxy() {
    dest, err := url.Parse("http://localhost:6969")
	if err != nil {
		log.Fatal("URL Parsing: ", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(dest)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        proxy.ServeHTTP(w, r)
    })

    server := &http.Server{
        Addr: ":1337",
        Handler: http.DefaultServeMux,
    }

    err = server.ListenAndServe()
    if err != nil {
        log.Fatal("Error staring server: ", err)
    }
}
