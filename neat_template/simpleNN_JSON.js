class NeuralNetwork {
  constructor(networkConfig) {
    // Initialize the network from a JSON configuration
    this.networkConfig = networkConfig;
    this.initializeWeights();
  }

  initializeWeights() {
    // Initialize weights and biases for each layer based on the network configuration
    this.weights = {};
    this.biases = {};

    // Loop through each layer and initialize weights and biases randomly
    for (let i = 0; i < this.networkConfig.layers.length - 1; i++) {
      const layer = this.networkConfig.layers[i];
      const nextLayer = this.networkConfig.layers[i + 1];

      // Initialize weights matrix for current layer to next layer
      this.weights[layer.name + "_to_" + nextLayer.name] = Array.from(
        { length: layer.nodes },
        () =>
          Array.from({ length: nextLayer.nodes }, () => Math.random() * 2 - 1)
      );

      // Initialize biases for next layer
      this.biases[nextLayer.name] = Array.from(
        { length: nextLayer.nodes },
        () => Math.random() * 2 - 1
      );
    }

    console.log("================weights===============");
    console.log(this.weights);
    console.log("================biases===============");
    console.log(this.biases);
  }

  sigmoid(x) {
    // Sigmoid activation function to normalize inputs between 0 and 1
    return 1 / (1 + Math.exp(-x));
  }

  feedforward(inputArray) {
    let inputs = inputArray;

    // Process each layer
    for (let i = 0; i < this.networkConfig.layers.length - 1; i++) {
      const layer = this.networkConfig.layers[i];
      const nextLayer = this.networkConfig.layers[i + 1];
      const weights = this.weights[layer.name + "_to_" + nextLayer.name];
      const biases = this.biases[nextLayer.name];

      // Calculate outputs for the next layer
      let outputs = new Array(nextLayer.nodes).fill(0).map((_, nodeIndex) => {
        let sum = inputs.reduce(
          (acc, input, inputIndex) =>
            acc + input * weights[inputIndex][nodeIndex],
          0
        );
        sum += biases[nodeIndex];
        console.log("================sum===============");
        console.log(sum);
        return this.sigmoid(sum); // Apply activation function
      });

      console.log("================outputs===============");
      console.log(outputs);

      inputs = outputs; // Set outputs as inputs for the next layer
    }

    return inputs; // Return the output of the last layer
  }
}

module.exports = {
  NeuralNetwork,
};
