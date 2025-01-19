package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
// @Description Error response with a message
type ErrorResponse struct {
	Error string `json:"error" example:"invalid request"`
}

// SuccessResponse represents a success response
// @Description Success response with optional message and data
type SuccessResponse struct {
	Message string      `json:"message,omitempty" example:"operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// UserResponse represents a user in responses
// @Description User profile information
type UserResponse struct {
	ID    string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email string `json:"email" example:"user@example.com"`
	Name  string `json:"name" example:"John Doe"`
}

// HealthResponse represents a health check response
// @Description Health check response with status
type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Version string `json:"version,omitempty" example:"1.0.0"`
}

func Error(c *gin.Context, status int, err error) {
	c.JSON(status, ErrorResponse{Error: err.Error()})
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{Data: data})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *gin.Context, err error) {
	Error(c, http.StatusBadRequest, err)
}

func Unauthorized(c *gin.Context, err error) {
	Error(c, http.StatusUnauthorized, err)
}

func InternalError(c *gin.Context, err error) {
	Error(c, http.StatusInternalServerError, err)
}
