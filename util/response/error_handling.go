package responseUtil

import (
	"coinbit-wallet/util/logger"

	"github.com/gin-gonic/gin"
)

func ErrorHandling(c *gin.Context) {
	if r := recover(); r != nil {
		logger.Error("Recovered from panic: %+v", r)
		var errorMessage string
		if err, ok := r.(error); ok {
			errorMessage = err.Error()
		} else if err, ok := r.(string); ok {
			errorMessage = err
		} else {
			errorMessage = "Internal Error"
		}
		Fail(c, errorMessage)
		return
	}
}
