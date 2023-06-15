package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/shortener/internal/app/service/mocks"
	"github.com/keyjin88/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
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
		want    string
		wantErr error
	}{
		{
			name: "success",
			args: args{
				id: "SHORTSTRING",
			},
			want:    "https://example.com/1",
			wantErr: errors.New("test error"),
		},
		{
			name: "not success",
			args: args{
				id: "SHORTSTRING",
			},
			want:    "",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockURLRepository := mocks.NewMockURLRepository(ctrl)
			mockURLRepository.EXPECT().FindByShortenedURL(tt.args.id).Return(tt.want, tt.wantErr)
			s := &ShortenService{
				urlRepository: mockURLRepository,
			}
			got, got1 := s.GetShortenedURLByID(tt.args.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, got1)
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
			mockURLRepository.EXPECT().FindByShortenedURL(gomock.Any()).Return("any string", errors.New("not found url"))
			mockURLRepository.EXPECT().Save(gomock.Any()).Times(tt.args.repositoryCallCount)
			s := &ShortenService{
				urlRepository: mockURLRepository,
				config:        &Config{},
			}

			got, err := s.ShortenURL(tt.args.serviceArgs, "any string")
			assert.Equal(t, tt.wantErr, err)
			assert.IsType(t, "String", got)
		})
	}
}

func TestShortenService_ShortenURLBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name                     string
		serviceArgs              storage.ShortenURLBatchRequest
		serviceError             error
		repositoryReturn         string
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
			serviceError:             nil,
			repositoryReturn:         "",
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
			serviceError:             errors.New("while saving batch"),
			repositoryReturn:         "",
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
			assert.Equal(t, tt.serviceError, err)
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
		error  error
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
				error: nil,
			},
		},
		{
			name: "success",
			repositoryReturn: repositoryReturn{
				result: nil,
				error:  errors.New("test error"),
			},
			serviceReturn: serviceReturn{
				result: nil,
				error:  errors.New("test error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockURLRepository := mocks.NewMockURLRepository(ctrl)
			mockURLRepository.EXPECT().FindAllByUserId(gomock.Any()).Return(tt.repositoryReturn.result, tt.repositoryReturn.error)
			s := &ShortenService{
				urlRepository: mockURLRepository,
				config:        &Config{},
			}
			got, err := s.GetShortenedURLByUserID("userId")
			assert.Equal(t, tt.serviceReturn.error, err)
			assert.Equal(t, tt.serviceReturn.result, got)
		})
	}
}
