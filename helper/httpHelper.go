package helper

import "net/http"

func CodeToStatus(code int) string {
	switch code {
	case http.StatusOK:
		return "ok"
	case http.StatusBadRequest:
		return "bad request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusNotFound:
		return "not found"
	default:
		return "internal server error"
	}
}
