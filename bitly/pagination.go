package bitly

type Paginate struct {
	Total int    `json:"total"`
	Size  int    `json:"size"`
	Prev  string `json:"prev"`
	Page  int    `json:"page"`
	Next  int    `json:"next"`
}

type Paginator interface {
	Next() error
	Prev() error
}
