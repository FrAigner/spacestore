package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type APIKeys struct {
	Keys map[string]string `json:"api_keys"`
}

// LoadAPIKeys lädt API-Schlüssel und zugehörige Ordner aus einer JSON-Datei
func LoadAPIKeys(filePath string) (map[string]string, error) {
	var apiKeys APIKeys

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &apiKeys)
	if err != nil {
		return nil, err
	}

	return apiKeys.Keys, nil
}
