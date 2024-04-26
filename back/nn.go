package main

import (
	"math"
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
