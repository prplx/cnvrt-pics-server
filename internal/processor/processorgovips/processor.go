package processorgovips

import (
	"context"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/prplx/lighter.pics/internal/helpers"
	"github.com/prplx/lighter.pics/internal/models"
	"github.com/prplx/lighter.pics/internal/repositories"
	"github.com/prplx/lighter.pics/internal/services"
	"github.com/prplx/lighter.pics/internal/types"
)

type Processor struct {
	communicator services.Communicator
	logger       services.Logger
	repositories *repositories.Repositories
	config       *types.Config
}

var Formats = map[string]vips.ImageType{
	"jpeg": vips.ImageTypeJPEG,
	"jpg":  vips.ImageTypeJPEG,
	"webp": vips.ImageTypeWEBP,
	"png":  vips.ImageTypePNG,
	"tiff": vips.ImageTypeTIFF,
	"gif":  vips.ImageTypeGIF,
	"heif": vips.ImageTypeHEIF,
	"avif": vips.ImageTypeAVIF,
}

func init() {
	vips.Startup(nil)
	// defer vips.Shutdown()
}

func NewProcessor(config *types.Config, r *repositories.Repositories, c services.Communicator, l services.Logger) *Processor {
	return &Processor{
		communicator: c,
		logger:       l,
		repositories: r,
		config:       config,
	}
}

func (p *Processor) Process(ctx context.Context, input types.ImageProcessInput) {
	jobID := input.JobID
	fileID := input.FileID
	fileName := input.FileName
	format := input.Format
	width := input.Width
	height := input.Height
	quality := input.Quality
	buffer := input.Buffer
	var resultFileName string

	reportError := func(err error) {
		p.communicator.SendErrorProcessing(jobID, fileID, fileName)
		p.logger.PrintError(err, types.AnyMap{
			"job_id": jobID,
			"file":   fileName,
		})
	}

	p.communicator.SendStartProcessing(jobID, fileID, fileName)

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
		image, err := vips.NewImageFromBuffer(buffer)
		if width != 0 && height != 0 {
			if err != nil {
				reportError(errors.Wrap(err, "error creating image from buffer"))
				return
			}

			if err := image.Resize(float64(width/height), vips.KernelLanczos3); err != nil {
				reportError(errors.Wrap(err, "error resizing image"))
				return
			}

		} else {
			width = image.Width()
			height = image.Height()
		}

		imageBytes, _, err := image.Export(&vips.ExportParams{
			Format:  Formats[format],
			Quality: quality,
		})
		if err != nil {
			reportError(errors.Wrap(err, "error exporting image"))
			return
		}

		resultFileName = uuid.NewString() + "." + format
		path := helpers.BuildPath(p.config.Process.UploadDir, jobID, resultFileName)
		err = os.WriteFile(path, imageBytes, 0644)
		if err != nil {
			reportError(errors.Wrap(err, "error writing image"))
			return
		}

	}

	sourceInfo, err := os.Stat(helpers.BuildPath(p.config.Process.UploadDir, jobID, fileName))
	if err != nil {
		reportError(err)
		return
	}

	targetInfo, err := os.Stat(helpers.BuildPath(p.config.Process.UploadDir, jobID, resultFileName))
	if err != nil {
		reportError(errors.Wrap(err, "error getting target file info"))
		return
	}

	operation := models.Operation{JobID: jobID, FileID: fileID, Format: format, Quality: quality, Width: width, Height: height, FileName: resultFileName}

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