package service

type News interface {
	Parse(query string) error
}
