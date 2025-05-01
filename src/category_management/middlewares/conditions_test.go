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
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/common_utils/go/serialization"
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

			conditions, exists := ctx.Get("conditions")
			if !exists {
				t.Error("conditions not set in context")
				return false
			}

			conds, ok := conditions.(*serialization.QueryConditions)
			if !ok {
				t.Error("conditions not parsed correctly")
				return false
			}

			if ctx.IsAborted() {
				t.Error("middleware aborted context")
				return false
			}

			return (*conds)[columnName] == value
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
			initialConditions := &serialization.QueryConditions{"initialColumn": "initialValue"}
			ctx.Set("conditions", initialConditions)
			middleware(ctx)
			conditions, exists := ctx.Get("conditions")
			if !exists {
				t.Error("conditions not set in context")
				return false
			}

			conds, ok := conditions.(*serialization.QueryConditions)
			if !ok {
				t.Error("conditions not parsed correctly")
				return false
			}

			if ctx.IsAborted() {
				t.Error("middleware aborted context")
				return false
			}

			return (*conds)[columnName] == value && (*conds)["initialColumn"] == "initialValue"
		}

		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error("failed checks", err)
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

			var unmarshalledBody serialization.ErrorResponse
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
			var unmarshalledBody serialization.ErrorResponse
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

	t.Run("should not be able to parse page param when it is not integers", func(t *testing.T) {
		assertion := func(page string, pageSize int) bool {
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
			var unmarshalledBody serialization.ErrorResponse
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

func TestParseFiltersMiddleware(t *testing.T) {
	t.Run("should be able to parse filters", func(t *testing.T) {
		assertion := func(nameFilter, idFilter string) bool {
			nameFilter = formatQueryValue(nameFilter)
			idFilter = formatQueryValue(idFilter)
			ctx, _ := getContext()
			ctx.Request = &http.Request{
				URL: &url.URL{
					Path:     "/test",
					RawQuery: fmt.Sprintf("filters=%s:%s,%s:%s", "name", nameFilter, "id", idFilter),
				},
			}
			middleware := middlewares.ParseFiltersMiddleware()
			middleware(ctx)

			conditions, exists := ctx.Get("conditions")
			if !exists {
				t.Error("filters not set in context")
				return false
			}

			conds, ok := conditions.(*serialization.QueryConditions)
			if !ok {
				t.Error("filters not in the correct format: ", conditions)
				return false
			}
			expected := map[string]string{
				"name": nameFilter,
				"id":   idFilter,
			}

			if !reflect.DeepEqual((*conds)["filters"], expected) {
				t.Errorf("expected %v, got %v", expected, (*conds)["filters"])
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
			var unmarshalledBody serialization.ErrorResponse
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

func formatQueryValue(queryValue string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(queryValue, "&", ""),
				" ",
				"",
			),
			":",
			"",
		),
		",",
		"",
	)

}
