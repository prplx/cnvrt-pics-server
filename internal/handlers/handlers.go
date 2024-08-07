package handlers

import (
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/prplx/cnvrt/internal/services"
	"github.com/prplx/cnvrt/internal/types"
)

type Handlers struct {
	services *services.Services
}

func NewHandlers(s *services.Services) *Handlers {
	return &Handlers{
		services: s,
	}
}

func (h *Handlers) Healthcheck(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status": "ok",
	})
}

func (h *Handlers) readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error opening file",
			"path":    path,
		})
		return nil, OpenFileError
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error reading file",
			"path":    path,
		})
		return nil, ReadingFileError
	}

	return buffer, nil
}
