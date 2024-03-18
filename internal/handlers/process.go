package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/pkg/errors"
	"github.com/prplx/cnvrt/internal/helpers"
	"github.com/prplx/cnvrt/internal/types"
	"github.com/prplx/cnvrt/internal/validator"
)

func (h *Handlers) HandleProcessJob(ctx *fiber.Ctx) error {
	session, err := h.getSession(ctx)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	session.Regenerate()

	v := validator.NewValidator()
	if validateRequestQueryParams(v, ctx, "format", "quality"); !v.Valid() {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": v.Errors,
		})
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing multipart form",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}
	jobID, err := h.services.Repositories.Jobs.Create(context.Background(), session.ID())
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error creating job",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	err = session.Save()
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error saving session",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	path := fmt.Sprintf(h.services.Config.Process.UploadDir+"/%d", jobID)
	reqFormat := ctx.Query("format")
	reqQuality := ctx.Query("quality")
	quality, err := strconv.Atoi(reqQuality)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing quality param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			h.services.Logger.PrintError(err, types.AnyMap{
				"message": "error creating directory",
			})
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	fileNames := []string{}
	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			err = ctx.SaveFile(fileHeader, helpers.BuildPath(h.services.Config.Process.UploadDir, jobID, fileHeader.Filename))
			if err != nil {
				h.services.Logger.PrintError(err, types.AnyMap{
					"message": "error saving file",
				})
				return ctx.SendStatus(http.StatusInternalServerError)
			}
			fileNames = append(fileNames, fileHeader.Filename)
		}
	}

	dbFiles, err := h.services.Repositories.Files.CreateBulk(context.Background(), jobID, fileNames)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error creating file records",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	fileNameToID := map[string]int64{}
	for _, file := range dbFiles {
		fileNameToID[file.Name] = file.ID
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	buffers := [][]byte{}
	for _, file := range files {
		filePath := helpers.BuildPath(h.services.Config.Process.UploadDir, jobID, file.Name())
		buffer, err := h.readFile(filePath)
		if err != nil {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		buffers = append(buffers, buffer)
	}

	if len(buffers) == len(files) {
		for idx, buffer := range buffers {
			name := files[idx].Name()
			fileID := fileNameToID[name]
			go h.services.Processor.Process(context.Background(), types.ImageProcessInput{JobID: jobID, FileID: fileID, FileName: name, Format: reqFormat, Quality: quality, Buffer: buffer})
		}
	}

	return ctx.Status(http.StatusAccepted).JSON(fiber.Map{
		"job_id": jobID,
	})
}

func (h *Handlers) HandleProcessFile(ctx *fiber.Ctx) error {
	session, err := h.getSession(ctx)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	reqJobID := ctx.Params("jobID")
	if reqJobID == "" {
		h.services.Logger.PrintError(JobIDIsNotFound)
		return ctx.SendStatus(http.StatusBadRequest)
	}

	v := validator.NewValidator()
	if validateRequestQueryParams(v, ctx, "format", "quality", "file_id", "width", "height"); !v.Valid() {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": v.Errors,
		})
	}

	format := ctx.Query("format")
	reqQuality := ctx.Query("quality")
	reqFileID := ctx.Query("file_id")
	reqFileIDInt, err := strconv.ParseInt(reqFileID, 10, 64)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing file_id param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}
	reqFileWidth, err := strconv.Atoi(ctx.Query("width"))
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing width param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	reqFileHeight, err := strconv.Atoi(ctx.Query("height"))
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing height param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	file, err := h.services.Repositories.Files.GetWithJobByID(context.Background(), reqFileIDInt)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error getting file by id",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	if file.Job.Session != session.ID() {
		h.services.Logger.PrintError(SessionIDDoesNotMatch, types.AnyMap{
			"message": "session id does not match",
		})
		return ctx.SendStatus(http.StatusForbidden)
	}
	filePath := helpers.BuildPath(h.services.Config.Process.UploadDir, reqJobID, file.Name)
	buffer, err := h.readFile(filePath)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	quality, err := strconv.Atoi(reqQuality)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing quality param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	fileID, err := strconv.ParseInt(reqFileID, 10, 64)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing file_id param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	jobID, err := strconv.ParseInt(reqJobID, 10, 64)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing job_id param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	go h.services.Processor.Process(context.Background(), types.ImageProcessInput{JobID: jobID, FileID: fileID, FileName: file.Name, Format: format, Quality: quality, Width: reqFileWidth, Height: reqFileHeight, Buffer: buffer})

	return nil
}

