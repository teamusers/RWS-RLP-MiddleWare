package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func LoadTestData[T any](path string) T {
	var data T

	filePath := filepath.Join("..", "..", "..", "..", "..", "testdata", path)

	_, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("FILE DOES NOT EXIST:", filePath)
		return data
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return data
	}
	err = json.Unmarshal(fileBytes, &data)
	return data
}
