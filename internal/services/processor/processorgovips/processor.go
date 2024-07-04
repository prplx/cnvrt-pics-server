package processorgovips

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/prplx/cnvrt/internal/helpers"
	"github.com/prplx/cnvrt/internal/models"
	"github.com/prplx/cnvrt/internal/repositories"
	"github.com/prplx/cnvrt/internal/services"
	"github.com/prplx/cnvrt/internal/types"
)

type Processor struct {
	communicator         services.Communicator
	logger               services.Logger
	operationsRepository repositories.Operations
	config               *types.Config
	scheduler            services.Scheduler
}

func NewProcessor(config *types.Config, or repositories.Operations, c services.Communicator, l services.Logger, s services.Scheduler) *Processor {
	return &Processor{
		communicator:         c,
		logger:               l,
		operationsRepository: or,
		config:               config,
		scheduler:            s,
	}
}

func (p *Processor) Startup() {
	vips.Startup(&vips.Config{
		ConcurrencyLevel: 8,
	})
}

func (p *Processor) Shutdown() {
	vips.Shutdown()
}

func (p *Processor) Process(ctx context.Context, input types.ImageProcessInput) {
	jobID := input.JobID
	fileID := input.FileID
	fileName := input.FileName
	format := input.Format
	width := input.Width
	height := input.Height
	quality := input.Quality
	buffer, err := io.ReadAll(input.Buffer)
	if err != nil {
		p.communicator.SendErrorProcessing(jobID, fileID, fileName)
		p.logger.PrintError(err, types.AnyMap{
			"job_id":  jobID,
			"file":    fileName,
			"message": "error reading buffer",
		})
		return
	}
	var resultFileName string
	var existingJobFileExists bool
	var originalWidth int
	var originalHeight int

	reportError := func(err error) {
		p.communicator.SendErrorProcessing(jobID, fileID, fileName)
		p.logger.PrintError(err, types.AnyMap{
			"job_id": jobID,
			"file":   fileName,
		})
	}

	p.communicator.SendStartProcessing(jobID, fileID, fileName)

	possiblyExistingOperation, err := p.operationsRepository.GetByParams(ctx, models.Operation{
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
		if _, err := os.Stat(helpers.BuildPath(p.config.Process.UploadDir, jobID, possiblyExistingOperation.FileName)); !os.IsNotExist(err) {
			existingJobFileExists = true
		}
	}

	image, err := vips.NewImageFromBuffer(buffer)
	if err != nil {
		reportError(errors.Wrap(err, "error creating image from buffer"))
		return
	}

	originalWidth = image.Width()
	originalHeight = image.Height()

	if existingJobFileExists {
		resultFileName = possiblyExistingOperation.FileName
		image.Close()
	} else {
		if width != 0 && height != 0 {
			if err := image.Resize(float64(width)/float64(originalWidth), vips.KernelLanczos3); err != nil {
				reportError(errors.Wrap(err, "error resizing image"))
				return
			}
		} else {
			width = image.Width()
			height = image.Height()
		}

		var imageBytes []byte

		switch format {
		case "jpeg", "jpg":
			exportParams := vips.NewDefaultJPEGExportParams()
			exportParams.Quality = quality
			exportParams.StripMetadata = true
			imageBytes, _, err = image.Export(exportParams)
		case "webp":
			exportParams := vips.NewDefaultWEBPExportParams()
			exportParams.Quality = quality
			exportParams.StripMetadata = true
			imageBytes, _, err = image.Export(exportParams)
		case "png":
			exportParams := vips.NewPngExportParams()
			exportParams.Compression = 8
			exportParams.Quality = quality
			exportParams.Palette = true
			exportParams.StripMetadata = true
			imageBytes, _, err = image.ExportPng(exportParams)
		case "avif":
			exportParams := vips.NewAvifExportParams()
			exportParams.Quality = quality
			exportParams.StripMetadata = true
			imageBytes, _, err = image.ExportAvif(exportParams)
		default:
			reportError(errors.New("unsupported format"))
			return
		}
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

		err = p.scheduler.ScheduleFlush(jobID, time.Duration(p.config.App.JobFlushTimeout)*time.Second)
		if err != nil {
			p.logger.PrintError(err, types.AnyMap{
				"job_id": jobID,
			})
		}
		image.Close()
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

	_, err = p.operationsRepository.Create(ctx, operation)
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
		Format:         format,
		Quality:        quality,
		OriginalWidth:  originalWidth,
		OriginalHeight: originalHeight,
	})
	if err != nil {
		p.logger.PrintError(err, types.AnyMap{
			"job_id":    jobID,
			"file_name": fileName,
			"file_id":   fileID,
			"operation": "send_success_processing",
		})
	}

}
