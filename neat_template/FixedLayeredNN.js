class FIXEDLayeredNeuralNetwork {
  constructor(networkConfig) {
    this.networkConfig = networkConfig;
    this.neurons = {};
  }

  activate(activationType, input) {
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
    // Initialize input layer neurons with input values
    for (let inputId in this.networkConfig.layers.input.neurons) {
      this.neurons[inputId] = inputValues[inputId];
    }

    console.log("-----------this.neurons----------");
    console.log(this.neurons);

    // Process hidden layers
    for (let layer of this.networkConfig.layers.hidden) {
      for (let nodeId in layer.neurons) {
        const node = layer.neurons[nodeId];
        let sum = 0;
        for (let inputId in node.connections) {
          const connection = node.connections[inputId];
          sum += this.neurons[inputId] * connection.weight;
        }
        sum += node.bias;
        this.neurons[nodeId] = this.activate(node.activationType, sum);

        console.log("-----------this.neurons" + nodeId + "----------");
        console.log(this.neurons);
      }
    }

    // Process output layer
    const outputs = {};
    for (let nodeId in this.networkConfig.layers.output.neurons) {
      const node = this.networkConfig.layers.output.neurons[nodeId];
      let sum = 0;
      for (let inputId in node.connections) {
        const connection = node.connections[inputId];
        sum += this.neurons[inputId] * connection.weight;
      }
      sum += node.bias;
      outputs[nodeId] = this.activate(node.activationType, sum);
    }

    return outputs;
  }
}

module.exports = {
  FIXEDLayeredNeuralNetwork,
};
