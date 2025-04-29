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

func (category *CategoryResponse) Marshal() serialization.JSONResponse {
	return serialization.JSONResponse{
		Status: "success",
		Data:   category,
	}
}

type UpdateCategoryModel struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Color       string  `json:"color"`
	Budget      float64 `json:"budget"`
}
