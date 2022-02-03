package responses

type PaginatedResponse struct {
	Start    int64       `json:"start"`
	PageSize int64       `json:"pageSize"`
	Total    int64       `json:"total"`
	Data     interface{} `json:"data"`
}
