package middleware

import (
	"fmt"
	"net/http"

	"github.com/FrAigner/spacestore/utils"
)

// APIKeyAuth überprüft den API-Key und weist den entsprechenden Upload-Ordner zu
func APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Hole den API-Key aus den Headern
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, "API Key required", http.StatusUnauthorized)
			return
		}

		// Lade die API-Schlüssel aus der JSON-Datei
		apiKeys, err := utils.LoadAPIKeys("api_keys.json")
		if err != nil {
			http.Error(w, fmt.Sprintf("Fehler beim Laden der API-Schlüssel: %v", err), http.StatusInternalServerError)
			return
		}

		// Überprüfe, ob der API-Key existiert
		uploadDir, validKey := apiKeys[apiKey]
		if !validKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Setze den Upload-Ordner in den Header, damit er im Upload-Handler verwendet werden kann
		r.Header.Set("Upload-Dir", uploadDir)

		// API-Key ist gültig, fahre fort
		next.ServeHTTP(w, r)
	})
}
