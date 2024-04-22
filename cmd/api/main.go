package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prplx/cnvrt/internal/config"
	"github.com/prplx/cnvrt/internal/handlers"
	"github.com/prplx/cnvrt/internal/pg"
	"github.com/prplx/cnvrt/internal/repositories"
	"github.com/prplx/cnvrt/internal/router"
	"github.com/prplx/cnvrt/internal/services"
	"github.com/prplx/cnvrt/internal/services/archiver"
	communicator "github.com/prplx/cnvrt/internal/services/communicator/communicatorwebsocket"
	"github.com/prplx/cnvrt/internal/services/jsonlog"
	processor "github.com/prplx/cnvrt/internal/services/processor/processorgovips"
	"github.com/prplx/cnvrt/internal/services/scheduler"
	"github.com/prplx/cnvrt/internal/types"
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
	scheduler := scheduler.NewScheduler(config, logger, communicator, repositories.Jobs, repositories.PlannedFlushes)
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

	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		logger.PrintInfo("Caught signal", types.AnyMap{"signal": s})
		ctx, cancel := context.WithTimeout(ctx, time.Duration(config.Server.ShutdownTimeout)*time.Second)
		defer cancel()

		shutdownError <- app.ShutdownWithContext(ctx)
	}()

	logger.PrintInfo("Starting server", types.AnyMap{"host": config.Server.Host, "port": config.Server.Port})
	app.Listen(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port))

	if err = <-shutdownError; err != nil {
		log.Fatal(err)
	}

	logger.PrintInfo("Server stopped", nil)
}
