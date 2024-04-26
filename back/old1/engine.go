package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
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

	// Ensure base directory and projects subdirectory exist
	/*baseDir := "./host"
	projectsDir := filepath.Join(baseDir, "projects")
	ensureDir(baseDir)
	ensureDir(projectsDir)

	// API endpoint to list files and folders in the base directory
	app.Get("/files", func(c *fiber.Ctx) error {
		return c.JSON(listFilesInDir(baseDir))
	})*/

	// Serve React static files - adjust "build" to the path of your React build directory
	//app.Static("/", "./front/build")

	// This route handler can be removed or adjusted if you want to use React Router for routing
	/*app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./front/build/index.html")
	})

	app.Get("/record", func(c *fiber.Ctx) error {
		return c.SendFile("./front/build/index.html")
	})

	app.Get("/preprocessing", func(c *fiber.Ctx) error {
		return c.SendFile("./front/build/index.html")
	})*/

	//setupRoutes(app)

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
