//go:build unit

package controllers_test

import (
	"reflect"
	"testing"
	"testing/quick"

	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/controllers"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/models"
	httpserialization "github.com/EdmilsonRodrigues/lilo-finance-manager/src/common_utils/go/serialization/http_serialization"
)

func TestCategoryResponseBindModel(t *testing.T) {
	assertion := func(id uint, budget, current float64, accountId, name, description, color string) bool {
		category := controllers.CategoryResponse{}
		model := models.Category{
			AccountID:   accountId,
			Name:        name,
			Description: description,
			Color:       color,
			Budget:      budget,
			Current:     current,
		}
		model.ID = id
		err := category.BindModel(model)
		if err != nil {
			t.Errorf("Error binding model: %v", err)
			return false
		}
		if category.ID != id {
			t.Errorf("Expected ID to be %d, got %d", id, category.ID)
			return false
		}
		if category.AccountID != accountId {
			t.Errorf("Expected AccountID to be %s, got %s", accountId, category.AccountID)
			return false
		}
		if category.Name != name {
			t.Errorf("Expected Name to be %s, got %s", name, category.Name)
			return false
		}
		if category.Description != description {
			t.Errorf("Expected Description to be %s, got %s", description, category.Description)
			return false
		}
		if category.Color != color {
			t.Errorf("Expected Color to be %s, got %s", color, category.Color)
			return false
		}
		if category.Budget != budget {
			t.Errorf("Expected Budget to be %f, got %f", budget, category.Budget)
			return false
		}
		if category.Current != current {
			t.Errorf("Expected Current to be %f, got %f", current, category.Current)
			return false
		}
		return true
	}
	if err := quick.Check(assertion, &quick.Config{
		MaxCount: 1000,
	}); err != nil {
		t.Error(err)
	}
}

func TestCategoryResponseMarshal(t *testing.T) {
	assertion := func(id uint, budget, current float64, accountId, name, description, color string) bool {
		category := controllers.CategoryResponse{
			ID:          id,
			AccountID:   accountId,
			Name:        name,
			Description: description,
			Color:       color,
			Budget:      budget,
			Current:     current,
		}
		response, err := category.Marshal([]string{})

		expected := httpserialization.JSONResponse{
			Status: "success",
			Data: map[string]interface{}{
				"id":          float64(id),
				"account_id":  accountId,
				"name":        name,
				"description": description,
				"color":       color,
				"budget":      budget,
				"current":     current,
			},
		}

		if err != nil {
			t.Errorf("Error marshaling: %v", err)
			return false
		}

		if !reflect.DeepEqual(response, expected) {
			t.Errorf("Expected %v, got %v", expected, response)
			return false
		}

		return true
	}
	if err := quick.Check(assertion, &quick.Config{
		MaxCount: 1000,
	}); err != nil {
		t.Error(err)
	}
}
