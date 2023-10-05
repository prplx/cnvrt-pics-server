package processor

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/prplx/lighter.pics/internal/imageProcessor"
	"github.com/prplx/lighter.pics/internal/repositories"
	"github.com/prplx/lighter.pics/internal/services"
	"github.com/prplx/lighter.pics/internal/types"
	"github.com/prplx/lighter.pics/internal/validator"
)

const (
	UploadDir = "./uploads"
)

type Processor struct {
	communicator   services.Communicator
	logger         services.Logger
	repositories   repositories.Repositories
	imageProcessor services.ImageProcessor
}

func NewProcessor(services *services.Services) *Processor {
	return &Processor{
		communicator:   services.Communicator,
		logger:         services.Logger,
		repositories:   services.Repositories,
		imageProcessor: services.ImageProcessor,
	}
}

func (p *Processor) HandleProcessJob(ctx *fiber.Ctx) error {
	v := validator.New()
	if validateRequestQueryParams(v, ctx, "format", "quality"); !v.Valid() {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": v.Errors,
		})
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		p.logger.PrintError(err, map[string]string{
			"message": "error parsing multipart form",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	jobID, err := p.repositories.Jobs.Create()
	if err != nil {
		p.logger.PrintError(err, map[string]string{
			"message": "error creating job",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	path := fmt.Sprintf(UploadDir+"/%s", jobID)
	reqFormat := ctx.Query("format")
	reqQuality := ctx.Query("quality")
	quality, err := strconv.Atoi(reqQuality)
	if err != nil {
		p.logger.PrintError(err, map[string]string{
			"message": "error parsing quality param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			p.logger.PrintError(err, map[string]string{
				"message": "error creating directory",
			})
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	fileNames := []string{}
	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			err = ctx.SaveFile(fileHeader, fmt.Sprintf(UploadDir+"/%s/%s", jobID, fileHeader.Filename))
			if err != nil {
				p.logger.PrintError(err, map[string]string{
					"message": "error saving file",
				})
				return ctx.SendStatus(http.StatusInternalServerError)
			}
			fileNames = append(fileNames, fileHeader.Filename)
		}
	}

	if err := p.repositories.Files.CreateBulk(jobID, fileNames); err != nil {
		p.logger.PrintError(err, map[string]string{
			"message": "error creating file records",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	buffers := [][]byte{}
	for _, file := range files {
		filePath := fmt.Sprintf(UploadDir+"/%s/%s", jobID, file.Name())
		buffer, err := p.readFile(filePath)
		if err != nil {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		buffers = append(buffers, buffer)
	}

	if len(buffers) == len(files) {
		for idx, buffer := range buffers {
			go p.process(jobID, files[idx].Name(), reqFormat, quality, buffer)
		}
	}

	return ctx.Status(http.StatusAccepted).JSON(fiber.Map{
		"job_id": jobID,
	})
}

func (p *Processor) HandleProcessFile(ctx *fiber.Ctx) error {
	jobID := ctx.Params("jobID")
	if jobID == "" {
		p.logger.PrintError(errors.New("jobID param does not exist for the existing job"), map[string]string{})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	v := validator.New()
	if validateRequestQueryParams(v, ctx, "format", "quality", "file_name"); !v.Valid() {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": v.Errors,
		})
	}

	reqFormat := ctx.Query("format")
	reqQuality := ctx.Query("quality")
	reqFileName := ctx.Query("file_name")
	filePath := fmt.Sprintf(UploadDir+"/%s/%s", jobID, reqFileName)
	buffer, err := p.readFile(filePath)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	quality, err := strconv.Atoi(reqQuality)
	if err != nil {
		p.logger.PrintError(err, map[string]string{
			"message": "error parsing quality param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	go p.process(jobID, reqFileName, reqFormat, quality, buffer)

	return nil
}

func (p *Processor) process(jobID, fileName, format string, quality int, buffer []byte) {
	p.communicator.SendStartProcessing(jobID, fileName)
	reportError := func(err error) {
		p.communicator.SendErrorProcessing(jobID, fileName)
		p.logger.PrintError(err, map[string]string{
			"job_id": jobID,
			"file":   fileName,
		})
	}

	file, err := p.repositories.GetByJobIDAndName(jobID, fileName)
	if err != nil {
		reportError(err)
		return
	}

	converted, err := p.imageProcessor.NewImage(buffer).Convert(imageProcessor.Formats[format])
	if err != nil {
		reportError(err)
		return
	}

	processed, err := p.imageProcessor.NewImage(converted).Process(imageProcessor.Options{Quality: quality})
	if err != nil {
		reportError(err)
		return
	}

	resultFileName := uuid.NewString() + "." + format
	writerError := p.imageProcessor.Write(UploadDir+"/"+jobID+"/"+resultFileName, processed)
	if writerError != nil {
		reportError(writerError)
		return
	}

	sourceInfo, err := os.Stat(UploadDir + "/" + jobID + "/" + fileName)
	if err != nil {
		reportError(err)
		return
	}

	targetInfo, err := os.Stat(UploadDir + "/" + jobID + "/" + resultFileName)
	if err != nil {
		reportError(err)
		return
	}

	_, err = p.repositories.Operations.Create(jobID, file.ID, format, quality, resultFileName, 0, 0)
	if err != nil {
		reportError(err)
		return
	}

	err = p.communicator.SendSuccessProcessing(jobID, types.SuccessResult{
		SourceFileName: fileName,
		TargetFileName: resultFileName,
		SourceFileSize: sourceInfo.Size(),
		TargetFileSize: targetInfo.Size(),
	})
	if err != nil {
		p.logger.PrintError(err, map[string]string{
			"job_id": jobID,
			"file":   fileName,
		})
	}
}

func validateRequestQueryParams(v *validator.Validator, ctx *fiber.Ctx, requiredParams ...string) {
	for _, param := range requiredParams {
		v.Check(ctx.Query(param) != "", param, param+" is required")
	}
}

func (p Processor) readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		p.logger.PrintError(err, map[string]string{
			"message": "error opening file",
			"path":    path,
		})
		return nil, OpenFileError
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		p.logger.PrintError(err, map[string]string{
			"message": "error reading file",
			"path":    path,
		})
		return nil, ReadingFileError
	}

	return buffer, nil
}
