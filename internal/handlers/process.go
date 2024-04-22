package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/prplx/cnvrt/internal/helpers"
	"github.com/prplx/cnvrt/internal/types"
)

func (h *Handlers) HandleProcessJob(ctx *fiber.Ctx) error {
	ctxb := context.Background()

	session, err := h.getSession(ctx)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	session.Regenerate()

	form, err := ctx.MultipartForm()
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error parsing multipart form",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}
	jobID, err := h.services.Repositories.Jobs.Create(ctxb, session.ID())
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

	var fileNameToBuffer = map[string][]byte{}

	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			fileName := fileHeader.Filename
			buffer, err := fileHeaderToBytes(fileHeader)
			if err != nil {
				h.services.Logger.PrintError(err, types.AnyMap{
					"message": "error reading file",
				})
				return ctx.SendStatus(http.StatusInternalServerError)
			}

			if err := h.validateFileType(buffer); err != nil {
				h.services.Logger.PrintError(err, types.AnyMap{
					"message": "file type is not supported",
				})
				return ctx.SendStatus(http.StatusBadRequest)
			}

			fileNameToBuffer[fileName] = buffer
			err = ctx.SaveFile(fileHeader, helpers.BuildPath(h.services.Config.Process.UploadDir, jobID, fileName))
			if err != nil {
				h.services.Logger.PrintError(err, types.AnyMap{
					"message": "error saving file",
				})
				return ctx.SendStatus(http.StatusInternalServerError)
			}
		}
	}

	dbFiles, err := h.services.Repositories.Files.CreateBulk(ctxb, jobID, helpers.GetMapKeys(fileNameToBuffer))
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

	for fileName, buffer := range fileNameToBuffer {
		fileID := fileNameToID[fileName]
		go h.services.Processor.Process(ctxb, types.ImageProcessInput{JobID: jobID, FileID: fileID, FileName: fileName, Format: reqFormat, Quality: quality, Buffer: bytes.NewReader(buffer)})
	}

	return ctx.Status(http.StatusAccepted).JSON(fiber.Map{
		"job_id": jobID,
	})
}

func (h *Handlers) HandleProcessFile(ctx *fiber.Ctx) error {
	ctxb := context.Background()

	session, err := h.getSession(ctx)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	reqJobID := ctx.Params("jobID")
	if reqJobID == "" {
		h.services.Logger.PrintError(JobIDIsNotFound)
		return ctx.SendStatus(http.StatusBadRequest)
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

	file, err := h.services.Repositories.Files.GetWithJobByID(ctxb, reqFileIDInt)
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

	go h.services.Processor.Process(ctxb, types.ImageProcessInput{JobID: jobID, FileID: fileID, FileName: file.Name, Format: format, Quality: quality, Width: reqFileWidth, Height: reqFileHeight, Buffer: bytes.NewReader(buffer)})

	return nil
}

func (h *Handlers) HandleAddFileToJob(ctx *fiber.Ctx) error {
	ctxb := context.Background()

	reqJobID := ctx.Params("jobID")
	if reqJobID == "" {
		h.services.Logger.PrintError(JobIDIsNotFound)
		return ctx.SendStatus(http.StatusBadRequest)
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

	fileHeader := form.File["image"][0]
	fileName := fileHeader.Filename
	buffer, err := fileHeaderToBytes(fileHeader)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error reading file",
		})
	}

	if err := h.validateFileType(buffer); err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "file type is not supported",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := h.verifyJobSession(ctx, jobID); err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error verifying job session",
		})
		return ctx.SendStatus(http.StatusForbidden)
	}

	dbFiles, err := h.services.Repositories.Files.GetByJobID(ctxb, jobID)
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

	if _, ok := fileNameToID[fileName]; ok {
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

	dbFile, err := h.services.Repositories.Files.AddToJob(ctxb, jobID, fileName)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error adding file to job",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	go h.services.Processor.Process(ctxb, types.ImageProcessInput{JobID: jobID, FileID: dbFile.ID, FileName: fileName, Format: reqFormat, Quality: quality, Buffer: bytes.NewReader(buffer)})

	return nil
}

func (h *Handlers) HandleDeleteFileFromJob(ctx *fiber.Ctx) error {
	ctxb := context.Background()

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

	err = h.verifyJobSession(ctx, jobID)
	if err != nil {
		return ctx.SendStatus(http.StatusForbidden)
	}

	err = h.services.Repositories.Files.DeleteFromJob(ctxb, jobID, fileID)
	if err != nil {
		h.services.Logger.PrintError(err, types.AnyMap{
			"message": "error deleting file",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return nil
}
