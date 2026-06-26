package httpcommon

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeBadRequest          = 400000
	CodeValidationError     = 400001
	CodeUnauthorized        = 401000
	CodeForbidden           = 403000
	CodeDuplicatePhone      = 409101
	CodeProfileNotFound     = 404101
	CodeCartItemNotFound    = 404201
	CodeInternalServerError = 500000
)

type Response struct {
	Success bool         `json:"success"`
	Payload any          `json:"payload"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

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
