package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
)

var (
	Config         DBConfig
	DBJB           *sql.DB
	scrappingLink  string
	dbNameSync     string
	collectionSync string
	tableNameSync  string
	computerSpecs  string
)

// Define a global variable to hold the connected WebSocket clients
var clients = make(map[*websocket.Conn]bool)

func main() {
	computerSpecs, pcErr := GetSystemInfo()
	if pcErr != nil {
		log.Fatal(pcErr)
	}
	fmt.Println(computerSpecs)
	/*start := time.Now() // Start timing
	testing()
	duration := time.Since(start)                                  // Calculate duration
	fmt.Printf("Execution time: %v ms\n", duration.Milliseconds()) // Print time in milliseconds
	fmt.Printf("Execution time: %v s\n", duration.Seconds())       // Print time in seconds
	*/
	setupDB()
	//startEvaluation()
	startWebServer()
}

func setupDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("un", os.Getenv("dbun"))
	Config = DBConfig{
		Username: os.Getenv("dbun"),
		Password: os.Getenv("dbpw"),
		Hostname: os.Getenv("dbhost"),
	}
	DBJB, err = ConnectToDB(Config)
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	lstTmpDb, _ := ListDatabases(DBJB)
	fmt.Println(lstTmpDb)
	if strings.Contains(strings.Join(lstTmpDb, ","), "openaudata") {
		fmt.Printf("Database '%s' exists in the list.\n")

	} else {
		fmt.Printf("Database '%s' does not exist in the list.\n")
		CreateDatabase(DBJB, "openaudata")
	}

	_, err = DBJB.Exec("USE openaudata")
	if err != nil {
		log.Printf("Failed to select database: %v", err)
		return
	}

	testdataerr := createTestTableAndData(DBJB)
	if err != nil {
		log.Fatalf("Failed to create and populate test table: %v", testdataerr)
	}

}

func startWebServer() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	//app.Get("/traffic", websocket.New(handleWebSocket))

	//str := "Hello, World! 123  !@#$%^&*()"
	//decimals := transformString(str)
	//fmt.Printf("String: %s\nDecimals: %s\n", str, decimals)

	// Set global variables
	dbNameSync = "old_your_db_name"
	collectionSync = "old_your_collection_name"
	tableNameSync = "old_your_table_name"

	app.Get("/msg", websocket.New(handleWebSocket))

	// Example route
	app.Get("/", func(c *fiber.Ctx) error {
		// Set new global variables
		dbNameSync = "new_db_name"
		collectionSync = "new_collection_name"
		tableNameSync = "new_table_name"

		// Send the new values through the WebSocket
		updateFrontend()

		// Send response
		return c.SendString("New values set and sent through WebSocket.")
	})

	/*app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})*/

	setupRoutes(app)

	go func() {
		if err := app.Listen(":4123"); err != nil {
			fmt.Println("Error running server:", err)
		}
	}()

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	fmt.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		fmt.Println("Error shutting down Fiber app:", err)
	}
	fmt.Println("Server shut down.")
}

func handleWebSocket(c *websocket.Conn) {
	clients[c] = true
	defer func() {
		c.Close()
		delete(clients, c)
	}()

	rowCount, errCount := GetRowCount(DBJB, tableNameSync)
	if errCount != nil {
		rowCount = 0
	}

	modelCount, _ := getModelCount(dbNameSync, collectionSync)

	rowCountStr := strconv.Itoa(rowCount)
	modelCountStr := strconv.Itoa(modelCount)

	err := c.WriteJSON(map[string]string{
		"dbNameSync":     dbNameSync,
		"collectionSync": collectionSync,
		"tableNameSync":  tableNameSync,
		"datasetSize":    rowCountStr,
		"modelCount":     modelCountStr,
	})
	if err != nil {
		log.Println("Error writing JSON to WebSocket:", err)
		return
	}

	// Introduce a for loop to keep the connection alive and handle any incoming messages.
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received: %s", msg)
	}
}

// Call this function whenever you need to update the frontend with new data
func updateFrontend() {

	rowCount, err := GetRowCount(DBJB, tableNameSync)
	if err != nil {
		rowCount = 0
	}

	rowCountStr := strconv.Itoa(rowCount)
	modelCount, _ := getModelCount(dbNameSync, collectionSync)
	modelCountStr := strconv.Itoa(modelCount)

	broadcastData(map[string]string{
		"dbNameSync":     dbNameSync,
		"collectionSync": collectionSync,
		"tableNameSync":  tableNameSync,
		"datasetSize":    rowCountStr,
		"modelCount":     modelCountStr,
	})
}

func broadcastData(data map[string]string) {
	for client := range clients {
		err := client.WriteJSON(data)
		if err != nil {
			log.Println("Error writing JSON to WebSocket:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func testing() {
	// Load the neural network configuration from a JSON file
	jsonData, err := ioutil.ReadFile("fixing_layered_network_config.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var networkConfig NetworkConfig
	err = json.Unmarshal(jsonData, &networkConfig)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Example input values - adjust based on your actual input configuration
	inputValues := map[string]float64{
		"1": 1,
		"2": 0.5,
		"3": 0.75,
	}

	outputs := feedforward(&networkConfig, inputValues)
	fmt.Println("Network outputs:", outputs)
}
