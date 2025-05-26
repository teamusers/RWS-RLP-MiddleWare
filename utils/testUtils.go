package utils

import (
	"encoding/json"
	"fmt"
	"lbe/config"
	"os"
	"path/filepath"
)

func LoadTestData[T any](path string) T {
	var data T

	basePaths := [][]string{
		{"..", "..", "..", "..", "..", "testdata"}, // current working test path for controllers in v1
		{"..", "..", "..", "..", "testdata"},       // alternate test path for controllers
		{"..", "..", "..", "testdata"},             // alternate test path for http
	}

	var fileBytes []byte
	var fileFound bool

	for _, base := range basePaths {
		filePath := filepath.Join(append(base, path)...)
		if _, err := os.Stat(filePath); err == nil {
			b, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println("FAILED TO READ FILE:", filePath)
				return data
			}
			fileBytes = b
			fileFound = true
			break
		}
	}

	if !fileFound {
		fmt.Println("FILE DOES NOT EXIST in any known path for:", path)
		return data
	}

	if err := json.Unmarshal(fileBytes, &data); err != nil {
		fmt.Println("FAILED TO UNMARSHAL JSON:", err)
	}
	return data
}

// test utils

func BuildCiamEmailFilter(email string) string {
	return fmt.Sprintf("mail eq '%s'", email)
}

func BuildCiamGrIdFilter(grId string) string {
	return fmt.Sprintf("%s/grid eq '%s'", config.GetConfig().Api.Eeid.UserIdLinkExtensionKey, grId)
}
