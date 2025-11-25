package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func AppendJSON[T any](fileName string, data T) error {
	var existingData []T
	fileContent, err := os.ReadFile(fileName)

	if err == nil {
		if err := json.Unmarshal(fileContent, &existingData); err != nil {
			return fmt.Errorf("could not parse existing JSON: %w", err)
		}
	}

	existingData = append(existingData, data)

	jsonData, err := json.MarshalIndent(existingData, "", "  ")
	if err != nil {
		return fmt.Errorf("could not serialize to json: %w", err)
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
