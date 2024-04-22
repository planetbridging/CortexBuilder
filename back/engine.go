package main

import "fmt"

func main() {
	network, err := NewNetworkFromFile("network_config.json")
	if err != nil {
		fmt.Println("Error creating network:", err)
		return
	}

	// Example multidimensional integer input
	input := [][]int{{1, 2}, {3, 4}}
	output := network.Forward(input)
	fmt.Println("Network output:", output)
}
