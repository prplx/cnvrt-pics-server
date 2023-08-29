package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type application struct{}

func main() {
	fiberApp := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024,
	})
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app := &application{}
	app.routes(fiberApp)

	log.Fatal(fiberApp.Listen(":3001"))
}
