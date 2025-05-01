package controllers_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"testing/quick"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/controllers"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/database"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/models"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/common_utils/go/serialization"
	"github.com/gin-gonic/gin"
)

func TestGetCategories(t *testing.T) {
	db := openDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	t.Run("should get all categories", func(t *testing.T) {
		assertion := func(categories [10]models.Category, accountId string) bool {
			ctx, body := getContext()

			for _, category := range categories {
				category.AccountID = accountId
				db.Create(&category)
			}

			conditions := serialization.QueryConditions{
				"account_id": accountId,
			}

			ctx.Set("conditions", conditions)

			controllers.GetCategories(ctx)

			if status := ctx.Writer.Status(); status != http.StatusOK {
				t.Errorf("expected status code 200, got %d", status)
				return false
			}

			var categoriesResponse gin.H
			if err := json.Unmarshal(*body, &categories); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if categoriesResponse["status"] != "success" {
				t.Errorf("expected status to be success, got %s", categoriesResponse["status"])
				return false
			}

			data, ok := categoriesResponse["data"].(gin.H)
			if !ok {
				t.Errorf("expected data to be of type gin.H, got %T", categoriesResponse["data"])
				return false
			}

			items := data["items"].([]gin.H)

			if len(items) != len(categories) {
				t.Errorf("expected %d categories, got %d", len(categories), len(categoriesResponse.Data.Items))
				return false
			}

			if !reflect.DeepEqual(data["filters"], conditions) {
				t.Errorf("expected filters to be %v, got %v", conditions, categoriesResponse.Data.Filters)
				return false
			}

			if data["total_items"] != len(categories) {
				t.Errorf("expected total to be %d, got %d", len(categories), categoriesResponse.Data.Total)
				return false
			}

			if data["total_pages"] != 1 {
				t.Errorf("expected total to be %d, got %d", 1, categoriesResponse.Data.Total)
				return false
			}

			if categoriesResponse.Data.Page != 1 {
				t.Errorf("expected page to be %d, got %d", 1, categoriesResponse.Data.Page)
				return false
			}

			if categoriesResponse.Data.Size != 10 {
				t.Errorf("expected size to be %d, got %d", 10, categoriesResponse.Data.Size)
				return false
			}

			for _, item :=  range categoriesResponse.Data.Items {
				category, ok := item.(controllers.CategoryResponse)
				if !ok {
					t.Errorf("expected item to be of type CategoryResponse, got %T", item)
					return false
				}
				if category.ID == 0 {
					t.Errorf("expected id to be %d, got %d", 0, category.ID)
					return false
				}
				if category.AccountID != accountId {
					t.Errorf("expected account id to be %s, got %s", accountId, category.AccountID)
					return false
				}
				if category.Name != "test" {
					t.Errorf("expected name to be %s, got %s", "test", category.Name)
					return false
				}
				if category.Description != "test" {
					t.Errorf("expected description to be %s, got %s", "test", category.Description)
					return false
				}
				if category.Color != "test" {
					t.Errorf("expected color to be %s, got %s", "test", category.Color)
					return false
				}
				if category.Budget != 100 {
					t.Errorf("expected budget to be %f, got %f", 100, category.Budget)
					return false
				}
			}

			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 100,
		}); err != nil {
			t.Error(err)
		}
	})


}

// func TestCategoryController(t *testing.T) {




// 	sampleCategory := models.Category{
// 		AccountID:   "67c4782c-7085-44aa-8c89-b5fee34cf19c",
// 		Name:        "test",
// 		Description: "test",
// 		Color:       "test",
// 		Budget:      100,
// 		Current:     50,
// 	}
// 	db.Create(&sampleCategory)

// 	t.Run("should get all categories", func(t *testing.T) {
// 		ctx, body := getContext()
// 		controllers.GetCategories(ctx)

// 		if  status := ctx.Writer.Status(); status != 200 {
// 			t.Errorf("expected status code 200, got %d", status)
// 		}

// 		if !reflect.DeepEqual(ctx.Writer.Header()["Content-Type"], jsonContentType) {
// 			t.Errorf("expected content type %s, got %s", jsonContentType, ctx.Writer.Header()["Content-Type"])
// 		}

