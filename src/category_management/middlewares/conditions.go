package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	customerrors "github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/custom_errors"
	"github.com/gin-gonic/gin"
)

// AddParamToConditionsMiddleware adds a param to the query conditions in the context.
// It takes two params, paramName and columnName, and uses them to set a key-value pair in the query conditions.
// If the key is already present, it will be overwritten.
// If the "conditions" key is not present in the context, it will be created.
//
// Parameters:
//   - paramName: the name of the param to add to the query conditions
//   - columnName: the name of the key to add to the query conditions
//
// Returns:
//   - gin.HandlerFunc: the middleware function
//
// Example:
//
//	group.GET("/:id", AddParamToConditionsMiddleware("id", "id"), categoryController.GetCategoryById)
func AddParamToConditionsMiddleware(paramName, columnName string) gin.HandlerFunc {
	if paramName == ""  {
		log.Fatal("AddParamToConditionsMiddleware called with empty paramName or columnName")
	}

	return func(ctx *gin.Context) {
		paramValue := ctx.Param(paramName)
		if paramValue == "" {
			log.Printf("AddParamToConditionsMiddleware applied to route missing %s parameter: %s", paramName, ctx.Request.URL.Path)
			ctx.JSON(http.StatusUnprocessableEntity, customerrors.UnprocessableEntityResponse)
			ctx.Abort()
			return
		}

		queryConditions, queryValues := GetConditions(ctx)

		queryValues = append(queryValues, paramValue)
		condition := fmt.Sprintf("%s = ?", columnName)
		if queryConditions == "" {
			queryConditions = condition
		} else {
			queryConditions = strings.Join([]string{queryConditions, condition}, " AND ")
		}

		ctx.Set("conditions", queryConditions)
		ctx.Set("conditions_values", queryValues)
		ctx.Next()
	}
}

// ParsePaginationParamsMiddleware parses the page and page_size query parameters and sets them in the gin context.
// If the page parameter is not present, it sets it to 1.
// If the page_size parameter is not present, it sets it to 10.
// If the page_size parameter is greater than 100, it sets it to 100.
// If there is an error parsing the page or page_size parameters, it returns a 400 response.
//
// Returns:
//   - gin.HandlerFunc: the middleware function
//
// Example:
//
//	group.GET("/", ParsePaginationParamsMiddleware(), categoryController.GetCategories)
func ParsePaginationParamsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page, pageSize := ctx.Query("page"), ctx.Query("page_size")
		var (
			pageNumber, pageSizeNumber int
			err                        error
		)

		pageNumber, err = parsePageNumber(page)
		if err != nil {
			log.Printf("Error parsing page number for request %s, page: %s", ctx.Request.URL.Path, page)
			ctx.JSON(http.StatusUnprocessableEntity, customerrors.UnprocessableEntityResponse)
			ctx.Abort()
			return
		}

		pageSizeNumber, err = parsePageSize(pageSize)
		if err != nil {
			log.Printf("Error parsing page size for request %s, page_size: %s", ctx.Request.URL.Path, pageSize)
			ctx.JSON(http.StatusUnprocessableEntity, customerrors.UnprocessableEntityResponse)
			ctx.Abort()
			return
		}

		ctx.Set("page", pageNumber)
		ctx.Set("pageSize", pageSizeNumber)
		ctx.Next()
	}
}

// ParseFiltersMiddleware parses the filters query parameter and sets it in the gin context.
// It should be run after all the AddParamToConditionsMiddleware functions.
// If there is an error parsing the filters parameter, it returns a 422 response.
//
// Returns:
//   - gin.HandlerFunc: the middleware function
//
// Example:
//
//	group.GET("/", ParseFiltersMiddleware(), categoryController.GetCategories)
func ParseFiltersMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		queryConditions, queryValues := GetConditions(ctx)

		filters, err := parseQueryMap(ctx.Query("filters"))
		if err != nil {
			log.Printf("Error parsing filters for request %s, filters: %s", ctx.Request.URL.Path, ctx.Query("filters"))
			ctx.JSON(http.StatusUnprocessableEntity, customerrors.UnprocessableEntityResponse)
			ctx.Abort()
			return
		}
		ctx.Set("filters", filters)

		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("%s AND ", queryConditions))
		for key, value := range filters {
			builder.WriteString(fmt.Sprintf("%s = ? AND ", key))
			queryValues = append(queryValues, value)
		}
		queryConditions = builder.String()[:builder.Len()-5]
		ctx.Set("conditions", queryConditions)
		ctx.Set("conditions_values", queryValues)
		ctx.Next()
	}
}

