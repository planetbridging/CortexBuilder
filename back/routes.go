package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Post("/mountpopulation", mountPopulationHandler)
	app.Post("/mounttrainingdata", mountTrainingDataHandler)
	app.Post("/evaluation", runEvaluation)

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

func mountPopulationHandler(c *fiber.Ctx) error {
	// Define a struct to map your JSON data
	type RequestData struct {
		DbName         string `json:"dbName"`
		CollectionName string `json:"collectionName"`
	}

	// Instance of the struct to hold your POST data
	data := new(RequestData)

	// Parsing the JSON body to the struct
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Logging the data to the console
	//fmt.Printf("Received dbName: %s\n", data.DbName)
	//fmt.Printf("Received collectionName: %s\n", data.CollectionName)

	// Check if data is received properly
	if data.DbName == "" || data.CollectionName == "" {
		fmt.Println("Did not receive expected variables.")
	} else {
		//fmt.Println("success")
		dbNameSync = data.DbName
		collectionSync = data.CollectionName
		updateFrontend()
	}

	// Return a response to the client
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data received successfully",
	})
}

func mountTrainingDataHandler(c *fiber.Ctx) error {
	// Define a struct to map your JSON data
	type RequestData struct {
		TableName string `json:"tableName"` // Notice TableName is capitalized
	}

	// Instance of the struct to hold your POST data
	data := new(RequestData)

	// Parsing the JSON body to the struct
	if err := c.BodyParser(data); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Logging the data to the console
	fmt.Printf("Received tableName: %s\n", data.TableName)

	// Check if data is received properly
	if data.TableName == "" {
		fmt.Println("Did not receive expected variables.")
	} else {
		fmt.Println("Received tableName:", data.TableName)
		tableNameSync = data.TableName // Assuming tableNameSync is a global variable
		updateFrontend()
	}

	// Return a response to the client
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data received successfully",
	})
}

func runEvaluation(c *fiber.Ctx) error {
	// Define a struct to map your JSON data
	type RequestData struct {
		TableName      string `json:"tableName"` // Notice TableName is capitalized
		DbName         string `json:"dbName"`
		CollectionName string `json:"collectionName"`
		Batch          string `json:"batchNumber"` // Corrected capitalization
	}

	// Instance of the struct to hold your POST data
	data := new(RequestData)

	// Parsing the JSON body to the struct
	if err := c.BodyParser(data); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Logging the data to the console
	fmt.Printf("Received tableName: %s\n", data.TableName)
	fmt.Printf("Received dbName: %s\n", data.DbName)
	fmt.Printf("Received collectionName: %s\n", data.CollectionName)
	fmt.Printf("Received batch: %s\n", data.Batch)

	startEvaluation(data.DbName, data.CollectionName, data.Batch)

	// Return a response to the client
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data received successfully",
	})
}

func getModelCount(dbName, collectionName string) (int, error) {
	url := fmt.Sprintf("http://localhost:1789/countListModels?dbName=%s&collectionName=%s", dbName, collectionName)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var count int
	_, err = fmt.Sscanf(string(bodyBytes), `{"count":%d}`, &count)
	if err != nil {
		return 0, fmt.Errorf("error parsing response: %v", err)
	}

	return count, nil
}
