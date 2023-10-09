package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/prplx/lighter.pics/internal/imageProcessor"
	"github.com/prplx/lighter.pics/internal/jsonlog"
	"github.com/prplx/lighter.pics/internal/pg"
	"github.com/prplx/lighter.pics/internal/processor"
	"github.com/prplx/lighter.pics/internal/repositories"
	"github.com/prplx/lighter.pics/internal/router"
	"github.com/prplx/lighter.pics/internal/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := pg.NewPG(context.Background(), os.Getenv("DB_DSN"))
	defer db.Close()

	fiberApp := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024,
	})
	fiberApp.Static("/uploads", "./uploads")
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	services := services.NewServices(services.Deps{
		Logger:         jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo),
		Repositories:   *repositories.NewRepositories(db.Pool),
		ImageProcessor: imageProcessor.NewImageProcessor(),
	})
	processor := processor.NewProcessor(services)
	router.Register(fiberApp, processor)

	log.Fatal(fiberApp.Listen(":3001"))
}
