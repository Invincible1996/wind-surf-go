package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response codes
const (
	CodeSuccess = 0
	CodeError   = 1
)

// Response is the standard API response structure
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse returns a success response with data
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// CreatedResponse returns a success response for resource creation
func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "created successfully",
		Data:    data,
	})
}

// ErrorResponse returns an error response
func ErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeError,
		Message: message,
	})
}

// ServerErrorResponse returns a server error response
func ServerErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeError,
		Message: message,
	})
}
