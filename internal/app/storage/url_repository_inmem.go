package storage

type URLRepositoryInMem struct {
}

func NewURLRepositoryInMem() *URLRepositoryInMem {
	return &URLRepositoryInMem{}
}

var (
	inMemStorage map[string]string
)

func init() {
	inMemStorage = make(map[string]string)
}

func (ur *URLRepositoryInMem) Create(uuidStr string, url string) {
	inMemStorage[uuidStr] = url
}

func (ur *URLRepositoryInMem) FindByShortenedString(id string) (string, bool) {
	url, ok := inMemStorage[id]
	return url, ok
}
