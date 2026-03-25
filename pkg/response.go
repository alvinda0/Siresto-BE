package pkg

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Status    int         `json:"status"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
	Meta      interface{} `json:"meta,omitempty"`
	Error     string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Success:   true,
		Message:   message,
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	})
}

func SuccessResponseWithMeta(c *gin.Context, status int, message string, data interface{}, meta interface{}) {
	c.JSON(status, Response{
		Success:   true,
		Message:   message,
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
		Meta:      meta,
	})
}

func ErrorResponse(c *gin.Context, status int, message string, err string) {
	c.JSON(status, Response{
		Success:   false,
		Message:   message,
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Error:     err,
	})
}
