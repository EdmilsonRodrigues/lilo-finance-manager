package middlewares

import (
	"encoding/json"
	"log"
	"net/http"

	customerrors "github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/custom_errors"
	"github.com/gin-gonic/gin"
)

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
			ctx.JSON(http.StatusBadRequest, customerrors.BadRequestResponse)
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
		ctx.Set("accountId", requiredAccountID)
		ctx.Next()
	}
}

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
