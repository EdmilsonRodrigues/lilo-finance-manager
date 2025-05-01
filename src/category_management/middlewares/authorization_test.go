//go:build unit

package middlewares_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"testing/quick"

	customerrors "github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/custom_errors"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/middlewares"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/common_utils/go/serialization"
	"github.com/gin-gonic/gin"
)

func TestAuthorizationMiddleware(t *testing.T) {
	t.Run("should authenticate user if the user has access to the account", func(t *testing.T) {
		assertion := func(accountId, userId, otherAccountId string) bool {
			ctx, _ := getContext()
			ctx.Request = &http.Request{
				Header: map[string][]string{
					"X-User-Id":       {userId},
					"X-User-Accounts": {fmt.Sprintf(`{"%s": "admin", "%s": "user"}`, accountId, otherAccountId)},
				},
			}
			ctx.Params = []gin.Param{{Key: "accountId", Value: accountId}}
			middleware := middlewares.AuthorizationMiddleware()
			middleware(ctx)

			if ctx.IsAborted() {
				t.Error("middleware aborted ")
				return false
			}

			if ctx.GetString("role") != "admin" {
				t.Error("role not set in context")
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

	t.Run("should not authenticate user if the user does not have access to the account", func(t *testing.T) {
		assertion := func(accountId, userId, otherAccountId string) bool {
			ctx, body := getContext()
			ctx.Request = &http.Request{
				Header: map[string][]string{
					"X-User-Id":       {userId},
					"X-User-Accounts": {fmt.Sprintf(`{"%s": "user"}`, otherAccountId)},
				},
				URL: &url.URL{Path: fmt.Sprintf("/accounts/%s", accountId)},
			}
			ctx.Params = []gin.Param{{Key: "accountId", Value: accountId}}
			middleware := middlewares.AuthorizationMiddleware()
			middleware(ctx)

			if !ctx.IsAborted() {
				t.Error("middleware did not abort")
				return false
			}

			if _, present := ctx.Get("role"); present {
				t.Error("role set in context")
				return false
			}

			if ctx.Writer.Status() != http.StatusForbidden {
				t.Error("status code is not forbidden")
				return false
			}

			var unmarshalledBody serialization.ErrorResponse
			if err := json.Unmarshal(*body, &unmarshalledBody); err != nil {
				t.Error(err)
				return false
			}

			if !reflect.DeepEqual(unmarshalledBody, customerrors.ForbiddenResponse) {
				t.Errorf("expected %v, got %v", customerrors.ForbiddenResponse, unmarshalledBody)
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

	t.Run("should not authenticate user if the user id is not set in the headers", func(t *testing.T) {
		assertion := func(accountId, userId, otherAccountId string) bool {
			ctx, body := getContext()
			ctx.Request = &http.Request{
				Header: map[string][]string{
					"X-User-Accounts": {fmt.Sprintf(`{"%s": "user"}`, otherAccountId)},
				},
				URL: &url.URL{Path: fmt.Sprintf("/accounts/%s", accountId)},
			}
			ctx.Params = []gin.Param{{Key: "accountId", Value: accountId}}
			middleware := middlewares.AuthorizationMiddleware()
			middleware(ctx)
			if !ctx.IsAborted() {
				t.Error("middleware did not abort")
				return false
			}

			if _, present := ctx.Get("role"); present {
				t.Error("role set in context")
				return false
			}

			if ctx.Writer.Status() != http.StatusForbidden {
				t.Error("status code is not forbidden")
				return false
			}

			var unmarshalledBody serialization.ErrorResponse
			if err := json.Unmarshal(*body, &unmarshalledBody); err != nil {
				t.Error(err)
				return false
			}

			if !reflect.DeepEqual(unmarshalledBody, customerrors.ForbiddenResponse) {
				t.Errorf("expected %v, got %v", customerrors.ForbiddenResponse, unmarshalledBody)
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

	t.Run("should not authenticate user if the user account is not set in the headers", func(t *testing.T) {
		assertion := func(accountId, userId, otherAccountId string) bool {
			ctx, body := getContext()
			ctx.Request = &http.Request{
				Header: map[string][]string{
					"X-User-Id": {userId},
				},
				URL: &url.URL{Path: fmt.Sprintf("/accounts/%s", accountId)},
			}
			ctx.Params = []gin.Param{{Key: "accountId", Value: accountId}}
			middleware := middlewares.AuthorizationMiddleware()
			middleware(ctx)
			if !ctx.IsAborted() {
				t.Error("middleware did not abort")
				return false
			}

			if _, present := ctx.Get("role"); present {
				t.Error("role set in context")
				return false
			}

			if ctx.Writer.Status() != http.StatusForbidden {
				t.Error("status code is not forbidden")
				return false
			}
			var unmarshalledBody serialization.ErrorResponse
			if err := json.Unmarshal(*body, &unmarshalledBody); err != nil {
				t.Error(err)
				return false
			}
			if !reflect.DeepEqual(unmarshalledBody, customerrors.ForbiddenResponse) {
				t.Errorf("expected %v, got %v", customerrors.ForbiddenResponse, unmarshalledBody)
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

	t.Run("should not authenticate user if the account id is not set in the params", func(t *testing.T) {
		assertion := func(accountId, userId, otherAccountId string) bool {
			ctx, body := getContext()
			ctx.Request = &http.Request{
				Header: map[string][]string{
					"X-User-Id":       {userId},
					"X-User-Accounts": {fmt.Sprintf(`{"%s": "user"}`, otherAccountId)},
				},
				URL: &url.URL{Path: fmt.Sprintf("/accounts/%s", accountId)},
			}
			middleware := middlewares.AuthorizationMiddleware()
			middleware(ctx)
			if !ctx.IsAborted() {
				t.Error("middleware did not abort")
				return false
			}

			if _, present := ctx.Get("role"); present {
				t.Error("role set in context")
				return false
			}

			if ctx.Writer.Status() != http.StatusForbidden {
				t.Error("status code is not forbidden")
				return false
			}

			var unmarshalledBody serialization.ErrorResponse
			if err := json.Unmarshal(*body, &unmarshalledBody); err != nil {
				t.Error(err)
				return false
			}

			if !reflect.DeepEqual(unmarshalledBody, customerrors.ForbiddenResponse) {
				t.Errorf("expected %v, got %v", customerrors.ForbiddenResponse, unmarshalledBody)
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

	t.Run("should not authenticate user if the user accounts are in a wrong format", func(t *testing.T) {
		assertion := func(accountId, userId, otherAccountId string) bool {
			ctx, body := getContext()
			ctx.Request = &http.Request{
				Header: map[string][]string{
					"X-User-Id":       {userId},
					"X-User-Accounts": {fmt.Sprintf(`{"%s": "user"}`, otherAccountId)},
				},
				URL: &url.URL{Path: fmt.Sprintf("/accounts/%s", accountId)},
			}
			ctx.Params = []gin.Param{{Key: "accountId", Value: accountId}}
			middleware := middlewares.AuthorizationMiddleware()
			middleware(ctx)
			if !ctx.IsAborted() {
				t.Error("middleware did not abort")
				return false
			}

			if _, present := ctx.Get("role"); present {
				t.Error("role set in context")
				return false
			}

			if ctx.Writer.Status() != http.StatusForbidden {
				t.Error("status code is not forbidden")
				return false
			}

			var unmarshalledBody serialization.ErrorResponse
			if err := json.Unmarshal(*body, &unmarshalledBody); err != nil {
				t.Error(err)
				return false
			}

			if !reflect.DeepEqual(unmarshalledBody, customerrors.ForbiddenResponse) {
				t.Errorf("expected %v, got %v", customerrors.ForbiddenResponse, unmarshalledBody)
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

func TestAdminOnlyMiddleware(t *testing.T) {
	t.Run("should authenticate user if the user is admin", func(t *testing.T) {
		ctx, _ := getContext()
		ctx.Set("role", "admin")
		middleware := middlewares.AdminOnlyMiddleware()
		middleware(ctx)
		if ctx.IsAborted() {
			t.Error("middleware should not abort the context")
		}
	})
	t.Run("should not authenticate user if the user is not admin", func(t *testing.T) {
		assertion := func(role string) bool {
			ctx, body := getContext()
			ctx.Request = &http.Request{
				Header: map[string][]string{
					"X-User-Id": {"id"},
				},
				URL: &url.URL{Path: "/accounts"},
			}
			if role == "admin" {
				return true
			}
			ctx.Set("role", role)
			middleware := middlewares.AdminOnlyMiddleware()
			middleware(ctx)
			if !ctx.IsAborted() {
				t.Error("middleware should abort the context")
				return false
			}
			if ctx.Writer.Status() != http.StatusForbidden {
				t.Error("middleware should return forbidden status")
				return false
			}

			var unmarshalledBody serialization.ErrorResponse
			if err := json.Unmarshal(*body, &unmarshalledBody); err != nil {
				t.Error(err)
				return false
			}

			if !reflect.DeepEqual(unmarshalledBody, customerrors.ForbiddenResponse) {
				t.Errorf("expected %v, got %v", customerrors.ForbiddenResponse, unmarshalledBody)
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

func TestEditorOrAdminOnlyMiddleware(t *testing.T) {
	t.Run("should authenticate user if the user is admin or editor", func(t *testing.T) {
		ctx, _ := getContext()
		ctx.Set("role", "admin")
		middleware := middlewares.EditorOrAdminOnlyMiddleware()
		middleware(ctx)
		if ctx.IsAborted() {
			t.Error("middleware should not abort the context")
		}
	})
	t.Run("should authenticate user if the user is editor", func(t *testing.T) {
		ctx, _ := getContext()
		ctx.Request = &http.Request{
			Header: map[string][]string{
				"X-User-Id": {"id"},
			},
			URL: &url.URL{Path: "/accounts"},
		}
		ctx.Set("role", "editor")
		middleware := middlewares.EditorOrAdminOnlyMiddleware()
		middleware(ctx)
		if ctx.IsAborted() {
			t.Error("middleware should not abort the context")
		}
	})
	t.Run("should not authenticate user if the user is not admin or editor", func(t *testing.T) {
		assertion := func(role string) bool {
			ctx, body := getContext()
			ctx.Request = &http.Request{
				Header: map[string][]string{
					"X-User-Id": {"id"},
				},
				URL: &url.URL{Path: "/accounts"},
			}
			if role == "admin" || role == "editor" {
				return true
			}
			ctx.Set("role", role)
			middleware := middlewares.EditorOrAdminOnlyMiddleware()
			middleware(ctx)
			if !ctx.IsAborted() {
				t.Error("middleware should abort the context")
				return false
			}
			if ctx.Writer.Status() != http.StatusForbidden {
				t.Error("middleware should return forbidden status")
				return false
			}
			var unmarshalledBody serialization.ErrorResponse
			if err := json.Unmarshal(*body, &unmarshalledBody); err != nil {
				t.Error(err)
				return false
			}
			if !reflect.DeepEqual(unmarshalledBody, customerrors.ForbiddenResponse) {
				t.Errorf("expected %v, got %v", customerrors.ForbiddenResponse, unmarshalledBody)
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
