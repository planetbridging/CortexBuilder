package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/files/*", func(c *fiber.Ctx) error {
		// Extracting the subpath or handling root
		subPath := c.Params("*")
		if subPath == "" {
			subPath = "./host" // default path
		}

		files, err := listFilesInDir(subPath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Assuming the folder chain is a breadcrumb of paths
		folderChain := []map[string]string{{"id": "root", "name": "Home"}}
		// Add more logic here to build the folder chain based on the path

		response := map[string]interface{}{
			"files":       files,
			"folderChain": folderChain,
		}

		return c.JSON(response)
	})

	app.Post("/initialize", func(c *fiber.Ctx) error {
		var request struct {
			NetworkType      string `json:"networkType"`
			SpawnCount       int    `json:"spawnCount"`
			AdditionalParam1 string `json:"additionalParam1"`
			AdditionalParam2 string `json:"additionalParam2"`
			AdditionalParam3 string `json:"additionalParam3"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request"})
		}

		fmt.Printf("Received initialization request: %+v\n", request)

		// Here, you would add your logic to handle the training initialization
		// For example, setting up configurations, preparing datasets, etc.

		return c.Status(fiber.StatusOK).SendString("Training initialization started successfully")
	})
}
