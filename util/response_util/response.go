package response_util

import (
	"coinbit-wallet/dto/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Message   string      `json:"message,omitempty"`
	Status    int         `json:"status,omitempty"`
	FromCache bool        `json:"from_cache"`
}

func Ok(c *gin.Context, data interface{}, isFromCache bool) {
	respond(c, http.StatusOK, Response{Message: "SUCCESS", Data: data, Status: http.StatusOK, FromCache: isFromCache})
}

func Accepted(c *gin.Context, data interface{}, isFromCache bool) {
	respond(c, http.StatusAccepted, Response{Message: "SUCCESS", Data: data, Status: http.StatusAccepted, FromCache: isFromCache})
}

func Fail(c *gin.Context, response app.ErrorHttpResponse) {
	respond(c, response.HttpStatus, Response{Message: response.Message, Error: response.ErrorName, Status: response.HttpStatus})
}

func respond(c *gin.Context, code int, payload Response) {
	c.JSON(code, payload)
	c.Abort()
}
