package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/prplx/cnvrt/internal/mocks"
	"github.com/prplx/cnvrt/internal/models"
	"github.com/prplx/cnvrt/internal/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_HandleProcessJob__should_return_correct_response_when_all_conditions_are_met(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileName := "file.png"
	jobID := 555
	jobsRepo := mocks.NewMockJobs(ctrl)
	communicator := mocks.NewMockCommunicator(ctrl)
	logger := mocks.NewMockLogger(ctrl)
	filesRepo := mocks.NewMockFiles(ctrl)
	processor := mocks.NewMockProcessor(ctrl)
	mocks := &Mocks{
		jobsRepo:     jobsRepo,
		filesRepo:    filesRepo,
		communicator: communicator,
		logger:       logger,
		processor:    processor,
	}
	jobsRepo.EXPECT().Create(gomock.Any()).Return(jobID, nil)
	logger.EXPECT().PrintError(gomock.Any()).AnyTimes()
	filesRepo.EXPECT().CreateBulk(gomock.Any(), jobID, []string{fileName}).Return([]models.File{
		{
			ID:   1,
			Name: fileName,
		},
	}, nil)
	processor.EXPECT().Process(gomock.Any(), gomock.Any())
	body, contentType := createFormFile(t, "file", fileName)
	app, services := setup(t, mocks)

	r := httptest.NewRequest(http.MethodPost, processEndpoint+"?format=webp&quality=80", body)
	r.Header.Add("Content-Type", contentType)

	time.Sleep(1 * time.Second)

	resp, _ := app.Test(r, -1)
	got, _ := io.ReadAll(resp.Body)
	want := `{"job_id":555}`

	assert.Equal(t, want, string(got))
	assert.Equal(t, 202, resp.StatusCode)

	cleanUp(t, services)
}

func Test_HandleProcessJob__should_return_400_when_required_params_is_missing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mocks := &Mocks{}
	body, contentType := createFormFile(t, "file", "file.png")
	app, services := setup(t, mocks)

	testCases := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "format",
			url:  processEndpoint + "?quality=80",
			want: `{"errors":{"format":"format is required"}}`,
		},
		{
			name: "quality",
			url:  processEndpoint + "?format=webp",
			want: `{"errors":{"quality":"quality is required"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tc.url, body)
			r.Header.Add("Content-Type", contentType)
			resp, _ := app.Test(r, -1)
			got, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tc.want, string(got))
			assert.Equal(t, 400, resp.StatusCode)
		})
	}

	cleanUp(t, services)
}

func Test_HandleProcessFile__should_return_correct_response_when_all_conditions_are_met(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := 555
	fileID := 1
	fileName := "file.png"
	logger := mocks.NewMockLogger(ctrl)
	filesRepo := mocks.NewMockFiles(ctrl)
	processor := mocks.NewMockProcessor(ctrl)
	mocks := &Mocks{
		filesRepo: filesRepo,
		logger:    logger,
		processor: processor,
	}
	filesRepo.EXPECT().GetByID(gomock.Any(), fileID).Return(&models.File{
		Name: fileName,
	}, nil).Times(1)

	app, services := setup(t, mocks)
	fileDir := fmt.Sprintf("%s/%d", services.Config.Process.UploadDir, jobID)
	filePath := fmt.Sprintf("%s/%s", fileDir, fileName)
	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	buffer, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	processor.EXPECT().Process(gomock.Any(), gomock.Eq(types.ImageProcessInput{
		JobID:    jobID,
		FileID:   fileID,
		FileName: fileName,
		Format:   "webp",
		Quality:  80,
		Width:    100,
		Height:   100,
		Buffer:   buffer,
	}))

	url := fmt.Sprintf("%s/%d?format=webp&quality=80&file_id=%d&width=100&height=100", processEndpoint, jobID, fileID)
	r := httptest.NewRequest(http.MethodPost, url, nil)

	resp, _ := app.Test(r, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cleanUp(t, services)
}

func Test_HandleProcessFile__should_return_400_when_required_param_is_missing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mocks := &Mocks{}
	jobID := 555
	app, services := setup(t, mocks)

	testCases := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "file_id",
			url:  fmt.Sprintf("%s/%d?format=webp&quality=80&width=100&height=100", processEndpoint, jobID),
			want: `{"errors":{"file_id":"file_id is required"}}`,
		},
		{
			name: "format",
			url:  fmt.Sprintf("%s/%d?quality=80&file_id=1&width=100&height=100", processEndpoint, jobID),
			want: `{"errors":{"format":"format is required"}}`,
		},
		{
			name: "quality",
			url:  fmt.Sprintf("%s/%d?format=webp&file_id=1&width=100&height=100", processEndpoint, jobID),
			want: `{"errors":{"quality":"quality is required"}}`,
		},
		{
			name: "width",
			url:  fmt.Sprintf("%s/%d?format=webp&quality=80&file_id=1&height=100", processEndpoint, jobID),
			want: `{"errors":{"width":"width is required"}}`,
		},
		{
			name: "height",
			url:  fmt.Sprintf("%s/%d?format=webp&quality=80&file_id=1&width=100", processEndpoint, jobID),
			want: `{"errors":{"height":"height is required"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tc.url, nil)
			resp, _ := app.Test(r, -1)
			got, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tc.want, string(got))
			assert.Equal(t, 400, resp.StatusCode)
		})
	}

	cleanUp(t, services)
}
