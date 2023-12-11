package router

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/prplx/lighter.pics/internal/handlers"
)

func Register(app *fiber.App, handlers *handlers.Handlers) {
	v1 := app.Group("/api/v1")
	v1.Post("/process", handlers.HandleProcessJob)
	v1.Post("/process/:jobID", handlers.HandleProcessFile)
	v1.Post("/archive/:jobID", handlers.HandleArchiveJob)
	v1.Get("/healthcheck", handlers.Healthcheck)
	v1.Get("/ws/:jobID", websocket.New(handlers.HandleWebsocket))
}
