package main

import "github.com/gofiber/fiber/v2"

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
}
