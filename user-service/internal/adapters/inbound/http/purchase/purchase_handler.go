package httppurchase

import (
	"errors"
	"net/http"

	httpcommon "user-service/internal/adapters/inbound/http/common"
	"user-service/internal/application/services"
	"user-service/internal/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	purchaseService *services.PurchaseService
}

func NewHandler(purchaseService *services.PurchaseService) *Handler {
	return &Handler{purchaseService: purchaseService}
}

type CreatePurchaseRequest struct {
	Items []CreatePurchaseItemRequest `json:"items" binding:"required,min=1,dive"`
}

type CreatePurchaseItemRequest struct {
	ClothesID int64  `json:"clothes_id" binding:"required,min=1"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	Size      string `json:"size"`
}

type PurchaseResponse struct {
	ID          int64                  `json:"id"`
	UserID      int64                  `json:"user_id"`
	TotalAmount float64                `json:"total_amount"`
	Items       []PurchaseItemResponse `json:"items"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

type PurchaseItemResponse struct {
	ID              int64   `json:"id"`
	PurchaseID      int64   `json:"purchase_id"`
	ClothesID       int64   `json:"clothes_id"`
	ProductName     string  `json:"product_name"`
	ProductImageURL *string `json:"product_image_url"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	Size            string  `json:"size"`
}

func (h *Handler) CreatePurchase(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "unauthorized", nil)
		return
	}

	var req CreatePurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpcommon.RespondError(c, http.StatusBadRequest, httpcommon.CodeValidationError, "invalid request body", map[string]any{"error": err.Error()})
		return
	}

	items := make([]domain.PurchaseItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.PurchaseItem{
			ClothesID: item.ClothesID,
			Quantity:  item.Quantity,
			Size:      item.Size,
		})
	}

	purchase, err := h.purchaseService.CreatePurchase(c.Request.Context(), userID, items)
	if err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusCreated, toResponse(purchase))
}

func (h *Handler) ListPurchases(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "unauthorized", nil)
		return
	}

	purchases, err := h.purchaseService.ListPurchases(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	response := make([]PurchaseResponse, 0, len(purchases))
	for _, purchase := range purchases {
		purchase := purchase
		response = append(response, toResponse(&purchase))
	}

	httpcommon.RespondSuccess(c, http.StatusOK, response)
}

func authUserID(c *gin.Context) (int64, bool) {
	value, exists := c.Get("auth.user_id")
	if !exists {
		return 0, false
	}

	userID, ok := value.(int64)
	return userID, ok
}

func toResponse(purchase *domain.Purchase) PurchaseResponse {
	items := make([]PurchaseItemResponse, 0, len(purchase.Items))
	for _, item := range purchase.Items {
		items = append(items, PurchaseItemResponse{
			ID:              item.ID,
			PurchaseID:      item.PurchaseID,
			ClothesID:       item.ClothesID,
			ProductName:     item.ProductName,
			ProductImageURL: item.ProductImageURL,
			Quantity:        item.Quantity,
			UnitPrice:       item.UnitPrice,
			Size:            item.Size,
		})
	}

	return PurchaseResponse{
		ID:          purchase.ID,
		UserID:      purchase.UserID,
		TotalAmount: purchase.TotalAmount,
		Items:       items,
		CreatedAt:   purchase.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   purchase.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidPurchaseItem):
		httpcommon.RespondError(c, http.StatusBadRequest, httpcommon.CodeValidationError, "invalid purchase item", nil)
	default:
		httpcommon.RespondError(c, http.StatusInternalServerError, httpcommon.CodeInternalServerError, "internal server error", map[string]any{"error": err.Error()})
	}
}
