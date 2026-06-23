package httpcustomer

import (
	"errors"
	"net/http"
	"time"

	httpcommon "user-service/internal/adapters/inbound/http/common"
	"user-service/internal/application/services"
	"user-service/internal/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	customerService *services.CustomerService
}

func NewHandler(customerService *services.CustomerService) *Handler {
	return &Handler{customerService: customerService}
}

type CustomerProfileResponse struct {
	UserID         int64   `json:"user_id"`
	FirstName      *string `json:"first_name"`
	LastName       *string `json:"last_name"`
	Phone          *string `json:"phone"`
	DateOfBirth    *string `json:"date_of_birth"`
	MarketingOptIn bool    `json:"marketing_opt_in"`
}

type UpdateCustomerProfileRequest struct {
	FirstName      *string `json:"first_name"`
	LastName       *string `json:"last_name"`
	Phone          *string `json:"phone"`
	DateOfBirth    *string `json:"date_of_birth"`
	MarketingOptIn bool    `json:"marketing_opt_in"`
}

func (h *Handler) GetMe(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "unauthorized", nil)
		return
	}

	profile, err := h.customerService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, toResponse(profile))
}

func (h *Handler) UpdateMe(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "unauthorized", nil)
		return
	}

	var req UpdateCustomerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpcommon.RespondError(c, http.StatusBadRequest, httpcommon.CodeValidationError, "invalid request body", map[string]any{"error": err.Error()})
		return
	}

	dateOfBirth, err := parseDate(req.DateOfBirth)
	if err != nil {
		httpcommon.RespondError(c, http.StatusBadRequest, httpcommon.CodeValidationError, "invalid date_of_birth", nil)
		return
	}

	profile := &domain.CustomerProfile{
		UserID:         userID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Phone:          req.Phone,
		DateOfBirth:    dateOfBirth,
		MarketingOptIn: req.MarketingOptIn,
	}

	if err := h.customerService.UpdateProfile(c.Request.Context(), profile); err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, toResponse(profile))
}

func authUserID(c *gin.Context) (int64, bool) {
	value, exists := c.Get("auth.user_id")
	if !exists {
		return 0, false
	}

	userID, ok := value.(int64)
	return userID, ok
}

func parseDate(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	parsed, err := time.Parse("2006-01-02", *value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func toResponse(profile *domain.CustomerProfile) CustomerProfileResponse {
	var dateOfBirth *string
	if profile.DateOfBirth != nil {
		formatted := profile.DateOfBirth.Format("2006-01-02")
		dateOfBirth = &formatted
	}

	return CustomerProfileResponse{
		UserID:         profile.UserID,
		FirstName:      profile.FirstName,
		LastName:       profile.LastName,
		Phone:          profile.Phone,
		DateOfBirth:    dateOfBirth,
		MarketingOptIn: profile.MarketingOptIn,
	}
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrCustomerProfileNotFound):
		httpcommon.RespondError(c, http.StatusNotFound, httpcommon.CodeProfileNotFound, "customer profile not found", nil)
	default:
		httpcommon.RespondError(c, http.StatusInternalServerError, httpcommon.CodeInternalServerError, "internal server error", map[string]any{"error": err.Error()})
	}
}
