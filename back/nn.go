package main

import (
	"encoding/json"
	"math"
	"strconv"
)

type Neuron struct {
	ID             int
	Incoming       []Connection
	Bias           float64
	ActivationType string
	Output         float64
	IsInput        bool
	IsOutput       bool
}

type Connection struct {
	Weight float64
	From   *Neuron
}

type NeuralNetwork struct {
	Neurons       map[int]*Neuron
	InputNeurons  []*Neuron
	OutputNeurons []*Neuron
}

func NewNetworkFromJSON(jsonData []byte) (NeuralNetwork, error) {
	var config struct {
		Neurons map[string]struct {
			Bias           float64            `json:"bias"`
			ActivationType string             `json:"activationType"`
			Connections    map[string]float64 `json:"connections"`
			IsInput        bool               `json:"isInput"`
			IsOutput       bool               `json:"isOutput"`
		} `json:"neurons"`
	}

	err := json.Unmarshal(jsonData, &config)
	if err != nil {
		return NeuralNetwork{}, err
	}

	network := NeuralNetwork{Neurons: make(map[int]*Neuron)}
	for id, neuronConfig := range config.Neurons {
		neuronID, _ := strconv.Atoi(id)
		neuron := &Neuron{
			ID:             neuronID,
			Bias:           neuronConfig.Bias,
			ActivationType: neuronConfig.ActivationType,
			Incoming:       []Connection{},
			IsInput:        neuronConfig.IsInput,
			IsOutput:       neuronConfig.IsOutput,
		}
		network.Neurons[neuronID] = neuron
		if neuron.IsInput {
			network.InputNeurons = append(network.InputNeurons, neuron)
		}
		if neuron.IsOutput {
			network.OutputNeurons = append(network.OutputNeurons, neuron)
		}
	}

	for id, neuronConfig := range config.Neurons {
		neuronID, _ := strconv.Atoi(id)
		neuron := network.Neurons[neuronID]
		for targetID, weight := range neuronConfig.Connections {
			targetNeuronID, _ := strconv.Atoi(targetID)
			targetNeuron := network.Neurons[targetNeuronID]
			connection := Connection{Weight: weight, From: neuron}
			targetNeuron.Incoming = append(targetNeuron.Incoming, connection)
		}
	}

	return network, nil
}

func (nn *NeuralNetwork) Forward(input [][]int) [][]float64 {
	// Set input neuron outputs
	for i, row := range input {
		for j, value := range row {
			index := i*len(row) + j
			if index < len(nn.InputNeurons) {
				nn.InputNeurons[index].Output = float64(value)
			}
		}
	}

	// Compute output for all neurons
	for _, neuron := range nn.Neurons {
		if !neuron.IsInput { // Skip input neurons since they are set directly
			sum := neuron.Bias
			for _, connection := range neuron.Incoming {
				sum += connection.From.Output * connection.Weight
			}
			neuron.Output = activate(sum, neuron.ActivationType)
		}
	}

	// Collect outputs
	output := make([][]float64, len(nn.OutputNeurons))
	for i, neuron := range nn.OutputNeurons {
		output[i] = []float64{neuron.Output} // Modify as needed to match your specific output structure
	}
	return output
}

func activate(x float64, activationType string) float64 {
	switch activationType {
	case "sigmoid":
		return 1 / (1 + math.Exp(-x))
	case "relu":
		if x > 0 {
			return x
		}
		return 0
	default:
		return x // Linear activation by default
	}
}
