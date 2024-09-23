package main

import (
	"log"

	handler "gofiber-api/httphandler"
	repo "gofiber-api/repository"
	service "gofiber-api/service"

	"github.com/gofiber/fiber/v2"
)

// refactor app

func main() {
	app := fiber.New()

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	db := repo.Db{}
	db.Init()

	threadService := service.NewThread(&db)
	threadHandler := handler.NewThreadHandler(threadService)

	api := app.Group("/api")
	api.Get("/threads", threadHandler.GetAllThreads)
	api.Post("/threads", threadHandler.CreateThread)
	api.Put("/threads/:id", threadHandler.EditThread)
	api.Delete("/threads/:id", threadHandler.DeleteThread)

	log.Fatal(app.Listen(":3001"))
}
