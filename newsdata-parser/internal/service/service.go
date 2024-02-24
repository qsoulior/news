package service

type News interface {
	Parse(query string, page string) (string, error)
}

type Page interface {
	Get() (string, error)
	Set(page string) error
}
