// Import the custom module you created
const { readFileAsJson } = require("./access");
const { ONN } = require("./ONN");

(async () => {
  console.log("welcome to neat template node/js to golang to compute shader");

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
