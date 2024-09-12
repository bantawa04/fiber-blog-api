package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"my-fiber-project/config"
	"my-fiber-project/controllers"
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func main() {
	app := fiber.New()
	// Connect to the database
	if err := config.Connect(); err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		response := Response{
			Message: "hello world",
			Status:  200,
		}
		return c.JSON(response)
	})

	postController := controllers.NewPostController()
	categoryController := controllers.NewCategoryController()

	app.Get("/posts", postController.Index)
	app.Post("/posts", postController.Store)
	app.Get("/posts/:id", postController.Show)
	app.Put("/posts/:id", postController.Update)
	app.Delete("/posts/:id", postController.Destroy)

	app.Get("/categories", categoryController.Index)
	app.Post("/categories", categoryController.Store)
	app.Get("/categories/:id", categoryController.Show)
	app.Put("/categories/:id", categoryController.Update)
	app.Delete("/categories/:id", categoryController.Destroy)

	log.Fatal(app.Listen(":3000"))
}
