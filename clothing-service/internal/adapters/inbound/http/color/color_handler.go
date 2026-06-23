package httpcolor

import (
	httpcommon "clothing-service/internal/adapters/inbound/http/common"
	"clothing-service/internal/application/services"
	"clothing-service/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.ColorService
}

func NewHandler(service *services.ColorService) *Handler {
	return &Handler{
		service: service,
	}
}

type ColorResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	HexCode string `json:"hex_code"`
}

func (h *Handler) ColorDropdown(c *gin.Context) {
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

func toResponseList(colors []domain.Color) []ColorResponse {
	responses := make([]ColorResponse, 0, len(colors))
	for _, color := range colors {
		responses = append(responses, ColorResponse{
			ID:      color.ID,
			Name:    color.Name,
			HexCode: color.HexCode,
		})
	}

	return responses
}
