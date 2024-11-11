package main

import (
	"log"
	"net/http"

	"github.com/FrAigner/spacestore/handlers"
	"github.com/FrAigner/spacestore/middleware"
	"github.com/FrAigner/spacestore/utils"

	"github.com/gorilla/mux"
)

func main() {
	utils.LoadEnv() // ENV laden

	r := mux.NewRouter()

	// Datei-Upload Endpoint mit API-Key Authentifizierung
	r.Handle("/upload", middleware.APIKeyAuth(http.HandlerFunc(handlers.FileUploadHandler))).Methods("POST")

	// Server starten
	log.Println("Server läuft auf Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
