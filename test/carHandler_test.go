package test

import (
	"cobaApp/customError"
	"cobaApp/handler"
	"cobaApp/helper"
	"cobaApp/model/dto"
	mck "cobaApp/test/mock"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInsertCarHandler(t *testing.T) {
	t.Run("test insert error bad request", func(t *testing.T) {
		app := fiber.New()
		carService := mck.NewCarServiceMock()
		carHandler := handler.NewCarHandler(carService)

		app.Post("/", carHandler.InsertData)

		// create request
		request := httptest.NewRequest(http.MethodPost, "/", nil)
		request.Header.Add("Content-Type", "application/json")

		response, err := app.Test(request)

		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		body, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		var responseBody map[string]any
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, "bad request", responseBody["status"].(string))
	})
	t.Run("test insert error not found", func(t *testing.T) {
		carService := mck.NewCarServiceMock()
		carHandler := handler.NewCarHandler(carService)
		app := fiber.New()
		app.Post("/", carHandler.InsertData)

		// mock
		errMessage := "record not found"
		carService.Mock.On("Insert", mock.Anything, mock.Anything).Return(nil, customError.NewNotFoundError(errMessage))

		// create request
		reqBody := dto.InsertCarRequest{
			Name:        "Toyota",
			Price:       1,
			ReleaseDate: "2020-10-10",
		}
		reqJson, _ := json.Marshal(&reqBody)

		// create request
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqJson)))
		request.Header.Add("Content-Type", "application/json")

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)

		body, _ := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusNotFound, response.StatusCode)
		assert.Equal(t, http.StatusNotFound, int(responseBody["status_code"].(float64)))
		assert.Equal(t, "not found", responseBody["status"].(string))
		assert.Equal(t, errMessage, responseBody["message"].(string))
		carService.Mock.AssertExpectations(t)
	})
}

func TestGetAllCarHandler(t *testing.T) {
	t.Run("test get all error not found", func(t *testing.T) {
		carService := mck.NewCarServiceMock()
		carHandler := handler.NewCarHandler(carService)

		app := fiber.New()
		app.Get("/", carHandler.GetAll)

		// mock
		errMessage := "record not found"
		carService.Mock.On("GetAll", mock.Anything).
			Return(nil, customError.NewNotFoundError(errMessage))

		// create request
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		// receive response
		response, err := app.Test(request)

		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		// receive response_body
		body, err := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusNotFound, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusNotFound), responseBody["status"].(string))
	})
	t.Run("test get all error bad request", func(t *testing.T) {
		carService := mck.NewCarServiceMock()
		carHandler := handler.NewCarHandler(carService)

		app := fiber.New()
		app.Get("/", carHandler.GetAll)

		// mock
		errMessage := "error bad request"
		carService.Mock.On("GetAll", mock.Anything).Return(nil, customError.NewBadRequestError(errMessage))

		// create request
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		// receive response_body
		body, err := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusBadRequest), responseBody["status"].(string))
		carService.Mock.AssertExpectations(t)
	})
	t.Run("test get all error internal server", func(t *testing.T) {
		carService := mck.NewCarServiceMock()
		carHandler := handler.NewCarHandler(carService)

		app := fiber.New()
		app.Get("/", carHandler.GetAll)

		// mock
		errorMessage := "error internal server error"
		carService.Mock.On("GetAll", mock.Anything).Return(nil, customError.NewInternalServerError(errorMessage))

		// create request
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

		// receive response_body
		body, err := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		// test
		assert.Equal(t, http.StatusInternalServerError, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusInternalServerError), responseBody["status"].(string))
		assert.Equal(t, errorMessage, responseBody["message"].(string))
		carService.Mock.AssertExpectations(t)
	})
	t.Run("test get all success", func(t *testing.T) {
		carService := mck.NewCarServiceMock()
		carHandler := handler.NewCarHandler(carService)

		app := fiber.New()
		app.Get("/", carHandler.GetAll)

		// mock
		carService.Mock.On("GetAll", mock.Anything).Return([]dto.InsertCarResponse{
			{
				Id:          1,
				Name:        "Toyota",
				Price:       614000000,
				ReleaseDate: "2020-10-10",
			},
		}, nil)

		// create request
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// receive response_body
		body, err := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		// test dulu
		assert.Equal(t, http.StatusOK, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusOK), responseBody["status"].(string))
		assert.Equal(t, "success get all data cars", responseBody["message"].(string))
		carService.Mock.AssertExpectations(t)
	})
}

func TestGetDetailHandler(t *testing.T) {
	t.Run("test get detail failed convert int", func(t *testing.T) {
		carService := mck.NewCarServiceMock()
		carHandler := handler.NewCarHandler(carService)

		app := fiber.New()
		app.Get("/:id", carHandler.GetDetail)

		// create request
		request := httptest.NewRequest(http.MethodGet, "/wasd", nil)

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		// receive response_body
		body, err := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusBadRequest), responseBody["status"].(string))
		assert.Equal(t, "cant convert id to int", responseBody["message"].(string))
		carService.Mock.AssertExpectations(t)
	})
	t.Run("test get detail not found", func(t *testing.T) {
		carService := mck.NewCarServiceMock()
		carHandler := handler.NewCarHandler(carService)

		app := fiber.New()
		app.Get("/:id", carHandler.GetDetail)

		// mock
		errMessage := "record not found"
		carService.Mock.On("GetDetail", mock.Anything, mock.Anything).Return(nil, customError.NewNotFoundError(errMessage))

		// create request
		request := httptest.NewRequest(http.MethodGet, "/999", nil)

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		// receive response_body
		body, err := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusNotFound, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusNotFound), responseBody["status"].(string))
		assert.Equal(t, errMessage, responseBody["message"].(string))
		carService.Mock.AssertExpectations(t)
	})
	t.Run("test get detail internal server error", func(t *testing.T) {
		carService := mck.NewCarServiceMock()
		carHandler := handler.NewCarHandler(carService)

		app := fiber.New()
		app.Get("/:id", carHandler.GetDetail)

		// mock
		errMessage := "error internal server error"
		carService.Mock.On("GetDetail", mock.Anything, mock.Anything).
			Return(nil, customError.NewInternalServerError(errMessage))

		// create request
		request := httptest.NewRequest(http.MethodGet, "/1", nil)

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

		// receive response_body
		body, err := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusInternalServerError, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusInternalServerError), responseBody["status"].(string))
		assert.Equal(t, errMessage, responseBody["message"].(string))
		carService.Mock.AssertExpectations(t)
	})
	t.Run("test get detail success", func(t *testing.T) {
		carService := mck.NewCarServiceMock()
		carHandler := handler.NewCarHandler(carService)

		app := fiber.New()
		app.Get("/:id", carHandler.GetDetail)

		// mock
		carService.Mock.On("GetDetail", mock.Anything, mock.Anything).Return(&dto.InsertCarResponse{
			Id:          1,
			Name:        "Toyota",
			Price:       12345,
			ReleaseDate: "2020-10-10",
		}, nil)

		// test -> create request
		request := httptest.NewRequest(http.MethodGet, "/1", nil)

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// receive response_body
		body, err := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusOK, int(responseBody["status_code"].(float64)))
		assert.Equal(t, "ok", helper.CodeToStatus(int(responseBody["status_code"].(float64))))
		carService.Mock.AssertExpectations(t)
	})
}
