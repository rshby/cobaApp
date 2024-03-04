package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
)

func LoggerMiddleware(log *logrus.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		startTime := time.Now()

		// get request
		body := ctx.Request().Body()
		requestBody := map[string]any{}
		json.Unmarshal(body, &requestBody)

		ctx.Next()

		response := ctx.Response().Body()
		responseBody := map[string]any{}
		json.Unmarshal(response, &responseBody)
		log.WithFields(map[string]interface{}{
			"url":           string(ctx.Request().URI().PathOriginal()),
			"method":        string(ctx.Request().Header.Method()),
			"request":       requestBody,
			"response":      responseBody,
			"status_code":   ctx.Response().StatusCode(),
			"response_time": fmt.Sprintf("%vms", time.Since(startTime).Milliseconds()),
		}).Info("request coming")
		return nil
	}
}
