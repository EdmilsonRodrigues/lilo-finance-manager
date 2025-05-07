//go:build unit

package middlewares_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"testing/quick"

	customerrors "github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/custom_errors"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/middlewares"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/common_utils/go/serialization/http_serialization"
	"github.com/gin-gonic/gin"
)

func TestAddParamToConditionsMiddleware(t *testing.T) {
	t.Run("should add param to conditions when param value is not empty and conditions not initially set", func(t *testing.T) {
		assertion := func(paramName, columnName, value string) bool {
			if value == "" || columnName == "" || paramName == "" {
				return true
			}
			ctx, _ := getContext()
			ctx.Params = gin.Params{
				gin.Param{
					Key:   paramName,
					Value: value,
				},
			}
			middleware := middlewares.AddParamToConditionsMiddleware(paramName, columnName)
			middleware(ctx)

			conditions := ctx.GetString("conditions")
			conditionsValues, exists := ctx.Get("conditions_values")
			if !exists {
				t.Error("conditions values not set in context")
				return false
			}

			cValues := make([]string, len(conditionsValues.([]interface{})))
			for i, v := range conditionsValues.([]interface{}) {
				cValues[i] = v.(string)
			}

			if ctx.IsAborted() {
				t.Error("middleware aborted context")
				return false
			}

			if conditions != fmt.Sprintf("%s = ?", columnName) {
				t.Errorf("expected %v, got %v", fmt.Sprintf("%s = ?", columnName), conditions)
				return false
			}

			if !reflect.DeepEqual(cValues, []string{value}) {
				t.Errorf("expected %v, got %v", []string{value}, cValues)
				return false
			}

			return true
		}

		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error("failed checks", err)
		}
	})

	t.Run("should add param to conditions when param value is not empty and conditions already set", func(t *testing.T) {
		assertion := func(paramName, columnName, value string) bool {
			if value == "" || columnName == "" || paramName == "" {
				return true
			}
			ctx, _ := getContext()
			ctx.Params = gin.Params{
				gin.Param{
					Key:   paramName,
					Value: value,
				},
			}
			middleware := middlewares.AddParamToConditionsMiddleware(paramName, columnName)
			initialConditions := "initialColumn = ?"
			ctx.Set("conditions", initialConditions)
			ctx.Set("conditions_values", []interface{}{"initialValue"})

			middleware(ctx)

			conditions := ctx.GetString("conditions")
			conditionsValues, exists := ctx.Get("conditions_values")
			if !exists {
				t.Error("conditions values not set in context")
				return false
			}

			cValues := make([]string, len(conditionsValues.([]interface{})))
			for i, v := range conditionsValues.([]interface{}) {
				cValues[i] = v.(string)
			}

			if ctx.IsAborted() {
				t.Error("middleware aborted context")
				return false
			}

			if !reflect.DeepEqual(cValues, []string{"initialValue", value}) {
				t.Errorf("expected %v, got %v", []string{"initialValue", value}, cValues)
				return false
			}

			if conditions != fmt.Sprintf("%s AND %s = ?", initialConditions, columnName) {
				t.Errorf("expected %v, got %v", fmt.Sprintf("%s AND %s = ?", initialConditions, columnName), conditions)
				return false
			}

			return true
		}

		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error("failed checks for Run ", err)
		}
	})

	t.Run("should not add param to conditions when param value is empty", func(t *testing.T) {
		assertion := func(paramName, columnName string) bool {
			if columnName == "" || paramName == "" {
				return true
			}
			ctx, body := getContext()
			ctx.Params = gin.Params{
				gin.Param{
					Key:   paramName,
					Value: "",
				},
			}
			ctx.Request = &http.Request{
				URL: &url.URL{
					Path: "/test",
				},
			}

			middleware := middlewares.AddParamToConditionsMiddleware(paramName, columnName)
			middleware(ctx)
			if ctx.Writer.Status() != http.StatusUnprocessableEntity {
				t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, ctx.Writer.Status())
				return false
			}

			var unmarshalledBody httpserialization.ErrorResponse
			err := json.Unmarshal(*body, &unmarshalledBody)
			if err != nil {
				t.Errorf("error unmarshalling body: %v", err)
				return false
			}

			if !reflect.DeepEqual(unmarshalledBody, customerrors.UnprocessableEntityResponse) {
				t.Errorf("expected %v, got %v", customerrors.UnprocessableEntityResponse, unmarshalledBody)
				return false
			}

			if !ctx.IsAborted() {
				t.Error("middleware aborted context")
				return false
			}

			_, exists := ctx.Get("conditions")
			return !exists
		}

		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error("failed checks", err)
		}
	})

}

