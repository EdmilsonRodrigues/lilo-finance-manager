//go:build unit

package models_test

import (
	"testing"

	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/models"
)

func TestCategoryModel(t *testing.T) {
	category := models.Category{
		AccountID:   "1234567890",
		Name:        "Test Category",
		Description: "This is a test category",
		Color:       "#FF0000",
		Budget:      1000.0,
		Current:     500.0,
	}

	if category.AccountID != "1234567890" {
		t.Errorf("Expected AccountID to be '1234567890', but got '%s'", category.AccountID)
	}

	if category.Name != "Test Category" {
		t.Errorf("Expected Name to be 'Test Category', but got '%s'", category.Name)
	}

	if category.Description != "This is a test category" {
		t.Errorf("Expected Description to be 'This is a test category', but got '%s'", category.Description)
	}

	if category.Color != "#FF0000" {
		t.Errorf("Expected Color to be '#FF0000', but got '%s'", category.Color)
	}

	if category.Budget != 1000.0 {
		t.Errorf("Expected Budget to be 1000.0, but got %f", category.Budget)
	}

	if category.Current != 500.0 {
		t.Errorf("Expected Current to be 500.0, but got %f", category.Current)
	}
}
