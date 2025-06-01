package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/foto-leistenschneider/admin-panel/internal/config"
)

var (
	server *http.Server
)

func Start() {
	ln, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to listen: %v", err)
		os.Exit(1)
	}

	tcpListener := ln.(*net.TCPListener)

	log.Info("Listening for TCP connections", "address", tcpListener.Addr().String())

	registerRoutes()

	server = &http.Server{
		Addr:    config.ServerAddress,
		Handler: http.DefaultServeMux,
	}

	if err := server.Serve(tcpListener); err != nil {
		log.Error("Server error", "error", err)
	}
}

func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error", "error", err)
	}
}
