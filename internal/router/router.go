package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prplx/lighter.pics/internal/processor"
)

func Register(app *fiber.App, processor *processor.Processor) {
	v1 := app.Group("/api/v1")
	v1.Post("/process", processor.Handle)
	v1.Get("/healthcheck", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status": "ok",
		})
	})
}
