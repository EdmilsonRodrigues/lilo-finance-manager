//go:build integration

package router_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"testing/quick"

	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/controllers"
	customerrors "github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/custom_errors"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/database"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/models"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/router"
	httpserialization "github.com/EdmilsonRodrigues/lilo-finance-manager/src/common_utils/go/serialization/http_serialization"
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

	t.Run("should not be able to do any request without the headers", func(t *testing.T) {
		response, err := http.Get(ts.URL + "/categories/1")
		if err != nil {
			t.Errorf("failed to make request: %v", err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusForbidden {
			t.Errorf("expected status code %d, but got %d", http.StatusForbidden, response.StatusCode)
		}

		var categoriesErrorResponse httpserialization.ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&categoriesErrorResponse); err != nil {
			t.Errorf("failed to decode response body: %v", err)
		}

		if !reflect.DeepEqual(categoriesErrorResponse, customerrors.ForbiddenResponse) {
			t.Errorf("expected error response %+v, but got %+v", customerrors.ForbiddenResponse, categoriesErrorResponse)
			return
		}

		response, err = http.Get(ts.URL + "/categories/1/1")
		if err != nil {
			t.Errorf("failed to make request: %v", err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusForbidden {
			t.Errorf("expected status code %d, but got %d", http.StatusForbidden, response.StatusCode)
		}

		var categoryErrorResponse httpserialization.ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&categoryErrorResponse); err != nil {
			t.Errorf("failed to decode response body: %v", err)
		}

		if !reflect.DeepEqual(categoryErrorResponse, customerrors.ForbiddenResponse) {
			t.Errorf("expected error response %+v, but got %+v", customerrors.ForbiddenResponse, categoryErrorResponse)
			return
		}

		response, err = http.Post(ts.URL + "/categories/1", "application/json", nil)
		if err != nil {
			t.Errorf("failed to make request: %v", err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusForbidden {
			t.Errorf("expected status code %d, but got %d", http.StatusForbidden, response.StatusCode)
		}

		var createCategoryErrorResponse httpserialization.ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&createCategoryErrorResponse); err != nil {
			t.Errorf("failed to decode response body: %v", err)
		}

		if !reflect.DeepEqual(createCategoryErrorResponse, customerrors.ForbiddenResponse) {
			t.Errorf("expected error response %+v, but got %+v", customerrors.ForbiddenResponse, createCategoryErrorResponse)
			return
		}

		req, err := http.NewRequest(http.MethodPatch, ts.URL + "/categories/1/1", nil)
		if err != nil {
			t.Errorf("failed to create request: %v", err)

		}
		req.Header.Set("Content-Type", "application/json")
		response, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("failed to make request: %v", err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusForbidden {
			t.Errorf("expected status code %d, but got %d", http.StatusForbidden, response.StatusCode)
		}

		var updateCategoryErrorResponse httpserialization.ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&updateCategoryErrorResponse); err != nil {
			t.Errorf("failed to decode response body: %v", err)
			return
		}

		if !reflect.DeepEqual(updateCategoryErrorResponse, customerrors.ForbiddenResponse) {
			t.Errorf("expected error response %+v, but got %+v", customerrors.ForbiddenResponse, updateCategoryErrorResponse)
			return
		}

		req, err = http.NewRequest(http.MethodDelete, ts.URL + "/categories/1/1", nil)
		if err != nil {
			t.Errorf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		response, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("failed to make request: %v", err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusForbidden {
			t.Errorf("expected status code %d, but got %d", http.StatusForbidden, response.StatusCode)
		}

		var deleteCategoryErrorResponse httpserialization.ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&deleteCategoryErrorResponse); err != nil {
			t.Errorf("failed to decode response body: %v", err)
			return
		}

		if !reflect.DeepEqual(deleteCategoryErrorResponse, customerrors.ForbiddenResponse) {
			t.Errorf("expected error response %+v, but got %+v", customerrors.ForbiddenResponse, deleteCategoryErrorResponse)
			return
		}
	})

	t.Run("should get five categories of page 2", func(t *testing.T) {
		assertion := func(accountId string, names, descriptions, colors [10]string, budgets, currents [10]float64) bool {
			accountId = formatQueryValue(accountId)
			if accountId == "" {
				return true
			}
			size := len(names)
			categories := make([]models.Category, size)
			for i := range size {
				categories[i] = models.Category{Name: names[i], AccountID: accountId, Description: descriptions[i], Color: colors[i], Budget: budgets[i], Current: currents[i]}
				if err := db.Create(&categories[i]).Error; err != nil {
					t.Errorf("failed to create test category: %v", err)
					return false
				}
			}

			page := 2
			pageSize := page / 2

			reqUrl := ts.URL + "/categories/" + accountId + fmt.Sprintf("?page=%d&page_size=%d", page, pageSize)

			req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
			if err != nil {
				t.Errorf("failed to create GET request: %v", err)
				return false
			}

			req.Header.Set("X-User-Accounts", fmt.Sprintf(`{"%s": "admin"}`, accountId))
			req.Header.Set("X-User-Id", "test-user-id")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("error getting categories: %v", err)
				return false
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
				return false
			}

			if resp.Header.Get("Content-Type") != jsonContentType[0] {
				t.Errorf("expected content type %s, got %s", jsonContentType[0], resp.Header.Get("Content-Type"))
				return false
			}

			var returnedCategories httpserialization.PaginatedJSONResponse
			if err := json.NewDecoder(resp.Body).Decode(&returnedCategories); err != nil {
				t.Errorf("error decoding response body: %v", err)
				return false
			}

			if len(returnedCategories.Data.Items) != pageSize {
				t.Errorf("expected %d categories, got %d", pageSize, len(returnedCategories.Data.Items))
				return false
			}

			categoryResponses := make([]httpserialization.Serializer, pageSize)
			for i := range pageSize {
				categoryResponse := controllers.CategoryResponse{}
				if err := categoryResponse.BindModel(categories[pageSize+i]); err != nil {
					t.Errorf("error binding category to response: %v", err)
					return false
				}
				categoryResponses[i] = &categoryResponse
			}

			expectedResponse, err := httpserialization.CreatePaginatedResponse(page, pageSize, size, map[string]string{}, categoryResponses, []string{})
			if err != nil {
				t.Errorf("error creating expected response: %v", err)
				return false
			}

			if !reflect.DeepEqual(returnedCategories, expectedResponse) {
				t.Errorf("expected categories %v, got %v", expectedResponse, returnedCategories)
				return false
			}

			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 100,
		}); err != nil {
			t.Error(err)
		}
	})

	t.Run("should get name and account_id of categories filtered by name", func(t *testing.T) {
		assertion := func(accountId string, names, descriptions, colors [10]string, budgets, currents [10]float64) bool {
			accountId = formatQueryValue(accountId)
			if accountId == "" {
				return true
			}
			size := len(names)
			categories := make([]models.Category, size)
			for i := range size {
				categories[i] = models.Category{Name: names[i], AccountID: accountId, Description: descriptions[i], Color: colors[i], Budget: budgets[i], Current: currents[i]}
				if err := db.Create(&categories[i]).Error; err != nil {
					t.Errorf("failed to create test category: %v", err)
					return false
				}
			}

			page := 1
			pageSize := 10
			name := names[0]

			reqUrl := ts.URL + "/categories/" + accountId + fmt.Sprintf("?page=%d&page_size=%d&filters=name:%s&return_fields=name,account_id", page, pageSize, url.QueryEscape(name))

			req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
			if err != nil {
				t.Errorf("failed to create GET request: %v", err)
				return false
			}

			req.Header.Set("X-User-Accounts", fmt.Sprintf(`{"%s": "admin"}`, accountId))
			req.Header.Set("X-User-Id", "test-user-id")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("error getting categories: %v", err)
				return false
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
				return false
			}

			if resp.Header.Get("Content-Type") != jsonContentType[0] {
				t.Errorf("expected content type %s, got %s", jsonContentType[0], resp.Header.Get("Content-Type"))
				return false
			}

			var returnedCategories httpserialization.PaginatedJSONResponse
			if err := json.NewDecoder(resp.Body).Decode(&returnedCategories); err != nil {
				t.Errorf("error decoding response body: %v", err)
				return false
			}

			categoryResponses := []httpserialization.Serializer{}
			for i := range pageSize {
				categoryResponse := controllers.CategoryResponse{}
				if err := categoryResponse.BindModel(categories[i]); err != nil {
					t.Errorf("error binding category to response: %v", err)
					return false
				}
				if categoryResponse.Name == name {
					categoryResponses = append(categoryResponses, &categoryResponse)
				}
			}

			expectedResponse, err := httpserialization.CreatePaginatedResponse(page, len(categoryResponses), len(categoryResponses), map[string]string{"name": name}, categoryResponses, []string{"name", "account_id"})
			if err != nil {
				t.Errorf("error creating expected response: %v", err)
				return false
			}

			if !reflect.DeepEqual(returnedCategories, expectedResponse) {
				t.Errorf("expected categories %v, got %v", expectedResponse, returnedCategories)
				return false
			}

			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 100,
		}); err != nil {
			t.Error(err)
		}
	})

	t.Run("should get an account by id", func(t *testing.T) {
		assertion := func(accountId string, name, description, color string, budget, current float64) bool {
			accountId = formatQueryValue(accountId)
			if accountId == "" {
				return true
			}
			category := models.Category{Name: name, AccountID: accountId, Description: description, Color: color, Budget: budget, Current: current}
			if err := db.Create(&category).Error; err != nil {
				t.Errorf("failed to create test category: %v", err)
				return false
			}

			reqUrl := ts.URL + "/categories/" + accountId + "/" + strconv.Itoa(int(category.ID))

			req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
			if err != nil {
				t.Errorf("failed to create GET request: %v", err)
				return false
			}

			req.Header.Set("X-User-Accounts", fmt.Sprintf(`{"%s": "admin"}`, accountId))
			req.Header.Set("X-User-Id", "test-user-id")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("error getting category by id: %v", err)
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
				return false
			}

			if resp.Header.Get("Content-Type") != jsonContentType[0] {
				t.Errorf("expected content type %s, got %s", jsonContentType[0], resp.Header.Get("Content-Type"))
				return false
			}

			var returnedCategory httpserialization.JSONResponse
			if err := json.NewDecoder(resp.Body).Decode(&returnedCategory); err != nil {
				t.Errorf("error decoding response body: %v", err)
				return false
			}

			categoryResponse := controllers.CategoryResponse{}
			if err := categoryResponse.BindModel(category); err != nil {
				t.Errorf("error binding category to response: %v", err)
				return false
			}

			expectedResponse, err := categoryResponse.Marshal([]string{})
			if err != nil {
				t.Errorf("error creating expected response: %v", err)
				return false
			}

			if !reflect.DeepEqual(returnedCategory, expectedResponse) {
				t.Errorf("expected category %v, got %v", expectedResponse, returnedCategory)
				return false
			}

			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 100,
		}); err != nil {
			t.Error(err)
		}
	})

	t.Run("should not get a category by id from another account", func(t *testing.T) {
		assertion := func(accountId string, name, description, color string, budget, current float64) bool {
			accountId = formatQueryValue(accountId)
			if accountId == "" {
				return true
			}
			category := models.Category{Name: name, AccountID: "2", Description: description, Color: color, Budget: budget, Current: current}
			if err := db.Create(&category).Error; err != nil {
				t.Errorf("failed to create test category: %v", err)
				return false
			}

			reqUrl := ts.URL + "/categories/" + accountId + "/" + strconv.Itoa(int(category.ID))

			req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
			if err != nil {
				t.Errorf("failed to create GET request: %v", err)
				return false
			}

			req.Header.Set("X-User-Accounts", fmt.Sprintf(`{"%s": "admin"}`, accountId))
			req.Header.Set("X-User-Id", "test-user-id")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("error getting category by id: %v", err)
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusNotFound {
				t.Errorf("expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
				return false
			}

			var errorResponse httpserialization.ErrorResponse
			if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
				t.Errorf("error decoding response body: %v", err)
				return false
			}

			if !reflect.DeepEqual(errorResponse, controllers.CategoryNotFoundResponse) {
				t.Errorf("expected error response %v, got %v", controllers.CategoryNotFoundResponse, errorResponse)
				return false
			}

			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 100,
		}); err != nil {
			t.Error(err)
		}
	})

	t.Run("should get name and account id from a category", func(t *testing.T) {
		assertion := func(accountId string, name, description, color string, budget, current float64) bool {
			accountId = formatQueryValue(accountId)
			if accountId == "" {
				return true
			}
			category := models.Category{Name: name, AccountID: accountId, Description: description, Color: color, Budget: budget, Current: current}
			if err := db.Create(&category).Error; err != nil {
				t.Errorf("failed to create test category: %v", err)
				return false
			}

			reqUrl := ts.URL + "/categories/" + accountId + "/" + strconv.Itoa(int(category.ID)) + "?return_fields=name,account_id"

			req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
			if err != nil {
				t.Errorf("failed to create GET request: %v", err)
				return false
			}

			req.Header.Set("X-User-Accounts", fmt.Sprintf(`{"%s": "admin"}`, accountId))
			req.Header.Set("X-User-Id", "test-user-id")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("error getting category by id: %v", err)
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
				return false
			}

			var categoryResponse httpserialization.JSONResponse
			if err := json.NewDecoder(resp.Body).Decode(&categoryResponse); err != nil {
				t.Errorf("error decoding response body: %v", err)
				return false
			}

			expected := httpserialization.JSONResponse{Status: "success", Data: httpserialization.DataItem{"name": name, "account_id": accountId}}

			if !reflect.DeepEqual(categoryResponse, expected) {
				t.Errorf("expected category %v, got %v", expected, categoryResponse)
				return false
			}
			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 100,
		}); err != nil {
			t.Error(err)
		}
	})

	t.Run("should update a category and get name and budget", func(t *testing.T) {
		assertion := func(accountId string, name, description, color string, budget, current float64, newName string) bool {
			newName = formatQueryValue(newName)
			if newName == "" {
				return true
			}
			accountId = formatQueryValue(accountId)
			if accountId == "" {
				return true
			}
			category := models.Category{Name: name, AccountID: accountId, Description: description, Color: color, Budget: budget, Current: current}
			if err := db.Create(&category).Error; err != nil {
				t.Errorf("failed to create test category: %v", err)
				return false
			}

			reqUrl := ts.URL + "/categories/" + accountId + "/" + strconv.Itoa(int(category.ID)) + "?return_fields=name,budget"

			reqBody := strings.NewReader(fmt.Sprintf(`{"name": "%s"}`, newName))

			req, err := http.NewRequest(http.MethodPatch, reqUrl, reqBody)
			if err != nil {
				t.Errorf("failed to create PATCH request: %v", err)
			}

			req.Header.Set("X-User-Accounts", fmt.Sprintf(`{"%s": "admin"}`, accountId))
			req.Header.Set("X-User-Id", "test-user-id")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("error updating category: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
				return false
			}

			var categoryResponse httpserialization.JSONResponse
			if err := json.NewDecoder(resp.Body).Decode(&categoryResponse); err != nil {
				t.Errorf("error decoding response body: %v", err)
				return false
			}

			expected := httpserialization.JSONResponse{Status: "success", Data: httpserialization.DataItem{"name": newName, "budget": budget}}

			if !reflect.DeepEqual(categoryResponse, expected) {
				t.Errorf("expected category %v, got %v", expected, categoryResponse)
				return false
			}
			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 100,
		}); err != nil {
			t.Error(err)
		}
	})

	t.Run("should not update a category from another account", func(t *testing.T) {
		assertion := func(accountId string, name, description, color string, budget, current float64, newName string) bool {
			accountId = formatQueryValue(accountId)
			if accountId == "" {
				return true
			}
			category := models.Category{Name: name, AccountID: "2", Description: description, Color: color, Budget: budget, Current: current}
			if err := db.Create(&category).Error; err != nil {
				t.Errorf("failed to create test category: %v", err)
				return false
			}

			reqUrl := ts.URL + "/categories/" + accountId + "/" + strconv.Itoa(int(category.ID))

			reqBody := strings.NewReader(fmt.Sprintf(`{"name": "%s"}`, newName))

			req, err := http.NewRequest(http.MethodPatch, reqUrl, reqBody)
			if err != nil {
				t.Errorf("failed to create PATCH request: %v", err)
			}

			req.Header.Set("X-User-Accounts", fmt.Sprintf(`{"%s": "admin"}`, accountId))
			req.Header.Set("X-User-Id", "test-user-id")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("error updating category: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusNotFound {
				t.Errorf("expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
				return false
			}

			var errorResponse httpserialization.ErrorResponse
			if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
				t.Errorf("error decoding response body: %v", err)
				return false
			}

			if !reflect.DeepEqual(errorResponse, controllers.CategoryNotFoundResponse) {
				t.Errorf("expected error response %v, got %v", controllers.CategoryNotFoundResponse, errorResponse)
				return false
			}

			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 100,
		}); err != nil {
			t.Error(err)
		}
	})

	t.Run("should delete a category", func(t *testing.T) {
		assertion := func(accountId string, name, description, color string, budget, current float64) bool {
			accountId = formatQueryValue(accountId)
			if accountId == "" {
				return true
			}
			category := models.Category{Name: name, AccountID: accountId, Description: description, Color: color, Budget: budget, Current: current}
			if err := db.Create(&category).Error; err != nil {
				t.Errorf("failed to create test category: %v", err)
				return false
			}

			reqUrl := ts.URL + "/categories/" + accountId + "/" + strconv.Itoa(int(category.ID))

			req, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
			if err != nil {
				t.Errorf("failed to create DELETE request: %v", err)
				return false
			}

			req.Header.Set("X-User-Accounts", fmt.Sprintf(`{"%s": "admin"}`, accountId))
			req.Header.Set("X-User-Id", "test-user-id")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("error deleting category: %v", err)
				return false
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusNoContent {
				t.Errorf("expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
				return false
			}

			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 100,
		}); err != nil {
			t.Error(err)
		}
	})

	t.Run("should not delete a category from another account", func(t *testing.T) {
		assertion := func(accountId string, name, description, color string, budget, current float64) bool {
			accountId = formatQueryValue(accountId)
			if accountId == "" {
				return true
			}
			category := models.Category{Name: name, AccountID: "2", Description: description, Color: color, Budget: budget, Current: current}
			if err := db.Create(&category).Error; err != nil {
				t.Errorf("failed to create test category: %v", err)
				return false
			}

			reqUrl := ts.URL + "/categories/" + accountId + "/" + strconv.Itoa(int(category.ID))

			req, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
			if err != nil {
				t.Errorf("failed to create DELETE request: %v", err)
				return false
			}

			req.Header.Set("X-User-Accounts", fmt.Sprintf(`{"%s": "admin"}`, accountId))
			req.Header.Set("X-User-Id", "test-user-id")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("error deleting category: %v", err)
				return false
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusNotFound {
				t.Errorf("expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
				return false
			}

			var errorResponse httpserialization.ErrorResponse
			if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
				t.Errorf("error decoding response body: %v", err)
				return false
			}

			if !reflect.DeepEqual(errorResponse, controllers.CategoryNotFoundResponse) {
				t.Errorf("expected error response %v, got %v", controllers.CategoryNotFoundResponse, errorResponse)
				return false
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

// t.Run("should not get a category from another account", func(t *testing.T) {
// 	// Create a category in the database for the test.
// 	category := models.Category{Name: "CategoryX", AccountID: "1", Description: "DescriptionX", Color: "ColorX", Budget: 150, Current: 75}
// 	if err := db.Create(&category).Error; err != nil {
// 		t.Fatalf("failed to create test category: %v", err)
// 	}

// 	resp, err := http.Get(ts.URL + "/categories/2/" + strconv.Itoa(int(category.ID)))
// 	if err != nil {
// 		t.Errorf("error getting category by id: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusNotFound {
// 		t.Errorf("expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
// 	}
// })

// t.Run("should get a category by id", func(t *testing.T) {
// 	// Create a category in the database for the test.
// 	createdCategory := models.Category{Name: "CategoryX", AccountID: "1", Description: "DescriptionX", Color: "ColorX", Budget: 150, Current: 75}
// 	if err := db.Create(&createdCategory).Error; err != nil {
// 		t.Fatalf("failed to create test category: %v", err)
// 	}

// 	resp, err := http.Get(ts.URL + "/categories/1/" + strconv.Itoa(int(createdCategory.ID)))
// 	if err != nil {
// 		t.Errorf("error getting category by id: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
// 	}

// 	if resp.Header.Get("Content-Type") != jsonContentType[0] {
// 		t.Errorf("expected content type %s, got %s", jsonContentType[0], resp.Header.Get("Content-Type"))
// 	}

// 	var category models.Category
// 	if err := json.NewDecoder(resp.Body).Decode(&category); err != nil {
// 		t.Errorf("error decoding response body: %v", err)
// 	}

// 	normalizeDateTimes(&createdCategory, &category)

// 	if !reflect.DeepEqual(createdCategory, category) {
// 		t.Errorf("expected category %+v, got %+v", createdCategory, category)
// 	}
// })

// t.Run("should create a category", func(t *testing.T) {
// 	newCategory := models.Category{Name: "NewCategory", AccountID: "2", Description: "NewDescription", Color: "NewColor", Budget: 250, Current: 125}
// 	categoryJSON, err := json.Marshal(newCategory)
// 	if err != nil {
// 		t.Fatalf("failed to marshal new category: %v", err)
// 	}

// 	resp, err := http.Post(ts.URL+"/categories/2", "application/json", bytes.NewReader(categoryJSON))
// 	if err != nil {
// 		t.Errorf("error creating category: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusCreated {
// 		t.Errorf("expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
// 	}

// 	if resp.Header.Get("Content-Type") != jsonContentType[0] {
// 		t.Errorf("expected content type %s, got %s", jsonContentType[0], resp.Header.Get("Content-Type"))
// 	}

// 	var createdCategory models.Category
// 	if err := json.NewDecoder(resp.Body).Decode(&createdCategory); err != nil {
// 		t.Errorf("error decoding response body: %v", err)
// 	}

// 	// Check if the category was created in the database.
// 	var dbCategory models.Category
// 	if err := db.First(&dbCategory, "id = ?", createdCategory.ID).Error; err != nil {
// 		t.Errorf("error querying database for created category: %v", err)
// 	}

// 	// Normalize time zones to UTC for comparison.
// 	normalizeDateTimes(&createdCategory, &dbCategory)

// 	if !reflect.DeepEqual(createdCategory, dbCategory) {
// 		t.Errorf("expected %+v, got %+v", dbCategory, createdCategory)
// 	}
// })

// t.Run("should update a category", func(t *testing.T) {
// 	// Create a category to update.
// 	categoryToUpdate := models.Category{Name: "OldName", AccountID: "3", Description: "OldDescription", Color: "OldColor", Budget: 300, Current: 150}
// 	if err := db.Create(&categoryToUpdate).Error; err != nil {
// 		t.Fatalf("failed to create category to update: %v", err)
// 	}

// 	updatedCategory := models.Category{Name: "NewName", Description: "NewDescription", Color: "NewColor", Budget: 400, Current: 200}
// 	updatedCategoryJSON, err := json.Marshal(updatedCategory)
// 	if err != nil {
// 		t.Fatalf("failed to marshal updated category: %v", err)
// 	}

// 	req, err := http.NewRequest(http.MethodPatch, ts.URL+"/categories/3/"+strconv.Itoa(int(categoryToUpdate.ID)), bytes.NewReader(updatedCategoryJSON))
// 	if err != nil {
// 		t.Fatalf("failed to create PATCH request: %v", err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	client := &http.Client{}
// 	resp, err := client.Do(req)

// 	if err != nil {
// 		t.Errorf("error updating category: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
// 	}

// 	if resp.Header.Get("Content-Type") != jsonContentType[0] {
// 		t.Errorf("expected content type %s, got %s", jsonContentType[0], resp.Header.Get("Content-Type"))
// 	}

// 	var updatedCategoryResponse models.Category
// 	if err := json.NewDecoder(resp.Body).Decode(&updatedCategoryResponse); err != nil {
// 		t.Errorf("error decoding response body: %v", err)
// 	}

// 	// Check the category in the database.
// 	var dbCategory models.Category
// 	if err := db.First(&dbCategory, "id = ?", categoryToUpdate.ID).Error; err != nil {
// 		t.Errorf("error querying database for updated category: %v", err)
// 	}

// 	// Normalize time zones to UTC for comparison.
// 	normalizeDateTimes(&updatedCategoryResponse, &dbCategory)

// 	//  update the fields, that were changed
// 	categoryToUpdate.Name = updatedCategory.Name
// 	categoryToUpdate.Description = updatedCategory.Description
// 	categoryToUpdate.Color = updatedCategory.Color
// 	categoryToUpdate.Budget = updatedCategory.Budget

// 	if !reflect.DeepEqual(updatedCategoryResponse, dbCategory) {
// 		t.Errorf("expected %+v, got %+v", dbCategory, updatedCategoryResponse)
// 	}

// 	if updatedCategoryResponse.Current != categoryToUpdate.Current {
// 		t.Errorf("expected Current to be unaltered, got %v", updatedCategoryResponse.Current)
// 	}

// 	if updatedCategoryResponse.AccountID != categoryToUpdate.AccountID {
// 		t.Errorf("expected AccountID to be unaltered, got %v", updatedCategoryResponse.AccountID)
// 	}
// })

// t.Run("should delete a category", func(t *testing.T) {
// 	// Create a category to delete.
// 	categoryToDelete := models.Category{Name: "ToDelete", AccountID: "4", Description: "ToDeleteDescription", Color: "ToDeleteColor", Budget: 450, Current: 225}
// 	if err := db.Create(&categoryToDelete).Error; err != nil {
// 		t.Fatalf("failed to create category to delete: %v", err)
// 	}

// 	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/categories/4/"+strconv.Itoa(int(categoryToDelete.ID)), nil)
// 	if err != nil {
// 		t.Fatalf("failed to create DELETE request: %v", err)
// 	}
// 	client := &http.Client{}
// 	resp, err := client.Do(req)

// 	if err != nil {
// 		t.Errorf("error deleting category: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusNoContent {
// 		t.Errorf("expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
// 	}

// 	// Check if the category is actually deleted (soft deleted).
// 	var deletedCategory models.Category
// 	if err := db.Unscoped().First(&deletedCategory, "id = ?", categoryToDelete.ID).Error; err != nil {
// 		t.Errorf("error querying database for deleted category: %v", err)
// 	}

// 	if deletedCategory.ID != categoryToDelete.ID {
// 		t.Errorf("expected deleted category id to be %d, got %d", categoryToDelete.ID, deletedCategory.ID)
// 	}

// 	if deletedCategory.Model.DeletedAt.Valid == false {
// 		t.Errorf("expected deleted category to be soft deleted, but it was not")
// 	}
// })

func normalizeDateTimes(category *models.Category) {
	category.Model.CreatedAt = category.Model.CreatedAt.UTC()
	category.Model.UpdatedAt = category.Model.UpdatedAt.UTC()
}

func setupRouter() (*gin.Context, *gin.Engine, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(recorder)
	router.HandleRequests(engine)
	return ctx, engine, recorder
}

func formatQueryValue(queryValue string) string {
	escaped := url.QueryEscape(queryValue)
	return strings.ReplaceAll(strings.ReplaceAll(escaped, "+", "%20"), "%", "")
}
