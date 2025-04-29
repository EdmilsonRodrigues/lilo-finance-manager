package serialization

type QueryConditions map[string]interface{}

type ErrorResponse struct {
	Details ErrorDetails `json:"details"`
}

type ErrorDetails struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type JSONResponse struct {
	Status string    `json:"status"`
	Data   Marshaler `json:"data"`
}

type PaginatedJSONResponse struct {
	Status string            `json:"status"`
	Data   PaginatedResponse `json:"data"`
}

type PaginatedResponse struct {
	Page       int                 `json:"page"`
	Total      int                 `json:"total_pages"`
	Size       int                 `json:"page_size"`
	TotalItems int                 `json:"total_items"`
	Filters    QueryConditions     `json:"filters"`
	Items      []Marshaler `json:"items"`
}
