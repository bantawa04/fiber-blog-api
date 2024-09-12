package dto

import (
	"my-fiber-project/model"
	"time"
)

type CategoryResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewCategoryResponse(category *model.Category) CategoryResponse {
	response := CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format(time.RFC3339),
		UpdatedAt: category.UpdatedAt.Format(time.RFC3339),
	}

	return response
}

func CategoryListResponse(categories []model.Category) []CategoryResponse {
	responses := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = NewCategoryResponse(&category)
	}
	return responses
}