// 		var categories []models.Category
// 		json.Unmarshal(*body, &categories)
// 		expected := []models.Category{sampleCategory}
// 		if !reflect.DeepEqual(categories, expected) {
// 			t.Errorf("expected %v, got %v", expected, categories)
// 		}
// 	})

// 	t.Run("should get one category", func(t *testing.T) {
// 		ctx, body := getContext()
// 		ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.FormatUint(uint64(sampleCategory.ID), 10)}}
// 		controllers.GetCategory(ctx)

// 		if  status := ctx.Writer.Status(); status != 200 {
// 			t.Errorf("expected status code 200, got %d", status)
// 		}

// 		if !reflect.DeepEqual(ctx.Writer.Header()["Content-Type"], jsonContentType) {
// 			t.Errorf("expected content type %s, got %s", jsonContentType, ctx.Writer.Header()["Content-Type"])
// 		}

// 		var categories models.Category
// 		json.Unmarshal(*body, &categories)
// 		expected := sampleCategory
// 		if !reflect.DeepEqual(categories, expected) {
// 			t.Errorf("expected %v, got %v", expected, categories)
// 		}
// 	})

// 	t.Run("should return not found when getting a category that does not exist", func(t *testing.T) {
// 		ctx, body := getContext()
// 		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1000"}}
// 		controllers.GetCategory(ctx)

// 		if  status := ctx.Writer.Status(); status != 404 {
// 			t.Errorf("expected status code 404, got %d", status)
// 		}

// 		if !reflect.DeepEqual(*body, notFoundResponse) {
// 			t.Errorf("expected %s, got %s", notFoundResponse, *body)
// 		}
// 	})

// 	t.Run("should create a category", func(t *testing.T) {
// 		ctx, body := getContext()
// 		newCategory := models.NewCategory("test", "test", "test", "test", 100, 50)
// 		categoryJSON, err := json.Marshal(newCategory)
// 		if err != nil {
// 			t.Errorf("error marshalling new category: %v", err)
// 		}
// 		ctx.Request = &http.Request{
// 			Body: io.NopCloser(bytes.NewReader(categoryJSON)),
// 		}
// 		controllers.CreateCategory(ctx)

// 		if status := ctx.Writer.Status(); status != 201 {
// 			t.Errorf("expected status code 201, got %d", status)
// 		}

// 		if !reflect.DeepEqual(ctx.Writer.Header()["Content-Type"], jsonContentType) {
// 			t.Errorf("expected content type %s, got %s", jsonContentType, ctx.Writer.Header()["Content-Type"])
// 		}

// 		var category models.Category
// 		json.Unmarshal(*body, &category)

// 		var dbCategory models.Category
// 		db.First(&dbCategory, "id = ?", category.ID)

// 		// Normalize time zones to UTC for comparison.
// 		category.Model.CreatedAt = category.Model.CreatedAt.UTC()
// 		category.Model.UpdatedAt = category.Model.UpdatedAt.UTC()
// 		dbCategory.Model.CreatedAt = dbCategory.Model.CreatedAt.UTC()
// 		dbCategory.Model.UpdatedAt = dbCategory.Model.UpdatedAt.UTC()

// 		if !reflect.DeepEqual(category, dbCategory) {
// 			t.Errorf("expected %+v, got %+v", dbCategory, category)
// 		}
// 	})

// 	t.Run("should update a category, but only fields that are present", func(t *testing.T) {
// 		ctx, body := getContext()
// 		updateBody := map[string]string{
// 			"name":        "test-abc",
// 		}
// 		categoryJSON, err := json.Marshal(updateBody)
// 		if err != nil {
// 			t.Errorf("error marshalling new category: %v", err)
// 		}
// 		ctx.Request = &http.Request{
// 			Body: io.NopCloser(bytes.NewReader(categoryJSON)),
// 		}
// 		ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.FormatUint(uint64(sampleCategory.ID), 10)}}
// 		controllers.UpdateCategory(ctx)

