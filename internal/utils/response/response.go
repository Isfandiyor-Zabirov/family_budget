package response

type ResponseModel struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Pagination struct {
	TotalRows   uint64 `json:"total_rows"`
	TotalPages  uint64 `json:"total_pages"`
	CurrentPage uint64 `json:"current_page"`
}
