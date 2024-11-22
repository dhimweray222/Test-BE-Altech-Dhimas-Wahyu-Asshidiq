package web

type WebResponse struct {
	Code    int         `json:"code"`
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type WebResponsePagination struct {
	Code      int         `json:"code"`
	Status    bool        `json:"status"`
	Page      int         `json:"page"`
	TotalPage int         `json:"total_page"`
	TotalData int         `json:"total_data"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}
