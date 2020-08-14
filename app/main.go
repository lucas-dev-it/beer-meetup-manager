package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
)

var port = internal.GetEnv("PORT", "3999")

func main() {
	builder := builder{}
	err := builder.injectDependencies()
	if err != nil {
		log.Fatalf("and error ocurrend during startup, got: %v", err)
	}
	defer builder.close()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      builder.controller,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Printf("Starting HTTP Server. Listening at %q", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("%v", err)
		} else {
			log.Println("Server closed!")
		}
	}()

	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)
	sig := <-sigquit
	log.Printf("caught sig: %+v", sig)
	log.Printf("Gracefully shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("Unable to shut down server: %v", err)
	} else {
		log.Println("Server stopped")
	}
}
