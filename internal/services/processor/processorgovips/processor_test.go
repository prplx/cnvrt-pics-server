package processorgovips_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/prplx/cnvrt/internal/helpers"
	"github.com/prplx/cnvrt/internal/mocks"
	processor "github.com/prplx/cnvrt/internal/services/processor/processorgovips"

	"math/rand"

	"bytes"
	"image"
	"image/png"

	"github.com/prplx/cnvrt/internal/config"
	"github.com/prplx/cnvrt/internal/models"
	"github.com/prplx/cnvrt/internal/types"
	"go.uber.org/mock/gomock"
)

func TestProcessor_Process_should_communicate_about_processing_error_when_getting_operations_by_params_returns_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocks.NewMockLogger(ctrl)
	comm := mocks.NewMockCommunicator(ctrl)
	operationsRepo := mocks.NewMockOperations(ctrl)
	scheduler := mocks.NewMockScheduler(ctrl)
	processor := processor.NewProcessor(config.TestConfig(), operationsRepo, comm, logger, scheduler)
	ctx := context.Background()
	input := processInput()

	logger.EXPECT().PrintError(gomock.Any(), gomock.Any()).AnyTimes()
	firstCall := comm.EXPECT().SendStartProcessing(input.JobID, input.FileID, input.FileName).Return(nil).Times(1)
	secondCall := operationsRepo.EXPECT().GetByParams(ctx, gomock.Any()).Return(nil, errors.New("Network error")).Times(1).After(firstCall)
	comm.EXPECT().SendErrorProcessing(input.JobID, input.FileID, input.FileName).Return(nil).Times(1).After(secondCall)

	processor.Process(ctx, input)
}

func TestProcessor_Process_should_communicate_about_processing_success_if_operation_and_file_already_exists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocks.NewMockLogger(ctrl)
	comm := mocks.NewMockCommunicator(ctrl)
	operationsRepo := mocks.NewMockOperations(ctrl)
	scheduler := mocks.NewMockScheduler(ctrl)
	processor := processor.NewProcessor(config.TestConfig(), operationsRepo, comm, logger, scheduler)
	ctx := context.Background()
	input := processInput()
	resultFileName := uuid.NewString() + "." + "webp"

	firstCall := comm.EXPECT().SendStartProcessing(input.JobID, input.FileID, input.FileName).Return(nil).Times(1)
	secondCall := operationsRepo.EXPECT().GetByParams(ctx, gomock.Any()).Return(&models.Operation{
		FileName: resultFileName,
	}, nil).Times(1).After(firstCall)
	thirdCall := operationsRepo.EXPECT().Create(ctx, gomock.Any()).Return(jobID(), nil).Times(1).After(secondCall)
	comm.EXPECT().SendSuccessProcessing(input.JobID, gomock.Any()).Return(nil).Times(1).After(thirdCall)

	filePath := helpers.BuildPath(config.TestConfig().Process.UploadDir, input.JobID, input.FileName)
	resultFilePath := helpers.BuildPath(config.TestConfig().Process.UploadDir, input.JobID, resultFileName)
	os.MkdirAll(helpers.BuildPath(config.TestConfig().Process.UploadDir, input.JobID), os.ModePerm)
	os.Create(resultFilePath)
	os.Create(filePath)

	processor.Process(ctx, input)

	os.RemoveAll(helpers.BuildPath(config.TestConfig().Process.UploadDir))
}

func TestProcessor_Process_should_convert_file_if_operation_exists_but_the_file_locally_doesnt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocks.NewMockLogger(ctrl)
	comm := mocks.NewMockCommunicator(ctrl)
	operationsRepo := mocks.NewMockOperations(ctrl)
	scheduler := mocks.NewMockScheduler(ctrl)
	processor := processor.NewProcessor(config.TestConfig(), operationsRepo, comm, logger, scheduler)
	ctx := context.Background()
	input := processInput()
	resultFileName := uuid.NewString() + "." + "webp"
	input.Buffer = getImageBuffer()

	firstCall := comm.EXPECT().SendStartProcessing(input.JobID, input.FileID, input.FileName).Return(nil).Times(1)
	secondCall := operationsRepo.EXPECT().GetByParams(ctx, gomock.Any()).Return(&models.Operation{
		FileName: resultFileName,
	}, nil).Times(1).After(firstCall)
	thirdCall := scheduler.EXPECT().ScheduleFlush(input.JobID, gomock.Any()).Return(nil).Times(1)
	fourthCall := operationsRepo.EXPECT().Create(ctx, gomock.Any()).Return(jobID(), nil).Times(1).After(secondCall).After(thirdCall)
	comm.EXPECT().SendSuccessProcessing(input.JobID, gomock.Any()).Return(nil).Times(1).After(fourthCall)

	os.MkdirAll(helpers.BuildPath(config.TestConfig().Process.UploadDir, input.JobID), os.ModePerm)
	filePath := helpers.BuildPath(config.TestConfig().Process.UploadDir, input.JobID, input.FileName)
	os.Create(filePath)

	processor.Process(ctx, input)

	os.RemoveAll(helpers.BuildPath(config.TestConfig().Process.UploadDir))
}

func TestProcessor_Process_should_convert_file_if_neither_operation_nor_file_locally_exist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocks.NewMockLogger(ctrl)
	comm := mocks.NewMockCommunicator(ctrl)
	operationsRepo := mocks.NewMockOperations(ctrl)
	scheduler := mocks.NewMockScheduler(ctrl)
	processor := processor.NewProcessor(config.TestConfig(), operationsRepo, comm, logger, scheduler)
	ctx := context.Background()
	input := processInput()

	myImage := image.NewRGBA(image.Rect(0, 0, 100, 200))
	var buff bytes.Buffer
	png.Encode(&buff, myImage)
	input.Buffer = bytes.NewReader(buff.Bytes())

	firstCall := comm.EXPECT().SendStartProcessing(input.JobID, input.FileID, input.FileName).Return(nil).Times(1)
	secondCall := operationsRepo.EXPECT().GetByParams(ctx, gomock.Any()).Return(nil, nil).Times(1).After(firstCall)
	thirdCall := scheduler.EXPECT().ScheduleFlush(input.JobID, gomock.Any()).Return(nil).Times(1)
	fourthCall := operationsRepo.EXPECT().Create(ctx, gomock.Any()).Return(jobID(), nil).Times(1).After(secondCall).After(thirdCall)
	comm.EXPECT().SendSuccessProcessing(input.JobID, gomock.Any()).Return(nil).Times(1).After(fourthCall)

	os.MkdirAll(helpers.BuildPath(config.TestConfig().Process.UploadDir, input.JobID), os.ModePerm)
	filePath := helpers.BuildPath(config.TestConfig().Process.UploadDir, input.JobID, input.FileName)
	os.Create(filePath)

	processor.Process(ctx, input)

	os.RemoveAll(helpers.BuildPath(config.TestConfig().Process.UploadDir))
}

func jobID() int64 {
	return int64(rand.Intn(2e8))
}

func processInput() types.ImageProcessInput {
	return types.ImageProcessInput{
		JobID:    jobID(),
		FileID:   jobID(),
		FileName: "test.jpg",
		Format:   "webp",
		Width:    100,
		Height:   100,
		Quality:  100,
		Buffer:   getImageBuffer(),
	}
}

func getImageBuffer() io.Reader {
	myImage := image.NewRGBA(image.Rect(0, 0, 100, 200))
	var buff bytes.Buffer
	png.Encode(&buff, myImage)
	return bytes.NewReader(buff.Bytes())
}
