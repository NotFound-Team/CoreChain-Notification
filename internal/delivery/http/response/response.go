package response

import (
	"github.com/gin-gonic/gin"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// JSON sends a success response with data
func JSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// JSONWithMessage sends a success response with data and message
func JSONWithMessage(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, errorMsg string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   errorMsg,
	})
}

// ErrorWithCode sends an error response with error code
func ErrorWithCode(c *gin.Context, statusCode int, errorMsg string, code string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   errorMsg,
		Code:    code,
	})
}
