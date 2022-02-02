package responses

type ErrorResponse struct {
	ErrorId string            `json:"errorId,omitempty"`
	Message string            `json:"message"`
	Status  int               `json:"status"`
	Title   string            `json:"title,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}
