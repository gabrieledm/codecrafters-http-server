package main

import (
	"fmt"
	"io"
	"os"
)

func readFile(baseDir string, fileName string) []byte {
	tmp := "/tmp/"
	if len(baseDir) != 0 {
		tmp = baseDir
	}

	contents, err := os.ReadFile(tmp + fileName)

	if err != nil {
		fmt.Println("File reading error", err)
		return nil
	}

	return contents
}

func writeFile(length int, reader io.Reader, fileName string) {
	fileContent := make([]byte, length)

	_, err := io.ReadFull(reader, fileContent)
	if err != nil {
		fmt.Println("Error reading file content:", err)
	}

	completePath := tmpDirectory + "/" + fileName
	err = os.WriteFile(completePath, fileContent, os.FileMode(0660))
	if err != nil {
		fmt.Printf("Error writing file to '%s': %s", completePath, err)
	}
	fmt.Printf("Succesfully written file at '%s'\n", completePath)
}
