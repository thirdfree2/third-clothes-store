package httpcategory

import (
	httpcommon "clothing-service/internal/adapters/inbound/http/common"
	"clothing-service/internal/application/services"
	"clothing-service/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.CategoryService
}

func NewHandler(service *services.CategoryService) *Handler {
	return &Handler{
		service: service,
	}
}

type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CategoryDropdown(c *gin.Context) {
	colors, err := h.service.List(c.Request.Context())
	if err != nil {
		httpcommon.RespondError(
			c,
			http.StatusInternalServerError,
			httpcommon.CodeInternalServerError,
			"internal server error",
			map[string]any{
				"error": err.Error(),
			},
		)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, toResponseList(colors))
}

func toResponseList(colors []domain.Category) []CategoryResponse {
	responses := make([]CategoryResponse, 0, len(colors))
	for _, color := range colors {
		responses = append(responses, CategoryResponse{
			ID:   color.ID,
			Name: color.Name,
		})
	}

	return responses
}
