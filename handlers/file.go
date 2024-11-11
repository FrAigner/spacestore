package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	// Logge den Beginn der Anfrage
	log.Println("Received file upload request")

	// Parsing der Multipart-Formulardaten
	err := r.ParseMultipartForm(10 << 20) // Maximal 10MB
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "File Error", http.StatusBadRequest)
		return
	}

	// Hole die Datei aus der Anfrage
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file from request: %v", err)
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Logge den Dateinamen, der hochgeladen wird
	log.Printf("Received file: %s", handler.Filename)

	// Bestimme den Zielordner und Dateipfad
	uploadDir := "uploads"
	filePath := filepath.Join(uploadDir, handler.Filename)

	// Überprüfen, ob der Ordner existiert, und ggf. erstellen
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		log.Printf("Error creating upload directory: %v", err)
		http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
		return
	}
	log.Printf("Upload directory '%s' created or already exists", uploadDir)

	// Erstelle die Datei im Zielordner
	f, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file at path '%s': %v", filePath, err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Kopiere den Inhalt der hochgeladenen Datei in die Zieldatei
	_, err = io.Copy(f, file)
	if err != nil {
		log.Printf("Error copying file contents to '%s': %v", filePath, err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Erfolgreiche Antwort
	log.Printf("File successfully uploaded to '%s'", filePath)
	w.Write([]byte("File uploaded successfully"))
}
