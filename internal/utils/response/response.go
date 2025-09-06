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

func calculateTotalPages(totalRows int64, pageLimit int) uint64 {
	if pageLimit == 0 {
		pageLimit = 1
	}
	return uint64(math.Ceil(float64(totalRows) / float64(pageLimit)))
}

func SetResponseData(resp interface{}, data interface{}, message string, success bool, pageLimit int, totalRows int64, currentPage int) {
	type dataStruct struct {
		List       interface{} `json:"list"`
		Pagination Pagination  `json:"pagination"`
	}

	newData := dataStruct{
		List:       data,
		Pagination: setPagination(calculateTotalPages(totalRows, pageLimit), totalRows, currentPage),
	}
	newRespModel := ResponseModel{
		Success: success,
		Message: message,
		Data:    newData,
	}

	// to avoid nil pointer dereference
	resp = newRespModel
}

func setPagination(totalPages uint64, totalRows int64, currentPage int) Pagination {
	return Pagination{
		TotalRows:   uint64(totalRows),
		TotalPages:  totalPages,
		CurrentPage: uint64(currentPage),
	}
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}
