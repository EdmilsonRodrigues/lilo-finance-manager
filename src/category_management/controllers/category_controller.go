package controllers

import (
	"log"
	"net/http"

	customerrors "github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/custom_errors"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/database"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/middlewares"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/models"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/common_utils/go/serialization/http_serialization"
	"github.com/gin-gonic/gin"
)

var CategoryNotFoundResponse = httpserialization.ErrorResponse{
	Details: httpserialization.ErrorDetails{
		Status:  404,
		Message: "Category not found",
	},
}

// GetCategories gets all categories that match the conditions set in the context.
//
// The method queries the database using the conditions and values set in the context.
// It uses the page and pageSize values from the context to limit the number of categories returned.
// It also counts the total number of categories that match the conditions.
// The method uses the BindArray function to convert the categories to an array of CategoryResponse.
// It then creates a PaginatedResponse using the page, pageSize, totalCount, filters and array of CategoryResponse.
// The method returns the PaginatedResponse as a JSON response.
func GetCategories(ctx *gin.Context) {
	categories := &[]models.Category{}
	conditions, conditionsValues := middlewares.GetConditions(ctx)

	page, pageSize := ctx.GetInt("page"), ctx.GetInt("pageSize")
	if page == 0 {
		log.Println("Page or pageSize not set, using default values")
		page = 1
	}
	if pageSize == 0 {
		log.Println("Page or pageSize not set, using default values")
		pageSize = 10
	}

	var totalCount int64
	log.Println("Getting categories for conditions: ", conditions, conditionsValues)
	database.DB.Where(conditions, conditionsValues...).Limit(pageSize).Offset((page - 1) * pageSize).Find(categories)
	database.DB.Where(conditions, conditionsValues...).Model(categories).Count(&totalCount)

	categoryResponses, err := httpserialization.BindArray[*CategoryResponse](*categories)
	if err != nil {
		log.Println("Error binding categories: ", err)
		ctx.JSON(http.StatusInternalServerError, customerrors.InternalServerError("Error binding categories"))
		return
	}

	fields, filters, serializerArray := getReturnFields(ctx), getFilters(ctx), convertToSerializer(categoryResponses)

	response, err := httpserialization.CreatePaginatedResponse(page, len(categoryResponses), int(totalCount), filters, serializerArray, fields)
	if err != nil {
		log.Println("Error creating response: ", err)
		ctx.JSON(http.StatusInternalServerError, customerrors.InternalServerError("Error creating response"))
		return
	}

	ctx.JSON(200, response)
}

// GetCategory gets a category by its ID.
//
// The method first gets the conditions and values from the context.
// It then queries the database using the conditions and values.
// If the category is not found, it returns a 404 response.
// If the category is found, it binds the category to a CategoryResponse and marshalls it to a JSON response.
// The method returns the JSON response.
func GetCategory(ctx *gin.Context) {
	conditions, conditionsValues := middlewares.GetConditions(ctx)
	category := models.Category{}

	log.Println("Getting category for conds: ", conditions, conditionsValues)
	database.DB.Where(conditions, conditionsValues...).First(&category)
	if category.ID == 0 {
		log.Println("Category with id ", ctx.Param("id"), " not found")
		ctx.JSON(404, CategoryNotFoundResponse)
		return
	}

	var categoryResponse CategoryResponse
	err := categoryResponse.BindModel(category)
	if err != nil {
		log.Println("Error binding category: ", err)
		ctx.JSON(http.StatusInternalServerError, customerrors.InternalServerError("Error binding category"))
		return
	}

	response, err := categoryResponse.Marshal(getReturnFields(ctx))
	if err != nil {
		log.Println("Error marshalling category: ", err)
		ctx.JSON(http.StatusInternalServerError, customerrors.InternalServerError("Error marshalling category"))
		return
	}

	ctx.JSON(200, response)
}

