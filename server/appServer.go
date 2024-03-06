package server

import (
	"cobaApp/config"
	"cobaApp/handler"
	"cobaApp/repository"
	"cobaApp/router"
	"cobaApp/service"
	"database/sql"
	"fmt"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AppServer struct {
	Router *fiber.App
	Config config.IConfig
}

func NewAppServer(db *sql.DB, config config.IConfig, log *logrus.Logger) IServer {
	// register validate
	validate := validator.New()

	// register repository
	carRepo := repository.NewCarRepository(db)

	// register service
	carService := service.NewCarService(db, validate, carRepo)

	// register handler
	carHandler := handler.NewCarHandler(carService, log)

	app := fiber.New(fiber.Config{
		Prefork: false,
	})

	prometheus := fiberprometheus.New("cobaApp-metrics")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	v1 := app.Group("/v1")

	// car router
	router.GenerateCarRouter(v1, carHandler)

	return &AppServer{
		Router: app,
		Config: config,
	}
}

func (a *AppServer) RunServer() error {
	err := a.Router.Listen(fmt.Sprintf(":%v", a.Config.GetConfig().App.Port))
	return err
}
