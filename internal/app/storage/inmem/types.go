package inmem

type Config struct {
	PathToStorageFile string //путь до фпйла для резервного хранения
}

type URLRepositoryInMem struct {
	config       Config
	inMemStorage map[string]string
}

func NewURLRepositoryInMem(pathToStorageFile string) *URLRepositoryInMem {
	return &URLRepositoryInMem{
		config: Config{
			PathToStorageFile: pathToStorageFile,
		},
		inMemStorage: make(map[string]string),
	}
}
