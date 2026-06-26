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
	UserID         int64                    `json:"user_id"`
	FirstName      *string                  `json:"first_name"`
	LastName       *string                  `json:"last_name"`
	Phone          *string                  `json:"phone"`
	DateOfBirth    *string                  `json:"date_of_birth"`
	MarketingOptIn bool                     `json:"marketing_opt_in"`
	Address        *CustomerAddressResponse `json:"address"`
}

type CustomerAddressResponse struct {
	ID            int64   `json:"id"`
	RecipientName string  `json:"recipient_name"`
	Phone         string  `json:"phone"`
	AddressLine1  string  `json:"address_line1"`
	AddressLine2  *string `json:"address_line2"`
	Subdistrict   *string `json:"subdistrict"`
	District      string  `json:"district"`
	Province      string  `json:"province"`
	PostalCode    string  `json:"postal_code"`
	CountryCode   string  `json:"country_code"`
	IsDefault     bool    `json:"is_default"`
}

type UpdateCustomerProfileRequest struct {
	FirstName      *string                       `json:"first_name"`
	LastName       *string                       `json:"last_name"`
	Phone          *string                       `json:"phone"`
	DateOfBirth    *string                       `json:"date_of_birth"`
	MarketingOptIn bool                          `json:"marketing_opt_in"`
	Address        *UpdateCustomerAddressRequest `json:"address"`
}

type UpdateCustomerAddressRequest struct {
	RecipientName string  `json:"recipient_name"`
	Phone         string  `json:"phone"`
	AddressLine1  string  `json:"address_line1"`
	AddressLine2  *string `json:"address_line2"`
	Subdistrict   *string `json:"subdistrict"`
	District      string  `json:"district"`
	Province      string  `json:"province"`
	PostalCode    string  `json:"postal_code"`
	CountryCode   string  `json:"country_code"`
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

	address, err := h.customerService.GetDefaultAddress(c.Request.Context(), userID)
	if err != nil && !errors.Is(err, domain.ErrCustomerAddressNotFound) {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, toResponse(profile, address))
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

	var address *domain.CustomerAddress
	if req.Address != nil {
		nextAddress, err := requestToAddress(userID, req.Address)
		if err != nil {
			httpcommon.RespondError(c, http.StatusBadRequest, httpcommon.CodeValidationError, err.Error(), nil)
			return
		}

		if err := h.customerService.UpdateDefaultAddress(c.Request.Context(), nextAddress); err != nil {
			respondError(c, err)
			return
		}
		address, _ = h.customerService.GetDefaultAddress(c.Request.Context(), userID)
	} else {
		address, _ = h.customerService.GetDefaultAddress(c.Request.Context(), userID)
	}

	httpcommon.RespondSuccess(c, http.StatusOK, toResponse(profile, address))
}

func (h *Handler) GetDefaultAddress(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "unauthorized", nil)
		return
	}

	address, err := h.customerService.GetDefaultAddress(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, toAddressResponse(address))
}

func (h *Handler) UpdateDefaultAddress(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "unauthorized", nil)
		return
	}

	var req UpdateCustomerAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpcommon.RespondError(c, http.StatusBadRequest, httpcommon.CodeValidationError, "invalid request body", map[string]any{"error": err.Error()})
		return
	}

	address, err := requestToAddress(userID, &req)
	if err != nil {
		httpcommon.RespondError(c, http.StatusBadRequest, httpcommon.CodeValidationError, err.Error(), nil)
		return
	}

	if err := h.customerService.UpdateDefaultAddress(c.Request.Context(), address); err != nil {
		respondError(c, err)
		return
	}

	address, err = h.customerService.GetDefaultAddress(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, toAddressResponse(address))
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

func toResponse(profile *domain.CustomerProfile, address *domain.CustomerAddress) CustomerProfileResponse {
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
		Address:        toAddressResponse(address),
	}
}

func toAddressResponse(address *domain.CustomerAddress) *CustomerAddressResponse {
	if address == nil {
		return nil
	}

	return &CustomerAddressResponse{
		ID:            address.ID,
		RecipientName: address.RecipientName,
		Phone:         address.Phone,
		AddressLine1:  address.AddressLine1,
		AddressLine2:  address.AddressLine2,
		Subdistrict:   address.Subdistrict,
		District:      address.District,
		Province:      address.Province,
		PostalCode:    address.PostalCode,
		CountryCode:   address.CountryCode,
		IsDefault:     address.IsDefault,
	}
}

func requestToAddress(userID int64, req *UpdateCustomerAddressRequest) (*domain.CustomerAddress, error) {
	countryCode := req.CountryCode
	if countryCode == "" {
		countryCode = "TH"
	}

	if req.RecipientName == "" {
		return nil, errors.New("recipient_name is required")
	}
	if req.Phone == "" {
		return nil, errors.New("address phone is required")
	}
	if req.AddressLine1 == "" {
		return nil, errors.New("address_line1 is required")
	}
	if req.District == "" {
		return nil, errors.New("district is required")
	}
	if req.Province == "" {
		return nil, errors.New("province is required")
	}
	if req.PostalCode == "" {
		return nil, errors.New("postal_code is required")
	}
	if len(countryCode) != 2 {
		return nil, errors.New("country_code must be 2 characters")
	}

	return &domain.CustomerAddress{
		UserID:        userID,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		AddressLine1:  req.AddressLine1,
		AddressLine2:  req.AddressLine2,
		Subdistrict:   req.Subdistrict,
		District:      req.District,
		Province:      req.Province,
		PostalCode:    req.PostalCode,
		CountryCode:   countryCode,
		IsDefault:     true,
	}, nil
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrCustomerProfileNotFound):
		httpcommon.RespondError(c, http.StatusNotFound, httpcommon.CodeProfileNotFound, "customer profile not found", nil)
	case errors.Is(err, domain.ErrCustomerAddressNotFound):
		httpcommon.RespondError(c, http.StatusNotFound, httpcommon.CodeProfileNotFound, "customer address not found", nil)
	case errors.Is(err, domain.ErrCustomerPhoneDuplicated):
		httpcommon.RespondError(c, http.StatusConflict, httpcommon.CodeDuplicatePhone, "phone number already exists", nil)
	default:
		httpcommon.RespondError(c, http.StatusInternalServerError, httpcommon.CodeInternalServerError, "internal server error", map[string]any{"error": err.Error()})
	}
}
