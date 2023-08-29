package main

import (
	"fmt"
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func (app *application) processHandler(ctx *fiber.Ctx) error {
	// Here will go all validation logic
	image, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	jobID, err := exec.Command("uuidgen").Output()
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.SaveFile(image, fmt.Sprintf("./%s/%s", jobID, image.Filename))
}
