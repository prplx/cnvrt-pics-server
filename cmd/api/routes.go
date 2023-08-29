package main

import (
	"github.com/gofiber/fiber/v2"
)

func (app *application) routes(fiberApp *fiber.App) {
	v1 := fiberApp.Group("/api/v1")
	v1.Post("/process", app.processHandler)
	v1.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "pong",
		})
	})
}
