// Package session handles session management for the BOHECO2 bill inquiry proxy.
// This file contains token validation functionality.
package session

import (
	"net/url"
	"time"
)

// HasToken checks if we have a valid, non-expired session token.
// It searches through stored cookies for a "session_token" cookie
// that hasn't expired yet.
// Returns true if a valid token exists, false otherwise.
func HasToken() bool {
	// Parse the API URL to check cookies for this domain
	u, _ := url.Parse(APIBaseURL)
	
	// Iterate through all cookies for the API domain
	for _, ck := range Jar.Cookies(u) {
		// Check if this is our session token and it's still valid
		if ck.Name == "session_token" && ck.Expires.After(time.Now()) {
			return true
		}
	}
	
	// No valid session token found
	return false
}