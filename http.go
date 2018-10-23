package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func httpInit(listen string) {
	r := http.NewServeMux()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Alive and well\n")
	})

	r.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:    listen,
		Handler: r,
	}

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP: Failed to listen on %s: %s", listen, err)
	}
}
