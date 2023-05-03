package storage

type URLRepository interface {
	FindByShortenedString(id string) (string, bool)
	Create(uuidStr string, url string)
}
