//go:build e2e

package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/database"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/models"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/router"
	"github.com/gin-gonic/gin"
)



func TestCategoryCRUD(t *testing.T) {
	os.Setenv("DB_DSN", "root:root@tcp(127.0.0.1:3306)/category_management?charset=utf8mb4&parseTime=True&loc=Local")
	database.StartDB()
	_, engine, _ := setupRouter()
	ts := httptest.NewServer(engine)
	defer ts.Close()

	// Create a new category
	createCategoryRequest := `{"name": "Test Category", "account_id": "123", "description": "Test Description", "color": "red", "budget": 100.0, "current": 50.0}`
	resp, err := http.Post(ts.URL + "/categories/123", "application/json", bytes.NewBufferString(createCategoryRequest))
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
	}

	var category models.Category
	err = json.NewDecoder(resp.Body).Decode(&category)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	id := category.ID


	// Get the category by ID
	resp, err = http.Get(ts.URL + "/categories/123/" + strconv.Itoa(int(id)))
	if err != nil {
		t.Fatalf("Failed to get category: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	var retrievedCategory models.Category
	err = json.NewDecoder(resp.Body).Decode(&retrievedCategory)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	if retrievedCategory.ID != id {
		t.Errorf("Expected category ID %v, but got %v", id, retrievedCategory.ID)
	}
	if !reflect.DeepEqual(category, retrievedCategory) {
		t.Errorf("Expected category %+v, but got %+v", category, retrievedCategory)
	}

	// Update the category
	updateCategoryRequest := `{"name": "Updated Category", "color": "blue", "budget": 200.0, "current": 100.0}`
	req, err := http.NewRequest(http.MethodPatch, ts.URL+"/categories/123/"+strconv.Itoa(int(id)), bytes.NewBufferString(updateCategoryRequest))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to update category: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	var updatedCategory models.Category
	err = json.NewDecoder(resp.Body).Decode(&updatedCategory)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if updatedCategory.ID != id {
		t.Errorf("Expected category ID %v, but got %v", id, updatedCategory.ID)
	}

	if !updatedCategory.UpdatedAt.After(category.UpdatedAt) {
		t.Errorf("Expected updatedAt to be after %s, but got %s", category.UpdatedAt, updatedCategory.UpdatedAt)
	}

	if updatedCategory.Name != "Updated Category" {
		t.Errorf("Expected name to be 'Updated Category', but got '%s'", updatedCategory.Name)
	}

	if updatedCategory.Color != "blue" {
		t.Errorf("Expected color to be 'blue', but got '%s'", updatedCategory.Color)
	}

	if updatedCategory.Budget != 200.0 {
		t.Errorf("Expected budget to be 200.0, but got %f", updatedCategory.Budget)
	}

	if updatedCategory.Current != category.Current {
		t.Errorf("Expected current to be 50.0, but got %f", updatedCategory.Current)
	}

	if updatedCategory.Description != category.Description {
		t.Errorf("Expected description to be '%s', but got '%s'", category.Description, updatedCategory.Description)
	}

	if updatedCategory.AccountID != category.AccountID {
		t.Errorf("Expected accountID to be '%s', but got '%s'", category.AccountID, updatedCategory.AccountID)
	}

	// Get all categories
	resp, err = http.Get(ts.URL + "/categories/123")
	if err != nil {
		t.Fatalf("Failed to get categories: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	var categories []models.Category
	err = json.NewDecoder(resp.Body).Decode(&categories)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	if len(categories) != 1 {
		t.Errorf("Expected 1 category, but got %d", len(categories))
	}
	if categories[0].ID != id {
		t.Errorf("Expected category ID %v, but got %v", id, categories[0].ID)
	}
	if !reflect.DeepEqual(categories[0], updatedCategory) {
		t.Errorf("Expected category %+v, but got %+v", updatedCategory, categories[0])
	}

	// Delete the category
	req, err = http.NewRequest(http.MethodDelete, ts.URL+"/categories/123/"+strconv.Itoa(int(id)), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err = http.DefaultClient.Do(req)

	if err != nil {
		t.Fatalf("Failed to delete category: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code %d, but got %d", http.StatusNoContent, resp.StatusCode)
	}

	// Check if the category is deleted
	resp, err = http.Get(ts.URL + "/categories/123/" + strconv.Itoa(int(id)))
	if err != nil {
		t.Fatalf("Failed to get category: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func setupRouter() (*gin.Context, *gin.Engine, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(recorder)
	router.HandleRequests(engine)
	return ctx, engine, recorder
}