// 		if status := ctx.Writer.Status(); status != 200 {
// 			t.Errorf("expected status code 200, got %d", status)
// 		}

// 		if !reflect.DeepEqual(ctx.Writer.Header()["Content-Type"], jsonContentType) {
// 			t.Errorf("expected content type %s, got %s", jsonContentType, ctx.Writer.Header()["Content-Type"])
// 		}

// 		var category models.Category
// 		json.Unmarshal(*body, &category)

// 		var dbCategory models.Category
// 		db.First(&dbCategory, "id = ?", category.ID)

// 		// Normalize time zones to UTC for comparison.
// 		category.Model.CreatedAt = category.Model.CreatedAt.UTC()
// 		category.Model.UpdatedAt = category.Model.UpdatedAt.UTC()
// 		dbCategory.Model.CreatedAt = dbCategory.Model.CreatedAt.UTC()
// 		dbCategory.Model.UpdatedAt = dbCategory.Model.UpdatedAt.UTC()

// 		if !reflect.DeepEqual(category, dbCategory) {
// 			t.Errorf("expected %+v, got %+v", dbCategory, category)
// 		}

// 		if dbCategory.Name != "test-abc" {
// 			t.Errorf("expected name to be test-abc, got %s", dbCategory.Name)
// 		}

// 		if dbCategory.Description != sampleCategory.Description {
// 			t.Errorf("expected description to be %s, got %s", sampleCategory.Description, dbCategory.Description)
// 		}
// 	})

// 	t.Run("should return not found when updating a category that does not exist", func(t *testing.T) {
// 		ctx, body := getContext()
// 		updateBody := map[string]string{
// 			"name":        "test-abc",
// 		}
// 		categoryJSON, err := json.Marshal(updateBody)
// 		if err != nil {
// 			t.Errorf("error marshalling new category: %v", err)
// 		}
// 		ctx.Request = &http.Request{
// 			Body: io.NopCloser(bytes.NewReader(categoryJSON)),
// 		}
// 		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1000"}}
// 		controllers.UpdateCategory(ctx)

// 		if status := ctx.Writer.Status(); status != 404 {
// 			t.Errorf("expected status code 404, got %d", status)
// 		}
// 		if !reflect.DeepEqual(*body, notFoundResponse) {
// 			t.Errorf("expected %s, got %s", notFoundResponse, *body)
// 		}
// 	})

// 	t.Run("should delete a category", func(t *testing.T) {
// 		ctx, body := getContext()
// 		ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.FormatUint(uint64(sampleCategory.ID), 10)}}
// 		controllers.DeleteCategory(ctx)

// 		if  status := ctx.Writer.Status(); status != 204 {
// 			t.Errorf("expected status code 204, got %d", status)
// 		}

// 		if !reflect.DeepEqual(*body, []byte{}) {
// 			t.Errorf("expected empty body, got %+v", *body)
// 		}

// 		var category models.Category
// 		db.First(&category, "id = ?", sampleCategory.ID)

// 		if category.ID != 0 {
// 			t.Errorf("expected category to be deleted, got %+v", category)
// 		}
// 	})

// 	t.Run("should return not found when deleting a category that does not exist", func(t *testing.T) {
// 		ctx, body := getContext()
// 		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1000"}}
// 		controllers.DeleteCategory(ctx)

// 		if status := ctx.Writer.Status(); status != 404 {
// 			t.Errorf("expected status code 404, got %d", status)
// 		}

// 		if !reflect.DeepEqual(*body, notFoundResponse) {
// 			t.Errorf("expected %s, got %s", notFoundResponse, *body)
// 		}
// 	})
// }


func getContext() (*gin.Context, *[]byte) {
	writer := &FakeWriter{
		HeadersMapping: make(http.Header),
		Body:           []byte{},
	}
	ctx, _ := gin.CreateTestContext(writer)
	return ctx, &writer.Body
}

func openDB() *gorm.DB {
	dsn := ":memory:"
	opener := sqlite.Open

	db := database.OpenDBConnection(dsn, opener)
	database.MakeMigrations(db)
	database.DB = db
	return db
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
