// internal/adapters/inbound/http/auth/handler.go

package httpauth

import (
	"errors"
	"net/http"

	httpcommon "auth-service/internal/adapters/inbound/http/common"
	"auth-service/internal/application/services"
	"auth-service/internal/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService *services.AuthService
}

func NewHandler(authService *services.AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	UserType string `json:"userType" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type RegisterResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	Status   string `json:"status"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpcommon.RespondError(
			c,
			http.StatusBadRequest,
			httpcommon.CodeValidationError,
			"invalid request body",
			map[string]any{"error": err.Error()},
		)
		return
	}

	user, err := h.authService.RegisterCustomer(
		c.Request.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusCreated, RegisterResponse{
		ID:       user.ID,
		Email:    user.Email,
		UserType: string(user.UserType),
		Status:   string(user.Status),
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpcommon.RespondError(
			c,
			http.StatusBadRequest,
			httpcommon.CodeValidationError,
			"invalid request body",
			map[string]any{"error": err.Error()},
		)
		return
	}

	token, err := h.authService.Login(
		c.Request.Context(),
		req.Email,
		req.Password,
		req.UserType,
	)
	if err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	})
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidEmailOrPassword):
		httpcommon.RespondError(
			c,
			http.StatusUnauthorized,
			httpcommon.CodeUnauthorized,
			"invalid email or password",
			nil,
		)

	case errors.Is(err, domain.ErrUserInactive):
		httpcommon.RespondError(
			c,
			http.StatusForbidden,
			httpcommon.CodeForbidden,
			"user inactive",
			nil,
		)

	case errors.Is(err, domain.ErrInvalidUserType):
		httpcommon.RespondError(
			c,
			http.StatusBadRequest,
			httpcommon.CodeValidationError,
			"invalid user type",
			nil,
		)

	case errors.Is(err, domain.ErrInvalidPassword):
		httpcommon.RespondError(
			c,
			http.StatusBadRequest,
			httpcommon.CodeValidationError,
			"password must be at least 8 characters",
			nil,
		)

	case errors.Is(err, domain.ErrEmailDuplicated):
		httpcommon.RespondError(
			c,
			http.StatusConflict,
			httpcommon.CodeEmailDuplicated,
			"email already exists",
			nil,
		)

	default:
		httpcommon.RespondError(
			c,
			http.StatusInternalServerError,
			httpcommon.CodeInternalServerError,
			"internal server error",
			map[string]any{"error": err.Error()},
		)
	}
}
