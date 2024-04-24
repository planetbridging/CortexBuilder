package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func NewNetworkFromFile(filename string) (NeuralNetwork, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return NeuralNetwork{}, err
	}

	return NewNetworkFromJSON(data)
}

func ensureDir(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		os.MkdirAll(dirName, os.ModePerm)
	}
}

// listFilesInDir lists the files and folders in a given directory
func OldlistFilesInDir(dirPath string) []map[string]interface{} {
	var fileList []map[string]interface{}
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip the root directory itself
		if path == dirPath {
			return nil
		}
		fileList = append(fileList, map[string]interface{}{
			"id":    info.Name(),
			"name":  info.Name(),
			"isDir": info.IsDir(),
			"path":  path,
		})
		return nil
	})
	return fileList
}

// Modified to include childrenIds
func Old2listFilesInDir(dirPath string) []map[string]interface{} {
	fileMap := make(map[string]map[string]interface{})
	fileList := make([]map[string]interface{}, 0)

	// Walk the directory and build the file map
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip the root directory itself
		if path == dirPath {
			return nil
		}

		// Create a unique ID for each file/folder as the map key
		id := filepath.Base(path)
		fileEntry := map[string]interface{}{
			"id":    id,
			"name":  info.Name(),
			"isDir": info.IsDir(),
			"path":  path,
		}
		if info.IsDir() {
			fileEntry["childrenIds"] = []string{}
		}

		fileMap[id] = fileEntry
		return nil
	})

	// Build the parent-child relationships
	for _, file := range fileMap {
		if file["isDir"].(bool) {
			dirPath := file["path"].(string)
			filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
				if path != dirPath { // Avoid adding the directory itself
					childId := filepath.Base(path)
					children := file["childrenIds"].([]string)
					children = append(children, childId)
					file["childrenIds"] = children
				}
				return nil
			})
		}
		fileList = append(fileList, file)
	}

	return fileList
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
