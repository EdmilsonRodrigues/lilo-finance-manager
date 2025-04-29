package router

import (
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/controllers"
	"github.com/gin-gonic/gin"
)

// HandleRequests sets up the routes and starts the Gin server.
func HandleRequests(engine *gin.Engine) {
	baseCategoryPath := "/categories/:accountId"
	engine.GET(baseCategoryPath, AddInitialCondsDecorator(controllers.GetCategories))
	engine.GET(baseCategoryPath+"/:id", AddInitialCondsDecorator(controllers.GetCategory))
	engine.POST(baseCategoryPath, controllers.CreateCategory)
	engine.PATCH(baseCategoryPath+"/:id", AddInitialCondsDecorator(controllers.UpdateCategory))
	engine.DELETE(baseCategoryPath+"/:id", AddInitialCondsDecorator(controllers.DeleteCategory))
}

func AddInitialCondsDecorator(function func (ctx *gin.Context, conds controllers.QueryConditions)) (func (ctx *gin.Context)) {
	return func(ctx *gin.Context) {
		conds := make(controllers.QueryConditions)
		conds["account_id"] = ctx.Param("accountId")
		function(ctx, conds)
	}
}
