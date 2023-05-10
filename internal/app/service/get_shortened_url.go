package service

func (s *ShortenService) GetShortenedURLByID(id string) (string, bool) {
	return s.urlRepository.FindByShortenedString(id)
}
