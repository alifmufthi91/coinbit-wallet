package error_util

import (
	"net/http"
)

func Unauthorized(message string) *CustomError {
	return &CustomError{
		Message:    message,
		HttpStatus: http.StatusUnauthorized,
		ErrorName:  "UNAUTHORIZED",
	}
}
