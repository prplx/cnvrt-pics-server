package router

import (
	"context"
	"log"
	"time"

	firebaseAdmin "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/appcheck"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gofiber/storage/postgres/v3"
	"github.com/prplx/cnvrt/internal/handlers"
	"github.com/prplx/cnvrt/internal/helpers"
	"github.com/prplx/cnvrt/internal/types"
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

func Register(app *fiber.App, handlers *handlers.Handlers, config *types.Config, dbPool *pgxpool.Pool) {
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
	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.Server.AllowOrigins,
		AllowHeaders:     config.Server.AllowHeaders,
		AllowMethods:     config.Server.AllowMethods,
		AllowCredentials: true,
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
		c.Locals("config", config)
		return c.Next()
	})
	app.Use(checkAppCheckToken)
	sessionConfig := session.Config{
		Expiration:     24 * time.Hour,
		KeyLookup:      "cookie:session_id",
		CookieDomain:   helpers.MustGetHostnameFromURL(config.Server.AllowOrigins),
		CookieHTTPOnly: true,
	}
	limiterConfig := limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return helpers.IsTest()
		},
		Max:               2,
		Expiration:        1 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}
	if !helpers.IsTest() {
		limiterConfig.Storage = postgres.New(postgres.Config{
			DB:    dbPool,
			Table: "ratelimit",
		})
		sessionConfig.Storage = postgres.New(postgres.Config{
			DB:    dbPool,
			Table: "session",
		})
	}
	if helpers.IsTest() {
		sessionConfig.KeyGenerator = func() string {
			return "session"
		}
	}
	app.Use(limiter.New(limiterConfig))
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	store := session.New(sessionConfig)
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("store", store)
		return c.Next()
	})
	app.Get(metricsEndpoint, monitor.New())
	app.Get(healthcheckEndpoint, handlers.Healthcheck)

	v1 := app.Group("/api/v1")
	v1.Post("/process", checkFormFileLength, handlers.HandleProcessJob)
	v1.Post("/process/:jobID", handlers.HandleProcessFile)
	v1.Put("/process/:jobID", handlers.HandleAddFileToJob)
	v1.Delete("/process/:jobID", handlers.HandleDeleteFileFromJob)
	v1.Post("/archive/:jobID", handlers.HandleArchiveJob)
	v1.Get("/ws/:jobID", websocket.New(func(c *websocket.Conn) {
		handlers.HandleWebsocket(c)
	}))
}
