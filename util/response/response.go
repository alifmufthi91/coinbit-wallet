package responseUtil

import (
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

func Success(c *gin.Context, data interface{}, isFromCache bool) {
	respond(c, http.StatusOK, Response{Message: "SUCCESS", Data: data, Status: http.StatusOK, FromCache: isFromCache})
}

func Fail(c *gin.Context, err string) {
	respond(c, http.StatusInternalServerError, Response{Message: err, Error: "Internal Server Error", Status: http.StatusInternalServerError})
}

func respond(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}
