class FULLNeuralNetwork {
  constructor(networkConfig) {
    this.networkConfig = networkConfig;
    this.neurons = {};
  }

  activate(activationType, input) {
    // Simple activation functions
    switch (activationType) {
      case "relu":
        return Math.max(0, input);
      case "sigmoid":
        return 1 / (1 + Math.exp(-input));
      default:
        return input;
    }
  }

  feedforward(inputValues) {
    // Set initial values for input neurons
    for (let inputId in inputValues) {
      this.neurons[inputId] = inputValues[inputId];
    }

    // Compute values for non-input neurons
    for (let nodeId in this.networkConfig.neurons) {
      let node = this.networkConfig.neurons[nodeId];
      if (!node.isInput) {
        let sum = 0;
        for (let inputId in node.connections) {
          let connection = node.connections[inputId];
          sum += this.neurons[inputId] * connection.weight; // Accumulate weighted input
        }
        sum += node.bias; // Add bias
        this.neurons[nodeId] = this.activate(node.activationType, sum); // Apply activation function
      }
    }

    // Extract and return outputs
    let outputs = {};
    for (let nodeId in this.networkConfig.neurons) {
      let node = this.networkConfig.neurons[nodeId];
      if (node.isOutput) {
        outputs[nodeId] = this.neurons[nodeId];
      }
    }
    return outputs;
  }
}

module.exports = {
  FULLNeuralNetwork,
};
