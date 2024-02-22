package service

type News interface {
	ParsePage(page string) error
	ParseQuery(query string) error
}
