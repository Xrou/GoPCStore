package main

import (
	"TestProject/controllers"
	"TestProject/model"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	model.Connect()
	defer model.Disconnect()

	app := fiber.New()

	apiGroup := app.Group("api")
	pcGroup := apiGroup.Group("computers")
	pcGroup.Get("/", controllers.GetComputers)
	pcGroup.Post("/", controllers.PostComputer)
	pcGroup.Put("/:id", controllers.PutComputer)
	pcGroup.Delete("/:id", controllers.DeleteComputer)

	log.Fatal(app.Listen(":3000"))
}
