package httpcommon

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 10
	MaxPerPage     = 100
)

type PaginationRequest struct {
	Page    int
	PerPage int
	Limit   int
	Offset  int
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"perPage"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse[T any] struct {
	Items      []T            `json:"items"`
	Pagination PaginationMeta `json:"pagination"`
}

func ParsePagination(c *gin.Context) PaginationRequest {
	page := parsePositiveInt(c.Query("page"), DefaultPage)
	perPage := parsePositiveInt(c.Query("perPage"), DefaultPerPage)

	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	offset := (page - 1) * perPage

	return PaginationRequest{
		Page:    page,
		PerPage: perPage,
		Limit:   perPage,
		Offset:  offset,
	}
}

func NewPaginatedResponse[T any](
	items []T,
	total int64,
	page int,
	perPage int,
) PaginatedResponse[T] {
	totalPages := 0
	if perPage > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(perPage)))
	}

	return PaginatedResponse[T]{
		Items: items,
		Pagination: PaginationMeta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}

func parsePositiveInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return defaultValue
	}

	return n
}
