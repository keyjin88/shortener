package service

//go:generate mockgen -destination=mocks/url_repository.go -package=mocks . URLRepository
type URLRepository interface {
	FindByShortenedString(id string) (string, bool)
	Create(uuidStr string, url string) error
}
