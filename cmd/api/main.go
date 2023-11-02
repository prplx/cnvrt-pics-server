package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/prplx/lighter.pics/internal/archiver"
	"github.com/prplx/lighter.pics/internal/communicator"
	"github.com/prplx/lighter.pics/internal/config"
	"github.com/prplx/lighter.pics/internal/handlers"
	"github.com/prplx/lighter.pics/internal/jsonlog"
	"github.com/prplx/lighter.pics/internal/pg"
	"github.com/prplx/lighter.pics/internal/processor"
	"github.com/prplx/lighter.pics/internal/repositories"
	"github.com/prplx/lighter.pics/internal/router"
	"github.com/prplx/lighter.pics/internal/services"
)

func main() {
	config, err := config.NewConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	db := pg.NewPG(context.Background(), config.DB.DSN)
	defer db.Close()

	fiberApp := fiber.New(fiber.Config{
		BodyLimit: config.Server.BodyLimit * 1024 * 1024,
	})
	fiberApp.Static("/uploads", config.Process.UploadDir, fiber.Static{
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
	archiver := archiver.NewArchiver(config, repositories, logger, communicator)
	processor := processor.NewProcessor(config, repositories, communicator, logger)

	services := services.NewServices(services.Deps{
		Logger:       logger,
		Repositories: repositories,
		Processor:    processor,
		Communicator: communicator,
		Archiver:     archiver,
		Config:       config,
	})

	handlers := handlers.NewHandlers(services)
	router.Register(fiberApp, handlers)

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)))
}
