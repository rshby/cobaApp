package router

import (
	"cobaApp/handler"
	"github.com/gofiber/fiber/v2"
)

func GenerateCarRouter(app fiber.Router, handler *handler.CarHandler) {
	app.Post("/car", handler.InsertData)
	app.Get("/cars", handler.GetAll)
	app.Get("/car/:id", handler.GetDetail)
}
