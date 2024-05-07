package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Connection struct {
	Weight float64 `json:"weight"`
}

type Neuron struct {
	ActivationType string                `json:"activationType"`
	Connections    map[string]Connection `json:"connections"`
	Bias           float64               `json:"bias"`
}

type Layer struct {
	Neurons map[string]Neuron `json:"neurons"`
}

type NetworkConfig struct {
	Layers struct {
		Input  Layer   `json:"input"`
		Hidden []Layer `json:"hidden"`
		Output Layer   `json:"output"`
	} `json:"layers"`
}

func activate(activationType string, input float64) float64 {
	switch activationType {
	case "relu":
		return math.Max(0, input)
	case "sigmoid":
		return 1 / (1 + math.Exp(-input))
	case "tanh":
		return math.Tanh(input)
	case "softmax":
		return math.Exp(input) // Should normalize later in the layer processing
	case "leaky_relu":
		if input > 0 {
			return input
		}
		return 0.01 * input
	case "swish":
		return input * (1 / (1 + math.Exp(-input))) // Beta set to 1 for simplicity
	case "elu":
		alpha := 1.0 // Alpha can be adjusted based on specific needs
		if input >= 0 {
			return input
		}
		return alpha * (math.Exp(input) - 1)
	case "selu":
		lambda := 1.0507    // Scale factor
		alphaSELU := 1.6733 // Alpha for SELU
		if input >= 0 {
			return lambda * input
		}
		return lambda * (alphaSELU * (math.Exp(input) - 1))
	case "softplus":
		return math.Log(1 + math.Exp(input))
	default:
		return input // Linear activation (no change)
	}
}

func feedforward(config *NetworkConfig, inputValues map[string]float64) map[string]float64 {
	neurons := make(map[string]float64)

	// Initialize input layer neurons with input values
	for inputID := range config.Layers.Input.Neurons {
		neurons[inputID] = inputValues[inputID]
	}

	// Process hidden layers
	for _, layer := range config.Layers.Hidden {
		for nodeID, node := range layer.Neurons {
			sum := 0.0
			for inputID, connection := range node.Connections {
				sum += neurons[inputID] * connection.Weight
			}
			sum += node.Bias
			neurons[nodeID] = activate(node.ActivationType, sum)
		}
	}

	// Process output layer
	outputs := make(map[string]float64)
	for nodeID, node := range config.Layers.Output.Neurons {
		sum := 0.0
		for inputID, connection := range node.Connections {
			sum += neurons[inputID] * connection.Weight
		}
		sum += node.Bias
		outputs[nodeID] = activate(node.ActivationType, sum)
	}

	return outputs
}

func randomizeModelOnlyLayer() string {
	rand.Seed(time.Now().UnixNano())
	activationTypes := []string{"relu", "sigmoid", "tanh", "softmax", "leaky_relu", "swish", "elu", "selu", "softplus"}
	activationType := activationTypes[rand.Intn(len(activationTypes))]

	// Randomize weights and bias for a single neuron
	weight1 := rand.NormFloat64() // Random weight from a normal distribution
	bias := rand.NormFloat64()

	// Constructing a JSON model with the randomized parameters
	model := map[string]interface{}{
		"layers": map[string]interface{}{
			"hidden": []map[string]interface{}{
				{
					"neurons": map[string]interface{}{
						"4": map[string]interface{}{
							"activationType": activationType,
							"connections": map[string]interface{}{
								"1": map[string]interface{}{
									"weight": weight1,
								},
							},
							"bias": bias,
						},
					},
				},
			},
		},
	}

	modelJSON, _ := json.Marshal(model)
	return string(modelJSON)
}

func randomWeight() float64 {
	return rand.NormFloat64() // Generate a Gaussian distribution random weight
}

func randomizeNetworkStaticTesting() string {
	model := map[string]interface{}{
		"layers": map[string]interface{}{
			"input": map[string]interface{}{
				"neurons": map[string]interface{}{
					"1": map[string]interface{}{},
					"2": map[string]interface{}{},
					"3": map[string]interface{}{},
				},
			},
			"hidden": []map[string]interface{}{
				{
					"neurons": map[string]interface{}{
						"4": map[string]interface{}{
							"activationType": "relu",
							"connections": map[string]interface{}{
								"1": map[string]interface{}{
									"weight": randomWeight(),
								},
							},
							"bias": rand.Float64(), // Random bias between 0 and 1
						},
					},
				},
			},
			"output": map[string]interface{}{
				"neurons": map[string]interface{}{
					"5": map[string]interface{}{
						"activationType": "sigmoid",
						"connections": map[string]interface{}{
							"4": map[string]interface{}{
								"weight": randomWeight(),
							},
						},
						"bias": rand.Float64(),
					},
					"6": map[string]interface{}{
						"activationType": "sigmoid",
						"connections": map[string]interface{}{
							"4": map[string]interface{}{
								"weight": randomWeight(),
							},
						},
						"bias": rand.Float64(),
					},
					"7": map[string]interface{}{
						"activationType": "sigmoid",
						"connections": map[string]interface{}{
							"4": map[string]interface{}{
								"weight": randomWeight(),
							},
						},
						"bias": rand.Float64(),
					},
				},
			},
		},
	}

	modelJSON, _ := json.MarshalIndent(model, "", "  ")
	return string(modelJSON)
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
