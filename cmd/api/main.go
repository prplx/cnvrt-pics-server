package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/prplx/lighter.pics/internal/communicator"
	"github.com/prplx/lighter.pics/internal/processor"
)

// type application struct {
// }

func main() {
	fiberApp := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024,
	})
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	communicator := communicator.NewCommunicator()
	processor := processor.NewProcessor(communicator)
	v1 := fiberApp.Group("/api/v1")
	v1.Post("/process", processor.Handle)
	v1.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "pong",
		})
	})

	log.Fatal(fiberApp.Listen(":3001"))
}
