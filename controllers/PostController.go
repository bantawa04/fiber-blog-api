package controllers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"my-fiber-project/config"
	"my-fiber-project/dto"
	"my-fiber-project/model"
	"my-fiber-project/utils"
)

type PostController struct {
	Validate     *validator.Validate
	ResponseUtil *utils.ResponseUtil
}

func NewPostController() *PostController {
	return &PostController{
		Validate:     validator.New(),
		ResponseUtil: &utils.ResponseUtil{},
	}
}

func (pc *PostController) Index(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search")

	offset := (page - 1) * limit

	query := config.Database.Model(&model.Post{}).Preload("Category")

	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var posts []model.Post
	result := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts)

	if result.Error != nil {
		return pc.ResponseUtil.SendError(c, "Could not fetch posts", result.Error.Error(), fiber.StatusInternalServerError)
	}

	postResponses := dto.NewPostsResponse(posts)

	return pc.ResponseUtil.SendPagination(c, postResponses, total, page, limit, "Posts fetched successfully")
}

func (pc *PostController) Store(c *fiber.Ctx) error {
	post := new(model.Post)

	if err := c.BodyParser(post); err != nil {
		return pc.ResponseUtil.SendError(c, "Cannot parse JSON", err.Error(), fiber.StatusBadRequest)
	}

	// Validate the post
	if err := pc.Validate.Struct(post); err != nil {
		// Convert validation errors to a string
		var errorsStr string
		for _, err := range err.(validator.ValidationErrors) {
			errorsStr += err.Field() + ": " + err.Tag() + "; "
		}
		return pc.ResponseUtil.SendError(c, "Validation failed", errorsStr, fiber.StatusBadRequest)
	}

	result := config.Database.Create(post)
	if result.Error != nil {
		return pc.ResponseUtil.SendError(c, "Could not create post", result.Error.Error(), fiber.StatusInternalServerError)
	}
	postResponse := dto.NewPostResponse(post)
	return pc.ResponseUtil.SendResponse(c, postResponse, "Post created successfully.")
}

func (pc *PostController) Show(c *fiber.Ctx) error {
	id := c.Params("id")
	var post model.Post

	result := config.Database.Preload("Category").First(&post, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return pc.ResponseUtil.SendError(c, "Post not found", result.Error.Error(), fiber.StatusNotFound)
		}
		return pc.ResponseUtil.SendError(c, "Error fetching post", result.Error.Error(), fiber.StatusNotFound)
	}
	postResponse := dto.NewPostResponse(&post)
	return pc.ResponseUtil.SendResponse(c, postResponse, "Post fetched successfully.")
}

func (pc *PostController) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	// First, check if the post exists
	var existingPost model.Post
	if err := config.Database.First(&existingPost, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pc.ResponseUtil.SendError(c, "Post not found", err.Error(), fiber.StatusNotFound)
		}
		return pc.ResponseUtil.SendError(c, "Error fetching post", err.Error(), fiber.StatusInternalServerError)
	}

	// Parse the request body
	updateData := new(model.Post)
	if err := c.BodyParser(updateData); err != nil {
		return pc.ResponseUtil.SendError(c, "Cannot parse JSON", err.Error(), fiber.StatusBadRequest)
	}

	// Validate the update data
	if err := pc.Validate.Struct(updateData); err != nil {
		// Convert validation errors to a string
		var errorsStr string
		for _, err := range err.(validator.ValidationErrors) {
			errorsStr += err.Field() + ": " + err.Tag() + "; "
		}
		return pc.ResponseUtil.SendError(c, "Validation failed", errorsStr, fiber.StatusBadRequest)
	}

	// Update the post
	if err := config.Database.Model(&existingPost).Updates(updateData).Error; err != nil {
		return pc.ResponseUtil.SendError(c, "Could not update post", err.Error(), fiber.StatusInternalServerError)
	}

	// Fetch the updated post to return in the dto
	if err := config.Database.First(&existingPost, id).Error; err != nil {
		return pc.ResponseUtil.SendError(c, "Error fetching updated post", err.Error(), fiber.StatusInternalServerError)
	}
	postResponse := dto.NewPostResponse(&existingPost)
	return pc.ResponseUtil.SendResponse(c, postResponse, "Post updated successfully.")
}

func (pc *PostController) Destroy(c *fiber.Ctx) error {
	id := c.Params("id")

	var post model.Post

	if err := config.Database.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pc.ResponseUtil.SendError(c, "Post not found", "The specified post does not exist", fiber.StatusNotFound)
		}
		return pc.ResponseUtil.SendError(c, "Error checking post", err.Error(), fiber.StatusInternalServerError)
	}

	result := config.Database.Delete(&post)

	if result.Error != nil {
		return pc.ResponseUtil.SendError(c, "Could not delete post", result.Error.Error(), fiber.StatusInternalServerError)
	}

	return pc.ResponseUtil.SendSuccess(c, "Post deleted successfully.")
}
