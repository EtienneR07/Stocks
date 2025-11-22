package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteJSON[T any](fileName string, data T) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("could not serialize to json: %w", err)
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
