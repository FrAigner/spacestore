package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/FrAigner/spacestore/utils"
	"github.com/gorilla/mux"
)

var apiKeys map[string]string

func main() {
	var err error
	// Lade die API-Schlüssel und Ordner aus der JSON-Datei
	apiKeys, err = utils.LoadAPIKeys("api_keys.json")
	if err != nil {
		log.Fatalf("Fehler beim Laden der API-Schlüssel: %v", err)
	}

	r := mux.NewRouter()

	// Middleware für API-Key-Überprüfung
	r.Use(APIKeyAuth)

	// Endpunkt zum Hochladen von Dateien
	r.HandleFunc("/upload", uploadFile).Methods("POST")

	log.Println("Server läuft auf Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Middleware für API-Key-Überprüfung und Zuordnung zum richtigen Ordner
func APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		log.Printf("Received API-Key: %s", apiKey)

		uploadDir, validKey := apiKeys[apiKey]
		if !validKey {
			log.Printf("Invalid API Key: %s", apiKey)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Setze den Upload-Ordner für die Anfrage
		r.Header.Set("Upload-Dir", uploadDir)
		log.Printf("API Key '%s' is valid. Upload folder: '%s'", apiKey, uploadDir)

		next.ServeHTTP(w, r)
	})
}

// uploadFile behandelt Datei-Uploads und speichert sie im entsprechenden Ordner
func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Logge den Beginn des Datei-Uploads
	log.Println("Received file upload request")

	// Parsing der Multipart-Formulardaten
	err := r.ParseMultipartForm(10 << 20) // Max. Größe: 10 MB
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "Fehler beim Parsen der Datei", http.StatusBadRequest)
		return
	}

	// Hole die Datei aus der Anfrage
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file from request: %v", err)
		http.Error(w, "Fehler beim Abrufen der Datei", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ordner aus den Header-Informationen holen
	uploadDir := r.Header.Get("Upload-Dir")
	log.Printf("Uploading file '%s' to folder '%s'", handler.Filename, uploadDir)

	// Prefix "uploads" zum Ordnernamen hinzufügen
	uploadDir = "./uploads/" + uploadDir

	// Pfad für die Datei festlegen
	filePath := uploadDir + "/" + handler.Filename

	// Überprüfen, ob der Ordner existiert, und ggf. erstellen
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		log.Printf("Error creating upload directory '%s': %v", uploadDir, err)
		http.Error(w, "Fehler beim Erstellen des Upload-Ordners", http.StatusInternalServerError)
		return
	}
	log.Printf("Upload directory '%s' created or already exists", uploadDir)

	// Datei im Upload-Ordner speichern
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file at path '%s': %v", filePath, err)
		http.Error(w, "Fehler beim Speichern der Datei", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Kopiere den Inhalt der hochgeladenen Datei
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Error writing file contents to '%s': %v", filePath, err)
		http.Error(w, "Fehler beim Schreiben der Datei", http.StatusInternalServerError)
		return
	}

	// Erfolgreiche Antwort
	log.Printf("File '%s' successfully uploaded to '%s'", handler.Filename, filePath)
	w.Write([]byte("Datei erfolgreich hochgeladen"))
}
