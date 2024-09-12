package dto

import (
	"my-fiber-project/model"
	"time"
)

type PostResponse struct {
	ID         uint           `json:"id"`
	Title      string         `json:"title"`
	Body       string         `json:"body"`
	CategoryID *uint          `json:"category_id"`
	Category   *CategoryBrief `json:"category,omitempty"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
}

type CategoryBrief struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewPostResponse(post *model.Post) PostResponse {
	response := PostResponse{
		ID:         post.ID,
		Title:      post.Title,
		Body:       post.Body,
		CategoryID: post.CategoryID,
		CreatedAt:  post.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  post.UpdatedAt.Format(time.RFC3339),
	}

	if post.Category != nil && post.Category.ID != 0 {
		response.Category = &CategoryBrief{
			ID:   post.Category.ID,
			Name: post.Category.Name,
		}
	}

	return response
}

func NewPostsResponse(posts []model.Post) []PostResponse {
	responses := make([]PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = NewPostResponse(&post)
	}
	return responses
}
