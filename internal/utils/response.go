package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"se69-lab04-starter/internal/dto"
)

// send an error as application/problem+json and stop the request
func WriteError(c *gin.Context, status int, detail string) {
	c.Header("Content-Type", "application/problem+json")
	c.AbortWithStatusJSON(status, dto.ErrorResponse{
		Status: status,
		Title:  http.StatusText(status),
		Detail: detail,
	})
}

// read :id from the URL path, for example /equipment/:id
func ParamID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(id), err == nil
}
