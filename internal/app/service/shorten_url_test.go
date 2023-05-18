package service

import (
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/shortener/internal/app/service/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShortenService_ShortenString(t *testing.T) {
	ctrl := gomock.NewController(t)
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
			mockURLRepository.EXPECT().FindByShortenedString(gomock.Any()).Return("any string", false)
			mockURLRepository.EXPECT().Create(gomock.Any(), tt.args.url).Times(1)
			s := &ShortenService{
				urlRepository: mockURLRepository,
			}

			got, err := s.ShortenURL(tt.args.url)
			assert.Equal(t, tt.wantErr, err)
			assert.IsType(t, "String", got)
		})
	}
}