func TestParsePaginationParamsMiddleware(t *testing.T) {
	t.Run("should be able to parse page and page_size params", func(t *testing.T) {
		assertion := func(page, pageSize int) bool {
			ctx, _ := getContext()
			ctx.Request = &http.Request{
				URL: &url.URL{
					Path:     "/test",
					RawQuery: "page=" + strconv.Itoa(page) + "&page_size=" + strconv.Itoa(pageSize),
				},
			}
			middleware := middlewares.ParsePaginationParamsMiddleware()
			middleware(ctx)

			pageNumber, exists := ctx.Get("page")

			t.Log("page number: ", pageNumber)
			switch {
			case !exists:
				t.Error("page number not set in context")
				return false
			case page <= 0 && pageNumber != 1:
				t.Error("page number should be 1 when page is less than or equal to 0, got ", pageNumber)
				return false
			case page > 0 && pageNumber != page:
				t.Error("page number should be set to ", page, " when page is greater than 0, got ", pageNumber)
				return false
			}

			t.Log("page size: ", pageSize)
			pageSizeNumber, exists := ctx.Get("pageSize")
			switch {
			case !exists:
				t.Error("page size not set in context")
				return false
			case pageSize <= 0 && pageSizeNumber != 10:
				t.Error("page size should be 10 when page size is less than or equal to 0, got ", pageSizeNumber)
				return false
			case pageSize >= 100 && pageSizeNumber != 100:
				t.Error("page size should be 100 when page size is greater than or equal to 100, got ", pageSizeNumber)
				return false
			case pageSize > 0 && pageSize < 100 && pageSizeNumber != pageSize:
				t.Error("page size should be set to ", pageSize, " got ", pageSizeNumber)
				return false
			}

			if ctx.IsAborted() {
				t.Error("middleware aborted context")
				return false
			}

			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error("failed checks", err)
		}
	})

	t.Run("should be able to parse page and page_size params when page and page_size are not set", func(t *testing.T) {
		assertion := func() bool {
			ctx, _ := getContext()
			ctx.Request = &http.Request{
				URL: &url.URL{
					Path: "/test",
				},
			}
			middleware := middlewares.ParsePaginationParamsMiddleware()
			middleware(ctx)

			pageNumber, exists := ctx.Get("page")
			if !exists {
				t.Error("page number not set in context")
				return false
			}
			pageSizeNumber, exists := ctx.Get("pageSize")
			if !exists {
				t.Error("page size not set in context")
				return false
			}

			if ctx.IsAborted() {
				t.Error("middleware aborted context")
				return false
			}

			return pageNumber == 1 && pageSizeNumber == 10
		}

		if !assertion() {
			t.Error("failed checks")
		}
	})

	t.Run("should not be able to parse page_size param when it is not integers", func(t *testing.T) {
		assertion := func(page int, pageSize string) bool {
			pageSize = formatQueryValue(pageSize)
			_, err := strconv.Atoi(pageSize)
			if pageSize == "" || err == nil {
				return true
			}
			ctx, body := getContext()
			ctx.Request = &http.Request{
				URL: &url.URL{
					Path:     "/test",
					RawQuery: "page=" + strconv.Itoa(page) + "&page_size=" + pageSize,
				},
			}
			middleware := middlewares.ParsePaginationParamsMiddleware()
			middleware(ctx)
			if ctx.Writer.Status() != http.StatusUnprocessableEntity {
				t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, ctx.Writer.Status())
				return false
			}
			var unmarshalledBody httpserialization.ErrorResponse

			if err := json.Unmarshal(*body, &unmarshalledBody); err != nil {
				t.Errorf("error unmarshalling body: %v", err)
				return false
			}

			if !reflect.DeepEqual(unmarshalledBody, customerrors.UnprocessableEntityResponse) {
				t.Errorf("expected %v, got %v", customerrors.UnprocessableEntityResponse, unmarshalledBody)
				return false
			}

			if !ctx.IsAborted() {
				t.Error("middleware didn't abort context")
				return false
			}

			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error("failed checks", err)
		}
	})

	t.Run("should not be able to parse page param when it is not integers", func(t *testing.T) {
		assertion := func(page string, pageSize int) bool {
			page = formatQueryValue(page)
			_, err := strconv.Atoi(page)
			if page == "" || err == nil {
				return true
			}
			ctx, body := getContext()
			ctx.Request = &http.Request{
				URL: &url.URL{
					Path:     "/test",
					RawQuery: "page=" + page + "&page_size=" + strconv.Itoa(pageSize),
				},
			}
			middleware := middlewares.ParsePaginationParamsMiddleware()
			middleware(ctx)
			if ctx.Writer.Status() != http.StatusUnprocessableEntity {
				t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, ctx.Writer.Status())
				return false
			}
			var unmarshalledBody httpserialization.ErrorResponse

			if err := json.Unmarshal(*body, &unmarshalledBody); err != nil {
				t.Errorf("error unmarshalling body: %v", err)
				return false
			}

			if !reflect.DeepEqual(unmarshalledBody, customerrors.UnprocessableEntityResponse) {
				t.Errorf("expected %v, got %v", customerrors.UnprocessableEntityResponse, unmarshalledBody)
				return false
			}

			if !ctx.IsAborted() {
				t.Error("middleware didn't abort context")
				return false
			}

			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error("failed checks", err)
		}
	})
}

