package handler

import (
	"cobaApp/customError"
	"cobaApp/helper"
	"cobaApp/model/dto"
	"cobaApp/service"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
	"strings"
)

type CarHandler struct {
	CarService service.ICarService
}

// function provider
func NewCarHandler(carService service.ICarService) *CarHandler {
	return &CarHandler{CarService: carService}
}

// handler insert data
func (c *CarHandler) InsertData(ctx *fiber.Ctx) error {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx.Context(), "Handler InsertData")
	defer span.Finish()
	// parsing body request
	var request dto.InsertCarRequest
	if err := ctx.BodyParser(&request); err != nil {
		statusCode := http.StatusBadRequest
		ctx.Status(statusCode)
		return ctx.JSON(&dto.ApiResponse{
			StatusCode: statusCode,
			Status:     "bad request",
			Message:    err.Error(),
		})
	}

	// log request to tracing
	reqJson, _ := json.Marshal(&request)
	span.LogFields(log.String("request", string(reqJson)))

	// call service
	newCar, err := c.CarService.Insert(ctxTracing, &request)

	var statusCode int
	if err != nil {
		// cek if error validator
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errMessage []string
			for _, errorField := range validationErrors {
				errMessage = append(errMessage, fmt.Sprintf("error on field [%v] with tag [%v]",
					errorField.Field(), errorField.ActualTag()))
			}

			statusCode = http.StatusBadRequest
			return ctx.JSON(&dto.ApiResponse{
				StatusCode: statusCode,
				Status:     helper.CodeToStatus(statusCode),
				Message:    strings.Join(errMessage, ". "),
			})
		}

		switch err.(type) {
		case *customError.BadRequestError:
			statusCode = http.StatusBadRequest
		case *customError.NotFoundError:
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}

		ctx.Status(statusCode)
		return ctx.JSON(&dto.ApiResponse{
			StatusCode: statusCode,
			Status:     helper.CodeToStatus(statusCode),
			Message:    err.Error(),
		})
	}

	// success insert
	statusCode = http.StatusOK
	ctx.Status(statusCode)
	return ctx.JSON(&dto.ApiResponse{
		StatusCode: statusCode,
		Status:     helper.CodeToStatus(statusCode),
		Message:    "success insert data",
		Data:       newCar,
	})
}

// handler get all data cars
func (c *CarHandler) GetAll(ctx *fiber.Ctx) error {
	// start span
	span, ctxTracing := opentracing.StartSpanFromContext(ctx.Context(), "Handler GetAll")
	defer span.Finish()

	// call service
	cars, err := c.CarService.GetAll(ctxTracing)
	if err != nil {
		var statusCode int
		switch err.(type) {
		case *customError.NotFoundError:
			statusCode = http.StatusNotFound
		case *customError.BadRequestError:
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}

		ctx.Status(statusCode)

		response := dto.ApiResponse{
			StatusCode: statusCode,
			Status:     helper.CodeToStatus(statusCode),
			Message:    err.Error(),
		}
		resJson, _ := json.Marshal(&response)
		span.LogFields(log.String("response", string(resJson)))
		return ctx.JSON(&response)
	}

	// success get data
	statusCode := http.StatusOK
	ctx.Status(statusCode)

	response := dto.ApiResponse{
		StatusCode: statusCode,
		Status:     helper.CodeToStatus(statusCode),
		Message:    "success get all data cars",
		Data:       cars,
	}
	resJson, _ := json.Marshal(&response)
	span.LogFields(log.String("response", string(resJson)))
	return ctx.JSON(&response)
}

// handler get detail
func (c *CarHandler) GetDetail(ctx *fiber.Ctx) error {
	// start span
	span, ctxTracing := opentracing.StartSpanFromContext(ctx.Context(), "Handler GetDetail")
	defer span.Finish()

	id, err := ctx.ParamsInt("id")
	if err != nil {
		statusCode := http.StatusBadRequest
		ctx.Status(statusCode)
		response := dto.ApiResponse{
			StatusCode: statusCode,
			Status:     helper.CodeToStatus(statusCode),
			Message:    "cant convert id to int",
		}
		resJson, _ := json.Marshal(&response)
		span.LogFields(log.String("response", string(resJson)))

		return ctx.JSON(response)
	}

	// call procedure in service
	var statusCode int
	car, err := c.CarService.GetDetail(ctxTracing, id)
	if err != nil {
		switch err.(type) {
		case *customError.NotFoundError:
			statusCode = http.StatusNotFound
		case *customError.BadRequestError:
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}

		ctx.Status(statusCode)

		response := dto.ApiResponse{
			StatusCode: statusCode,
			Status:     helper.CodeToStatus(statusCode),
			Message:    err.Error(),
		}

		resJson, _ := json.Marshal(&response)
		span.LogFields(log.String("response", string(resJson)))

		return ctx.JSON(&response)
	}

	// success get detail
	resJson, _ := json.Marshal(&car)
	span.LogFields(log.String("response", string(resJson)))
	statusCode = http.StatusOK
	ctx.Status(statusCode)
	return ctx.JSON(&dto.ApiResponse{
		StatusCode: statusCode,
		Status:     helper.CodeToStatus(statusCode),
		Message:    "success get data detail car",
		Data:       car,
	})
}
