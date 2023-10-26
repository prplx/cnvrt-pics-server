package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/prplx/lighter.pics/internal/archiver"
	"github.com/prplx/lighter.pics/internal/communicator"
	"github.com/prplx/lighter.pics/internal/config"
	"github.com/prplx/lighter.pics/internal/handlers"
	"github.com/prplx/lighter.pics/internal/imageProcessor"
	"github.com/prplx/lighter.pics/internal/jsonlog"
	"github.com/prplx/lighter.pics/internal/pg"
	"github.com/prplx/lighter.pics/internal/repositories"
	"github.com/prplx/lighter.pics/internal/router"
	"github.com/prplx/lighter.pics/internal/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config, err := config.NewConfig("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	db := pg.NewPG(context.Background(), os.Getenv("DB_DSN"))
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
	communicator := communicator.NewCommunicator()
	archiver := archiver.NewArchiver(repositories, logger, communicator)
	imageProcessor := imageProcessor.NewImageProcessor(config, repositories, communicator, logger)

	services := services.NewServices(services.Deps{
		Logger:         logger,
		Repositories:   repositories,
		ImageProcessor: imageProcessor,
		Communicator:   communicator,
		Archiver:       archiver,
		Config:         config,
	})

	handlers := handlers.NewHandlers(services)
	router.Register(fiberApp, handlers)

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)))
}
