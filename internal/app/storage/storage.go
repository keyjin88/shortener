package storage

type Storage struct {
	urlRepository *UrlRepository
}

func NewStorage() *Storage {
	return &Storage{}
}

// public repo for URLS
func (s *Storage) Urls() *UrlRepository {
	if s.urlRepository != nil {
		return s.urlRepository
	}
	a := &UrlRepository{
		storage: s,
	}
	s.urlRepository = a
	return a
}
