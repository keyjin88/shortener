package service

import (
	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/keyjin88/shortener/internal/app/service/mocks"
	"github.com/keyjin88/shortener/internal/app/storage"
	"github.com/keyjin88/shortener/internal/app/storage/inmem"
	"github.com/stretchr/testify/assert"
	"go/constant"
	"testing"
)

func TestShortenService_GetShortenedURLByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    storage.ShortenedURL
		wantErr any
	}{
		{
			name: "success",
			args: args{
				id: "SHORTSTRING",
			},
			want:    storage.ShortenedURL{OriginalURL: "https://example.com/1"},
			wantErr: nil,
		},
		{
			name: "not success",
			args: args{
				id: "SHORTSTRING",
			},
			want:    storage.ShortenedURL{},
			wantErr: errors.New("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockURLRepository := mocks.NewMockURLRepository(ctrl)
			mockURLRepository.EXPECT().FindByShortenedURL(tt.args.id).Return(tt.want, tt.wantErr)
			s := &ShortenService{
				urlRepository: mockURLRepository,
			}
			got, err := s.GetShortenedURLByID(tt.args.id)
			if err != nil {
				assert.NotNil(t, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestShortenService_ShortenString(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		serviceArgs         string
		repositoryCallCount int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "success",
			args: args{
				serviceArgs:         "https://example.com/1",
				repositoryCallCount: 1,
			},
			want:    "SHORTSTRING",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockURLRepository := mocks.NewMockURLRepository(ctrl)
			mockURLRepository.EXPECT().FindByShortenedURL(
				gomock.Any()).Return(storage.ShortenedURL{},
				errors.New("not found url"))
			mockURLRepository.EXPECT().Save(gomock.Any()).Times(tt.args.repositoryCallCount)
			s := &ShortenService{
				urlRepository: mockURLRepository,
				config:        &Config{},
			}

			got, err := s.ShortenURL(tt.args.serviceArgs, "any string")
			assert.Equal(t, tt.wantErr, err)
			assert.IsType(t, constant.String.String(), got)
		})
	}
}

func TestShortenService_ShortenURLBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name                     string
		serviceArgs              storage.ShortenURLBatchRequest
		serviceError             string
		repositoryReturn         storage.ShortenedURL
		repositoryError          error
		repositoryCallCount      int
		repositorySaveBatchError error
	}{
		{
			name: "success",
			serviceArgs: storage.ShortenURLBatchRequest{
				{CorrelationID: "AAA", OriginalURL: "https://example.com/1"},
				{CorrelationID: "BBB", OriginalURL: "https://yandex.com/1"},
			},
			serviceError:             "",
			repositoryReturn:         storage.ShortenedURL{},
			repositoryError:          errors.New("repository"),
			repositoryCallCount:      2,
			repositorySaveBatchError: nil,
		},
		{
			name: "error",
			serviceArgs: storage.ShortenURLBatchRequest{
				{CorrelationID: "AAA", OriginalURL: "https://example.com/1"},
				{CorrelationID: "BBB", OriginalURL: "https://yandex.com/1"},
			},
			serviceError:             "failed to save batch: while saving batch",
			repositoryReturn:         storage.ShortenedURL{},
			repositoryError:          errors.New("repository"),
			repositoryCallCount:      2,
			repositorySaveBatchError: errors.New("while saving batch"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockURLRepository := mocks.NewMockURLRepository(ctrl)
			mockURLRepository.EXPECT().
				FindByShortenedURL(gomock.Any()).
				Return(tt.repositoryReturn, tt.repositoryError).
				Times(tt.repositoryCallCount)
			mockURLRepository.EXPECT().SaveBatch(gomock.Any()).
				Return(tt.repositorySaveBatchError).
				Times(1)

			s := &ShortenService{
				urlRepository: mockURLRepository,
				config:        &Config{},
			}
			_, err := s.ShortenURLBatch(tt.serviceArgs, "any string")
			if err != nil {
				assert.Equal(t, tt.serviceError, err.Error())
			}
		})
	}
}

func TestShortenService_GetShortenedURLByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type repositoryReturn struct {
		result []storage.UsersURLResponse
		error  error
	}

	type serviceReturn struct {
		result []storage.UsersURLResponse
		error  string
	}

	tests := []struct {
		name             string
		repositoryReturn repositoryReturn
		serviceReturn    serviceReturn
	}{
		{
			name: "success",
			repositoryReturn: repositoryReturn{
				result: []storage.UsersURLResponse{
					{
						OriginalURL: "https://yandex.com/1",
						ShortURL:    "111111",
					},
					{
						OriginalURL: "https://mail.com/1",
						ShortURL:    "222222",
					},
				},
				error: nil,
			},
			serviceReturn: serviceReturn{
				result: []storage.UsersURLResponse{
					{
						OriginalURL: "https://yandex.com/1",
						ShortURL:    "/111111",
					},
					{
						OriginalURL: "https://mail.com/1",
						ShortURL:    "/222222",
					},
				},
				error: "",
			},
		},
		{
			name: "error",
			repositoryReturn: repositoryReturn{
				result: nil,
				error:  errors.New("test error"),
			},
			serviceReturn: serviceReturn{
				result: nil,
				error:  "failed to find URL from repository: test error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockURLRepository := mocks.NewMockURLRepository(ctrl)
			mockURLRepository.EXPECT().
				FindAllByUserID(gomock.Any()).
				Return(tt.repositoryReturn.result, tt.repositoryReturn.error)
			s := &ShortenService{
				urlRepository: mockURLRepository,
				config:        &Config{},
			}
			got, err := s.GetShortenedURLByUserID("userId")
			if err != nil {
				assert.Equal(t, tt.serviceReturn.error, err.Error())
			}
			assert.Equal(t, tt.serviceReturn.result, got)
		})
	}
}

func TestShortenService_DeleteURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name   string
		urls   []string
		userID string
		error  error
	}{
		{
			name:   "success",
			urls:   []string{"https://example.com/1", "https://example.com/2", "https://yandex.com/1"},
			userID: "userId",
			error:  nil,
		},
		{
			name:   "error",
			urls:   []string{"https://example.com/1", "https://example.com/2", "https://yandex.com/1"},
			userID: "userId",
			error:  errors.New("failed to delete URL: test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockURLRepository := mocks.NewMockURLRepository(ctrl)
			mockURLRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(tt.error)

			s := &ShortenService{
				urlRepository: mockURLRepository,
				config:        &Config{},
			}
			err := s.DeleteURLs(&tt.urls, tt.userID)
			if err != nil {
				actualErr := errors.Cause(err).Error()
				assert.Equal(t, actualErr, tt.error.Error())
			}
		})
	}
}

func BenchmarkGenerateShortURL(b *testing.B) {
	b.StopTimer() // останавливаем таймер
	var repo = inmem.NewURLRepositoryInMem()
	s := &ShortenService{
		urlRepository: repo,
	}
	b.StartTimer() // возобновляем таймер
	for i := 0; i < b.N; i++ {
		_, err := s.generateShortenURL()
		if err != nil {
			return
		}
	}
}
