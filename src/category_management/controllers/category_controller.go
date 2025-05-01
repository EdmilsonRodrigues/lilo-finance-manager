package controllers

import (
	"fmt"
	"log"
	"net/http"

	customerrors "github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/custom_errors"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/database"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/models"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/common_utils/go/serialization"
	"github.com/gin-gonic/gin"
)

var CategoryNotFoundResponse = serialization.ErrorResponse{Details: serialization.ErrorDetails{Status: 404, Message: "Category not found"}}

func GetCategories(ctx *gin.Context) {
	categories := &[]models.Category{}
	conditions, err := getConditions(ctx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, customerrors.InternalServerError("No conditions found in context"))
		return
	}

	page, pageSize := ctx.GetInt("page"), ctx.GetInt("pageSize")
	if page == 0 || pageSize == 0 {
		log.Println("Page or pageSize not set, using default values")
		page, pageSize = 1, 10
	}

	var totalCount int64
	log.Println("Getting categories for conditions: ", conditions)
	database.DB.Where(conditions).Limit(pageSize).Offset((page - 1) * pageSize).Find(categories)
	database.DB.Model(categories).Count(&totalCount)

	var r CategoryResponse
	categoryResponses, err := serialization.BindArray(*categories, &r)
	if err != nil {
		log.Println("Error binding categories: ", err)
		ctx.JSON(http.StatusInternalServerError, customerrors.InternalServerError("Error binding categories"))
		return
	}

	response := serialization.CreatePaginatedResponse(page, min(pageSize, len(categoryResponses)), int(totalCount), conditions, categoryResponses)

	ctx.JSON(200, response)
}

// func GetCategory(ctx *gin.Context) {
// 	conds["id"] = ctx.Param("id")
// 	category := models.Category{}

// 	log.Println("Getting category for conds: ", conds)
// 	database.DB.Where(conds).First(&category)
// 	if category.ID == 0 {
// 		log.Println("Category with id ", ctx.Param("id"), " not found")
// 		ctx.JSON(404, CategoryNotFoundResponse)
// 		return
// 	}
// 	ctx.JSON(200, category)
// }

// func CreateCategory(ctx *gin.Context) {
// 	category := models.Category{}
// 	ctx.BindJSON(&category)

// 	log.Println("Creating category: ", category)
// 	database.DB.Create(&category)

// 	ctx.JSON(201, category)
// }

// func UpdateCategory(ctx *gin.Context) {
// 	conds["id"] = ctx.Param("id")

// 	log.Println("Updating category for conds: ", conds)

// 	category := models.Category{}
// 	database.DB.Where(conds).First(&category)
// 	if category.ID == 0 {
// 		log.Println("Category with id ", ctx.Param("id"), " not found")
// 		ctx.JSON(404, CategoryNotFoundResponse)
// 		return
// 	}

// 	updateBodyJson := UpdateCategoryModel{}
// 	err := ctx.BindJSON(&updateBodyJson)
// 	if err != nil {
// 		log.Println("Error parsing request body: ", err)
// 		ctx.JSON(http.StatusUnprocessableEntity, UnprocessableEntityResponse)
// 		return
// 	}

// 	updateBody := createUpdateBody(&updateBodyJson)

// 	if len(updateBody) == 0 {
// 		log.Println("No fields to update for category with id ", conds["id"])
// 		ctx.JSON(200, category)
// 		return
// 	}

// 	log.Println("Updating category with id ", conds["id"], " with body: ", updateBody)
// 	database.DB.Model(&category).Updates(updateBody)

// 	ctx.JSON(200, category)
// }

// func DeleteCategory(ctx *gin.Context) {
// 	conds["id"] = ctx.Param("id")
// 	category := models.Category{}

// 	log.Println("Deleting category for conds: ", conds)
// 	database.DB.Where(conds).First(&category)

// 	if category.ID == 0 {
// 		log.Println("Category with id ", conds["id"], " not found")
// 		ctx.JSON(404, CategoryNotFoundResponse)
// 		return
// 	}

// 	log.Println("Deleting category with id ", conds["id"])
// 	database.DB.Delete(&category)

// 	ctx.JSON(204, http.NoBody)
// }

// func createUpdateBody(updateBodyJson *UpdateCategoryModel) gin.H {
// 	updateBody := make(gin.H)

// 	if updateBodyJson.Name != "" {
// 		updateBody["name"] = updateBodyJson.Name
// 	}
// 	if updateBodyJson.Description != "" {
// 		updateBody["description"] = updateBodyJson.Description
// 	}
// 	if updateBodyJson.Color != "" {
// 		updateBody["color"] = updateBodyJson.Color
// 	}
// 	if updateBodyJson.Budget != 0 {
// 		updateBody["budget"] = updateBodyJson.Budget
// 	}

// 	return updateBody
// }

func getConditions(ctx *gin.Context) (serialization.QueryConditions, error) {
	conditions, exists := ctx.Get("conditions")
	conds, ok := conditions.(serialization.QueryConditions)
	if !exists || !ok {
		return nil, fmt.Errorf("No conditions found in context, these must be set in the middleware")
	}
	return conds, nil
}
