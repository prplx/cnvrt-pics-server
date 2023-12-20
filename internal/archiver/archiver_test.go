package archiver_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/prplx/cnvrt/internal/mocks"

	"math/rand"

	"github.com/prplx/cnvrt/internal/archiver"
	"github.com/prplx/cnvrt/internal/config"
	"github.com/prplx/cnvrt/internal/helpers"
	"github.com/prplx/cnvrt/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestArchiver_Archive__should_report_and_return_when_communicator_returns_error_while_sending_start_archiving(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := jobID()
	logger := mocks.NewMockLogger(ctrl)
	comm := mocks.NewMockCommunicator(ctrl)
	filesRepo := mocks.NewMockFiles(ctrl)
	archiver := archiver.NewArchiver(config.TestConfig(), filesRepo, logger, comm)
	logger.EXPECT().PrintError(gomock.Any(), gomock.Any()).AnyTimes()

	comm.EXPECT().SendStartArchiving(jobID).Return(errors.New("Communication problem")).Times(1)
	comm.EXPECT().SendErrorArchiving(jobID).Return(nil).Times(1)

	err := archiver.Archive(jobID)

	if assert.Error(t, err) {
		assert.Equal(t, "error sending start archiving: Communication problem", err.Error())
	}
}

func TestArchiver_Archive__should_report_and_return_when_files_repository_returns_error_while_getting_files_with_latest_operations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := jobID()
	logger := mocks.NewMockLogger(ctrl)
	comm := mocks.NewMockCommunicator(ctrl)
	filesRepo := mocks.NewMockFiles(ctrl)
	archiver := archiver.NewArchiver(config.TestConfig(), filesRepo, logger, comm)
	logger.EXPECT().PrintError(gomock.Any(), gomock.Any()).AnyTimes()

	comm.EXPECT().SendStartArchiving(jobID).Return(nil).Times(1)
	filesRepo.EXPECT().GetWithLatestOperationsByJobID(jobID).Return(nil, errors.New("Database problem")).Times(1)
	comm.EXPECT().SendErrorArchiving(jobID).Return(nil).Times(1)

	err := archiver.Archive(jobID)

	if assert.Error(t, err) {
		assert.Equal(t, "error getting files with latest operations: Database problem", err.Error())
	}
}

func TestArchiver_Archive__should_successfully_zip_files_and_communicate_when_conditions_are_met(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := jobID()
	logger := mocks.NewMockLogger(ctrl)
	comm := mocks.NewMockCommunicator(ctrl)
	filesRepo := mocks.NewMockFiles(ctrl)
	archiver := archiver.NewArchiver(config.TestConfig(), filesRepo, logger, comm)

	logger.EXPECT().PrintError(gomock.Any(), gomock.Any()).AnyTimes()
	comm.EXPECT().SendStartArchiving(jobID).Return(nil).Times(1)
	filesRepo.EXPECT().GetWithLatestOperationsByJobID(jobID).Return([]*models.File{}, nil).Times(1)
	comm.EXPECT().SendSuccessArchiving(jobID, helpers.BuildPath(config.TestConfig().Process.UploadDir, jobID, fmt.Sprintf("%s.zip", config.TestConfig().App.Name))).Return(nil).Times(1)

	os.MkdirAll(config.TestConfig().Process.UploadDir+"/"+fmt.Sprint(jobID), os.ModePerm)

	err := archiver.Archive(jobID)

	assert.Nil(t, err)

	os.RemoveAll(config.TestConfig().Process.UploadDir)
}

func jobID() int {
	return rand.Intn(2e8)
}
