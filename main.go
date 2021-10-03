package main

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	connectioncontext "github.com/Dafaque/constanta/connection_context"
	"github.com/Dafaque/constanta/handlers"
)

func main() {
	var connections int = 0
	var mu sync.Mutex

	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("/", handlers.Root)

	var server *http.Server = &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		//? Time out with a margin for sterilization
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  10 * time.Second,
		//? Active connections counter
		ConnState: func(c net.Conn, cs http.ConnState) {
			switch cs {
			case http.StateActive:
				mu.Lock()
				connections++
				mu.Unlock()
			case http.StateClosed:
				mu.Lock()
				connections--
				mu.Unlock()
			default:
			}
		},
		//? Sending to every handler info about currently active connections
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return context.WithValue(
				ctx,
				connectioncontext.ConnectionsKeyCounter,
				connections,
			)
		},
	}

	serveErr := server.ListenAndServe()
	println("serve err:", serveErr)

	//? Write and Idle timeouts are setted in server config
	//? Cannot test on my machine
	server.Shutdown(context.Background())
}
