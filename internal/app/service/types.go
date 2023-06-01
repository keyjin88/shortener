package service

//go:generate mockgen -destination=mocks/url_repository.go -package=mocks . URLRepository
type URLRepository interface {
	FindByShortenedURL(shortURL string) (string, error)
	Save(shortURL string, url string) error
}
