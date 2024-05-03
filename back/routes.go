package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/files/*", func(c *fiber.Ctx) error {
		// Extracting the subpath or handling root
		subPath := c.Params("*")
		if subPath == "" {
			subPath = "./host" // default path
		}

		files, err := listFilesInDir(subPath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Assuming the folder chain is a breadcrumb of paths
		folderChain := []map[string]string{{"id": "root", "name": "Home"}}
		// Add more logic here to build the folder chain based on the path

		response := map[string]interface{}{
			"files":       files,
			"folderChain": folderChain,
		}

		return c.JSON(response)
	})

	app.Post("/initialize", func(c *fiber.Ctx) error {
		var request struct {
			NetworkType      string `json:"networkType"`
			SpawnCount       int    `json:"spawnCount"`
			AdditionalParam1 string `json:"additionalParam1"`
			AdditionalParam2 string `json:"additionalParam2"`
			AdditionalParam3 string `json:"additionalParam3"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request"})
		}

		fmt.Printf("Received initialization request: %+v\n", request)

		// Generate a random model
		//randomModel := randomizeNetworkStaticTesting()
		//fmt.Println("Generated Model:", randomModel)

		for i := 0; i < request.SpawnCount; i++ {
			model := randomizeNetworkStaticTesting()

			// Prepare payload to send to /saveModel
			payload := map[string]interface{}{
				"dbName":         request.AdditionalParam1,
				"collectionName": request.AdditionalParam2,
				"model":          json.RawMessage(model),
			}

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				fmt.Println("Error marshaling model:", err)
				continue
			}

			resp, err := http.Post("http://localhost:1789/saveModel", "application/json", bytes.NewBuffer(payloadBytes))
			if err != nil {
				fmt.Println("Failed to send model:", err)
				continue // Optionally handle this error more gracefully
			}

			fmt.Printf("Model sent, response status: %d\n", resp.StatusCode)
			resp.Body.Close() // Don't forget to close the response body
		}

		// Here, you would add your logic to handle the training initialization
		// For example, setting up configurations, preparing datasets, etc.

		return c.Status(fiber.StatusOK).SendString("Training initialization started successfully")
	})

	app.Get("/data", func(c *fiber.Ctx) error {
		dbName := c.Query("dbname")

		if dbName == "" {
			// List databases
			databaseNames, err := ListDatabases(DBJB)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error fetching databases: " + err.Error())
			}

			// Filter out system databases
			filteredDatabases := []string{}
			for _, dbName := range databaseNames {
				if dbName != "sys" && dbName != "performance_schema" && dbName != "mysql" && dbName != "information_schema" {
					filteredDatabases = append(filteredDatabases, dbName)
				}
			}

			// Format as Chonky files
			files := formatAsChonkyFiles(filteredDatabases, true) // Assuming isDir=true for databases

			folderChain := []map[string]string{{"id": "root", "name": "Home"}}
			response := map[string]interface{}{
				"files":       files,
				"folderChain": folderChain,
			}
			return c.JSON(response)

		} else {
			// Ensure the database exists (modify if needed)
			err := EnsureDatabaseExists(DBJB, dbName)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).SendString("Database not found: " + dbName)
			}

			// Get table names from the provided database
			tableNames, err := GetTableNames(DBJB, dbName)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error fetching tables: " + err.Error())
			}

			// Format as Chonky files
			files := formatAsChonkyFiles(tableNames, true) // Assuming isDir=true for tables

			// Adjust folderChain if needed (e.g., to include the dbName)
			folderChain := []map[string]string{
				{"id": "root", "name": "Home"},
				{"id": dbName, "name": dbName},
			}

			response := map[string]interface{}{
				"files":       files,
				"folderChain": folderChain,
			}
			return c.JSON(response)
		}
	})
}

func formatAsChonkyFiles(names []string, isDir bool) []map[string]string {
	files := make([]map[string]string, len(names))
	for i, name := range names {
		files[i] = map[string]string{
			"id":    name,
			"name":  name,
			"isDir": strconv.FormatBool(isDir), // Convert boolean to string
		}
	}
	return files
}
