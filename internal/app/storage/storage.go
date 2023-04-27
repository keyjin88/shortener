package storage

type Storage struct {
	urlRepository *URLRepository
}

func NewStorage() *Storage {
	return &Storage{}
}

// public repo for URLS
func (s *Storage) Urls() *URLRepository {
	if s.urlRepository != nil {
		return s.urlRepository
	}
	a := &URLRepository{
		storage: s,
	}
	s.urlRepository = a
	return a
}
