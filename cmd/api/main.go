package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/prplx/lighter.pics/internal/archiver"
	communicator "github.com/prplx/lighter.pics/internal/communicator/communicatorpusher"
	"github.com/prplx/lighter.pics/internal/config"
	"github.com/prplx/lighter.pics/internal/handlers"
	"github.com/prplx/lighter.pics/internal/jsonlog"
	"github.com/prplx/lighter.pics/internal/pg"
	processor "github.com/prplx/lighter.pics/internal/processor/processorgovips"
	"github.com/prplx/lighter.pics/internal/repositories"
	"github.com/prplx/lighter.pics/internal/router"
	"github.com/prplx/lighter.pics/internal/scheduler"
	"github.com/prplx/lighter.pics/internal/services"
)

const uploadsDir = "/uploads"

func main() {
	config, err := config.NewConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	db := pg.NewPG(ctx, config.DB.DSN)
	if err := db.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fiberApp := fiber.New(fiber.Config{
		BodyLimit: config.Server.BodyLimit * 1024 * 1024,
	})
	fiberApp.Static(uploadsDir, config.Process.UploadDir, fiber.Static{
		Download: true,
	})
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: config.Server.AllowOrigins,
		AllowHeaders: config.Server.AllowHeaders,
		AllowMethods: config.Server.AllowMethods,
	}))

	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)
	repositories := repositories.NewRepositories(db.Pool)
	communicator := communicator.NewCommunicator(config)
	archiver := archiver.NewArchiver(config, repositories.Files, logger, communicator)
	scheduler := scheduler.NewScheduler(config, logger, communicator)
	processor := processor.NewProcessor(config, repositories.Operations, communicator, logger, scheduler)

	processor.Startup()
	defer processor.Shutdown()

	services := services.NewServices(services.Deps{
		Logger:       logger,
		Repositories: repositories,
		Processor:    processor,
		Communicator: communicator,
		Scheduler:    scheduler,
		Archiver:     archiver,
		Config:       config,
	})

	handlers := handlers.NewHandlers(services)
	router.Register(fiberApp, handlers)

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)))
}
