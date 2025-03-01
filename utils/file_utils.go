package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadFile[T any](path string) ([]T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var result []T
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling file: %w", err)
	}

	return result, nil
}

func WriteFile[T any](path string, data []T) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	byteValue, err := json.MarshalIndent(data, "", "  ") // Indented for readability
	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}

	_, err = file.Write(byteValue)
	if err != nil {
		return fmt.Errorf("error writing data: %w", err)
	}
	return nil
}

