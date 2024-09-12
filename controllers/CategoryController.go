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

type CategoryController struct {
	Validate     *validator.Validate
	ResponseUtil *utils.ResponseUtil
}

func NewCategoryController() *CategoryController {
	return &CategoryController{
		Validate:     validator.New(),
		ResponseUtil: &utils.ResponseUtil{},
	}
}

func (controller *CategoryController) Index(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search")

	offset := (page - 1) * limit

	query := config.Database.Model(&model.Category{})

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var categories []model.Category

	result := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&categories)

	if result.Error != nil {
		return controller.ResponseUtil.SendError(c, "Could not fetch posts", result.Error.Error(), fiber.StatusInternalServerError)
	}

	postResponses := dto.CategoryListResponse(categories)

	return controller.ResponseUtil.SendPagination(c, postResponses, total, page, limit, "Posts fetched successfully")
}

func (controller *CategoryController) Store(c *fiber.Ctx) error {
	category := new(model.Category)

	if err := c.BodyParser(category); err != nil {
		return controller.ResponseUtil.SendError(c, "Cannot parse JSON", err.Error(), fiber.StatusBadRequest)
	}

	//request validation
	if err := controller.Validate.Struct(category); err != nil {
		// Convert validation errors to a string
		var errorsStr string
		for _, err := range err.(validator.ValidationErrors) {
			errorsStr += err.Field() + ": " + err.Tag() + "; "
		}
		return controller.ResponseUtil.SendError(c, "Validation failed", errorsStr, fiber.StatusBadRequest)
	}
	result := config.Database.Create(category)

	if result.Error != nil {
		return controller.ResponseUtil.SendError(c, "Could not create category", result.Error.Error(), fiber.StatusInternalServerError)
	}

	categoryResponse := dto.NewCategoryResponse(category)

	return controller.ResponseUtil.SendResponse(c, categoryResponse, "Category created successfully")
}

func (controller *CategoryController) Show(c *fiber.Ctx) error {
	id := c.Params("id")

	var category model.Category

	result := config.Database.First(&category, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return controller.ResponseUtil.SendError(c, "Category not found", result.Error.Error(), fiber.StatusNotFound)
		}
		return controller.ResponseUtil.SendError(c, "Error fetching category", result.Error.Error(), fiber.StatusNotFound)
	}
	catResponse := dto.NewCategoryResponse(&category)
	return controller.ResponseUtil.SendResponse(c, catResponse, "Category fetched successfully")
}

func (controller *CategoryController) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var existingCat model.Category

	if err := config.Database.First(&existingCat, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return controller.ResponseUtil.SendError(c, "Category not found", err.Error(), fiber.StatusNotFound)
		}
		return controller.ResponseUtil.SendError(c, "Error fetching category", err.Error(), fiber.StatusNotFound)
	}

	updateData := new(model.Category)

	if err := c.BodyParser(updateData); err != nil {
		return controller.ResponseUtil.SendError(c, "Cannot parse JSON", err.Error(), fiber.StatusBadRequest)
	}

	if err := controller.Validate.Struct(updateData); err != nil {
		var errorsStr string
		for _, err := range err.(validator.ValidationErrors) {
			errorsStr += err.Field() + ": " + err.Tag() + "; "
		}
		return controller.ResponseUtil.SendError(c, "Validation failed", errorsStr, fiber.StatusBadRequest)
	}

	if err := config.Database.Model(&existingCat).Updates(updateData).Error; err != nil {
		return controller.ResponseUtil.SendError(c, "Could not update category", err.Error(), fiber.StatusInternalServerError)
	}

	if err := config.Database.First(&existingCat, id).Error; err != nil {
		return controller.ResponseUtil.SendError(c, "Error fetching updated category", err.Error(), fiber.StatusInternalServerError)
	}
	postResponse := dto.NewCategoryResponse(&existingCat)
	return controller.ResponseUtil.SendResponse(c, postResponse, "Category updated successfully.")
}

func (controller *CategoryController) Destroy(c *fiber.Ctx) error {
	id := c.Params("id")

	var category model.Category

	// First, fetch the category
	if err := config.Database.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return controller.ResponseUtil.SendError(c, "Category not found", "The specified category does not exist", fiber.StatusNotFound)
		}
		return controller.ResponseUtil.SendError(c, "Error fetching category", err.Error(), fiber.StatusInternalServerError)
	}

	// Now delete the category
	result := config.Database.Delete(&category)

	if result.Error != nil {
		return controller.ResponseUtil.SendError(c, "Could not delete category", result.Error.Error(), fiber.StatusInternalServerError)
	}

	if result.RowsAffected == 0 {
		return controller.ResponseUtil.SendError(c, "No category deleted", "The operation did not delete any category", fiber.StatusNotFound)
	}

	return controller.ResponseUtil.SendSuccess(c, "Category deleted successfully.")
}
