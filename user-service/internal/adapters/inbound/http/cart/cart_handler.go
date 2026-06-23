package httpcart

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	httpcommon "user-service/internal/adapters/inbound/http/common"
	"user-service/internal/application/ports"
	"user-service/internal/application/services"
	"user-service/internal/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cartService   *services.CartService
	clothesClient ports.ClothesClient
}

func NewHandler(cartService *services.CartService, clothesClient ports.ClothesClient) *Handler {
	return &Handler{
		cartService:   cartService,
		clothesClient: clothesClient,
	}
}

type AddCartItemRequest struct {
	ClothesID int64  `json:"clothes_id" binding:"required,min=1"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	Size      string `json:"size"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type CartResponse struct {
	ID     int64              `json:"id"`
	UserID int64              `json:"user_id"`
	Status string             `json:"status"`
	Items  []CartItemResponse `json:"items"`
}

type CartItemResponse struct {
	ID        int64            `json:"id"`
	ClothesID int64            `json:"clothes_id"`
	Clothes   *ClothesResponse `json:"clothes,omitempty"`
	Quantity  int              `json:"quantity"`
	Size      string           `json:"size"`
}

type ClothesResponse struct {
	ID      int64                  `json:"id"`
	Name    string                 `json:"name"`
	Price   float64                `json:"price"`
	ColorID *int64                 `json:"color_id"`
	Images  []ClothesImageResponse `json:"images"`
}

type ClothesImageResponse struct {
	ID        int64  `json:"id"`
	ClothesID int64  `json:"clothes_id"`
	ImageURL  string `json:"image_url"`
}

func (h *Handler) GetCart(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "unauthorized", nil)
		return
	}

	cart, err := h.cartService.GetCart(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	if err := h.enrichCartClothes(c.Request.Context(), cart); err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, toResponse(cart))
}

func (h *Handler) AddItem(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "unauthorized", nil)
		return
	}

	var req AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpcommon.RespondError(c, http.StatusBadRequest, httpcommon.CodeValidationError, "invalid request body", map[string]any{"error": err.Error()})
		return
	}

	cart, err := h.cartService.AddItem(c.Request.Context(), userID, req.ClothesID, req.Quantity, req.Size)
	if err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, toResponse(cart))
}

func (h *Handler) UpdateItem(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "unauthorized", nil)
		return
	}

	itemID, ok := itemIDParam(c)
	if !ok {
		httpcommon.RespondError(c, http.StatusBadRequest, httpcommon.CodeValidationError, "invalid item id", nil)
		return
	}

	var req UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpcommon.RespondError(c, http.StatusBadRequest, httpcommon.CodeValidationError, "invalid request body", map[string]any{"error": err.Error()})
		return
	}

	cart, err := h.cartService.UpdateItemQuantity(c.Request.Context(), userID, itemID, req.Quantity)
	if err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, toResponse(cart))
}

func (h *Handler) DeleteItem(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "unauthorized", nil)
		return
	}

	itemID, ok := itemIDParam(c)
	if !ok {
		httpcommon.RespondError(c, http.StatusBadRequest, httpcommon.CodeValidationError, "invalid item id", nil)
		return
	}

	cart, err := h.cartService.DeleteItem(c.Request.Context(), userID, itemID)
	if err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, toResponse(cart))
}

func (h *Handler) ClearCart(c *gin.Context) {
	userID, ok := authUserID(c)
	if !ok {
		httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "unauthorized", nil)
		return
	}

	if err := h.cartService.ClearCart(c.Request.Context(), userID); err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, toResponse(&domain.Cart{
		UserID: userID,
		Status: domain.CartStatusActive,
		Items:  []domain.CartItem{},
	}))
}

func authUserID(c *gin.Context) (int64, bool) {
	value, exists := c.Get("auth.user_id")
	if !exists {
		return 0, false
	}

	userID, ok := value.(int64)
	return userID, ok
}

func itemIDParam(c *gin.Context) (int64, bool) {
	itemID, err := strconv.ParseInt(c.Param("itemId"), 10, 64)
	return itemID, err == nil && itemID > 0
}

func toResponse(cart *domain.Cart) CartResponse {
	items := make([]CartItemResponse, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, CartItemResponse{
			ID:        item.ID,
			ClothesID: item.ClothesID,
			Clothes:   toClothesResponse(item.Clothes),
			Quantity:  item.Quantity,
			Size:      item.Size,
		})
	}

	return CartResponse{
		ID:     cart.ID,
		UserID: cart.UserID,
		Status: string(cart.Status),
		Items:  items,
	}
}

func toClothesResponse(clothes *domain.Clothes) *ClothesResponse {
	if clothes == nil {
		return nil
	}

	images := make([]ClothesImageResponse, 0, len(clothes.Images))
	for _, image := range clothes.Images {
		images = append(images, ClothesImageResponse{
			ID:        image.ID,
			ClothesID: image.ClothesID,
			ImageURL:  image.ImageURL,
		})
	}

	return &ClothesResponse{
		ID:      clothes.ID,
		Name:    clothes.Name,
		Price:   clothes.Price,
		ColorID: clothes.ColorID,
		Images:  images,
	}
}

func (h *Handler) enrichCartClothes(ctx context.Context, cart *domain.Cart) error {
	if h.clothesClient == nil {
		return nil
	}

	clothesByID := make(map[int64]*domain.Clothes)
	for index := range cart.Items {
		item := &cart.Items[index]
		if clothes, ok := clothesByID[item.ClothesID]; ok {
			item.Clothes = clothes
			continue
		}

		clothes, err := h.clothesClient.GetByID(ctx, item.ClothesID)
		if err != nil {
			return err
		}

		clothesByID[item.ClothesID] = clothes
		item.Clothes = clothes
	}

	return nil
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrCartItemNotFound):
		httpcommon.RespondError(c, http.StatusNotFound, httpcommon.CodeCartItemNotFound, "cart item not found", nil)
	case errors.Is(err, domain.ErrInvalidCartItem):
		httpcommon.RespondError(c, http.StatusBadRequest, httpcommon.CodeValidationError, "invalid cart item", nil)
	default:
		httpcommon.RespondError(c, http.StatusInternalServerError, httpcommon.CodeInternalServerError, "internal server error", map[string]any{"error": err.Error()})
	}
}
