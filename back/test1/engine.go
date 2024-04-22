package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	// Create a network with 3 inputs, 3 neurons in the hidden layer, and 2 outputs
	network := NewNetwork([]int{3, 3, 2})

	// Example input
	input := []float64{0.5, 0.3, 0.2}
	output := network.Forward(input)
	fmt.Println("Output:", output)
}
