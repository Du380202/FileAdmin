package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"statuscode"`
}

func (e *Response) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.StatusCode, e.Status, e.Message)
}

func ErrorResponse(c *gin.Context, message string, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Status:     "error",
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	})
}
func SuccessResponse(c *gin.Context, message string, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Status:     "success",
		Message:    message,
		Data:       data,
		StatusCode: statusCode,
	})
}
