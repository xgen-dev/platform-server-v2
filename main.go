package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type Hello struct {
	Hello string `json:"hello"`
}

func main() {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	log.Fatal(app.Listen(":5050"))
}
