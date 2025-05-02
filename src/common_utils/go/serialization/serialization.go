package serialization

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Marshaler interface {
	Marshal() JSONResponse
}

type ModelBinder interface {
	BindModel(model interface{}) error
}

type Serializer interface {
	Marshaler
	ModelBinder
}

// BindArray binds a slice of models to a slice of response types by using the BindModel method of the Serializer interface.
// It iterates over each model, binds it to the response type, and populates the resulting slice of response types.
// If binding fails for any model, an error is returned along with the partially populated response slice.
//
// Type Parameters:
//   - ResponseType: a type that implements the Serializer interface
//   - ModelType: any type representing the model
//
// Parameters:
//   - models: a slice of models to be bound
//
// Returns:
//   - []ResponseType: a slice of bound response types
//   - error: an error if any model binding fails
func BindArray[ResponseType Serializer, ModelType any](models []ModelType) ([]ResponseType, error) {
	response := make([]ResponseType, len(models))
	elementType, err := getElementType[ResponseType]()
	if err != nil {
		return response, fmt.Errorf("failed to get element type: %v", err)
	}

	for i, model := range models {
		boundItem, err := bindReflect[ResponseType](model, elementType)
		if err != nil {
			return response, fmt.Errorf("failed to bind model: %v", err)
		}
		response[i] = boundItem
	}

	return response, nil
}

// CreatePaginatedResponse generates a PaginatedJSONResponse from the given parameters.
// It calculates the total number of pages based on the total items and page size.
// If the page size is zero, the total pages will be set to zero.
// The response includes the current page, total pages, page size, total items,
// filters applied, and the list of items.
//
// Parameters:
//   - page: the current page number
//   - size: the number of items per page
//   - total: the total number of items
//   - filters: the query conditions used to filter the data
//   - items: the list of items to be included in the response
//
// Returns:
//   - PaginatedJSONResponse: the structured response containing pagination details and items
//   - error: an error if any occurs during response creation
func CreatePaginatedResponse(page, size, total int, filters QueryConditions, items []DataItem) (PaginatedJSONResponse) {
	totalPages := 0
	if size > 0 {
		totalPages = (total + size - 1) / size
		if total == 0 {
			totalPages = 0
		}
	}
	return PaginatedJSONResponse{
		Status: "success",
		Data: PaginatedResponse{
			Page:       page,
			Total:      totalPages,
			Size:       size,
			TotalItems: total,
			Filters:    filters,
			Items:      items,
		},
	}
}

// FilterSerializerFields filters the fields of a serialized struct based on a given list of field names.
// It first marshals the provided serializer into JSON format, then unmarshals it into a map, allowing
// for selective field extraction. If the list of fields is empty, all fields from the serializer are returned.
//
// Parameters:
//   - serializer: the struct implementing the Marshaler interface to be filtered
//   - fields: a slice of strings representing the field names to be included in the result
//
// Returns:
//   - DataItem: a map containing only the specified fields from the serializer
//   - error: an error if any occurs during marshaling or unmarshaling
func FilterSerializerFields(serializer Marshaler, fields []string) (DataItem, error) {
	marshalled, err := json.Marshal(serializer)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal serializer: %v", err)
	}
	var unmarshalled DataItem
	if err := json.Unmarshal(marshalled, &unmarshalled); err != nil {
		return nil, fmt.Errorf("failed to unmarshal serializer: %v", err)
	}

	if len(fields) == 0 {
		return unmarshalled, nil
	}

	result := make(DataItem)
	for _, field := range fields {
		value, ok := unmarshalled[field]
		if !ok {
			continue
		}
		result[field] = value
	}

	return result, nil
}

func bindReflect[ResponseType Serializer](model any, elementType reflect.Type) (response ResponseType, err error) {
	newValue := reflect.New(elementType)

	itemAsInterface := newValue.Interface()

	itemAsSerializer, ok := itemAsInterface.(Serializer)
	if !ok {
		err = fmt.Errorf("created instance of type %s does not implement Serializer interface", newValue.Type())
		return
	}

	if bindErr := itemAsSerializer.BindModel(model); bindErr != nil {
		err = fmt.Errorf("failed to bind model: %v", bindErr)
		return
	}

	itemAsResponseType, ok := itemAsInterface.(ResponseType)
	if !ok {
		err = fmt.Errorf("created instance of type %s does not implement ResponseType interface", newValue.Type())
		return
	}

	return itemAsResponseType, nil
}

func getElementType[ResponseType Serializer]() (reflect.Type, error) {
	var zero ResponseType
	responseReflectType := reflect.TypeOf(zero)

	if responseReflectType.Kind() != reflect.Ptr || responseReflectType.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("ResponseType must be a pointer to a struct, got %s", responseReflectType.Kind())
	}

	return responseReflectType.Elem(), nil
}
