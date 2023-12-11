package router

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/prplx/lighter.pics/internal/handlers"
	svc "github.com/prplx/lighter.pics/internal/services"
	"github.com/prplx/lighter.pics/internal/types"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type jobCreatedResponse struct {
	JobID string `json:"job_id"`
}

type CommunicatorMock struct {
	mu         sync.Mutex
	StartCalls int
	ErrCalls   int
}

func (c *CommunicatorMock) AddClient(jobID int, connection *websocket.Conn) {}

func (c *CommunicatorMock) RemoveClient(jobID int) {}

func (c *CommunicatorMock) SendStartProcessing(jobID, fileID int, fileName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.StartCalls++
	return nil
}

func (c *CommunicatorMock) SendErrorProcessing(jobID, fileID int, fileName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ErrCalls++
	return nil
}

func (c *CommunicatorMock) SendSuccessProcessing(jobID int, result types.SuccessResult) error {
	return nil
}

func (c *CommunicatorMock) SendStartArchiving(jobID int) error {
	return nil
}

func (c *CommunicatorMock) SendErrorArchiving(jobID int) error {
	return nil
}

func (c *CommunicatorMock) SendSuccessArchiving(jobID int, path string) error {
	return nil
}

func (c *CommunicatorMock) SendSuccessFlushing(jobID int) error {
	return nil
}

func (c *CommunicatorMock) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.StartCalls = 0
	c.ErrCalls = 0
}

func (c *CommunicatorMock) GetStartCalls() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.StartCalls
}

type LoggerMock struct{}

func (l *LoggerMock) PrintInfo(message string, properties ...types.AnyMap) {}

func (l *LoggerMock) PrintError(err error, properties ...types.AnyMap) {}

func (l *LoggerMock) PrintFatal(err error, properties ...types.AnyMap) {}

func (l *LoggerMock) Write(message []byte) (n int, err error) { return 0, nil }

var communicator = &CommunicatorMock{}
var services = svc.Services{
	Communicator: communicator,
	Logger:       &LoggerMock{},
	Config: &types.Config{
		Process: struct {
			UploadDir string `yaml:"uploadDir"`
		}{
			UploadDir: "./uploads",
		},
	},
}

const (
	healthcheckEndpoint = "/api/v1/healthcheck"
	processEndpoint     = "/api/v1/process"
)

func Test_Healthcheck(t *testing.T) {
	app := setup(t)
	r := httptest.NewRequest(http.MethodGet, healthcheckEndpoint, nil)
	resp, _ := app.Test(r, -1)
	got, _ := io.ReadAll(resp.Body)
	want := `{"status":"ok"}`

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, want, string(got))
}

func Test_Process(t *testing.T) {
	body, contentType := createFormFile(t, "file", "file.png")
	app := setup(t)
	r := httptest.NewRequest(http.MethodPost, processEndpoint, body)
	r.Header.Add("Content-Type", contentType)

	resp, _ := app.Test(r, -1)
	got, _ := io.ReadAll(resp.Body)
	createdResponse := jobCreatedResponse{}
	json.Unmarshal(got, &createdResponse)
	jobID := createdResponse.JobID

	assert.NotEqual(t, "", jobID)
	assert.Equal(t, 202, resp.StatusCode)
	assert.Equal(t, 1, communicator.GetStartCalls())
	cleanUp(t)
}

func setup(t *testing.T) *fiber.App {
	t.Helper()
	app := fiber.New()
	handlers := handlers.NewHandlers(&services)

	Register(app, handlers)
	return app
}

func cleanUp(t *testing.T) {
	t.Helper()
	os.RemoveAll(services.Config.Process.UploadDir)
	communicator.Reset()
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
