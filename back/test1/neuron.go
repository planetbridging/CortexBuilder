package main

import (
	"math"
	"math/rand"
)

type Neuron struct {
	Weights []float64
	Bias    float64
}

type Layer struct {
	Neurons []Neuron
}

type NeuralNetwork struct {
	Layers []Layer
}

func NewNeuron(inputCount int) Neuron {
	neuron := Neuron{
		Weights: make([]float64, inputCount),
		Bias:    rand.Float64(),
	}
	for i := range neuron.Weights {
		neuron.Weights[i] = rand.Float64()*2 - 1 // Initialize with random weights between -1 and 1
	}
	return neuron
}

func NewLayer(neuronCount, inputCount int) Layer {
	layer := Layer{
		Neurons: make([]Neuron, neuronCount),
	}
	for i := range layer.Neurons {
		layer.Neurons[i] = NewNeuron(inputCount)
	}
	return layer
}

func NewNetwork(layerSizes []int) NeuralNetwork {
	network := NeuralNetwork{
		Layers: make([]Layer, len(layerSizes)-1),
	}
	for i := 0; i < len(network.Layers); i++ {
		network.Layers[i] = NewLayer(layerSizes[i+1], layerSizes[i])
	}
	return network
}

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func (nn *NeuralNetwork) Forward(input []float64) []float64 {
	currentResult := input
	for _, layer := range nn.Layers {
		nextResult := make([]float64, len(layer.Neurons))
		for i, neuron := range layer.Neurons {
			sum := neuron.Bias
			for j, input := range currentResult {
				sum += input * neuron.Weights[j]
			}
			nextResult[i] = sigmoid(sum) // Apply activation function
		}
		currentResult = nextResult
	}
	return currentResult
}
