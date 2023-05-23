package storage

type Config struct {
	PathToStorageFile string //путь до фпйла для резервного хранения
}

type URLJSON struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
