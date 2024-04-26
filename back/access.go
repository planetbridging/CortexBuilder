package main

import (
	"os"
)

func ensureDir(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		os.MkdirAll(dirName, os.ModePerm)
	}
}

// Function to list files in a directory and format the response
func listFilesInDir(dirPath string) ([]map[string]interface{}, error) {
	var files []map[string]interface{}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		_, err := entry.Info()
		if err != nil {
			return nil, err
		}
		file := map[string]interface{}{
			"id":    entry.Name(),
			"name":  entry.Name(),
			"isDir": entry.IsDir(),
		}
		if entry.IsDir() {
			// Recursively list files if necessary or simply mark as directory
			// Not expanding children here unless necessary
		}
		files = append(files, file)
	}

	return files, nil
}
