package dto

type WebResponse[T any] struct {
	Data  *T      `json:"data,omitempty"`
	Error *string `json:"error,omitempty"`
	Path  string  `json:"path"`
}
type WebPagination[T any] struct {
	Content      T   `json:"contents"`
	Page         int `json:"page"`
	Size         int `json:"size"`
	TotalPage    int `json:"total_page"`
	TotalElement int `json:"total_element"`
}
