package archiver

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/prplx/lighter.pics/internal/helpers"
	"github.com/prplx/lighter.pics/internal/repositories"
	"github.com/prplx/lighter.pics/internal/services"
	"github.com/prplx/lighter.pics/internal/types"
)

type Archiver struct {
	config       *types.Config
	repositories *repositories.Repositories
	logger       services.Logger
	communicator services.Communicator
}

func NewArchiver(config *types.Config, r *repositories.Repositories, l services.Logger, c services.Communicator) *Archiver {
	return &Archiver{
		config:       config,
		repositories: r,
		logger:       l,
		communicator: c,
	}
}

func (a *Archiver) Archive(jobID int) error {
	reportError := func(err error) {
		a.communicator.SendErrorArchiving(jobID)
		a.logger.PrintError(err, types.AnyMap{
			"message": "error while archiving files",
		})
	}

	err := a.communicator.SendStartArchiving(jobID)
	if err != nil {
		reportError(err)
		return errors.Wrap(err, "error sending start archiving")
	}

	filesWithOperaton, err := a.repositories.Files.GetWithLatestOperationsByJobID(jobID)
	if err != nil {
		reportError(err)
		return errors.Wrap(err, "error getting files with latest operations")
	}

	files := map[string]string{}
	for _, file := range filesWithOperaton {
		files[file.Name] = file.LatestOperation.FileName
	}

	archiveName := fmt.Sprintf("%s.zip", a.config.App.Name)

	err = zipFiles(archiveName, helpers.BuildPath(a.config.Process.UploadDir, jobID), files)
	if err != nil {
		reportError(err)
		return errors.Wrap(err, "error zipping files")
	}

	downloadPath := helpers.BuildPath(a.config.Process.UploadDir, jobID, archiveName)
	err = a.communicator.SendSuccessArchiving(jobID, downloadPath)
	if err != nil {
		reportError(err)
		return errors.Wrap(err, "error sending success archiving")
	}

	return nil
}

func zipFiles(zipFile string, dir string, files map[string]string) error {
	newZipFile, err := os.Create(helpers.BuildPath(dir, zipFile))
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for srcFile, dstFile := range files {
		filePath := filepath.Join(dir, dstFile)
		fileToZip, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		info, err := fileToZip.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = fmt.Sprintf("%s%s", helpers.FileNameWithoutExtension(srcFile), helpers.FileExtension(dstFile))
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, fileToZip)
		if err != nil {
			return err
		}
	}

	err = zipWriter.Close()
	if err != nil {
		return err
	}

	return nil
}
