package handlers

import (
	"context"
	"strings"

	"io"

	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/h2non/filetype"
	"github.com/prplx/cnvrt/internal/helpers"
	"github.com/prplx/cnvrt/internal/types"
)

func (h *Handlers) getSession(ctx *fiber.Ctx) (*session.Session, error) {
	sessionStore, ok := ctx.Locals("store").(*session.Store)
	if !ok {
		h.services.Logger.PrintError(StoreIsNotFoundInContext)
		return nil, StoreIsNotFoundInContext
	}
	session, err := sessionStore.Get(ctx)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error getting session",
		})
		return nil, err
	}

	return session, nil
}

func (h *Handlers) verifyJobSession(ctx *fiber.Ctx, jobID int64) error {
	session, err := h.getSession(ctx)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error getting session",
		})
		return SessionIsNotFoundInContext
	}

	job, err := h.services.Repositories.Jobs.GetByID(context.Background(), jobID)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error getting job by id",
		})
		return err
	}

	if job.Session != session.ID() {
		h.services.Logger.PrintError(SessionIDDoesNotMatch, types.AnyMap{
			"message": "session id does not match",
		})
		return SessionIDDoesNotMatch
	}

	return nil
}

func (h *Handlers) validateFileType(buf []byte) error {
	if helpers.IsTest() {
		return nil
	}

	kind, _ := filetype.Match(buf)
	allowedExtensions := strings.Split(h.services.Config.App.SupportedFileTypes, ", ")
	if !helpers.Contains(allowedExtensions, kind.Extension) {
		return FileTypeIsNotSupported
	}
	return nil
}

func fileHeaderToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
