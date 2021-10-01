package main

import (
	"net/http"
	"time"

	"github.com/Dafaque/constanta/handlers"
)

func main() {
	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("/", handlers.Root)

	var server *http.Server = &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       30 * time.Second,
	}

	panic(server.ListenAndServe())
	// TODO: GFST
}
