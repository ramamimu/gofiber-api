package main

import (
	"log"

	handler "gofiber-api/httphandler"
	repo "gofiber-api/repository"
	"gofiber-api/router"
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
	threadRouter := router.NewThreadRoute(threadHandler)

	api := app.Group("/api")
	threadRouter.Route(api)
	log.Fatal(app.Listen(":3001"))
}
