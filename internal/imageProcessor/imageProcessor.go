package imageProcessor

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/h2non/bimg"
	"github.com/pkg/errors"
	"github.com/prplx/lighter.pics/internal/helpers"
	"github.com/prplx/lighter.pics/internal/models"
	"github.com/prplx/lighter.pics/internal/repositories"
	"github.com/prplx/lighter.pics/internal/services"
	"github.com/prplx/lighter.pics/internal/types"
)

type ImageType int

const (
	// UNKNOWN represents an unknow image type value.
	UNKNOWN ImageType = iota
	// JPEG represents the JPEG image type.
	JPEG
	// WEBP represents the WEBP image type.
	WEBP
	// PNG represents the PNG image type.
	PNG
	// TIFF represents the TIFF image type.
	TIFF
	// GIF represents the GIF image type.
	GIF
	// PDF represents the PDF type.
	PDF
	// SVG represents the SVG image type.
	SVG
	// MAGICK represents the libmagick compatible genetic image type.
	MAGICK
	// HEIF represents the HEIC/HEIF/HVEC image type
	HEIF
	// AVIF represents the AVIF image type.
	AVIF
)

var Formats = map[string]ImageType{
	"jpeg":   JPEG,
	"jpg":    JPEG,
	"webp":   WEBP,
	"png":    PNG,
	"tiff":   TIFF,
	"gif":    GIF,
	"pdf":    PDF,
	"svg":    SVG,
	"magick": MAGICK,
	"heif":   HEIF,
	"avif":   AVIF,
}

type Options struct {
	Height  int
	Width   int
	Quality int
}

type Dimensions struct {
	Height int
	Width  int
}

type ImageProcessor struct {
	communicator services.Communicator
	logger       services.Logger
	repositories *repositories.Repositories
	config       *types.Config
}

type Image struct {
	buffer []byte
}

func NewImageProcessor(config *types.Config, r *repositories.Repositories, c services.Communicator, l services.Logger) *ImageProcessor {
	return &ImageProcessor{
		communicator: c,
		logger:       l,
		repositories: r,
		config:       config,
	}
}

func (p *ImageProcessor) NewImage(buffer []byte) *Image {
	return &Image{buffer}
}

func (p *ImageProcessor) Write(path string, buf []byte) error {
	return bimg.Write(path, buf)
}

func (i Image) Process(options Options) ([]byte, error) {
	bimgOptions := bimg.Options{
		Width:   options.Width,
		Height:  options.Height,
		Quality: options.Quality,
	}
	return bimg.NewImage(i.buffer).Process(bimgOptions)
}

func (i Image) Convert(t ImageType) ([]byte, error) {
	return bimg.NewImage(i.buffer).Convert(bimg.ImageType(t))
}

func (i Image) Resize(width, height int) ([]byte, error) {
	return bimg.NewImage(i.buffer).Resize(width, height)
}

func (i Image) Size() Dimensions {
	size, _ := bimg.NewImage(i.buffer).Size()
	return Dimensions{size.Height, size.Width}
}

func (p *ImageProcessor) Process(ctx context.Context, input types.ImageProcessInput) {
	jobID := input.JobID
	fileID := input.FileID
	fileName := input.FileName
	format := input.Format
	width := input.Width
	height := input.Height
	quality := input.Quality
	buffer := input.Buffer
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
			resized, err := p.NewImage(buffer).Resize(width, height)
			if err != nil {
				reportError(err)
				return
			}

			buffer = resized
		} else {
			dimensions := p.NewImage(buffer).Size()
			width = dimensions.Width
			height = dimensions.Height
		}

		converted, err := p.NewImage(buffer).Convert(Formats[format])
		if err != nil {
			reportError(err)
			return
		}

		processed, err := p.NewImage(converted).Process(Options{Quality: quality})
		if err != nil {
			reportError(err)
			return
		}

		resultFileName = uuid.NewString() + "." + format
		writerError := p.Write(helpers.BuildPath(p.config.Process.UploadDir, jobID, resultFileName), processed)
		if writerError != nil {
			reportError(writerError)
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
