package middleware

import (
	"coinbit-wallet/dto/app"
	"coinbit-wallet/util/error_util"
	"coinbit-wallet/util/logger"
	"coinbit-wallet/util/response_util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Recovered from panic: %+v", r)
				var errorMessage string
				if err, ok := r.(*error_util.CustomError); ok {
					response_util.Fail(c, app.ErrorHttpResponse{
						Message:    err.Error(),
						HttpStatus: err.HttpStatus,
						ErrorName:  err.ErrorName,
					})
					return
				} else if err, ok := r.(error); ok {
					errorMessage = err.Error()
				} else if err, ok := r.(string); ok {
					errorMessage = err
				} else {
					errorMessage = "Internal Error"
				}
				response_util.Fail(c, app.ErrorHttpResponse{
					Message:    errorMessage,
					HttpStatus: http.StatusInternalServerError,
					ErrorName:  "INTERNAL SERVER ERROR",
				})
			}
		}()
		c.Next()
	}
}
