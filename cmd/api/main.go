package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/prplx/lighter.pics/internal/processor"
	"github.com/prplx/lighter.pics/internal/services"
	"github.com/prplx/lighter.pics/pkg/jsonlog"
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

	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)
	services := services.NewServices(services.Deps{
		Logger: logger,
	})
	processor := processor.NewProcessor(services)
	v1 := fiberApp.Group("/api/v1")
	v1.Post("/process", processor.Handle)
	v1.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "pong",
		})
	})

	log.Fatal(fiberApp.Listen(":3001"))
}
