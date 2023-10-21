package imageProcessor

import (
	"github.com/h2non/bimg"
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

type ImageProcessor struct{}

type Image struct {
	buffer []byte
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func (p ImageProcessor) NewImage(buffer []byte) *Image {
	return &Image{buffer}
}

func (p ImageProcessor) Write(path string, buf []byte) error {
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
