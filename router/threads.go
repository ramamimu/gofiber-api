package router

import "github.com/gofiber/fiber/v2"

type RouterImplementation interface {
	GetAllThreads(c *fiber.Ctx) error
	CreateThread(c *fiber.Ctx) error
	EditThread(c *fiber.Ctx) error
	DeleteThread(c *fiber.Ctx) error
}

type ThreadRoute struct {
	RouterImplementation
}

func NewThreadRoute(r RouterImplementation) *ThreadRoute {
	return &ThreadRoute{
		RouterImplementation: r,
	}
}

func (tr *ThreadRoute) Route(app fiber.Router) {
	app.Get("/threads", tr.GetAllThreads)
	app.Post("/threads", tr.CreateThread)
	app.Put("/threads/:id", tr.EditThread)
	app.Delete("/threads/:id", tr.DeleteThread)
}
