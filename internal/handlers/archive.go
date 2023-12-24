package handlers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) HandleArchiveJob(ctx *fiber.Ctx) error {
	reqJobID := ctx.Params("jobID")

	jobID, err := strconv.Atoi(reqJobID)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": "jobID must be a number",
		})
	}

	go h.services.Archiver.Archive(jobID)

	return ctx.SendStatus(http.StatusOK)
}
