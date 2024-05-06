package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// SystemInfo holds data about system resources
type SystemInfo struct {
	CPUs        []cpu.InfoStat `json:"cpus"`
	TotalMemory uint64         `json:"total_memory"`
	FreeMemory  uint64         `json:"free_memory"`
}

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

func GetSystemInfo() (string, error) {
	// Getting CPU information
	cpus, err := cpu.Info()
	if err != nil {
		return "", fmt.Errorf("error getting CPU info: %v", err)
	}

	// Getting virtual memory stats
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", fmt.Errorf("error getting virtual memory stats: %v", err)
	}

	sysInfo := SystemInfo{
		CPUs:        cpus,
		TotalMemory: vmStat.Total,
		FreeMemory:  vmStat.Available,
	}

	// Convert system info to JSON
	jsonBytes, err := json.Marshal(sysInfo)
	if err != nil {
		return "", fmt.Errorf("error marshalling system info to JSON: %v", err)
	}

	return string(jsonBytes), nil
}
