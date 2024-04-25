// Import the fs module from Node.js
const fs = require("fs").promises;

// Async function to read a file and parse it as JSON
async function readFileAsJson(filePath) {
  try {
    // Read the file asynchronously
    const data = await fs.readFile(filePath, "utf8");
    // Parse the data as JSON
    return JSON.parse(data);
  } catch (error) {
    // Handle errors (e.g., file not found, invalid JSON)
    console.error("Error reading or parsing file:", error);
    throw error; // Re-throw the error if needed
  }
}

// Export the function so it can be used in other files
module.exports = {
  readFileAsJson,
};
