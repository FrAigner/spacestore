package handlers

import (
	"archive/tar"
	"archive/zip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Entpackt eine ZIP-Datei in den angegebenen Zielordner
func Unzip(src, dest string) error {
	// Füge den 'uploads' Prefix hinzu
	finalDest := filepath.Join("uploads", dest) // 'uploads' wird hier als Prefix vorangestellt

	// Öffne die ZIP-Datei
	zipFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Erstelle ein ZIP-Reader
	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	// Erstelle einen Ordner mit dem gleichen Namen wie die ZIP-Datei im Zielordner
	dirName := filepath.Join(finalDest, filepath.Base(src[:len(src)-4])) // Entferne die .zip-Endung
	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return err
	}

	// Entpacke jede Datei im ZIP-Archiv
	for _, file := range zipReader.File {
		err := unzipFile(file, dirName)
		if err != nil {
			return err
		}
	}

	// Lösche das ZIP-Archiv nach dem Entpacken
	err = os.Remove(src)
	if err != nil {
		return err
	}

	return nil
}

// Entpackt eine Datei aus einem ZIP-Archiv in den Zielordner
func unzipFile(file *zip.File, dest string) error {
	// Öffne die Datei im ZIP-Archiv
	zipFile, err := file.Open()
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Bestimme den Zielpfad der entpackten Datei
	destPath := filepath.Join(dest, file.Name)
	if file.FileInfo().IsDir() {
		err := os.MkdirAll(destPath, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		// Erstelle die entpackte Datei
		outFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Kopiere den Inhalt der Datei
		_, err = io.Copy(outFile, zipFile)
		if err != nil {
			return err
		}
	}
	return nil
}

// Entpackt eine TAR-Datei in den angegebenen Zielordner
func Untar(src, dest string) error {
	// Füge den 'uploads' Prefix hinzu
	finalDest := filepath.Join("uploads", dest) // 'uploads' wird hier als Prefix vorangestellt

	// Öffne die TAR-Datei
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	// Erstelle einen TAR-Reader
	tarReader := tar.NewReader(file)

	// Bestimme den Ordnernamen aus dem TAR-Archiv
	dirName := filepath.Join(finalDest, filepath.Base(src[:len(src)-4])) // Entferne die .tar-Endung
	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return err
	}

	// Entpacke jede Datei im TAR-Archiv
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Bestimme den Zielpfad der entpackten Datei
		destPath := filepath.Join(dirName, header.Name)
		if header.Typeflag == tar.TypeDir {
			err := os.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			outFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			// Kopiere den Inhalt der Datei
			_, err = io.Copy(outFile, tarReader)
			if err != nil {
				return err
			}
		}
	}

	// Lösche die TAR-Datei nach dem Entpacken
	err = os.Remove(src)
	if err != nil {
		return err
	}

	return nil
}

// uploadFile behandelt Datei-Uploads und speichert sie im entsprechenden Ordner
func UploadFile(w http.ResponseWriter, r *http.Request) {
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

	// Hier den Prefix 'uploads' hinzufügen
	finalUploadDir := filepath.Join("uploads", uploadDir)

	// Pfad für die Datei festlegen
	filePath := filepath.Join(finalUploadDir, handler.Filename)

	// Überprüfen, ob der Ordner existiert, und ggf. erstellen
	err = os.MkdirAll(finalUploadDir, os.ModePerm)
	if err != nil {
		log.Printf("Error creating upload directory '%s': %v", finalUploadDir, err)
		http.Error(w, "Fehler beim Erstellen des Upload-Ordners", http.StatusInternalServerError)
		return
	}
	log.Printf("Upload directory '%s' created or already exists", finalUploadDir)

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

	// Überprüfe, ob es sich um eine ZIP- oder TAR-Datei handelt und entpacke sie
	if strings.HasSuffix(handler.Filename, ".zip") {
		log.Printf("File '%s' is a ZIP file, unpacking...", handler.Filename)
		err = Unzip(filePath, uploadDir)
		if err != nil {
			log.Printf("Error unpacking ZIP file: %v", err)
			http.Error(w, "Fehler beim Entpacken der ZIP-Datei", http.StatusInternalServerError)
			return
		}
	} else if strings.HasSuffix(handler.Filename, ".tar") {
		log.Printf("File '%s' is a TAR file, unpacking...", handler.Filename)
		err = Untar(filePath, uploadDir)
		if err != nil {
			log.Printf("Error unpacking TAR file: %v", err)
			http.Error(w, "Fehler beim Entpacken der TAR-Datei", http.StatusInternalServerError)
			return
		}
	}

	// Erfolgreiche Antwort
	log.Printf("File '%s' successfully uploaded and unpacked", handler.Filename)
	w.Write([]byte("Datei erfolgreich hochgeladen und entpackt"))
}
