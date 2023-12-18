package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	firebaseAdmin "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/appcheck"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/prplx/lighter.pics/internal/handlers"
	"github.com/prplx/lighter.pics/internal/types"
)

const (
	uploadsDir            = "/uploads"
	metricsEndpoint       = "/metrics"
	healthcheckEndpoint   = "/healthcheck"
	firebaseAppCheckQuery = "appCheckToken"
)

var (
	appCheck *appcheck.Client
)

func Register(app *fiber.App, handlers *handlers.Handlers, config *types.Config) {
	firebaseApp, err := firebaseAdmin.NewApp(context.Background(), &firebaseAdmin.Config{
		ProjectID: config.Firebase.ProjectID,
	})
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	appCheck, err = firebaseApp.AppCheck(context.Background())
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	app.Static(uploadsDir, config.Process.UploadDir, fiber.Static{
		Download: true,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.Server.AllowOrigins,
		AllowHeaders: config.Server.AllowHeaders,
		AllowMethods: config.Server.AllowMethods,
	}))
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			config.App.MetricsUser: config.App.MetricsPassword,
		},
		Next: func(c *fiber.Ctx) bool {
			return c.Path() != metricsEndpoint
		},
	}))
	app.Use(func(c *fiber.Ctx) error {
		if strings.Contains(c.Path(), "/ws") || strings.Contains(c.Path(), healthcheckEndpoint) {
			return c.Next()
		}

		var appCheckToken string

		if strings.Contains(c.Path(), "/ws") {
			appCheckToken = c.Query(firebaseAppCheckQuery)
		} else {
			appCheckToken = c.Get(config.Firebase.AppCheckHeader)
		}

		if err := requireAppCheck(appCheck, appCheckToken); err != nil {
			return c.Status(http.StatusForbidden).SendString("Forbidden")
		}

		return c.Next()

	})
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get(metricsEndpoint, monitor.New())
	app.Get(healthcheckEndpoint, handlers.Healthcheck)

	v1 := app.Group("/api/v1")
	v1.Post("/process", handlers.HandleProcessJob)
	v1.Post("/process/:jobID", handlers.HandleProcessFile)
	v1.Post("/archive/:jobID", handlers.HandleArchiveJob)
	v1.Get("/ws/:jobID", websocket.New(handlers.HandleWebsocket))
}

func requireAppCheck(appCheck *appcheck.Client, appCheckToken string) error {
	if _, err := appCheck.VerifyToken(appCheckToken); err != nil {
		return fmt.Errorf("AppCheck token verification failed: %v", err)
	}

	return nil
}
