package utils

import (
	"encoding/json"
	"fmt"
	"gobot/models"
	"io"
	"os"
)

func ReadFile(path string) []models.NewsType {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return nil
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file: ", err)
		return nil
	}

	var news []models.NewsType
	err = json.Unmarshal(byteValue, &news)
	if err != nil {
		fmt.Println("Error unmarshalling file: ", err)
		return nil
	}
	return news
}

func WriteFile(path string, data []models.NewsType) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}
	defer file.Close()

	byteValue, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data: ", err)
		return
	}

	_, err = file.Write(byteValue)
	if err != nil {
		fmt.Println("Error writing data: ", err)
		return
	}
}
