package handlers_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gofiber/fiber/v2"
	"github.com/prplx/cnvrt/internal/handlers"
	"github.com/prplx/cnvrt/internal/mocks"
	"github.com/prplx/cnvrt/internal/repositories"
	"github.com/prplx/cnvrt/internal/router"
	svc "github.com/prplx/cnvrt/internal/services"
	"github.com/prplx/cnvrt/internal/types"
	"github.com/spf13/afero"
)

const (
	processEndpoint = "/api/v1/process"
	archiveEndpoint = "/api/v1/archive"
)

type Mocks struct {
	jobsRepo     *mocks.MockJobs
	filesRepo    *mocks.MockFiles
	communicator *mocks.MockCommunicator
	logger       *mocks.MockLogger
	processor    *mocks.MockProcessor
	archiver     *mocks.MockArchiver
}

func Test_Healthcheck(t *testing.T) {
	mocks := &Mocks{}
	app, _ := setup(t, mocks)
	r := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	resp, _ := app.Test(r, -1)
	got, _ := io.ReadAll(resp.Body)
	want := `{"status":"ok"}`

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, want, string(got))
}

func setup(t *testing.T, mocks *Mocks) (*fiber.App, *svc.Services) {
	t.Helper()
	app := fiber.New()
	services := getServices(t, mocks)
	handlers := handlers.NewHandlers(&services)

	router.Register(app, handlers, services.Config, nil)
	return app, &services
}

func cleanUp(t *testing.T, services *svc.Services) {
	t.Helper()
	os.RemoveAll(services.Config.Process.UploadDir)
}

func createFormFile(t *testing.T, fieldName, filePath string) (*bytes.Buffer, string) {
	t.Helper()
	body := new(bytes.Buffer)

	mw := multipart.NewWriter(body)
	appFS := afero.NewMemMapFs()
	afero.WriteFile(appFS, filePath, []byte("file"), 0644)

	file, err := appFS.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}

	w, err := mw.CreateFormFile(fieldName, filePath)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.Copy(w, file); err != nil {
		t.Fatal(err)
	}

	mw.Close()
	return body, mw.FormDataContentType()
}

func getServices(_ *testing.T, mocks *Mocks) svc.Services {
	return svc.Services{
		Communicator: mocks.communicator,
		Logger:       mocks.logger,
		Config: &types.Config{
			Process: struct {
				UploadDir string `yaml:"uploadDir"`
			}{
				UploadDir: "./uploads",
			},
		},
		Repositories: &repositories.Repositories{
			Jobs:  mocks.jobsRepo,
			Files: mocks.filesRepo,
		},
		Processor: mocks.processor,
		Archiver:  mocks.archiver,
	}
}