func TestParseFiltersMiddleware(t *testing.T) {
	t.Run("should be able to parse filters", func(t *testing.T) {
		assertion := func(nameFilter, idFilter string) bool {
			if nameFilter == "" && idFilter == "" {
				return true
			}
			nameFilter = formatQueryValue(nameFilter)
			idFilter = formatQueryValue(idFilter)
			ctx, _ := getContext()
			ctx.Request = &http.Request{
				URL: &url.URL{
					Path:     "/test",
					RawQuery: fmt.Sprintf("filters=%s:%s,%s:%s", "name", formatQueryValue(nameFilter), "id", formatQueryValue(idFilter)),
				},
			}
			ctx.Set("conditions", "account_id = ?")
			ctx.Set("conditions_values", []interface{}{"account_id"})
			middleware := middlewares.ParseFiltersMiddleware()
			middleware(ctx)

			conditions := ctx.GetString("conditions")
			expectedConditions := "account_id = ? AND id = ? AND name = ?"
			conditionsValues, exists := ctx.Get("conditions_values")
			if !exists {
				t.Error("conditions values not set in context")
				return false
			}

			cValues := make([]string, len(conditionsValues.([]interface{})))
			for i, v := range conditionsValues.([]interface{}) {
				cValues[i] = v.(string)
			}

			expectedValues := []string{"account_id", formatQueryValue(idFilter), formatQueryValue(nameFilter)}
			if !assertEqualStringSlice(cValues, expectedValues) {
				t.Errorf("expected %v, got %v", expectedValues, cValues)
				return false
			}

			if !assertEqualConditions(conditions, expectedConditions) {
				t.Errorf("expected %v, got %v", expectedConditions, conditions)
				return false
			}

			if ctx.IsAborted() {
				t.Error("middleware aborted context")
				return false
			}

			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error("failed checks", err)
		}
	})

	t.Run("should not be able to parse filters when filters are not in the correct format", func(t *testing.T) {
		assertion := func(nameFilter, idFilter string) bool {
			nameFilter = formatQueryValue(nameFilter)
			idFilter = formatQueryValue(idFilter)
			ctx, body := getContext()
			ctx.Request = &http.Request{
				URL: &url.URL{
					Path:     "/test",
					RawQuery: fmt.Sprintf("filters=%s:%s,,%s:%s", "name", nameFilter, "id", idFilter),
				},
			}
			middleware := middlewares.ParseFiltersMiddleware()
			middleware(ctx)
			if ctx.Writer.Status() != http.StatusUnprocessableEntity {
				t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, ctx.Writer.Status())
				return false
			}
			var unmarshalledBody httpserialization.ErrorResponse
			err := json.Unmarshal(*body, &unmarshalledBody)
			if err != nil {
				t.Errorf("error unmarshalling body: %v", err)
				return false
			}
			if !reflect.DeepEqual(unmarshalledBody, customerrors.UnprocessableEntityResponse) {
				t.Errorf("expected %v, got %v", customerrors.UnprocessableEntityResponse, unmarshalledBody)
				return false
			}

			if !ctx.IsAborted() {
				t.Error("middleware didn't abort context")
				return false
			}

			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error("failed checks", err)
		}
	})
}

