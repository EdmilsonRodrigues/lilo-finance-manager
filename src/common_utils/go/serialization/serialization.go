package serialization

import (
	"encoding/json"
	"fmt"
)

type Marshaler interface {
	Marshal() JSONResponse
}

type ModelBinder interface {
	BindModel(model interface{}) error
}

func BindArray[ResponseType ModelBinder, ModelType any](models []ModelType, responseType ResponseType) ([]ResponseType, error) {
	response := make([]ResponseType, len(models))
	for i, model := range models {
		var item ResponseType
		if err := item.BindModel(model); err != nil {
			return response, fmt.Errorf("failed to bind model: %v", err)
		}
		response[i] = item
	}
	return response, nil
}

func CreatePaginatedResponse[ResponseArray []Marshaler](page, size, total int, filters QueryConditions, items ResponseArray) PaginatedJSONResponse {
	return PaginatedJSONResponse{
		Status: "success",
		Data: PaginatedResponse{
			Page:       page,
			Total:      total / size,
			Size:       size,
			TotalItems: total,
			Filters:    filters,
			Items:      items,
		},
	}
}

// FilterSerializerFields filters the fields of a struct according to the given list of fields.
// It ignores any fields that are not present in the given list.
// It returns a map of the filtered fields.
func FilterSerializerFields(serializer Marshaler, fields []string) (map[string]any, error) {
	marshalled, err := json.Marshal(serializer)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal serializer: %v", err)
	}
	var unmarshalled map[string]any
	if err := json.Unmarshal(marshalled, &unmarshalled); err != nil {
		return nil, fmt.Errorf("failed to unmarshal serializer: %v", err)
	}

	result := make(map[string]any)
	for _, field := range fields {
		value, ok := unmarshalled[field]
		if !ok {
			continue
		}
		result[field] = value
	}

	return result, nil
}
