package router

import (
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/controllers"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/middlewares"
	"github.com/gin-gonic/gin"
)

// HandleRequests sets up the routes and starts the Gin server.
func HandleRequests(engine *gin.Engine) {
	group := engine.Group("/categories/:accountId")
	group.Use(
		middlewares.AuthorizationMiddleware(),
		middlewares.AddParamToConditionsMiddleware("accountId", "account_id"),
	)
	group.GET("/",
		middlewares.ParsePaginationParamsMiddleware(),
		middlewares.ParseFiltersMiddleware(),
		middlewares.ParseReturnFieldsMiddleware(true),
		controllers.GetCategories,
	)
	group.GET("/:id",
		middlewares.AddParamToConditionsMiddleware("id", "id"),
		middlewares.ParseReturnFieldsMiddleware(),
		controllers.GetCategory,
	)
	group.POST("/",
		middlewares.AdminOnlyMiddleware(),
		middlewares.ParseReturnFieldsMiddleware(),
		controllers.CreateCategory,
	)
	group.PATCH("/:id",
		middlewares.AdminOnlyMiddleware(),
		middlewares.AddParamToConditionsMiddleware("id", "id"),
		middlewares.ParseReturnFieldsMiddleware(),
		controllers.UpdateCategory,
	)
	group.DELETE("/:id",
		middlewares.AdminOnlyMiddleware(),
		middlewares.AddParamToConditionsMiddleware("id", "id"),
		controllers.DeleteCategory,
	)
}
