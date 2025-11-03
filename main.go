// Package main provides a proxy server for BOHECO2 bill inquiry API.
// This server acts as a middleware to handle CORS and session management
// for the BOHECO2 bill inquiry frontend application.
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kennyclark/boheco2-bill-inquiry-server/server"
	"github.com/kennyclark/boheco2-bill-inquiry-server/session"
)

// main starts the HTTP proxy server on port 3000.
// The server handles bill inquiry requests and forwards them to the upstream API.
func main() {
	// Initialize server configuration
	port, apiBaseURL := server.Initialize()

	// Configure session package with the API base URL
	session.SetAPIBaseURL(apiBaseURL)

	// Handle bill inquiry API endpoint
	http.HandleFunc("/api/v1/bill", func(w http.ResponseWriter, r *http.Request) {
		slog.Info(fmt.Sprintf("Received %s request to /api/v1/bill", r.Method))

		// Get the origin from the request
		requestOrigin := r.Header.Get("Origin")

		// Check if the origin is allowed and set CORS headers accordingly
		if server.IsOriginAllowed(requestOrigin) {
			w.Header().Set("Access-Control-Allow-Origin", requestOrigin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Check if we have a valid session token, initialize if needed
		if !session.HasToken() {
			if err := session.Initialize(); err != nil {
				http.Error(w, "Failed to init session: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Read the incoming request body to forward to upstream API
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// Create request to upstream API with the same body
		apiURL := apiBaseURL + "/api/v1/bill"
		req, err := http.NewRequest("POST", apiURL, bytes.NewReader(bodyBytes))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set required headers to mimic requests from BOHECO2 website
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Origin", "https://www.boheco2.com.ph")
		req.Header.Set("Referer", "https://www.boheco2.com.ph")
		req.Header.Set("Content-Type", "application/json")

		// Execute the request using the session client (with cookies)
		resp, err := session.Client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Forward the upstream response status and body to the client
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	// Create HTTP server with timeouts for better resource management
	srv := &http.Server{
		Addr:         "0.0.0.0:" + port,
		Handler:      nil, // Use default ServeMux
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine so it doesn't block signal handling
	go func() {
		slog.Info(fmt.Sprintf("BOHECO 2 API Proxy Server listening on port: %s", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(fmt.Sprintf("Failed to start server: %v", err))
			os.Exit(1)
		}
	}()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-sigChan
	// slog.Info("Received signal, initiating graceful shutdown...", "signal", sig)
	slog.Info(fmt.Sprintf("Received signal %v, initiating graceful shutdown...", sig))

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		slog.Info(fmt.Sprintf("Server forced to shutdown: %v", err))
	} else {
		slog.Info("Server gracefully stopped")
	}
}
