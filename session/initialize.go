// Package session handles session management for the BOHECO2 bill inquiry proxy.
// This file contains session initialization and HTTP client setup.
package session

import (
	"io"
	"net/http"
	"net/http/cookiejar"
)

// Global variables for session management
var (
	// Jar stores cookies for maintaining session state with the upstream API
	Jar    *cookiejar.Jar
	// Client is the HTTP client configured with cookie jar for session persistence
	Client *http.Client
	// APIBaseURL stores the base URL for the API (configurable via environment variables)
	APIBaseURL string
)


// init initializes the HTTP client with cookie jar support.
// This runs automatically when the package is imported.
func init() {
	// Create a new cookie jar to store session cookies
	Jar, _ = cookiejar.New(nil)
	// Create HTTP client that automatically handles cookies
	Client = &http.Client{Jar: Jar}
}

// SetAPIBaseURL configures the base URL for API requests.
// This should be called before using Initialize() or HasToken().
func SetAPIBaseURL(baseURL string) {
	APIBaseURL = baseURL
}

// Initialize establishes a session with the upstream API by calling the session-init endpoint.
// This function should be called when no valid session token exists.
// Returns an error if the session initialization fails.
func Initialize() error {
	// Create request to the session initialization endpoint
	apiURL := APIBaseURL + "/api/v1/session-init"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}

	// Set headers to mimic requests from the official BOHECO2 website
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Origin", "https://www.boheco2.com.ph")
	req.Header.Set("Referer", "https://www.boheco2.com.ph")

	// Execute the session initialization request
	resp, err := Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Discard response body (we only care about the session cookies)
	io.Copy(io.Discard, resp.Body)
	return nil
}