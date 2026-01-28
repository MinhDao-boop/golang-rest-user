package dto

type ListResponse struct {
	Data     interface{} `json:"data"`
	Page     *int        `json:"page"`
	PageSize *int        `json:"pageSize"`
	Total    int64       `json:"total"`
}
