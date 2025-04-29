package serialization

import "fmt"

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
