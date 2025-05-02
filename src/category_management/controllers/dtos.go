package controllers

import (
	"fmt"

	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/models"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/common_utils/go/serialization"
)

type CategoryResponse struct {
	ID          uint    `json:"id"`
	AccountID   string  `json:"account_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Color       string  `json:"color"`
	Budget      float64 `json:"budget"`
	Current     float64 `json:"current"`
}

// BindModel binds the given model to this CategoryResponse.
// The model must be a pointer to models.Category.
// If the model is not valid, an error is returned.
//
// Parameters:
//   - model (models.Category): the model to bind to this CategoryResponse
//
// Returns:
//   - error: an error if the model is not valid
//
// Example:
//
//	model := ctx.BindJSON(&category)
//	category := CategoryResponse{}
//	err := category.BindModel(model)
//	if err != nil {
//	  ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//	  return
//	}
func (category *CategoryResponse) BindModel(model interface{}) error {
	cat, ok := model.(models.Category)
	if !ok {
		return fmt.Errorf("invalid model type: %T, expected: *models.Category", model)
	}
	category.ID = cat.ID
	category.AccountID = cat.AccountID
	category.Name = cat.Name
	category.Description = cat.Description
	category.Color = cat.Color
	category.Budget = cat.Budget
	category.Current = cat.Current
	return nil
}

func (category *CategoryResponse) Marshal(fields []string) (serialization.JSONResponse, error) {
	data, err := serialization.FilterSerializerFields(category, fields)
	if err != nil {
		return serialization.JSONResponse{}, err
	}
	return serialization.JSONResponse{
		Status: "success",
		Data:   data,
	}, nil
}

type UpdateCategoryModel struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Color       string  `json:"color"`
	Budget      float64 `json:"budget"`
}
