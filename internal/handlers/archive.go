package handlers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/prplx/lighter.pics/internal/validator"
)

func (h *Handlers) HandleArchiveJob(ctx *fiber.Ctx) error {
	v := validator.NewValidator()
	reqJobID := ctx.Params("jobID")
	if v.Check(reqJobID != "", "jobID", " jobId must be provided"); !v.Valid() {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": v.Errors,
		})
	}

	jobID, err := strconv.Atoi(reqJobID)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	go h.services.Archiver.Archive(jobID)

	return ctx.SendStatus(http.StatusOK)
}
