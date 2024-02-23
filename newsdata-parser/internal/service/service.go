package service

type News interface {
	ParsePage(page string) (string, error)
	ParseQuery(query string) error
}

type Page interface {
	Get() (string, error)
	Set(page string) error
}
