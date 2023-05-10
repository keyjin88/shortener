package service

//go:generate mockgen -destination=mocks/url_repository.go -package=mocks . URLRepository
type URLRepository interface {
	FindByShortenedString(id string) (string, bool)
	Create(uuidStr string, url string)
}

type ShortenService struct {
	urlRepository URLRepository
}

func NewShortenService(urlRepository URLRepository) *ShortenService {
	return &ShortenService{urlRepository: urlRepository}
}
