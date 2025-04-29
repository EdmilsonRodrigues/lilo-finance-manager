//go:build integration

package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/database"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/models"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/router"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
)


func TestRouter(t *testing.T) {
	db := database.OpenDBConnection(":memory:", sqlite.Open)

	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	database.MakeMigrations(db)
	database.DB = db

	_, engine, _ := setupRouter()
	ts := httptest.NewServer(engine)
	defer ts.Close()

	jsonContentType := []string{"application/json; charset=utf-8"}

	t.Run("should get all categories", func(t *testing.T) {
		// Create some categories in the database for the test.
		category1 := models.Category{Name: "Category1", AccountID: "1", Description: "Description1", Color: "Color1", Budget: 100, Current: 50}
		if err := db.Create(&category1).Error; err != nil {
			t.Fatalf("failed to create test category: %v", err)
		}

		category2 := models.Category{Name: "Category2", AccountID: "1", Description: "Description2", Color: "Color2", Budget: 200, Current: 100}
		if err := db.Create(&category2).Error; err != nil {
			t.Fatalf("failed to create test category: %v", err)
		}

		resp, err := http.Get(ts.URL + "/categories/1")
		if err != nil {
			t.Errorf("error getting categories: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != jsonContentType[0] {
			t.Errorf("expected content type %s, got %s", jsonContentType[0], resp.Header.Get("Content-Type"))
		}

		var categories []models.Category
		if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
			t.Errorf("error decoding response body: %v", err)
		}
		if len(categories) != 2 {
			t.Errorf("expected 2 categories, got %d", len(categories))
		}

		databaseCategories := []models.Category{}
		if err := db.Find(&databaseCategories).Error; err != nil {
			t.Errorf("error getting categories from database: %v", err)
		}
		if len(databaseCategories) != 2 {
			t.Errorf("expected 2 categories in database, got %d", len(databaseCategories))
		}

		normalizeDateTimes(&categories[0], &categories[1], &databaseCategories[0], &databaseCategories[1])

		if !reflect.DeepEqual(categories, databaseCategories) {
			t.Errorf("expected categories to be equal, got %v and %v", categories, databaseCategories)
		}
	})

	t.Run("should not get a category from another account", func(t *testing.T) {
		// Create a category in the database for the test.
		category := models.Category{Name: "CategoryX", AccountID: "1", Description: "DescriptionX", Color: "ColorX", Budget: 150, Current: 75}
		if err := db.Create(&category).Error; err != nil {
			t.Fatalf("failed to create test category: %v", err)
		}

		resp, err := http.Get(ts.URL + "/categories/2/" + strconv.Itoa(int(category.ID)))
		if err != nil {
			t.Errorf("error getting category by id: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
		}
	})

	t.Run("should get a category by id", func(t *testing.T) {
		// Create a category in the database for the test.
		createdCategory := models.Category{Name: "CategoryX", AccountID: "1", Description: "DescriptionX", Color: "ColorX", Budget: 150, Current: 75}
		if err := db.Create(&createdCategory).Error; err != nil {
			t.Fatalf("failed to create test category: %v", err)
		}

		resp, err := http.Get(ts.URL + "/categories/1/" + strconv.Itoa(int(createdCategory.ID)))
		if err != nil {
			t.Errorf("error getting category by id: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != jsonContentType[0] {
			t.Errorf("expected content type %s, got %s", jsonContentType[0], resp.Header.Get("Content-Type"))
		}

		var category models.Category
		if err := json.NewDecoder(resp.Body).Decode(&category); err != nil {
			t.Errorf("error decoding response body: %v", err)
		}

		normalizeDateTimes(&createdCategory, &category)

		if !reflect.DeepEqual(createdCategory, category) {
			t.Errorf("expected category %+v, got %+v", createdCategory, category)
		}
	})

	t.Run("should create a category", func(t *testing.T) {
		newCategory := models.Category{Name: "NewCategory", AccountID: "2", Description: "NewDescription", Color: "NewColor", Budget: 250, Current: 125}
		categoryJSON, err := json.Marshal(newCategory)
		if err != nil {
			t.Fatalf("failed to marshal new category: %v", err)
		}

		resp, err := http.Post(ts.URL+"/categories/2", "application/json", bytes.NewReader(categoryJSON))
		if err != nil {
			t.Errorf("error creating category: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != jsonContentType[0] {
			t.Errorf("expected content type %s, got %s", jsonContentType[0], resp.Header.Get("Content-Type"))
		}

		var createdCategory models.Category
		if err := json.NewDecoder(resp.Body).Decode(&createdCategory); err != nil {
			t.Errorf("error decoding response body: %v", err)
		}

		// Check if the category was created in the database.
		var dbCategory models.Category
		if err := db.First(&dbCategory, "id = ?", createdCategory.ID).Error; err != nil {
			t.Errorf("error querying database for created category: %v", err)
		}

		// Normalize time zones to UTC for comparison.
		normalizeDateTimes(&createdCategory, &dbCategory)

		if !reflect.DeepEqual(createdCategory, dbCategory) {
			t.Errorf("expected %+v, got %+v", dbCategory, createdCategory)
		}
	})

	t.Run("should update a category", func(t *testing.T) {
		// Create a category to update.
		categoryToUpdate := models.Category{Name: "OldName", AccountID: "3", Description: "OldDescription", Color: "OldColor", Budget: 300, Current: 150}
		if err := db.Create(&categoryToUpdate).Error; err != nil {
			t.Fatalf("failed to create category to update: %v", err)
		}

		updatedCategory := models.Category{Name: "NewName", Description: "NewDescription", Color: "NewColor", Budget: 400, Current: 200}
		updatedCategoryJSON, err := json.Marshal(updatedCategory)
		if err != nil {
			t.Fatalf("failed to marshal updated category: %v", err)
		}

		req, err := http.NewRequest(http.MethodPatch, ts.URL+"/categories/3/"+strconv.Itoa(int(categoryToUpdate.ID)), bytes.NewReader(updatedCategoryJSON))
		if err != nil {
			t.Fatalf("failed to create PATCH request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Errorf("error updating category: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != jsonContentType[0] {
			t.Errorf("expected content type %s, got %s", jsonContentType[0], resp.Header.Get("Content-Type"))
		}

		var updatedCategoryResponse models.Category
		if err := json.NewDecoder(resp.Body).Decode(&updatedCategoryResponse); err != nil {
			t.Errorf("error decoding response body: %v", err)
		}

		// Check the category in the database.
		var dbCategory models.Category
		if err := db.First(&dbCategory, "id = ?", categoryToUpdate.ID).Error; err != nil {
			t.Errorf("error querying database for updated category: %v", err)
		}

		// Normalize time zones to UTC for comparison.
		normalizeDateTimes(&updatedCategoryResponse, &dbCategory)

		//  update the fields, that were changed
		categoryToUpdate.Name = updatedCategory.Name
		categoryToUpdate.Description = updatedCategory.Description
		categoryToUpdate.Color = updatedCategory.Color
		categoryToUpdate.Budget = updatedCategory.Budget

		if !reflect.DeepEqual(updatedCategoryResponse, dbCategory) {
			t.Errorf("expected %+v, got %+v", dbCategory, updatedCategoryResponse)
		}

		if updatedCategoryResponse.Current != categoryToUpdate.Current {
			t.Errorf("expected Current to be unaltered, got %v", updatedCategoryResponse.Current)
		}

		if updatedCategoryResponse.AccountID != categoryToUpdate.AccountID {
			t.Errorf("expected AccountID to be unaltered, got %v", updatedCategoryResponse.AccountID)
		}
	})

	t.Run("should delete a category", func(t *testing.T) {
		// Create a category to delete.
		categoryToDelete := models.Category{Name: "ToDelete", AccountID: "4", Description: "ToDeleteDescription", Color: "ToDeleteColor", Budget: 450, Current: 225}
		if err := db.Create(&categoryToDelete).Error; err != nil {
			t.Fatalf("failed to create category to delete: %v", err)
		}

		req, err := http.NewRequest(http.MethodDelete, ts.URL+"/categories/4/"+strconv.Itoa(int(categoryToDelete.ID)), nil)
		if err != nil {
			t.Fatalf("failed to create DELETE request: %v", err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Errorf("error deleting category: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
		}

		// Check if the category is actually deleted (soft deleted).
		var deletedCategory models.Category
		if err := db.Unscoped().First(&deletedCategory, "id = ?", categoryToDelete.ID).Error; err != nil {
			t.Errorf("error querying database for deleted category: %v", err)
		}

		if deletedCategory.ID != categoryToDelete.ID {
			t.Errorf("expected deleted category id to be %d, got %d", categoryToDelete.ID, deletedCategory.ID)
		}

		if deletedCategory.Model.DeletedAt.Valid == false {
			t.Errorf("expected deleted category to be soft deleted, but it was not")
		}
	})
}

func normalizeDateTimes(categories ...*models.Category) {
	for _, category := range categories {
		category.Model.CreatedAt = category.Model.CreatedAt.UTC()
		category.Model.UpdatedAt = category.Model.UpdatedAt.UTC()
	}
}

func setupRouter() (*gin.Context, *gin.Engine, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(recorder)
	router.HandleRequests(engine)
	return ctx, engine, recorder
}
