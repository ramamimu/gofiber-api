package httphandler

import (
	"context"
	repo "gofiber-api/repository"

	"github.com/gofiber/fiber/v2"
)

type HttpThreadHandlerRepo interface {
	GetAll(ctx context.Context) []repo.Thread
	Add(ctx context.Context, author string, content string) error
	Edit(ctx context.Context, id string, content string) error
	Delete(ctx context.Context, id string) error
}

type PositiveResponseType struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type NegativeResponseType struct {
	Data []string `json:"data"`
}

type CreateThreadRequestType struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

type EditThreadRequestType struct {
	NewContent string `json:"content"`
}

type ThreadHandler struct {
	HttpThreadHandlerRepo
}

func NewThreadHandler(threadService HttpThreadHandlerRepo) *ThreadHandler {
	return &ThreadHandler{
		HttpThreadHandlerRepo: threadService,
	}
}

func (th *ThreadHandler) GetAllThreads(c *fiber.Ctx) error {
	threads := th.GetAll(context.Background())
	return c.Status(fiber.StatusOK).JSON(PositiveResponseType{
		Status:  fiber.StatusOK,
		Message: "success get threads",
		Data:    threads,
	})
}

func (th *ThreadHandler) CreateThread(c *fiber.Ctx) error {
	threadRequest := new(CreateThreadRequestType)

	if err := c.BodyParser(threadRequest); err != nil {
		return c.JSON(NegativeResponseType{
			Data: []string{err.Error()},
		})
	}

	if err := th.Add(context.Background(), threadRequest.Author, threadRequest.Content); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(NegativeResponseType{
			Data: []string{err.Error()},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(PositiveResponseType{
		Status:  fiber.StatusCreated,
		Message: "success create thread",
		Data:    nil,
	})
}

func (th *ThreadHandler) EditThread(c *fiber.Ctx) error {
	threadRequest := new(EditThreadRequestType)

	if err := c.BodyParser(threadRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(NegativeResponseType{
			Data: []string{err.Error()},
		})
	}

	if err := th.Edit(context.Background(), c.Params("id"), threadRequest.NewContent); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(NegativeResponseType{
			Data: []string{err.Error()},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(PositiveResponseType{
		Status:  fiber.StatusCreated,
		Message: "success edit thread",
		Data:    nil,
	})
}

func (th *ThreadHandler) DeleteThread(c *fiber.Ctx) error {
	if err := th.Delete(context.Background(), c.Params("id")); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(NegativeResponseType{
			Data: []string{err.Error()},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(PositiveResponseType{
		Status:  fiber.StatusCreated,
		Message: "success delete thread",
		Data:    nil,
	})
}