// CreateCategory creates a new category.
//
// The method first binds the JSON body to a models.Category.
// It then checks if the account id in the path matches the account id in the body.
// If not, it returns a 400 response.
// Otherwise, it creates the category in the database.
// If the creation is successful, it binds the category to a CategoryResponse and marshalls it to a JSON response.
// The method returns the JSON response with a 201 status code.
func CreateCategory(ctx *gin.Context) {
	category := models.Category{}
	ctx.BindJSON(&category)

	accountId := ctx.Param("accountId")
	if accountId != category.AccountID {
		log.Println("Account id in path does not match account id in body")
		ctx.JSON(http.StatusBadRequest, customerrors.BadRequestResponse)
		return
	}

	log.Println("Creating category: ", category)
	database.DB.Create(&category)

	var categoryResponse CategoryResponse
	err := categoryResponse.BindModel(category)
	if err != nil {
		log.Println("Error binding category: ", err)
		ctx.JSON(http.StatusInternalServerError, customerrors.InternalServerError("Error binding category"))
		return
	}

	response, err := categoryResponse.Marshal(getReturnFields(ctx))
	if err != nil {
		log.Println("Error marshalling category: ", err)
		ctx.JSON(http.StatusInternalServerError, customerrors.InternalServerError("Error marshalling category"))
		return
	}

	ctx.JSON(201, response)
}

// UpdateCategory updates a category.
//
// The method first gets the conditions and values from the context.
// It then queries the database using the conditions and values.
// If the category is not found, it returns a 404 response.
// Otherwise, it binds the JSON body to an UpdateCategoryModel and updates the category in the database.
// If the update is successful, it binds the category to a CategoryResponse and marshalls it to a JSON response.
// The method returns the JSON response with a 200 status code.
func UpdateCategory(ctx *gin.Context) {
	conditions, conditionsValues := middlewares.GetConditions(ctx)
	category := models.Category{}

	log.Println("Updating category for conds: ", conditions, conditionsValues)
	database.DB.Where(conditions, conditionsValues...).First(&category)

	if category.ID == 0 {
		log.Println("Category with id ", ctx.Param("id"), " not found")
		ctx.JSON(404, CategoryNotFoundResponse)
		return
	}

	updateBody := UpdateCategoryModel{}
	err := ctx.BindJSON(&updateBody)
	if err != nil {
		log.Println("Error parsing request body: ", err)
		ctx.JSON(http.StatusUnprocessableEntity, customerrors.UnprocessableEntityResponse)
		return
	}

	log.Println("Updating category with id ", ctx.Param("id"), " with body: ", updateBody)
	database.DB.Model(&category).Updates(updateBody)

	var categoryResponse CategoryResponse
	err = categoryResponse.BindModel(category)
	if err != nil {
		log.Println("Error binding category: ", err)
		ctx.JSON(http.StatusInternalServerError, customerrors.InternalServerError("Error binding category"))
		return
	}

	response, err := categoryResponse.Marshal(getReturnFields(ctx))
	if err != nil {
		log.Println("Error marshalling category: ", err)
		ctx.JSON(http.StatusInternalServerError, customerrors.InternalServerError("Error marshalling category"))
		return
	}
	ctx.JSON(200, response)
}

func DeleteCategory(ctx *gin.Context) {
	conditions, conditionsValues := middlewares.GetConditions(ctx)
	category := models.Category{}

	log.Println("Updating category for conds: ", conditions, conditionsValues)
	database.DB.Where(conditions, conditionsValues...).First(&category)

	if category.ID == 0 {
		log.Println("Category with id ", ctx.Param("id"), " not found")
		ctx.JSON(404, CategoryNotFoundResponse)
		return
	}

	log.Println("Deleting category with id ", ctx.Param("id"))
	database.DB.Delete(&category)

	ctx.JSON(204, http.NoBody)
}

func getFilters(ctx *gin.Context) map[string]string {
	f, _ := ctx.Get("filters")
	filters, ok := f.(map[string]string)
	if !ok {
		filters = make(map[string]string)
	}
	return filters
}

func getReturnFields(ctx *gin.Context) []string {
	f, _ := ctx.Get("returnFields")
	fields, ok := f.([]string)
	if !ok {
		log.Println("Fields not set in context")
		fields = []string{}
	}
	return fields

}

func convertToSerializer(categoryResponses []*CategoryResponse) []httpserialization.Serializer {
	serializerArray := make([]httpserialization.Serializer, len(categoryResponses))
	for i, categoryResponse := range categoryResponses {
		serializerArray[i] = categoryResponse
	}
	return serializerArray
}
