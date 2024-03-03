package test

import (
	"cobaApp/customError"
	"cobaApp/handler"
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
