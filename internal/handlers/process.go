package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/prplx/lighter.pics/internal/helpers"
	"github.com/prplx/lighter.pics/internal/types"
	"github.com/prplx/lighter.pics/internal/validator"
)

func (h *Handlers) HandleProcessJob(ctx *fiber.Ctx) error {
	v := validator.New()
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

	jobID, err := h.services.Repositories.Jobs.Create(context.Background())
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error creating job",
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

	fileNameToId := map[string]int{}
	for _, file := range dbFiles {
		fileNameToId[file.Name] = file.ID
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
			fileId := fileNameToId[name]
			go h.services.Processor.Process(context.Background(), types.ImageProcessInput{JobID: jobID, FileID: fileId, FileName: name, Format: reqFormat, Quality: quality, Buffer: buffer})
		}
	}

	return ctx.Status(http.StatusAccepted).JSON(fiber.Map{
		"job_id": jobID,
	})
}

func (h *Handlers) HandleProcessFile(ctx *fiber.Ctx) error {
	reqJobID := ctx.Params("jobID")
	if reqJobID == "" {
		h.services.Logger.PrintError(JobIDIsNotFound)
		return ctx.SendStatus(http.StatusBadRequest)
	}

	v := validator.New()
	if validateRequestQueryParams(v, ctx, "format", "quality", "file_id", "width", "height"); !v.Valid() {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": v.Errors,
		})
	}

	format := ctx.Query("format")
	reqQuality := ctx.Query("quality")
	reqFileID := ctx.Query("file_id")
	reqFileIdInt, err := strconv.Atoi(reqFileID)
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

	file, err := h.services.Repositories.Files.GetById(context.Background(), reqFileIdInt)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error getting file by id",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
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

	fileID, err := strconv.Atoi(reqFileID)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing file_id param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	jobID, err := strconv.Atoi(reqJobID)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing job_id param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	go h.services.Processor.Process(context.Background(), types.ImageProcessInput{JobID: jobID, FileID: fileID, FileName: file.Name, Format: format, Quality: quality, Width: reqFileWidth, Height: reqFileHeight, Buffer: buffer})

	return nil
}
