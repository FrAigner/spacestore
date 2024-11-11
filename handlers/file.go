package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // Maximal 10MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File Error", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath := filepath.Join("uploads", handler.Filename)

	f, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	io.Copy(f, file)

	w.Write([]byte("File uploaded successfully"))
}
