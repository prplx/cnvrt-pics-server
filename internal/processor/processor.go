package processor

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/prplx/lighter.pics/internal/types"
)

type Processor struct {
	communicator types.Communicator
}

func NewProcessor(communicator types.Communicator) *Processor {
	return &Processor{
		communicator: communicator,
	}
}

func (p *Processor) Handle(ctx *fiber.Ctx) error {
	// Here will go all validation logic
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	jobID := strings.ReplaceAll(string(uuid), "\n", "")
	path := fmt.Sprintf("./uploads/%s", jobID)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			err = ctx.SaveFile(fileHeader, fmt.Sprintf("./uploads/%s/%s", jobID, fileHeader.Filename))
			if err != nil {
				return ctx.SendStatus(http.StatusInternalServerError)
			}
			// TODO: Process only all of the files were saved successfully
			go p.process(jobID, fileHeader.Filename)
		}
	}

	return ctx.Status(http.StatusAccepted).JSON(fiber.Map{
		"job_id": jobID,
	})
}

func (p *Processor) process(jobID, fileName string) {
	p.communicator.SendStartProcess(jobID, fileName)
}
