package server

import (
	"os"
	"strings"
)

var allowedOrigins string

func Initialize() (port, apiBaseURL string) {
	// Default values for development mode
	defaultPort := "3000"
	defaultAPIBaseURL := "https://bill-inquiry-api.onrender.com"
	defaultAllowedOrigins := "*"

	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "development"
	}

	if mode == "production" {
		// In production, only use environment variables (no fallbacks)
		port = os.Getenv("BOHECO2_PROXY_SERVER_PORT")
		apiBaseURL = os.Getenv("BOHECO2_API_BASE_URL")
		allowedOrigins = os.Getenv("BOHECO2_PROXY_SERVER_ALLOWED_ORIGINS")
		
		// Ensure required environment variables are set in production
		if port == "" || apiBaseURL == "" || allowedOrigins == "" {
			panic("Required environment variables BOHECO2_PROXY_SERVER_PORT and BOHECO2_API_BASE_URL must be set in production mode")
		}
	} else {
		// In development mode, use environment variables with fallbacks to defaults
		port = getEnvWithFallback("BOHECO2_PROXY_SERVER_PORT", defaultPort)
		apiBaseURL = getEnvWithFallback("BOHECO2_API_BASE_URL", defaultAPIBaseURL)
		allowedOrigins = getEnvWithFallback("BOHECO2_PROXY_SERVER_ALLOWED_ORIGINS", defaultAllowedOrigins)
	}

	return
}

// getEnvWithFallback returns the environment variable value or fallback if empty/non-existent
func getEnvWithFallback(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func IsOriginAllowed(requestOrigin string) bool {
	if allowedOrigins == "*" {
		return true
	}

	for origin := range strings.SplitSeq(allowedOrigins, ",") {
		if strings.TrimSpace(origin) == requestOrigin {
			return true
		}
	}
	return false
}