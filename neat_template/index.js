// Load the required modules
const fs = require("fs");

// Import the custom module you created
const { NeuralNetwork } = require("./simpleNN_json");

const { readFileAsJson } = require("./access");
const { ONN } = require("./ONN");

const { FULLNeuralNetwork } = require("./fullFFNN");
const { LayeredNeuralNetwork } = require("./layredfullFFNN");

const { FIXEDLayeredNeuralNetwork } = require("./FixedLayeredNN");

const { startHosting } = require("./hosting");

(async () => {
  console.log("welcome to neat template node/js to golang to compute shader");
  /*simpleTestingExample();
  try {
    // Use the function to read a file and parse it as JSON
    const jsonObj = await readFileAsJson("network_config.json");
    console.log("JSON Object:", jsonObj);
  } catch (error) {
    // Handle errors that may come from reading the file or parsing JSON
    console.error("Failed to read or parse JSON:", error);
  }
  */
  //testing();
  //testingfullffnnrun();
  //testingLayeredFullFFNNrun();
  fixedtestingLayeredFullFFNNrun();

  startHosting();
})();

function testing() {
  var objNN = new ONN();
}

function simpleTestingExample() {
  // JSON configuration for the neural network
  const networkConfig = {
    layers: [
      { name: "input", nodes: 2 },
      { name: "hidden1", nodes: 3 },
      { name: "output", nodes: 1 },
    ],
  };

  // Initialize the neural network with the JSON configuration
  const nn = new NeuralNetwork(networkConfig);

  // Example input
  const input = [1, 0];
  console.log("Output:", nn.feedforward(input));
}

function testingfullffnnrun() {
  console.log("==============starting full ff nn===========");
  // Load the neural network configuration and run it
  fs.readFile("network_config.json", (err, data) => {
    if (err) {
      console.error("Error reading file:", err);
      return;
    }
    const config = JSON.parse(data);
    console.log(config);
    const nn = new FULLNeuralNetwork(config);

    // Example: Set input values here as needed
    const outputs = nn.feedforward({ 1: 1, 2: 0.5, 3: 0.75 });
    console.log("Neural network outputs:", outputs);
  });
}

function testingLayeredFullFFNNrun() {
  console.log("==============starting LAYERED full ff nn===========");
  // Loading the neural network configuration from a JSON file and running it
  fs.readFile("layered_network_config.json", "utf8", (err, data) => {
    if (err) {
      console.error("Error reading file:", err);
      return;
    }

    const networkConfig = JSON.parse(data);
    const nn = new LayeredNeuralNetwork(networkConfig);

    // Example input values - adjust based on your actual input configuration
    const inputValues = { 1: 1, 2: 0.5, 3: 0.75 };
    const outputs = nn.feedforward(inputValues);
    console.log("Network outputs:", outputs);
  });
}

function fixedtestingLayeredFullFFNNrun() {
  console.log("==============starting LAYERED full ff nn===========");
  // Loading the neural network configuration from a JSON file and running it
  fs.readFile("fixing_layered_network_config.json", "utf8", (err, data) => {
    if (err) {
      console.error("Error reading file:", err);
      return;
    }

    const networkConfig = JSON.parse(data);
    const nn = new FIXEDLayeredNeuralNetwork(networkConfig);

    // Example input values - adjust based on your actual input configuration
    const inputValues = { 1: 1, 2: 0.5, 3: 0.75 };
    const outputs = nn.feedforward(inputValues);
    console.log("Network outputs:", outputs);
  });
}
