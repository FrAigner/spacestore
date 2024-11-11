package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type APIKeysContainer struct {
	Keys map[string]string `json:"keys"`
}

// LoadAPIKeys lädt die API-Schlüssel aus einer JSON-Datei
func LoadAPIKeys(filePath string) (map[string]string, error) {
	// Öffne die JSON-Datei
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening API keys file: %v", err)
		return nil, err
	}
	defer file.Close()

	// Lese den Inhalt der Datei
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Error reading API keys file: %v", err)
		return nil, err
	}

	// Parsen der JSON-Daten in die APIKeysContainer-Struktur
	var container APIKeysContainer
	err = json.Unmarshal(data, &container)
	if err != nil {
		log.Printf("Error unmarshalling API keys: %v", err)
		return nil, err
	}

	// Rückgabe der geladenen API-Schlüssel
	return container.Keys, nil
}
