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
	"gofiber-api/router"
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
	threadRouter := router.NewThreadRoute(threadHandler)

	api := s.app.Group("/api")
	threadRouter.Route(api)
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
	s.Equal(fiber.StatusCreated, resp.StatusCode)

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

func (s *ThreadHttpHandlerSuite) TestGetThreads() {
	s.Db.AddThread(context.Background(), "the-author-1", "the content 1")
	s.Db.AddThread(context.Background(), "the-author-2", "the content 2")
	s.Db.AddThread(context.Background(), "the-author-3", "the content 3")

	req := httptest.NewRequest(fiber.MethodGet, "/api/threads", nil)
	resp, err := s.app.Test(req)
	s.NoError(err)
	s.Equal(fiber.StatusOK, resp.StatusCode)

	bodyString, err := io.ReadAll(resp.Body)
	s.Nil(err)
	s.NotNil(bodyString)
	s.NotEmpty(bodyString)

	var positiveResponse handler.PositiveResponseType
	err = json.Unmarshal(bodyString, &positiveResponse)
	s.NoError(err)
	s.Equal(fiber.StatusOK, positiveResponse.Status)
	s.Equal("success get threads", positiveResponse.Message)

	// Assert that Data is a map and then unmarshal it into Thread
	threadsData, err := json.Marshal(positiveResponse.Data)
	s.NoError(err)
	s.NotEmpty(threadsData)

	var threads []repo.Thread
	err = json.Unmarshal(threadsData, &threads)
	s.NoError(err)
	s.Equal(3, len(threads))

	thread3 := threads[0]
	s.Equal("the-author-3", thread3.Author)
	s.Equal("the content 3", thread3.Content)
	thread2 := threads[1]
	s.Equal("the-author-2", thread2.Author)
	s.Equal("the content 2", thread2.Content)
	thread1 := threads[2]
	s.Equal("the-author-1", thread1.Author)
	s.Equal("the content 1", thread1.Content)
}

func (s *ThreadHttpHandlerSuite) TestEditThread() {
	s.Db.AddThread(context.Background(), "the-author-1", "the content 1")
	s.Db.AddThread(context.Background(), "the-author-2", "the content 2")
	s.Db.AddThread(context.Background(), "the-author-3", "the content 3")

	reqBody := `{"content":"some new contents"}`
	req := httptest.NewRequest(fiber.MethodPut, "/api/threads/1", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.NoError(err)
	s.Equal(fiber.StatusOK, resp.StatusCode)

	bodyString, err := io.ReadAll(resp.Body)
	s.Nil(err)
	s.NotNil(bodyString)
	s.NotEmpty(bodyString)

	var positiveResponse handler.PositiveResponseType
	err = json.Unmarshal(bodyString, &positiveResponse)
	s.NoError(err)
	s.Equal(fiber.StatusOK, positiveResponse.Status)
	s.Equal("success edit thread", positiveResponse.Message)
}
