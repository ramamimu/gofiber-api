package httphandler

import (
	"context"
	repo "gofiber-api/repository"
	"sort"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type HttpThreadHandlerRepo interface {
	GetAll(ctx context.Context) []repo.Thread
	Add(ctx context.Context, author string, content string) error
	Edit(ctx context.Context, id string, content string) error
	Delete(ctx context.Context, id string) error
}

type ResponseType struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CreateThreadRequestType struct {
	Content string `json:"content" validate:"required"`
	Author  string `json:"author" validate:"required"`
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
	sort.SliceStable(threads, func(i, j int) bool {
		return threads[i].LastUpdate.After(threads[j].LastUpdate)
	})

	return c.Status(fiber.StatusOK).JSON(ResponseType{
		Status:  fiber.StatusOK,
		Message: "success get threads",
		Data:    threads,
	})
}

var validate = validator.New()

func (th *ThreadHandler) CreateThread(c *fiber.Ctx) error {
	threadRequest := new(CreateThreadRequestType)

	if err := c.BodyParser(threadRequest); err != nil {
		return c.JSON(
			ResponseType{
				Status:  c.Response().StatusCode(),
				Message: c.Response().String(),
				Data: []string{
					err.Error(),
				},
			})
	}

	// empty body request validation
	// best practice is put all of the parameter in domain
	// Validate the parsed body using the validator
	if err := validate.Struct(threadRequest); err != nil {
		// If validation fails, return a bad request error with details
		return c.Status(fiber.StatusBadRequest).JSON(ResponseType{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Data: func(err error) []string {
				var errors []string
				for _, err := range err.(validator.ValidationErrors) {
					errors = append(errors, err.Field()+" is "+err.Tag())
				}
				return errors
			}(err),
		})
	}

	if err := th.Add(context.Background(), threadRequest.Author, threadRequest.Content); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseType{
			Status:  c.Response().StatusCode(),
			Message: c.Response().String(),
			Data: []string{
				err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(ResponseType{
		Status:  fiber.StatusCreated,
		Message: "success create thread",
		Data:    nil,
	})
}

func (th *ThreadHandler) EditThread(c *fiber.Ctx) error {
	threadRequest := new(EditThreadRequestType)

	if err := c.BodyParser(threadRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseType{
			Status:  fiber.StatusBadRequest,
			Message: "bad request",
			Data: []string{
				err.Error(),
			},
		})
	}

	if err := th.Edit(context.Background(), c.Params("id"), threadRequest.NewContent); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ResponseType{
			Status:  fiber.StatusNotFound,
			Message: "not found",
			Data: []string{
				err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(ResponseType{
		Status:  fiber.StatusOK,
		Message: "success edit thread",
		Data:    nil,
	})
}

func (th *ThreadHandler) DeleteThread(c *fiber.Ctx) error {
	if err := th.Delete(context.Background(), c.Params("id")); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ResponseType{
			Status:  fiber.StatusNotFound,
			Message: "not found",
			Data: []string{
				err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(ResponseType{
		Status:  fiber.StatusOK,
		Message: "success delete thread",
		Data:    nil,
	})
}
