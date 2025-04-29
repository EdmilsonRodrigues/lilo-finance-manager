//go:build unit

package controllers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"gorm.io/driver/sqlite"

	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/controllers"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/database"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/models"
	"github.com/gin-gonic/gin"
)

func TestCategoryController(t *testing.T) {
	notFoundResponse, err := json.Marshal(controllers.CategoryNotFoundResponse)
	if err != nil {
		t.Fatal(err)
	}

	dsn := ":memory:"
	opener := sqlite.Open

	db := database.OpenDBConnection(dsn, opener)

	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	database.MakeMigrations(db)
	database.DB = db



	sampleCategory := models.Category{
		AccountID:   "67c4782c-7085-44aa-8c89-b5fee34cf19c",
		Name:        "test",
		Description: "test",
		Color:       "test",
		Budget:      100,
		Current:     50,
	}
	db.Create(&sampleCategory)

	t.Run("should get all categories", func(t *testing.T) {
		ctx, body := getContext()
		controllers.GetCategories(ctx)

		if  status := ctx.Writer.Status(); status != 200 {
			t.Errorf("expected status code 200, got %d", status)
		}

		if !reflect.DeepEqual(ctx.Writer.Header()["Content-Type"], jsonContentType) {
			t.Errorf("expected content type %s, got %s", jsonContentType, ctx.Writer.Header()["Content-Type"])
		}

		var categories []models.Category
		json.Unmarshal(*body, &categories)
		expected := []models.Category{sampleCategory}
		if !reflect.DeepEqual(categories, expected) {
			t.Errorf("expected %v, got %v", expected, categories)
		}
	})

	t.Run("should get one category", func(t *testing.T) {
		ctx, body := getContext()
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.FormatUint(uint64(sampleCategory.ID), 10)}}
		controllers.GetCategory(ctx)

		if  status := ctx.Writer.Status(); status != 200 {
			t.Errorf("expected status code 200, got %d", status)
		}

		if !reflect.DeepEqual(ctx.Writer.Header()["Content-Type"], jsonContentType) {
			t.Errorf("expected content type %s, got %s", jsonContentType, ctx.Writer.Header()["Content-Type"])
		}

		var categories models.Category
		json.Unmarshal(*body, &categories)
		expected := sampleCategory
		if !reflect.DeepEqual(categories, expected) {
			t.Errorf("expected %v, got %v", expected, categories)
		}
	})

	t.Run("should return not found when getting a category that does not exist", func(t *testing.T) {
		ctx, body := getContext()
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1000"}}
		controllers.GetCategory(ctx)

		if  status := ctx.Writer.Status(); status != 404 {
			t.Errorf("expected status code 404, got %d", status)
		}

		if !reflect.DeepEqual(*body, notFoundResponse) {
			t.Errorf("expected %s, got %s", notFoundResponse, *body)
		}
	})

	t.Run("should create a category", func(t *testing.T) {
		ctx, body := getContext()
		newCategory := models.NewCategory("test", "test", "test", "test", 100, 50)
		categoryJSON, err := json.Marshal(newCategory)
		if err != nil {
			t.Errorf("error marshalling new category: %v", err)
		}
		ctx.Request = &http.Request{
			Body: io.NopCloser(bytes.NewReader(categoryJSON)),
		}
		controllers.CreateCategory(ctx)

		if status := ctx.Writer.Status(); status != 201 {
			t.Errorf("expected status code 201, got %d", status)
		}

		if !reflect.DeepEqual(ctx.Writer.Header()["Content-Type"], jsonContentType) {
			t.Errorf("expected content type %s, got %s", jsonContentType, ctx.Writer.Header()["Content-Type"])
		}

		var category models.Category
		json.Unmarshal(*body, &category)

		var dbCategory models.Category
		db.First(&dbCategory, "id = ?", category.ID)

		// Normalize time zones to UTC for comparison.
		category.Model.CreatedAt = category.Model.CreatedAt.UTC()
		category.Model.UpdatedAt = category.Model.UpdatedAt.UTC()
		dbCategory.Model.CreatedAt = dbCategory.Model.CreatedAt.UTC()
		dbCategory.Model.UpdatedAt = dbCategory.Model.UpdatedAt.UTC()

		if !reflect.DeepEqual(category, dbCategory) {
			t.Errorf("expected %+v, got %+v", dbCategory, category)
		}
	})

	t.Run("should update a category, but only fields that are present", func(t *testing.T) {
		ctx, body := getContext()
		updateBody := map[string]string{
			"name":        "test-abc",
		}
		categoryJSON, err := json.Marshal(updateBody)
		if err != nil {
			t.Errorf("error marshalling new category: %v", err)
		}
		ctx.Request = &http.Request{
			Body: io.NopCloser(bytes.NewReader(categoryJSON)),
		}
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.FormatUint(uint64(sampleCategory.ID), 10)}}
		controllers.UpdateCategory(ctx)

		if status := ctx.Writer.Status(); status != 200 {
			t.Errorf("expected status code 200, got %d", status)
		}

		if !reflect.DeepEqual(ctx.Writer.Header()["Content-Type"], jsonContentType) {
			t.Errorf("expected content type %s, got %s", jsonContentType, ctx.Writer.Header()["Content-Type"])
		}

		var category models.Category
		json.Unmarshal(*body, &category)

		var dbCategory models.Category
		db.First(&dbCategory, "id = ?", category.ID)

		// Normalize time zones to UTC for comparison.
		category.Model.CreatedAt = category.Model.CreatedAt.UTC()
		category.Model.UpdatedAt = category.Model.UpdatedAt.UTC()
		dbCategory.Model.CreatedAt = dbCategory.Model.CreatedAt.UTC()
		dbCategory.Model.UpdatedAt = dbCategory.Model.UpdatedAt.UTC()

		if !reflect.DeepEqual(category, dbCategory) {
			t.Errorf("expected %+v, got %+v", dbCategory, category)
		}

		if dbCategory.Name != "test-abc" {
			t.Errorf("expected name to be test-abc, got %s", dbCategory.Name)
		}

		if dbCategory.Description != sampleCategory.Description {
			t.Errorf("expected description to be %s, got %s", sampleCategory.Description, dbCategory.Description)
		}
	})

	t.Run("should return not found when updating a category that does not exist", func(t *testing.T) {
		ctx, body := getContext()
		updateBody := map[string]string{
			"name":        "test-abc",
		}
		categoryJSON, err := json.Marshal(updateBody)
		if err != nil {
			t.Errorf("error marshalling new category: %v", err)
		}
		ctx.Request = &http.Request{
			Body: io.NopCloser(bytes.NewReader(categoryJSON)),
		}
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1000"}}
		controllers.UpdateCategory(ctx)

		if status := ctx.Writer.Status(); status != 404 {
			t.Errorf("expected status code 404, got %d", status)
		}
		if !reflect.DeepEqual(*body, notFoundResponse) {
			t.Errorf("expected %s, got %s", notFoundResponse, *body)
		}
	})

	t.Run("should delete a category", func(t *testing.T) {
		ctx, body := getContext()
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.FormatUint(uint64(sampleCategory.ID), 10)}}
		controllers.DeleteCategory(ctx)

		if  status := ctx.Writer.Status(); status != 204 {
			t.Errorf("expected status code 204, got %d", status)
		}

		if !reflect.DeepEqual(*body, []byte{}) {
			t.Errorf("expected empty body, got %+v", *body)
		}

		var category models.Category
		db.First(&category, "id = ?", sampleCategory.ID)

		if category.ID != 0 {
			t.Errorf("expected category to be deleted, got %+v", category)
		}
	})

	t.Run("should return not found when deleting a category that does not exist", func(t *testing.T) {
		ctx, body := getContext()
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1000"}}
		controllers.DeleteCategory(ctx)

		if status := ctx.Writer.Status(); status != 404 {
			t.Errorf("expected status code 404, got %d", status)
		}

		if !reflect.DeepEqual(*body, notFoundResponse) {
			t.Errorf("expected %s, got %s", notFoundResponse, *body)
		}
	})
}

func getContext() (*gin.Context, *[]byte) {
	writer := &FakeWriter{
		HeadersMapping: make(http.Header),
		Body:           []byte{},
	}
	ctx, _ := gin.CreateTestContext(writer)
	return ctx, &writer.Body
}

type FakeWriter struct {
	gin.ResponseWriter

	StatusCode     int
	HeadersMapping http.Header
	Body           []byte
}

func (w *FakeWriter) WriteHeader(code int) {
	w.StatusCode = code
}

func (w *FakeWriter) Header() http.Header {
	return w.HeadersMapping
}

func (w *FakeWriter) Write(b []byte) (int, error) {
	w.Body = append(w.Body, b...)
	return len(b), nil
}

func (w *FakeWriter) Status() int {
	return w.StatusCode
}

var jsonContentType = []string{"application/json; charset=utf-8"}