// ParseReturnFieldsMiddleware parses the return_fields query parameter and sets it in the gin context.
// The return_fields parameter is a comma-separated string of field names.
// If the parameter is not present, it sets an empty array in the context.
//
// Returns:
//   - gin.HandlerFunc: the middleware function
//
// Example:
//
//	group.GET("/", ParseReturnFieldsMiddleware(), categoryController.GetCategories)
func ParseReturnFieldsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		returnFields := parseQueryArray(ctx.Query("return_fields"))
		ctx.Set("returnFields", returnFields)
		ctx.Next()
	}
}

// parseQueryMap takes a comma-separated string of key-value pairs in the format "key:value,key:value"
// and returns a map where each key is mapped to its corresponding value. If the query is empty,
// it returns an empty map. It returns an error if any pair does not adhere to the "key:value" format.
// It removes any leading or trailing whitespace from each item.
//
// Expected Query Format: "fields=key1:value1,key2:value2"
//
// Expected Response Format: map[string]string{"key1": "value1", "key2": "value2"}
func parseQueryMap(query string) (queryConditions map[string]string, err error) {
	queryConditions = make(map[string]string)
	query = strings.TrimSpace(query)
	if query == "" {
		return
	}

	pairs := strings.Split(query, ",")
	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			err = fmt.Errorf("invalid query pair: %s", pair)
			return
		}
		queryConditions[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return
}

// parseQueryArray takes a comma-separated string and splits it into a slice of strings.
// It returns an empty slice if the query is empty.
// It removes any leading or trailing whitespace from each item.
//
// Expected Query Format: "items=item1,item2,item3"
//
// Expected Result: []string{"item1", "item2", "item3"}
func parseQueryArray(query string) (queryConditions []string) {
	queryConditions = []string{}
	query = strings.TrimSpace(query)
	if query == "" {
		return
	}

	items := strings.Split(query, ",")
	for _, item := range items {
		if item = strings.TrimSpace(item); item == "" {
			continue
		}
		queryConditions = append(queryConditions, strings.TrimSpace(item))
	}
	return
}

// parsePageNumber converts a page string to an integer page number.
// If the page string is empty, it defaults to 1.
// Returns an error if the page string is not a valid integer.
//
// Parameters:
//   - page: a string representing the page number
//
// Returns:
//   - int: the page number
//   - error: an error if the page string is not a valid integer
//
// Example:
//
//	pageNumber, err := parsePageNumber("10")
//	if err != nil {
//	    log.Println("Error parsing page number for request ", ctx.Request.URL, " with params ", ctx.Params)
//	}
func parsePageNumber(page string) (int, error) {
	if page == "" {
		return 1, nil
	}
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return 0, fmt.Errorf("error parsing page number for request %s", page)
	}
	if pageNumber <= 0 {
		pageNumber = 1
	}
	return pageNumber, nil
}

// parsePageSize takes a string representing the page size and returns the page size as an integer or an error if the string does not represent a valid integer.
// If the string is empty, it returns 10 as the default page size.
// If the parsed page size is greater than 100, it returns 100 as the maximum allowed page size.
//
// Parameters:
//   - pageSize: a string representing the page size
//
// Returns:
//   - int: the page size
//   - error: an error if the string does not represent a valid integer or if the page size is greater than 100
//
// Example:
//
//	pageSize, err := parsePageSize("20")
func parsePageSize(pageSize string) (int, error) {
	if pageSize == "" {
		return 10, nil
	}
	pageSizeNumber, err := strconv.Atoi(pageSize)
	if err != nil {
		return 0, fmt.Errorf("error parsing page size for request %s", pageSize)
	}
	if pageSizeNumber <= 0 {
		pageSizeNumber = 10
	}
	if pageSizeNumber > 100 {
		pageSizeNumber = 100
	}
	return pageSizeNumber, nil
}

func GetConditions(ctx *gin.Context) (queryConditions string, queryValues []interface{}) {
	queryConditions = ctx.GetString("conditions")
	values, exists := ctx.Get("conditions_values")
	if !exists {
		queryValues = make([]interface{}, 0)
	} else {
		queryValues = values.([]interface{})
	}
	return
}
