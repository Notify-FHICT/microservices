package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Notes!")
	})

	app.Get("/api/notes/*", func(c *fiber.Ctx) error {
		return c.SendString("API path: " + c.Params("*"))
		// => API path: user/john
	})

	app.Listen(":3000")
}
