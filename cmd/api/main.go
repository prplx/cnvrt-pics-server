package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/prplx/cnvrt/internal/archiver"
	communicator "github.com/prplx/cnvrt/internal/communicator/communicatorwebsocket"
	"github.com/prplx/cnvrt/internal/config"
	"github.com/prplx/cnvrt/internal/handlers"
	"github.com/prplx/cnvrt/internal/jsonlog"
	"github.com/prplx/cnvrt/internal/pg"
	processor "github.com/prplx/cnvrt/internal/processor/processorgovips"
	"github.com/prplx/cnvrt/internal/repositories"
	"github.com/prplx/cnvrt/internal/router"
	"github.com/prplx/cnvrt/internal/scheduler"
	"github.com/prplx/cnvrt/internal/services"
)

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

	app := fiber.New(fiber.Config{
		BodyLimit: config.Server.BodyLimit * 1024 * 1024,
	})

	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)
	repositories := repositories.NewRepositories(db.Pool)
	communicator := communicator.NewCommunicator()
	archiver := archiver.NewArchiver(config, repositories.Files, logger, communicator)
	scheduler := scheduler.NewScheduler(config, logger, communicator, repositories.Jobs)
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
	router.Register(app, handlers, config, db.Pool)

	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)))
}
