package httphandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"

	handler "gofiber-api/httphandler"
	repo "gofiber-api/repository"
	service "gofiber-api/service"
)

type ThreadHttpHandlerSuite struct {
	suite.Suite
	app *fiber.App
	Db  repo.Db
}

func TestThreadHttpHandlerSuite(t *testing.T) {
	suite.Run(t, new(ThreadHttpHandlerSuite))
}

func (s *ThreadHttpHandlerSuite) SetupSuite() {
	s.app = fiber.New()

	s.app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	s.Db = repo.Db{}

	threadService := service.NewThread(&s.Db)
	threadHandler := handler.NewThreadHandler(threadService)

	api := s.app.Group("/api")
	api.Get("/threads", threadHandler.GetAllThreads)
	api.Post("/threads", threadHandler.CreateThread)
	api.Put("/threads/:id", threadHandler.EditThread)
	api.Delete("/threads/:id", threadHandler.DeleteThread)
}

func (s *ThreadHttpHandlerSuite) SetupTest() {
	s.Db.Clear()
	s.Db.Init()
}

func (s *ThreadHttpHandlerSuite) AfterTest() {
	s.Db.Clear()
}

func (s *ThreadHttpHandlerSuite) TestPing() {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	resp, err := s.app.Test(req)

	s.Nil(err)
	s.Equal(200, resp.StatusCode)

	bodyString, err := io.ReadAll(resp.Body)
	s.Nil(err)
	s.Equal("pong", string(bodyString))
}

func (s *ThreadHttpHandlerSuite) TestCreateNewThread() {
	reqBody := `{"content":"hello world","author":"ramamimu"}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Nil(err)

	bodyString, err := io.ReadAll(resp.Body)
	s.Nil(err)
	s.NotNil(bodyString)
	s.NotEmpty(bodyString)

	var positiveResponseType handler.PositiveResponseType
	err = json.Unmarshal(bodyString, &positiveResponseType)
	s.Nil(err)
	s.Equal(fiber.StatusCreated, positiveResponseType.Status)
	s.Equal("success create thread", positiveResponseType.Message)
	s.Nil(positiveResponseType.Data)

	thread, err := s.Db.GetThreadByID("0")
	s.Nil(err)
	s.NotEmpty(thread)
	s.Equal("hello world", thread.Content)
	s.Equal("ramamimu", thread.Author)

	threadsEntity := s.Db.GetThreadsEntity()
	s.Equal(1, len(threadsEntity))

	threads := s.Db.GetThreads(context.Background())
	s.Equal(1, len(threads))

}