func (h *Handlers) HandleAddFileToJob(ctx *fiber.Ctx) error {
	session, err := h.getSession(ctx)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	reqJobID := ctx.Params("jobID")
	if reqJobID == "" {
		h.services.Logger.PrintError(JobIDIsNotFound)
		return ctx.SendStatus(http.StatusBadRequest)
	}

	v := validator.NewValidator()
	if validateRequestQueryParams(v, ctx, "format", "quality"); !v.Valid() {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": v.Errors,
		})
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing multipart form",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	path := fmt.Sprintf(h.services.Config.Process.UploadDir+"/%s", reqJobID)
	reqFormat := ctx.Query("format")
	reqQuality := ctx.Query("quality")
	jobID, err := strconv.ParseInt(reqJobID, 10, 64)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing job_id param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}
	quality, err := strconv.Atoi(reqQuality)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing quality param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "job directory not found",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	dbJob, err := h.services.Repositories.Jobs.GetByID(context.Background(), jobID)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error getting job by id",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	if dbJob.Session != session.ID() {
		h.services.Logger.PrintError(SessionIDDoesNotMatch, types.AnyMap{
			"message": "session id does not match",
		})
		return ctx.SendStatus(http.StatusForbidden)
	}

	dbFiles, err := h.services.Repositories.Files.GetByJobID(context.Background(), jobID)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error getting files by job id",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	fileNameToID := map[string]int64{}
	for _, file := range dbFiles {
		fileNameToID[file.Name] = file.ID
	}

	var fileName string
	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			if _, ok := fileNameToID[fileHeader.Filename]; ok {
				h.services.Logger.PrintError(err, types.AnyMap{
					"message": "file already exists",
				})
				return ctx.SendStatus(http.StatusBadRequest)
			}
			err = ctx.SaveFile(fileHeader, helpers.BuildPath(h.services.Config.Process.UploadDir, jobID, fileHeader.Filename))
			if err != nil {
				h.services.Logger.PrintError(err, types.AnyMap{
					"message": "error saving file",
				})
				return ctx.SendStatus(http.StatusInternalServerError)
			}
			fileName = fileHeader.Filename
		}
	}

	dbFile, err := h.services.Repositories.Files.AddToJob(context.Background(), jobID, fileName)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error adding file to job",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	buffer, err := h.readFile(helpers.BuildPath(h.services.Config.Process.UploadDir, jobID, fileName))
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error reading file",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	go h.services.Processor.Process(context.Background(), types.ImageProcessInput{JobID: jobID, FileID: dbFile.ID, FileName: fileName, Format: reqFormat, Quality: quality, Buffer: buffer})

	return nil
}

func (h *Handlers) HandleDeleteFileFromJob(ctx *fiber.Ctx) error {
	session, err := h.getSession(ctx)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	reqJobID := ctx.Params("jobID")
	reqFileID := ctx.Query("file_id")
	if reqJobID == "" || reqFileID == "" {
		h.services.Logger.PrintError(JobIDIsNotFound)
		return ctx.SendStatus(http.StatusBadRequest)
	}

	jobID, err := strconv.ParseInt(reqJobID, 10, 64)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing job_id param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	fileID, err := strconv.ParseInt(reqFileID, 10, 64)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing file_id param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	dbJob, err := h.services.Repositories.Jobs.GetByID(context.Background(), jobID)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error getting job by id",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	if dbJob.Session != session.ID() {
		h.services.Logger.PrintError(SessionIDDoesNotMatch, types.AnyMap{
			"message": "session id does not match",
		})
		return ctx.SendStatus(http.StatusForbidden)
	}

	err = h.services.Repositories.Files.DeleteFromJob(context.Background(), jobID, fileID)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error deleting file",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return nil
}

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