func TestParseReturnFieldsMiddleware(t *testing.T) {
	t.Run("should be able to parse return fields", func(t *testing.T) {
		assertion := func(field1, field2 string) bool {
			field1, field2 = formatQueryValue(field1), formatQueryValue(field2)
			if field1 == "" || field2 == "" {
				return true
			}
			ctx, _ := getContext()
			ctx.Request = &http.Request{
				URL: &url.URL{
					Path:     "/test",
					RawQuery: fmt.Sprintf("return_fields=%s,%s", field1, field2),
				},
			}
			middleware := middlewares.ParseReturnFieldsMiddleware()
			middleware(ctx)
			returnFields, exists := ctx.Get("returnFields")
			if !exists {
				t.Error("returnFields not set in context")
				return false
			}
			if !reflect.DeepEqual(returnFields, []string{field1, field2}) {
				t.Errorf("expected %v, got %v", []string{field1, field2}, returnFields)
				return false
			}

			if ctx.IsAborted() {
				t.Error("middleware aborted context")
				return false
			}

			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error("failed checks", err)
		}
	})

	t.Run("should be able to parse return fields and ignore empty fields", func(t *testing.T) {
		assertion := func(field1, field2 string) bool {
			field1, field2 = formatQueryValue(field1), formatQueryValue(field2)
			if field1 == "" || field2 == "" {
				return true
			}
			ctx, _ := getContext()
			ctx.Request = &http.Request{
				URL: &url.URL{
					Path:     "/test",
					RawQuery: fmt.Sprintf("return_fields=%s,,%s", field1, field2),
				},
			}
			middleware := middlewares.ParseReturnFieldsMiddleware()
			middleware(ctx)
			returnFields, exists := ctx.Get("returnFields")
			if !exists {
				t.Error("returnFields not set in context")
				return false
			}
			if !reflect.DeepEqual(returnFields, []string{field1, field2}) {
				t.Errorf("expected %v, got %v", []string{field1, field2}, returnFields)
				return false
			}

			if ctx.IsAborted() {
				t.Error("middleware aborted context")
				return false
			}

			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error("failed checks", err)
		}
	})
}

func assertEqualConditions(conditions1, conditions2 string) bool {
	if conditions1 == conditions2 {
		return true
	}
	c1 := strings.Split(conditions1, " AND ")
	c2 := strings.Split(conditions2, " AND ")
	return assertEqualStringSlice(c1, c2)
}

func assertEqualStringSlice(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if !contains(slice2, slice1[i]) {
			return false
		}
	}
	return true
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
