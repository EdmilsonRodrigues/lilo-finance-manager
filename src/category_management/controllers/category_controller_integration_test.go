//go:build integration

package controllers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"strconv"
	"testing"
	"testing/quick"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/controllers"
	customerrors "github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/custom_errors"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/database"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/models"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/common_utils/go/serialization/http_serialization"
	"github.com/gin-gonic/gin"
)

func TestGetCategories(t *testing.T) {
	db := openDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	t.Run("should get 10 categories", func(t *testing.T) {
		assertion := func(names, descriptions, colors [10]string, budgets, currents [10]float64, accountId string) bool {
			if accountId == "" {
				return true
			}
			ctx, body := getContext()

			categories := make([]models.Category, 10)

			for index := range 10 {
				categories[index] = models.Category{
					AccountID:   accountId,
					Name:        names[index],
					Description: descriptions[index],
					Color:       colors[index],
					Budget:      budgets[index],
					Current:     currents[index],
				}
				db.Create(&categories[index])
			}

			ctx.Set("conditions", "account_id = ?")
			ctx.Set("conditions_values", []interface{}{accountId})

			controllers.GetCategories(ctx)

			if status := ctx.Writer.Status(); status != http.StatusOK {
				t.Errorf("expected status code 200, got %d", status)
				return false
			}

			var categoriesResponse httpserialization.PaginatedJSONResponse
			if err := json.Unmarshal(*body, &categoriesResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if categoriesResponse.Status != "success" {
				t.Errorf("expected status to be success, got %s", categoriesResponse.Status)
				return false
			}

			if len(categoriesResponse.Data.Items) != len(categories) {
				t.Errorf("expected %d categories, got %d", len(categories), len(categoriesResponse.Data.Items))
				return false
			}

			if len(categoriesResponse.Data.Filters) != 0 {
				t.Errorf("expected filters to be empty, got %v", categoriesResponse.Data.Filters)
				return false
			}

			if categoriesResponse.Data.TotalItems != len(categories) {
				t.Errorf("expected total to be %d, got %d", len(categories), categoriesResponse.Data.TotalItems)
				return false
			}

			if categoriesResponse.Data.Total != 1 {
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

			for i, category := range categoriesResponse.Data.Items {
				if category["id"] == 0 {
					t.Errorf("expected id to be %d, got %d", 0, category["id"])
					return false
				}

				if category["account_id"] != accountId {
					t.Errorf("expected account id to be %s, got %s", accountId, category["account_id"])
					return false
				}

				if category["name"] != names[i] {
					t.Errorf("expected name to be %s, got %s", names[i], category["name"])
					return false
				}

				if category["description"] != descriptions[i] {
					t.Errorf("expected description to be %s, got %s", descriptions[i], category["description"])
					return false
				}

				if category["color"] != colors[i] {
					t.Errorf("expected color to be %s, got %s", colors[i], category["color"])
					return false
				}

				if category["budget"] != budgets[i] {
					t.Errorf("expected budget to be %f, got %v", budgets[i], category["budget"])
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

	t.Run("should get 5 categories", func(t *testing.T) {
		assertion := func(names, descriptions, colors [10]string, budgets, currents [10]float64, accountId string) bool {
			if accountId == "" {
				return true
			}
			ctx, body := getContext()

			categories := make([]models.Category, 10)

			for index := range 10 {
				categories[index] = models.Category{
					AccountID:   accountId,
					Name:        names[index],
					Description: descriptions[index],
					Color:       colors[index],
					Budget:      budgets[index],
					Current:     currents[index],
				}
				db.Create(&categories[index])
			}

			ctx.Set("conditions", "account_id = ?")
			ctx.Set("conditions_values", []interface{}{accountId})
			ctx.Set("pageSize", 5)

			controllers.GetCategories(ctx)

			if status := ctx.Writer.Status(); status != http.StatusOK {
				t.Errorf("expected status code 200, got %d", status)
				return false
			}

			var categoriesResponse httpserialization.PaginatedJSONResponse
			if err := json.Unmarshal(*body, &categoriesResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if categoriesResponse.Status != "success" {
				t.Errorf("expected status to be success, got %s", categoriesResponse.Status)
				return false
			}

			if len(categoriesResponse.Data.Items) != 5 {
				t.Errorf("expected %d categories, got %d", 5, len(categoriesResponse.Data.Items))
				return false
			}

			if len(categoriesResponse.Data.Filters) != 0 {
				t.Errorf("expected filters to be empty, got %v", categoriesResponse.Data.Filters)
				return false
			}

			if categoriesResponse.Data.TotalItems != len(categories) {
				t.Errorf("expected total to be %d, got %d", len(categories), categoriesResponse.Data.TotalItems)
				return false
			}

			if categoriesResponse.Data.Total != 2 {
				t.Errorf("expected total to be %d, got %d", 2, categoriesResponse.Data.Total)
				return false
			}

			if categoriesResponse.Data.Page != 1 {
				t.Errorf("expected page to be %d, got %d", 1, categoriesResponse.Data.Page)
				return false
			}

			if categoriesResponse.Data.Size != 5 {
				t.Errorf("expected size to be %d, got %d", 5, categoriesResponse.Data.Size)
				return false
			}

			for i, category := range categoriesResponse.Data.Items {
				if category["id"] == 0 {
					t.Errorf("expected id to be %d, got %d", 0, category["id"])
					return false
				}

				if category["account_id"] != accountId {
					t.Errorf("expected account id to be %s, got %s", accountId, category["account_id"])
					return false
				}

				if category["name"] != names[i] {
					t.Errorf("expected name to be %s, got %s", names[i], category["name"])
					return false
				}

				if category["description"] != descriptions[i] {
					t.Errorf("expected description to be %s, got %s", descriptions[i], category["description"])
					return false
				}

				if category["color"] != colors[i] {
					t.Errorf("expected color to be %s, got %s", colors[i], category["color"])
					return false
				}

				if category["budget"] != budgets[i] {
					t.Errorf("expected budget to be %f, got %v", budgets[i], category["budget"])
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

	t.Run("should get categories filtered by name", func(t *testing.T) {
		assertion := func(names, descriptions, colors [10]string, budgets, currents [10]float64, accountId string) bool {
			if accountId == "" || names[0] == "" {
				return true
			}
			ctx, body := getContext()

			categories := make([]models.Category, 10)
			name := formatQueryValue(names[0])

			for index := range 10 {
				name := formatQueryValue(names[index])
				categories[index] = models.Category{
					AccountID:   accountId,
					Name:        name,
					Description: descriptions[index],
					Color:       colors[index],
					Budget:      budgets[index],
					Current:     currents[index],
				}
				db.Create(&categories[index])
			}

			ctx.Set("conditions", "account_id = ? AND name = ?")
			ctx.Set("conditions_values", []interface{}{accountId, name})
			ctx.Set("filters", map[string]string{"name": name})

			controllers.GetCategories(ctx)

			if status := ctx.Writer.Status(); status != http.StatusOK {
				t.Errorf("expected status code 200, got %d", status)
				return false
			}

			var categoriesResponse httpserialization.PaginatedJSONResponse
			if err := json.Unmarshal(*body, &categoriesResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if categoriesResponse.Status != "success" {
				t.Errorf("expected status to be success, got %s", categoriesResponse.Status)
				return false
			}

			if len(categoriesResponse.Data.Items) != 1 {
				t.Errorf("expected %d categories, got %d", 1, len(categoriesResponse.Data.Items))
				return false
			}

			if !reflect.DeepEqual(categoriesResponse.Data.Filters, map[string]string{"name": name}) {
				t.Errorf("expected filters to be %v, got %v", map[string]string{"name": name}, categoriesResponse.Data.Filters)
				return false
			}

			if categoriesResponse.Data.TotalItems != 1 {
				t.Errorf("expected total to be %d, got %d", 1, categoriesResponse.Data.TotalItems)
				return false
			}

			if categoriesResponse.Data.Total != 1 {
				t.Errorf("expected total to be %d, got %d", 1, categoriesResponse.Data.Total)
				return false
			}

			if categoriesResponse.Data.Page != 1 {
				t.Errorf("expected page to be %d, got %d", 1, categoriesResponse.Data.Page)
				return false
			}

			if categoriesResponse.Data.Size != 1 {
				t.Errorf("expected size to be %d, got %d", 1, categoriesResponse.Data.Size)
				return false
			}

			for i, category := range categoriesResponse.Data.Items {
				if category["id"] == 0 {
					t.Errorf("expected id not to be %d, got %d", 0, category["id"])
					return false
				}

				if category["account_id"] != accountId {
					t.Errorf("expected account id to be %s, got %s", accountId, category["account_id"])
					return false
				}

				if category["name"] != name {
					t.Errorf("expected name to be %s, got %s", name, category["name"])
					return false
				}

				if category["description"] != descriptions[i] {
					t.Errorf("expected description to be %s, got %s", descriptions[i], category["description"])
					return false
				}

				if category["color"] != colors[i] {
					t.Errorf("expected color to be %s, got %s", colors[i], category["color"])
					return false
				}

				if category["budget"] != budgets[i] {
					t.Errorf("expected budget to be %f, got %v", budgets[i], category["budget"])
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

	t.Run("should get categories without budget", func(t *testing.T) {
		assertion := func(names, descriptions, colors [10]string, budgets, currents [10]float64, accountId string) bool {
			if accountId == "" {
				return true
			}
			ctx, body := getContext()

			categories := make([]models.Category, 10)

			for index := range 10 {
				categories[index] = models.Category{
					AccountID:   accountId,
					Name:        names[index],
					Description: descriptions[index],
					Color:       colors[index],
					Budget:      budgets[index],
					Current:     currents[index],
				}
				db.Create(&categories[index])
			}

			ctx.Set("conditions", "account_id = ?")
			ctx.Set("conditions_values", []interface{}{accountId})
			ctx.Set("returnFields", []string{"id", "account_id", "name", "description", "color", "current"})

			controllers.GetCategories(ctx)

			if status := ctx.Writer.Status(); status != http.StatusOK {
				t.Errorf("expected status code 200, got %d", status)
				return false
			}

			var categoriesResponse httpserialization.PaginatedJSONResponse
			if err := json.Unmarshal(*body, &categoriesResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if categoriesResponse.Status != "success" {
				t.Errorf("expected status to be success, got %s", categoriesResponse.Status)
				return false
			}

			if len(categoriesResponse.Data.Items) != len(categories) {
				t.Errorf("expected %d categories, got %d", len(categories), len(categoriesResponse.Data.Items))
				return false
			}

			if len(categoriesResponse.Data.Filters) != 0 {
				t.Errorf("expected filters to be empty, got %v", categoriesResponse.Data.Filters)
				return false
			}

			if categoriesResponse.Data.TotalItems != len(categories) {
				t.Errorf("expected total to be %d, got %d", len(categories), categoriesResponse.Data.TotalItems)
				return false
			}

			if categoriesResponse.Data.Total != 1 {
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

			for i, category := range categoriesResponse.Data.Items {
				if category["id"] == 0 {
					t.Errorf("expected id to be %d, got %d", 0, category["id"])
					return false
				}

				if category["account_id"] != accountId {
					t.Errorf("expected account id to be %s, got %s", accountId, category["account_id"])
					return false
				}

				if category["name"] != names[i] {
					t.Errorf("expected name to be %s, got %s", names[i], category["name"])
					return false
				}

				if category["description"] != descriptions[i] {
					t.Errorf("expected description to be %s, got %s", descriptions[i], category["description"])
					return false
				}

				if category["color"] != colors[i] {
					t.Errorf("expected color to be %s, got %s", colors[i], category["color"])
					return false
				}

				if _, exists := category["budget"]; exists {
					t.Errorf("expected budget to not be in response, got %v", category["budget"])
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

func TestGetCategory(t *testing.T) {
	db := openDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	t.Run("should get a category", func(t *testing.T) {
		assertion := func(name, description, color string, budget, current float64, accountId string) bool {
			if accountId == "" {
				return true
			}
			ctx, body := getContext()

			category := models.Category{
				AccountID:   accountId,
				Name:        name,
				Description: description,
				Color:       color,
				Budget:      budget,
				Current:     current,
			}

			db.Create(&category)

			ctx.Set("conditions", "account_id = ? AND id = ?")
			ctx.Set("conditions_values", []interface{}{accountId, category.ID})

			controllers.GetCategory(ctx)

			if status := ctx.Writer.Status(); status != http.StatusOK {
				t.Errorf("expected status code 200, got %d", status)
				return false
			}

			var categoriesResponse httpserialization.JSONResponse
			if err := json.Unmarshal(*body, &categoriesResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if categoriesResponse.Status != "success" {
				t.Errorf("expected status to be success, got %s", categoriesResponse.Status)
				return false
			}

			expectedResponse := controllers.CategoryResponse{}
			expectedResponse.BindModel(category)
			expectedData, err := expectedResponse.Marshal([]string{})
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if !reflect.DeepEqual(categoriesResponse, expectedData) {
				t.Errorf("expected %+v categories, got %+v", expectedData, categoriesResponse)
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

	t.Run("should not get a non existing category", func(t *testing.T) {
		assertion := func(name, description, color string, budget, current float64, accountId string) bool {
			if accountId == "" {
				return true
			}
			ctx, body := getContext()

			category := models.Category{
				AccountID:   accountId,
				Name:        name,
				Description: description,
				Color:       color,
				Budget:      budget,
				Current:     current,
			}

			db.Create(&category)

			ctx.Set("conditions", "account_id = ? AND id = ?")
			ctx.Set("conditions_values", []interface{}{accountId, category.ID + 1})

			controllers.GetCategory(ctx)

			if status := ctx.Writer.Status(); status != http.StatusNotFound {
				t.Errorf("expected status code 404, got %d", status)
				return false
			}

			var errorResponse httpserialization.ErrorResponse
			if err := json.Unmarshal(*body, &errorResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if !reflect.DeepEqual(errorResponse, controllers.CategoryNotFoundResponse) {
				t.Errorf("expected %+v error, got %+v", controllers.CategoryNotFoundResponse, errorResponse)
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

	t.Run("should get a category, but only the selected fields", func(t *testing.T) {
		assertion := func(name, description, color string, budget, current float64, accountId string) bool {
			if accountId == "" {
				return true
			}
			ctx, body := getContext()

			category := models.Category{
				AccountID:   accountId,
				Name:        name,
				Description: description,
				Color:       color,
				Budget:      budget,
				Current:     current,
			}

			db.Create(&category)

			possibleFields := []string{"id", "account_id", "name", "description", "color", "budget", "current"}

			fields := []string{}
			for range 4 {
				p := possibleFields[rand.Intn(len(possibleFields))]
				if !contains(fields, p) {
					fields = append(fields, p)
				}
			}

			ctx.Set("conditions", "account_id = ? AND id = ?")
			ctx.Set("conditions_values", []interface{}{accountId, category.ID})
			ctx.Set("returnFields", fields)

			controllers.GetCategory(ctx)

			if status := ctx.Writer.Status(); status != http.StatusOK {
				t.Errorf("expected status code 200, got %d", status)
				return false
			}

			var categoriesResponse httpserialization.JSONResponse
			if err := json.Unmarshal(*body, &categoriesResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if categoriesResponse.Status != "success" {
				t.Errorf("expected status to be success, got %s", categoriesResponse.Status)
				return false
			}

			expectedResponse := controllers.CategoryResponse{}
			expectedResponse.BindModel(category)
			expectedData, err := expectedResponse.Marshal(fields)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if !reflect.DeepEqual(categoriesResponse, expectedData) {
				t.Errorf("expected %+v categories, got %+v", expectedData, categoriesResponse)
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

func TestCreateCategory(t *testing.T) {
	db := openDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	t.Run("should create a category", func(t *testing.T) {
		assertion := func(name, description, color string, budget, current float64, accountId string) bool {
			if accountId == "" {
				return true
			}
			ctx, body := getContext()

			category := models.Category{
				AccountID:   accountId,
				Name:        name,
				Description: description,
				Color:       color,
				Budget:      budget,
				Current:     current,
			}

			createReq := controllers.CategoryResponse{}
			createReq.BindModel(category)

			categoryJson, err := json.Marshal(createReq)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			ctx.Request = &http.Request{Body: io.NopCloser(bytes.NewReader(categoryJson))}
			ctx.Params = append(ctx.Params, gin.Param{Key: "accountId", Value: accountId})

			controllers.CreateCategory(ctx)

			if status := ctx.Writer.Status(); status != http.StatusCreated {
				t.Errorf("expected status code 201, got %d", status)
				return false
			}

			var categoriesResponse httpserialization.JSONResponse
			if err := json.Unmarshal(*body, &categoriesResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if categoriesResponse.Status != "success" {
				t.Errorf("expected status to be success, got %s", categoriesResponse.Status)
				return false
			}

			if categoriesResponse.Data["id"] == nil {
				t.Errorf("expected ID to not be nil")
				return false
			}

			if categoriesResponse.Data["account_id"] != accountId {
				t.Errorf("expected account_id to be %s, got %s", accountId, categoriesResponse.Data["account_id"])
				return false
			}

			if categoriesResponse.Data["name"] != name {
				t.Errorf("expected name to be %s, got %s", name, categoriesResponse.Data["name"])
				return false
			}

			if categoriesResponse.Data["description"] != description {
				t.Errorf("expected description to be %s, got %s", description, categoriesResponse.Data["description"])
				return false
			}

			if categoriesResponse.Data["color"] != color {
				t.Errorf("expected color to be %s, got %s", color, categoriesResponse.Data["color"])
				return false
			}

			if categoriesResponse.Data["budget"] != budget {
				t.Errorf("expected budget to be %f, got %f", budget, categoriesResponse.Data["budget"])
				return false
			}

			if categoriesResponse.Data["current"] != current {
				t.Errorf("expected current to be %f, got %f", current, categoriesResponse.Data["current"])
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

	t.Run("should not create a category if accountId is different from the one in the url", func(t *testing.T) {
		assertion := func(name, description, color string, budget, current float64, accountId string) bool {
			if accountId == "" {
				return true
			}
			ctx, body := getContext()

			category := models.Category{
				AccountID:   accountId,
				Name:        name,
				Description: description,
				Color:       color,
				Budget:      budget,
				Current:     current,
			}

			createReq := controllers.CategoryResponse{}
			createReq.BindModel(category)

			categoryJson, err := json.Marshal(createReq)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			ctx.Request = &http.Request{Body: io.NopCloser(bytes.NewReader(categoryJson))}
			ctx.Params = append(ctx.Params, gin.Param{Key: "accountId", Value: accountId + "1"})

			controllers.CreateCategory(ctx)

			if status := ctx.Writer.Status(); status != http.StatusBadRequest {
				t.Errorf("expected status code 400, got %d", status)
				return false
			}

			var errorResponse httpserialization.ErrorResponse
			if err := json.Unmarshal(*body, &errorResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if !reflect.DeepEqual(errorResponse, customerrors.BadRequestResponse) {
				t.Errorf("expected error response %+v, got %+v", customerrors.BadRequestResponse, errorResponse)
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

	t.Run("should create a category and return specific fields", func(t *testing.T) {
		assertion := func(name, description, color string, budget, current float64, accountId string) bool {
			if accountId == "" {
				return true
			}
			ctx, body := getContext()

			category := models.Category{
				AccountID:   accountId,
				Name:        name,
				Description: description,
				Color:       color,
				Budget:      budget,
				Current:     current,
			}

			createReq := controllers.CategoryResponse{}
			createReq.BindModel(category)

			categoryJson, err := json.Marshal(createReq)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			fields := []string{}
			possibleFields := []string{"account_id", "name", "description", "color", "budget", "current"}
			for range 4 {
				field := possibleFields[rand.Intn(len(possibleFields))]

				if !contains(fields, field) {
					fields = append(fields, field)
				}
			}

			ctx.Request = &http.Request{Body: io.NopCloser(bytes.NewReader(categoryJson))}
			ctx.Params = append(ctx.Params, gin.Param{Key: "accountId", Value: accountId})
			ctx.Set("returnFields", fields)

			controllers.CreateCategory(ctx)

			if status := ctx.Writer.Status(); status != http.StatusCreated {
				t.Errorf("expected status code 201, got %d", status)
				return false
			}

			var categoriesResponse httpserialization.JSONResponse
			if err := json.Unmarshal(*body, &categoriesResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if categoriesResponse.Status != "success" {
				t.Errorf("expected status to be success, got %s", categoriesResponse.Status)
				return false
			}

			categoryMarshalled, err := createReq.Marshal(fields)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if !reflect.DeepEqual(categoriesResponse.Data, categoryMarshalled.Data) {
				t.Errorf("expected %+v, got %+v", categoryMarshalled.Data, categoriesResponse.Data)
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

func TestUpdateCategory(t *testing.T) {
	db := openDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	t.Run("should update a category", func(t *testing.T) {
		assertion := func(name, description, color, accountId string, budget, current float64, newName, newDescription, newColor string, newBudget float64) bool {
			if accountId == "" || newName == "" || newDescription == "" || newColor == "" || newBudget == 0 {
				return true
			}
			ctx, body := getContext()

			category := models.Category{
				AccountID:   accountId,
				Name:        name,
				Description: description,
				Color:       color,
				Budget:      budget,
				Current:     current,
			}

			updateReq := map[string]interface{}{
				"name":        newName,
				"description": newDescription,
				"color":       newColor,
				"budget":      newBudget,
			}

			db.Create(&category)

			categoryJson, err := json.Marshal(updateReq)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			ctx.Params = append(ctx.Params, gin.Param{Key: "accountId", Value: accountId})
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(int(category.ID))})

			ctx.Set("conditions", "account_id = ? AND id = ?")
			ctx.Set("conditions_values", []interface{}{accountId, category.ID})

			ctx.Request = &http.Request{Body: io.NopCloser(bytes.NewReader(categoryJson))}
			controllers.UpdateCategory(ctx)

			if status := ctx.Writer.Status(); status != http.StatusOK {
				t.Errorf("expected status code 200, got %d", status)
				return false
			}

			var categoryResponse httpserialization.JSONResponse
			if err := json.Unmarshal(*body, &categoryResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if categoryResponse.Data["name"] != newName {
				t.Errorf("expected name to be %s, got %s", newName, categoryResponse.Data["name"])
				return false
			}
			if categoryResponse.Data["description"] != newDescription {
				t.Errorf("expected description to be %s, got %s", newDescription, categoryResponse.Data["description"])
				return false
			}
			if categoryResponse.Data["color"] != newColor {
				t.Errorf("expected color to be %s, got %s", newColor, categoryResponse.Data["color"])
				return false
			}
			if categoryResponse.Data["budget"] != newBudget {
				t.Errorf("expected budget to be %f, got %f", newBudget, categoryResponse.Data["budget"])
				return false
			}
			if categoryResponse.Data["current"] != float64(current) {
				t.Errorf("expected current to be %f, got %f", float64(current), categoryResponse.Data["current"])
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

	t.Run("should return 404 if category is not found", func(t *testing.T) {
		assertion := func(name, description, color, accountId string, budget, current float64, newName, newDescription, newColor string, newBudget float64) bool {
			if accountId == "" || newName == "" {
				return true
			}
			ctx, body := getContext()

			updateReq := map[string]interface{}{
				"name":        newName,
				"description": newDescription,
				"color":       newColor,
				"budget":      newBudget,
			}
			categoryJson, err := json.Marshal(updateReq)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			ctx.Params = append(ctx.Params, gin.Param{Key: "accountId", Value: accountId})
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "0"})

			ctx.Set("conditions", "account_id = ? AND id = ?")
			ctx.Set("conditions_values", []interface{}{accountId, 0})

			ctx.Request = &http.Request{Body: io.NopCloser(bytes.NewReader(json.RawMessage(categoryJson)))}
			controllers.UpdateCategory(ctx)

			if status := ctx.Writer.Status(); status != http.StatusNotFound {
				t.Errorf("expected status code 404, got %d", status)
				return false
			}

			var errResponse httpserialization.ErrorResponse
			if err := json.Unmarshal(*body, &errResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if !reflect.DeepEqual(errResponse, controllers.CategoryNotFoundResponse) {
				t.Errorf("expected %+v, got %+v", controllers.CategoryNotFoundResponse, errResponse)
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

	t.Run("should ignore current field", func(t *testing.T) {
		assertion := func(name, description, color, accountId string, budget, current, newCurrent float64) bool {
			if accountId == "" {
				return true
			}
			ctx, body := getContext()

			category := models.Category{
				AccountID:   accountId,
				Name:        name,
				Description: description,
				Color:       color,
				Budget:      budget,
				Current:     current,
			}

			updateReq := map[string]interface{}{
				"current": newCurrent,
			}

			db.Create(&category)

			categoryJson, err := json.Marshal(updateReq)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			ctx.Params = append(ctx.Params, gin.Param{Key: "accountId", Value: accountId})
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(int(category.ID))})

			ctx.Set("conditions", "account_id = ? AND id = ?")
			ctx.Set("conditions_values", []interface{}{accountId, category.ID})

			ctx.Request = &http.Request{Body: io.NopCloser(bytes.NewReader(categoryJson))}
			controllers.UpdateCategory(ctx)

			var expectedResponse controllers.CategoryResponse
			expectedResponse.BindModel(category)
			marshalledResponse, err := expectedResponse.Marshal([]string{})
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			var categoryResponse httpserialization.JSONResponse
			if err := json.Unmarshal(*body, &categoryResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if ctx.Writer.Status() != http.StatusOK {
				t.Errorf("expected status code 200, got %d", ctx.Writer.Status())
				return false
			}

			if !reflect.DeepEqual(categoryResponse, marshalledResponse) {
				t.Errorf("expected %+v, got %+v", marshalledResponse, categoryResponse)
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

	t.Run("should return requested fields", func(t *testing.T) {
		assertion := func(name, description, color, accountId string, budget, current float64, newName, newDescription, newColor string, newBudget float64) bool {
			if accountId == "" || newName == "" || newDescription == "" || newColor == "" || newBudget == 0 {
				return true
			}
			ctx, body := getContext()

			category := models.Category{
				AccountID:   accountId,
				Name:        name,
				Description: description,
				Color:       color,
				Budget:      budget,
				Current:     current,
			}

			db.Create(&category)

			possibleFields := []string{"name", "description", "color", "budget"}
			fields := []string{}
			for range 2 {
				field := possibleFields[rand.Intn(len(possibleFields))]
				if !contains(fields, field) {
					fields = append(fields, field)
				}
			}

			updateReq := map[string]interface{}{
				"name":        newName,
				"description": newDescription,
				"color":       newColor,
				"budget":      newBudget,
			}
			jsonReq, err := json.Marshal(updateReq)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			ctx.Params = append(ctx.Params, gin.Param{Key: "accountId", Value: accountId})
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(int(category.ID))})

			ctx.Set("conditions", "account_id = ? AND id = ?")
			ctx.Set("conditions_values", []interface{}{accountId, category.ID})
			ctx.Set("returnFields", fields)

			ctx.Request = &http.Request{Body: io.NopCloser(bytes.NewReader(json.RawMessage(jsonReq)))}
			controllers.UpdateCategory(ctx)

			var expectedResponse controllers.CategoryResponse
			expectedResponse.BindModel(category)
			marshalledResponse, err := expectedResponse.Marshal(fields)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			var categoryResponse httpserialization.JSONResponse
			if err := json.Unmarshal(*body, &categoryResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if ctx.Writer.Status() != http.StatusOK {
				t.Errorf("expected status code 200, got %d", ctx.Writer.Status())
				return false
			}

			for key := range marshalledResponse.Data {
				marshalledResponse.Data[key] = updateReq[key]
			}

			if !reflect.DeepEqual(categoryResponse, marshalledResponse) {
				t.Errorf("expected %+v, got %+v", marshalledResponse, categoryResponse)
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

func TestDeleteCategory(t *testing.T) {
	db := openDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	t.Run("should delete category", func(t *testing.T) {
		assertion := func(name, description, color, accountId string, budget, current float64) bool {
			if accountId == "" {
				return true
			}
			ctx, _ := getContext()

			category := models.Category{
				AccountID:   accountId,
				Name:        name,
				Description: description,
				Color:       color,
				Budget:      budget,
				Current:     current,
			}

			db.Create(&category)

			ctx.Params = append(ctx.Params, gin.Param{Key: "accountId", Value: accountId})
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(int(category.ID))})
			ctx.Set("conditions", "account_id = ? AND id = ?")
			ctx.Set("conditions_values", []interface{}{accountId, category.ID})
			controllers.DeleteCategory(ctx)
			if ctx.Writer.Status() != http.StatusNoContent {
				t.Errorf("expected status code 204, got %d", ctx.Writer.Status())
				return false
			}

			var nonExistentCategory models.Category
			db.Where("account_id = ? AND id = ?", accountId, category.ID).First(&nonExistentCategory)
			if nonExistentCategory.ID != 0 {
				t.Errorf("expected category to be deleted, got %+v", nonExistentCategory)
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

	t.Run("should return 404 if category does not exist", func(t *testing.T) {
		assertion := func(accountId string, randomId uint) bool {
			if accountId == "" {
				return true
			}
			ctx, body := getContext()

			ctx.Params = append(ctx.Params, gin.Param{Key: "accountId", Value: accountId})
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(int(randomId))})

			ctx.Set("conditions", "account_id = ? AND id = ?")
			ctx.Set("conditions_values", []interface{}{accountId, randomId})
			controllers.DeleteCategory(ctx)

			if ctx.Writer.Status() != http.StatusNotFound {
				t.Errorf("expected status code 404, got %d", ctx.Writer.Status())
				return false
			}

			var errorResponse httpserialization.ErrorResponse
			if err := json.Unmarshal(*body, &errorResponse); err != nil {
				t.Errorf("expected no error, got %s", err)
				return false
			}

			if !reflect.DeepEqual(errorResponse, controllers.CategoryNotFoundResponse) {
				t.Errorf("expected %+v, got %+v", controllers.CategoryNotFoundResponse, errorResponse)
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

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func formatQueryValue(queryValue string) string {
	escaped := url.QueryEscape(queryValue)
	return strings.ReplaceAll(strings.ReplaceAll(escaped, "+", "%20"), "%", "")
}
