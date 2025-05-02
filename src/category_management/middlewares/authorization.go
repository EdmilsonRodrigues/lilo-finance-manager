package middlewares

import (
	"encoding/json"
	"log"
	"net/http"

	customerrors "github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/custom_errors"
	"github.com/gin-gonic/gin"
)

// AuthorizationMiddleware is a middleware that checks if the user has access to the account with the given accountId.
// It checks the X-User-Id and X-User-Accounts headers and verifies that the user has access to the account.
// If the user does not have access, it returns a 403 response.
// If the headers are missing or invalid, it returns a 400 or 403 response.
// If the accountId parameter is missing, it returns a 403 response.
// If the user has access, it sets the role and accountId in the gin context and calls the next handler.
//
// Returns:
//   - gin.HandlerFunc: the middleware function
//
// Example:
//   group.GET("/:accountId", AuthorizationMiddleware(), categoryController.GetCategory)
func AuthorizationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, userAccounts := ctx.GetHeader("X-User-Id"), ctx.GetHeader("X-User-Accounts")
		if userId == "" || userAccounts == "" {
			log.Println("Missing headers X-User-Id or X-User-Accounts for request ", ctx.Request.URL, " with params ", ctx.Params)
			ctx.JSON(http.StatusForbidden, customerrors.ForbiddenResponse)
			ctx.Abort()
			return
		}

		requiredAccountID := ctx.Param("accountId")
		if requiredAccountID == "" {
			log.Printf("AuthorizationMiddleware applied to route missing accountId parameter: %s", ctx.Request.URL.Path)
			ctx.JSON(http.StatusForbidden, customerrors.ForbiddenResponse)
			ctx.Abort()
			return
		}

		accountRolesMap := make(map[string]string)
		err := json.Unmarshal([]byte(userAccounts), &accountRolesMap)
		if err != nil {
			log.Printf("Failed to parse X-User-Accounts header as JSON: %v for request %s, header value: %s", err, ctx.Request.URL.Path, userAccounts)
			ctx.JSON(http.StatusUnprocessableEntity, customerrors.UnprocessableEntityResponse)
			ctx.Abort()
			return
		}

		role, hasAccess := accountRolesMap[requiredAccountID]
		if !hasAccess {
			log.Printf("User with id %s does not have access to account with id %s", userId, requiredAccountID)
			ctx.JSON(http.StatusForbidden, customerrors.ForbiddenResponse)
			ctx.Abort()
			return
		}

		ctx.Set("role", role)
		ctx.Next()
	}
}

// AdminOnlyMiddleware is a middleware that checks if the user has admin role.
// If the user does not have admin role, it returns a 403 response.
// If the user has admin role, it calls the next handler.
//
// Returns:
//   - gin.HandlerFunc: the middleware function
//
// Example:
//   group.GET("/", AdminOnlyMiddleware(), categoryController.GetCategories)
func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.GetString("role")
		if role != "admin" {
			log.Println("User with id ", ctx.GetHeader("X-User-Id"), " does not have admin role")
			ctx.JSON(403, customerrors.ForbiddenResponse)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// EditorOrAdminOnlyMiddleware is a middleware that checks if the user has admin or editor role.
// If the user does not have either role, it returns a 403 response.
// If the user has either admin or editor role, it calls the next handler.
//
// Returns:
//   - gin.HandlerFunc: the middleware function
//
// Example:
//   group.GET("/", EditorOrAdminOnlyMiddleware(), categoryController.GetCategories)
func EditorOrAdminOnlyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.GetString("role")
		if role != "admin" && role != "editor" {
			log.Println("User with id ", ctx.GetHeader("X-User-Id"), " does not have admin or editor role")
			ctx.JSON(403, customerrors.ForbiddenResponse)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
