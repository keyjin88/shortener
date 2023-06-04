package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/shortener/internal/app/service/mocks"
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
		url string
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
				url: "https://example.com/1",
			},
			want:    "SHORTSTRING",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockURLRepository := mocks.NewMockURLRepository(ctrl)
			mockURLRepository.EXPECT().FindByShortenedURL(gomock.Any()).Return("any string", errors.New("not found url"))
			mockURLRepository.EXPECT().Save(gomock.Any(), tt.args.url).Times(1)
			s := &ShortenService{
				urlRepository: mockURLRepository,
				config:        &Config{PathToStorageFile: ""},
			}

			got, err := s.ShortenURL(tt.args.url)
			assert.Equal(t, tt.wantErr, err)
			assert.IsType(t, "String", got)
		})
	}
}
