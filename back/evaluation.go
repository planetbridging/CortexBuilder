package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	// other imports...
)

type EvaluationMetrics struct {
	TrainingError   float64
	ValidationError float64
	Accuracy        float64
	Precision       float64
	Recall          float64
	F1Score         float64
	AUCROC          float64
	MAE             float64
	MSE             float64
	RMSE            float64
}

func startEvaluation(dbName string, collectionName string, batchSize string) {
	//dbName := "db_df6d2eb2-1890-4674-ab4e-a4009947574c"
	//collectionName := "col_c8f4f784-85b5-4c4e-b1e5-614f1bdda83d"
	//batchSize := 100

	batchSizeInt, err := strconv.Atoi(batchSize)
	if err != nil {
		fmt.Println("Error converting batchSize to integer:", err)
		return
	}

	models, err := getListModels(dbName, collectionName, batchSizeInt)
	if err != nil {
		log.Fatal(err)
	}

	// Process models concurrently
	var wg sync.WaitGroup
	results := make(chan map[string]interface{}, len(models)) // Channel to collect results
	for _, model := range models {
		wg.Add(1)
		go func(m map[string]interface{}) {
			defer wg.Done()
			result := processModel(m, dbName, collectionName)
			results <- result // Send the result to the channel
		}(model)
	}
	wg.Wait()
	close(results) // Close the channel after all goroutines finish

	fmt.Println(results)

	// Collect results from the channel
	var processedModels []map[string]interface{}
	for result := range results {
		processedModels = append(processedModels, result)
	}

	// Print the processed models
	for _, m := range processedModels {
		fmt.Printf("Model ID: %s, outputs: %s\n", m["modelID"], m["outputs"])
	}
}

func getListModels(dbName, collectionName string, batchSize int) ([]map[string]interface{}, error) {
	requestURL := fmt.Sprintf("http://localhost:1789/listModels?dbName=%s&collectionName=%s", dbName, collectionName)
	fmt.Println(requestURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response (assuming it's JSON)
	var models []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&models)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Split models into batches
	var batches [][]map[string]interface{}
	for i := 0; i < len(models); i += batchSize {
		end := i + batchSize
		if end > len(models) {
			end = len(models)
		}
		batches = append(batches, models[i:end])
	}

	// Flatten the batches into a single slice
	var flattenedModels []map[string]interface{}
	for _, batch := range batches {
		flattenedModels = append(flattenedModels, batch...)
	}

	return flattenedModels, nil
}

func processModel(model map[string]interface{}, dbName string, collectionName string) map[string]interface{} {
	// Process the model here
	// For example, insert into your database or perform other operations
	// based on the model data.
	// You can also add error handling for each model if needed.
	// Modify this part according to your actual processing logic.
	modelID := model["_id"].(string)

	// Make an additional HTTP request to getModel
	getModelURL := fmt.Sprintf("http://localhost:1789/getModel?dbName=%s&collectionName=%s&modelId=%s", dbName, collectionName, modelID)
	resp, err := http.Get(getModelURL)
	if err != nil {
		// Handle the error (e.g., log it or return an error map)
		fmt.Println("Error fetching model details:", err)
		//return nil
	}
	defer resp.Body.Close()

	// Parse the response (assuming it's JSON)
	/*var modelDetails map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&modelDetails)
	if err != nil {
		// Handle the error (e.g., log it or return an error map)
		fmt.Println("Error decoding model details:", err)
		//return nil
	}*/

	// Parse the response (assuming it's JSON)
	var modelDetails NetworkConfig
	err = json.NewDecoder(resp.Body).Decode(&modelDetails)
	if err != nil {
		fmt.Println("Error decoding model details:", err)
		//return nil
	}

	// Example input values - adjust based on your actual input configuration
	inputValues := map[string]float64{
		"1": 1,
		"2": 0.5,
		"3": 0.75,
	}

	outputs := feedforward(&modelDetails, inputValues)

	// Return the processed model
	return map[string]interface{}{
		"modelID": modelID,
		"outputs": outputs,
	}
}
