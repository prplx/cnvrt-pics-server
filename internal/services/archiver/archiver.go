package archiver

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/prplx/cnvrt/internal/helpers"
	"github.com/prplx/cnvrt/internal/repositories"
	"github.com/prplx/cnvrt/internal/services"
	"github.com/prplx/cnvrt/internal/types"
)

type Archiver struct {
	config          *types.Config
	filesRepository repositories.Files
	logger          services.Logger
	communicator    services.Communicator
}

func NewArchiver(config *types.Config, fr repositories.Files, l services.Logger, c services.Communicator) *Archiver {
	return &Archiver{
		config:          config,
		filesRepository: fr,
		logger:          l,
		communicator:    c,
	}
}

func (a *Archiver) Archive(jobID int64) error {
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

	filesWithOperaton, err := a.filesRepository.GetWithLatestOperationsByJobID(context.Background(), jobID)
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

	downloadPath := helpers.BuildPath("/uploads", jobID, archiveName)
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
		fmt.Println(err)
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
