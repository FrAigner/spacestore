package main

import (
	"log"
	"net/http"

	"github.com/FrAigner/spacestore/handlers"   // Importiere den handlers-Package
	"github.com/FrAigner/spacestore/middleware" // Importiere das middleware-Package
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// APIKeyAuth als Middleware registrieren
	r.Use(middleware.APIKeyAuth)

	// Endpunkt zum Hochladen von Dateien
	r.HandleFunc("/upload", handlers.UploadFile).Methods("POST")

	log.Println("Server läuft auf Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
