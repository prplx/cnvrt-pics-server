package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/prplx/lighter.pics/internal/imageProcessor"
	"github.com/prplx/lighter.pics/internal/processor"
	"github.com/prplx/lighter.pics/internal/repositories"
	"github.com/prplx/lighter.pics/internal/router"
	"github.com/prplx/lighter.pics/internal/services"
	"github.com/prplx/lighter.pics/pkg/jsonlog"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_DSN"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

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
		Repositories:   *repositories.NewRepositories(conn),
		ImageProcessor: imageProcessor.NewImageProcessor(),
	})
	processor := processor.NewProcessor(services)
	router.Register(fiberApp, processor)

	log.Fatal(fiberApp.Listen(":3001"))
}
