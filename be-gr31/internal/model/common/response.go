package common

// Response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

// PaginatedData
type PaginatedData struct {
	Items    interface{} `json:"items"`
	Page     int         `json:"page"`
	Limit    int         `json:"limit"`
	Total    int         `json:"total"`
	HasMore  bool        `json:"hasMore"`
	NextPage *string     `json:"nextPage,omitempty"`
}

// Success
func OK(data interface{}, msg string) Response {
	return Response{Success: true, Message: msg, Data: data}
}

// Fail
func Fail(err string) Response {
	return Response{Success: false, Error: err}
}

// Paginated
func Paginated(data interface{}, page, limit, total int, hasMore bool, msg string) PaginatedResponse {
	return PaginatedResponse{
		Success: true,
		Message: msg,
		Data: PaginatedData{
			Items:   data,
			Page:    page,
			Limit:   limit,
			Total:   total,
			HasMore: hasMore,
		},
	}
}
