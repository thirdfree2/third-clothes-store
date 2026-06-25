package httpclothes

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	httpcommon "clothing-service/internal/adapters/inbound/http/common"
	"clothing-service/internal/application/ports"
	"clothing-service/internal/application/services"
	"clothing-service/internal/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service            *services.ClothesService
	imageService       *services.ClothesImageService
	minIOPublicBaseURL string
}

func NewHandler(service *services.ClothesService, imageService *services.ClothesImageService, minIOPublicBaseURL string) *Handler {
	return &Handler{
		service:            service,
		imageService:       imageService,
		minIOPublicBaseURL: minIOPublicBaseURL,
	}
}

type CreateClothesRequest struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"min=0"`
	ColorID     *int64  `json:"color_id"`
	CategoryIDs []int64 `json:"category_ids"`
}

type UpdateClothesRequest struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"min=0"`
	ColorID     *int64  `json:"color_id"`
	CategoryIDs []int64 `json:"category_ids"`
}

type ColorResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	HexCode string `json:"hex_code"`
}

type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ClothesResponse struct {
	ID         int64                  `json:"id"`
	Name       string                 `json:"name"`
	Price      float64                `json:"price"`
	ColorID    *int64                 `json:"color_id"`
	Color      *ColorResponse         `json:"color,omitempty"`
	Categories []CategoryResponse     `json:"categories"`
	Images     []ClothesImageResponse `json:"images"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

type ClothesImageResponse struct {
	ID               int64     `json:"id"`
	ClothesID        int64     `json:"clothes_id"`
	Bucket           string    `json:"bucket"`
	ObjectKey        string    `json:"object_key"`
	ImageURL         string    `json:"image_url"`
	OriginalFilename *string   `json:"original_filename,omitempty"`
	ContentType      string    `json:"content_type"`
	SizeBytes        int64     `json:"size_bytes"`
	CreatedAt        time.Time `json:"created_at"`
}

func toImageResponse(image *domain.ClothesImage, minIOPublicBaseURL string) ClothesImageResponse {
	return ClothesImageResponse{
		ID:               image.ID,
		ClothesID:        image.ClothesID,
		Bucket:           image.Bucket,
		ObjectKey:        image.ObjectKey,
		ImageURL:         buildImageURL(minIOPublicBaseURL, image.Bucket, image.ObjectKey),
		OriginalFilename: image.OriginalFilename,
		ContentType:      image.ContentType,
		SizeBytes:        image.SizeBytes,
		CreatedAt:        image.CreatedAt,
	}
}

func toImageResponseList(images []domain.ClothesImage, minIOPublicBaseURL string) []ClothesImageResponse {
	responses := make([]ClothesImageResponse, 0, len(images))

	for _, image := range images {
		item := image
		responses = append(responses, toImageResponse(&item, minIOPublicBaseURL))
	}

	return responses
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateClothesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpcommon.RespondError(
			c,
			http.StatusBadRequest,
			httpcommon.CodeValidationError,
			"invalid request body",
			map[string]any{
				"error": err.Error(),
			},
		)
		return
	}

	clothes, err := h.service.Create(c.Request.Context(), req.Name, req.Price, req.ColorID, req.CategoryIDs)
	if err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusCreated, toResponse(clothes, h.minIOPublicBaseURL))
}

func (h *Handler) FindByID(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		respondError(c, domain.ErrInvalidClothesID)
		return
	}

	clothes, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, toResponse(clothes, h.minIOPublicBaseURL))
}

func (h *Handler) List(c *gin.Context) {
	p := httpcommon.ParsePagination(c)
	categoryID, err := parseOptionalID(c.Query("category_id"))
	if err != nil {
		respondError(c, domain.ErrInvalidCategoryID)
		return
	}

	clothesList, total, err := h.service.List(
		c.Request.Context(),
		ports.Pagination{
			Limit:  p.Limit,
			Offset: p.Offset,
		},
		ports.ClothesListFilter{
			CategoryID: categoryID,
		},
	)

	if err != nil {
		respondError(c, err)
		return
	}

	responses := toResponseList(clothesList, h.minIOPublicBaseURL)
	payload := httpcommon.NewPaginatedResponse(
		responses,
		total,
		p.Page,
		p.PerPage,
	)

	httpcommon.RespondSuccess(c, http.StatusOK, payload)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		respondError(c, domain.ErrInvalidClothesID)
		return
	}

	var req UpdateClothesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpcommon.RespondError(
			c,
			http.StatusBadRequest,
			httpcommon.CodeValidationError,
			"invalid request body",
			map[string]any{
				"error": err.Error(),
			},
		)
		return
	}

	clothes, err := h.service.Update(c.Request.Context(), id, req.Name, req.Price, req.ColorID, req.CategoryIDs)
	if err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, toResponse(clothes, h.minIOPublicBaseURL))
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		respondError(c, domain.ErrInvalidClothesID)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, map[string]bool{
		"deleted": true,
	})
}

func (h *Handler) UploadImages(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		respondError(c, domain.ErrInvalidClothesID)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		httpcommon.RespondError(
			c,
			http.StatusBadRequest,
			httpcommon.CodeValidationError,
			"invalid multipart form",
			map[string]any{
				"error": err.Error(),
			},
		)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		httpcommon.RespondError(
			c,
			http.StatusBadRequest,
			httpcommon.CodeValidationError,
			"files are required",
			nil,
		)
		return
	}

	responses := make([]ClothesImageResponse, 0, len(files))

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			httpcommon.RespondError(
				c,
				http.StatusBadRequest,
				httpcommon.CodeValidationError,
				"cannot open uploaded file",
				map[string]any{
					"filename": fileHeader.Filename,
					"error":    err.Error(),
				},
			)
			return
		}

		contentType, err := detectContentType(file)
		if err != nil {
			file.Close()

			httpcommon.RespondError(
				c,
				http.StatusBadRequest,
				httpcommon.CodeValidationError,
				"cannot detect file content type",
				map[string]any{
					"filename": fileHeader.Filename,
					"error":    err.Error(),
				},
			)
			return
		}

		image, err := h.imageService.Upload(
			c.Request.Context(),
			id,
			fileHeader.Filename,
			contentType,
			fileHeader.Size,
			file,
		)

		file.Close()

		if err != nil {
			respondError(c, err)
			return
		}

		responses = append(responses, toImageResponse(image, h.minIOPublicBaseURL))
	}

	httpcommon.RespondSuccess(c, http.StatusCreated, responses)
}

func (h *Handler) DeleteImage(c *gin.Context) {
	clothesID, err := parseID(c.Param("id"))
	if err != nil {
		respondError(c, domain.ErrInvalidClothesID)
		return
	}

	imageID, err := parseID(c.Param("image_id"))
	if err != nil {
		httpcommon.RespondError(
			c,
			http.StatusBadRequest,
			httpcommon.CodeValidationError,
			"invalid image id",
			nil,
		)
		return
	}

	if err := h.imageService.Delete(c.Request.Context(), clothesID, imageID); err != nil {
		respondError(c, err)
		return
	}

	httpcommon.RespondSuccess(c, http.StatusOK, map[string]bool{
		"deleted": true,
	})
}

func parseID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	if id <= 0 {
		return 0, domain.ErrInvalidClothesID
	}

	return id, nil
}

func parseOptionalID(value string) (*int64, error) {
	if value == "" {
		return nil, nil
	}

	id, err := parseID(value)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func toResponse(clothes *domain.Clothes, minIOPublicBaseURL string) ClothesResponse {
	response := ClothesResponse{
		ID:         clothes.ID,
		Name:       clothes.Name,
		Price:      clothes.Price,
		ColorID:    clothes.ColorID,
		Categories: make([]CategoryResponse, 0, len(clothes.Categories)),
		Images:     make([]ClothesImageResponse, 0, len(clothes.Images)),
		CreatedAt:  clothes.CreatedAt,
		UpdatedAt:  clothes.UpdatedAt,
	}

	if clothes.Color != nil {
		response.Color = &ColorResponse{
			ID:      clothes.Color.ID,
			Name:    clothes.Color.Name,
			HexCode: clothes.Color.HexCode,
		}
	}

	for _, category := range clothes.Categories {
		response.Categories = append(response.Categories, CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	for _, image := range clothes.Images {
		item := image
		response.Images = append(
			response.Images,
			toImageResponse(&item, minIOPublicBaseURL),
		)
	}

	return response
}

func detectContentType(file io.ReadSeeker) (string, error) {
	buffer := make([]byte, 512)

	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return http.DetectContentType(buffer[:n]), nil
}

func toResponseList(clothesList []domain.Clothes, minIOPublicBaseURL string) []ClothesResponse {
	responses := make([]ClothesResponse, 0, len(clothesList))

	for _, clothes := range clothesList {
		item := clothes
		responses = append(responses, toResponse(&item, minIOPublicBaseURL))
	}

	return responses
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidClothesID):
		httpcommon.RespondError(
			c,
			http.StatusBadRequest,
			httpcommon.CodeInvalidClothesID,
			"invalid clothes id",
			nil,
		)

	case errors.Is(err, domain.ErrInvalidClothesName):
		httpcommon.RespondError(
			c,
			http.StatusBadRequest,
			httpcommon.CodeValidationError,
			"invalid clothes name",
			nil,
		)

	case errors.Is(err, domain.ErrInvalidClothesPrice):
		httpcommon.RespondError(
			c,
			http.StatusBadRequest,
			httpcommon.CodeValidationError,
			"invalid clothes price",
			nil,
		)

	case errors.Is(err, domain.ErrClothesNotFound):
		httpcommon.RespondError(
			c,
			http.StatusNotFound,
			httpcommon.CodeClothesNotFound,
			"clothes not found",
			nil,
		)
	case errors.Is(err, domain.ErrInvalidColorID):
		httpcommon.RespondError(
			c,
			http.StatusBadRequest,
			httpcommon.CodeInvalidColorID,
			"invalid color id",
			nil,
		)

	case errors.Is(err, domain.ErrInvalidCategoryID):
		httpcommon.RespondError(
			c,
			http.StatusBadRequest,
			httpcommon.CodeInvalidCategoryID,
			"invalid category id",
			nil,
		)
	case errors.Is(err, domain.ErrClothesNameAlreadyExists):
		httpcommon.RespondError(
			c,
			http.StatusConflict,
			httpcommon.CodeValidationError,
			"clothes name already exists",
			nil,
		)
	default:
		httpcommon.RespondError(
			c,
			http.StatusInternalServerError,
			httpcommon.CodeInternalServerError,
			"internal server error",
			map[string]any{
				"error": err.Error(),
			},
		)
	}
}

func buildImageURL(baseURL string, bucket string, objectKey string) string {
	baseURL = strings.TrimRight(baseURL, "/")
	objectKey = strings.TrimLeft(objectKey, "/")

	return baseURL + "/" + bucket + "/" + objectKey
}
