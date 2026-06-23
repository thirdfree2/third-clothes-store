package httpcommon

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// Generic 400 errors
	CodeBadRequest      = 400000
	CodeValidationError = 400001

	// Clothes 400 errors
	CodeInvalidClothesID   = 400101
	CodeInvalidClothesName = 400102
	CodeInvalidColorID     = 400103
	CodeInvalidCategoryID  = 400104
	CodeForbidden          = 400105
	CodeUnauthorized       = 400106

	// Clothes 404 errors
	CodeClothesNotFound = 404101

	// Clothes 409 errors
	CodeClothesDuplicated = 409101

	// Generic 500 errors
	CodeInternalServerError = 500000
	CodeBadGateway          = 500002
)

// Response is the standard HTTP response envelope.
type Response struct {
	Success bool         `json:"success" example:"true"`
	Payload any          `json:"payload" swaggertype:"object"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

// ErrorDetail is the standard error body inside Response.
type ErrorDetail struct {
	Code    int            `json:"code" example:"400001"`
	Message string         `json:"message" example:"invalid request body"`
	Details map[string]any `json:"details,omitempty" swaggertype:"object"`
}

// HealthResponse is the response body for health check.
type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

// Health godoc
// @Summary Health check
// @Description Check whether the API is running.
// @Tags system
// @Produce json
// @Success 200 {object} Response
// @Router /health [get]
func Health(c *gin.Context) {
	RespondSuccess(c, http.StatusOK, HealthResponse{Status: "ok"})
}

func RespondSuccess(c *gin.Context, status int, payload any) {
	c.JSON(status, Response{
		Success: true,
		Payload: payload,
	})
}

func RespondError(c *gin.Context, status int, code int, message string, details map[string]any) {
	c.JSON(status, Response{
		Success: false,
		Payload: nil,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

type ResponseOf[T any] struct {
	Success bool         `json:"success"`
	Payload T            `json:"payload"`
	Error   *ErrorDetail `json:"error,omitempty"`
}
