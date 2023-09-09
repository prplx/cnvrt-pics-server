package processor

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/h2non/bimg"
	"github.com/pkg/errors"
	"github.com/prplx/lighter.pics/internal/services"
)

type Processor struct {
	communicator services.Communicator
	logger       services.Logger
}

func NewProcessor(services *services.Services) *Processor {
	return &Processor{
		communicator: services.Communicator,
		logger:       services.Logger,
	}
}

func (p *Processor) Handle(ctx *fiber.Ctx) error {
	// Here will go all validation logic
	form, err := ctx.MultipartForm()
	if err != nil {
		p.logger.PrintError(err, map[string]string{
			"message": "error parsing multipart form",
		})
		return ctx.SendStatus(http.StatusBadRequest)
	}

	jobID := uuid.New().String()
	path := fmt.Sprintf("./uploads/%s", jobID)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			p.logger.PrintError(err, map[string]string{
				"message": "error creating directory",
			})
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			err = ctx.SaveFile(fileHeader, fmt.Sprintf("./uploads/%s/%s", jobID, fileHeader.Filename))
			if err != nil {
				p.logger.PrintError(err, map[string]string{
					"message": "error saving file",
				})
				return ctx.SendStatus(http.StatusInternalServerError)
			}
		}
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	buffers := [][]byte{}
	for _, file := range files {
		filePath := fmt.Sprintf("./uploads/%s/%s", jobID, file.Name())
		file, err := os.Open(filePath)
		if err != nil {
			p.logger.PrintError(err, map[string]string{
				"message": "error opening file",
			})
			return ctx.SendStatus(http.StatusInternalServerError)
		}
		defer file.Close()

		buffer, err := io.ReadAll(file)
		if err != nil {
			p.logger.PrintError(err, map[string]string{
				"message": "error reading file",
			})
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		buffers = append(buffers, buffer)
	}

	if len(buffers) == len(files) {
		for idx, buffer := range buffers {
			go p.process(jobID, files[idx].Name(), buffer)
		}
	}

	return ctx.Status(http.StatusAccepted).JSON(fiber.Map{
		"job_id": jobID,
	})
}

func (p *Processor) process(jobID, fileName string, buffer []byte) {
	p.communicator.SendStartProcessing(jobID, fileName)
	reportError := func(err error) {
		p.communicator.SendErrorProcessing(jobID, fileName)
		p.logger.PrintError(err, map[string]string{
			"job_id": jobID,
			"file":   fileName,
		})
	}

	converted, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	if err != nil {
		reportError(err)
		return
	}

	processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: 80})
	if err != nil {
		reportError(err)
		return
	}

	writerError := bimg.Write("./uploads/"+jobID+"/"+fileName+".webp", processed)
	if writerError != nil {
		reportError(writerError)
	}
}
