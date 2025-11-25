package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadFile[T any](fileName string) ([]T, error) {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return []T{}, err
	}

	var data []T
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		fmt.Printf("Error deserializing JSON: %s\n", err)
		return []T{}, err
	}

	return data, nil
}
