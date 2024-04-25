// Import the custom module you created
const { NeuralNetwork } = require("./simpleNN_json");

const { readFileAsJson } = require("./access");
const { ONN } = require("./ONN");

(async () => {
  console.log("welcome to neat template node/js to golang to compute shader");
  simpleTestingExample();
  try {
    // Use the function to read a file and parse it as JSON
    const jsonObj = await readFileAsJson("network_config.json");
    console.log("JSON Object:", jsonObj);
  } catch (error) {
    // Handle errors that may come from reading the file or parsing JSON
    console.error("Failed to read or parse JSON:", error);
  }

  testing();
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
