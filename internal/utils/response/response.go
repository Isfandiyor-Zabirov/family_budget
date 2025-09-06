package response

import "math"

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

func CalculateTotalPages(totalRows int64, pageLimit int) uint64 {
	return uint64(math.Ceil(float64(totalRows) / float64(pageLimit)))
}

func SetResponseData(data interface{}, message string, success bool) ResponseModel {
	return ResponseModel{
		Success: success,
		Message: message,
		Data:    data,
	}
}

func SetPagination(totalPages uint64, totalRows int64, currentPage int) Pagination {
	return Pagination{
		TotalRows:   uint64(totalRows),
		TotalPages:  totalPages,
		CurrentPage: uint64(currentPage),
	}
}
