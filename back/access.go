package main

import "io/ioutil"

func NewNetworkFromFile(filename string) (NeuralNetwork, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return NeuralNetwork{}, err
	}

	return NewNetworkFromJSON(data)
}
