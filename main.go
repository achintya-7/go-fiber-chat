package main

import (
	"github.com/achintya-7/go-fiber-chat/configs"
	"github.com/achintya-7/go-fiber-chat/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func main() {

	app := fiber.New()

	// adding cache middleware
	app.Use(cache.New())

	configs.ConnectDB()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"Status":         "200",
			"API is running": "Successfully",
		})
	})

	routes.UserRoute(app)
	routes.ChatRoute(app)

	app.Listen("127.0.0.1:4000")

}
