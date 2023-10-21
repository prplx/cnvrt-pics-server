package processor

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/prplx/lighter.pics/internal/helpers"
	"github.com/prplx/lighter.pics/internal/imageProcessor"
	"github.com/prplx/lighter.pics/internal/models"
	"github.com/prplx/lighter.pics/internal/repositories"
	"github.com/prplx/lighter.pics/internal/services"
	"github.com/prplx/lighter.pics/internal/types"
	"github.com/prplx/lighter.pics/internal/validator"
)

const (
	UploadDir = "./uploads"
	format    = "format"
	quality   = "quality"
	fileName  = "file_name"
	fileID    = "file_id"
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

type processInput struct {
	jobID    int
	fileID   int
	fileName string
	format   string
	quality  int
	width    int
	height   int
	buffer   []byte
}

func (p *Processor) HandleProcessJob(ctx *fiber.Ctx) error {
	v := validator.New()
	if validateRequestQueryParams(v, ctx, format, quality); !v.Valid() {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": v.Errors,
		})
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		p.logger.PrintError(err, types.AnyMap{
			"message": "error parsing multipart form",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	jobID, err := p.repositories.Jobs.Create(context.Background())
	if err != nil {
		p.logger.PrintError(err, types.AnyMap{
			"message": "error creating job",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	path := fmt.Sprintf(UploadDir+"/%d", jobID)
	reqFormat := ctx.Query(format)
	reqQuality := ctx.Query(quality)
	quality, err := strconv.Atoi(reqQuality)
	if err != nil {
		p.logger.PrintError(err, types.AnyMap{
			"message": "error parsing quality param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			p.logger.PrintError(err, types.AnyMap{
				"message": "error creating directory",
			})
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	fileNames := []string{}
	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			err = ctx.SaveFile(fileHeader, helpers.BuildPath(UploadDir, jobID, fileHeader.Filename))
			if err != nil {
				p.logger.PrintError(err, types.AnyMap{
					"message": "error saving file",
				})
				return ctx.SendStatus(http.StatusInternalServerError)
			}
			fileNames = append(fileNames, fileHeader.Filename)
		}
	}

	fileIds, err := p.repositories.Files.CreateBulk(context.Background(), jobID, fileNames)
	if err != nil {
		p.logger.PrintError(err, types.AnyMap{
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
		filePath := helpers.BuildPath(UploadDir, jobID, file.Name())
		buffer, err := p.readFile(filePath)
		if err != nil {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		buffers = append(buffers, buffer)
	}

	if len(buffers) == len(files) {
		for idx, buffer := range buffers {
			go p.process(context.Background(), processInput{jobID: jobID, fileID: fileIds[idx], fileName: files[idx].Name(), format: reqFormat, quality: quality, buffer: buffer})
		}
	}

	return ctx.Status(http.StatusAccepted).JSON(fiber.Map{
		"job_id": jobID,
	})
}

func (p *Processor) HandleProcessFile(ctx *fiber.Ctx) error {
	reqJobID := ctx.Params("jobID")
	if reqJobID == "" {
		p.logger.PrintError(JobIDIsNotFound)
		return ctx.SendStatus(http.StatusBadRequest)
	}

	v := validator.New()
	if validateRequestQueryParams(v, ctx, format, quality, fileID, "width", "height"); !v.Valid() {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": v.Errors,
		})
	}

	format := ctx.Query(format)
	reqQuality := ctx.Query(quality)
	reqFileID := ctx.Query(fileID)
	reqFileIdInt, err := strconv.Atoi(reqFileID)
	if err != nil {
		p.logger.PrintError(err, types.AnyMap{
			"message": "error parsing file_id param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	reqFileWidth, err := strconv.Atoi(ctx.Query("width"))
	if err != nil {
		p.logger.PrintError(err, types.AnyMap{
			"message": "error parsing width param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	reqFileHeight, err := strconv.Atoi(ctx.Query("height"))
	if err != nil {
		p.logger.PrintError(err, types.AnyMap{
			"message": "error parsing height param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	file, err := p.repositories.Files.GetById(context.Background(), reqFileIdInt)
	if err != nil {
		p.logger.PrintError(err, types.AnyMap{
			"message": "error getting file by id",
		})
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	filePath := helpers.BuildPath(UploadDir, reqJobID, file.Name)
	buffer, err := p.readFile(filePath)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	quality, err := strconv.Atoi(reqQuality)
	if err != nil {
		p.logger.PrintError(err, types.AnyMap{
			"message": "error parsing quality param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	fileID, err := strconv.Atoi(reqFileID)
	if err != nil {
		p.logger.PrintError(err, types.AnyMap{
			"message": "error parsing file_id param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	jobID, err := strconv.Atoi(reqJobID)
	if err != nil {
		p.logger.PrintError(err, types.AnyMap{
			"message": "error parsing job_id param",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	go p.process(context.Background(), processInput{jobID: jobID, fileID: fileID, fileName: file.Name, format: format, quality: quality, width: reqFileWidth, height: reqFileHeight, buffer: buffer})

	return nil
}

func (p *Processor) process(ctx context.Context, input processInput) {
	jobID := input.jobID
	fileID := input.fileID
	fileName := input.fileName
	format := input.format
	width := input.width
	height := input.height
	quality := input.quality
	buffer := input.buffer
	var resultFileName string

	p.communicator.SendStartProcessing(jobID, fileID, fileName)
	reportError := func(err error) {
		p.communicator.SendErrorProcessing(jobID, fileID, fileName)
		p.logger.PrintError(err, types.AnyMap{
			"job_id": jobID,
			"file":   fileName,
		})
	}

	possiblyExistingOperation, err := p.repositories.Operations.GetByParams(ctx, models.Operation{
		JobID:   jobID,
		FileID:  fileID,
		Format:  format,
		Quality: quality,
		Width:   width,
		Height:  height,
	})

	if err != nil {
		reportError(errors.Wrap(err, "error getting operation by params"))
		return
	}

	if possiblyExistingOperation != nil {
		resultFileName = possiblyExistingOperation.FileName
	} else {
		if width != 0 && height != 0 {
			resized, err := p.imageProcessor.NewImage(buffer).Resize(width, height)
			if err != nil {
				reportError(err)
				return
			}

			buffer = resized
		} else {
			dimensions := p.imageProcessor.NewImage(buffer).Size()
			width = dimensions.Width
			height = dimensions.Height
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

		resultFileName = uuid.NewString() + "." + format
		writerError := p.imageProcessor.Write(helpers.BuildPath(UploadDir, jobID, resultFileName), processed)
		if writerError != nil {
			reportError(writerError)
			return
		}
	}

	sourceInfo, err := os.Stat(helpers.BuildPath(UploadDir, jobID, fileName))
	if err != nil {
		reportError(err)
		return
	}

	targetInfo, err := os.Stat(helpers.BuildPath(UploadDir, jobID, resultFileName))
	if err != nil {
		reportError(errors.Wrap(err, "error getting target file info"))
		return
	}

	operation := models.Operation{JobID: jobID, FileID: fileID, Format: format, Quality: quality, Width: width, Height: height, FileName: resultFileName}

	err = p.repositories.Operations.UnsetLatest(ctx)
	if err != nil {
		reportError(errors.Wrap(err, "error unsetting latest operation"))
		return
	}
	_, err = p.repositories.Operations.Create(ctx, operation)
	if err != nil {
		reportError(errors.Wrap(err, "error creating operation"))
		return
	}

	err = p.communicator.SendSuccessProcessing(jobID, types.SuccessResult{
		SourceFileName: fileName,
		SourceFileID:   fileID,
		TargetFileName: resultFileName,
		SourceFileSize: sourceInfo.Size(),
		TargetFileSize: targetInfo.Size(),
		Width:          width,
		Height:         height,
	})
	if err != nil {
		p.logger.PrintError(err, types.AnyMap{
			"job_id":    jobID,
			"file_name": fileName,
			"file_id":   fileID,
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
		p.logger.PrintError(err, types.AnyMap{
			"message": "error opening file",
			"path":    path,
		})
		return nil, OpenFileError
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		p.logger.PrintError(err, types.AnyMap{
			"message": "error reading file",
			"path":    path,
		})
		return nil, ReadingFileError
	}

	return buffer, nil
}
