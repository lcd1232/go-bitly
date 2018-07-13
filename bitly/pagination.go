package bitly

type Paginate struct {
	Total int    `json:"total"`
	Size  int    `json:"size"`
	Prev  string `json:"prev"`
	Page  int    `json:"page"`
	Next  string `json:"next"`
}

// Paginator is a simple interface which need for pagination
//
type Paginator interface {
	Next() bool
	Prev() bool
	Get() error
}
