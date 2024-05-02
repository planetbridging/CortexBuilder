package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var (
	Config        DBConfig
	DBJB          *sql.DB
	scrappingLink string
)

func main() {
	/*start := time.Now() // Start timing
	testing()
	duration := time.Since(start)                                  // Calculate duration
	fmt.Printf("Execution time: %v ms\n", duration.Milliseconds()) // Print time in milliseconds
	fmt.Printf("Execution time: %v s\n", duration.Seconds())       // Print time in seconds
	*/
	setupDB()
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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

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
