package service

import (
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/shortener/internal/app/service/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShortenService_GetShortenedURLByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		id string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "success",
			args: args{
				id: "SHORTSTRING",
			},
			want:  "https://example.com/1",
			want1: true,
		},
		{
			name: "not success",
			args: args{
				id: "SHORTSTRING",
			},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockURLRepository := mocks.NewMockURLRepository(ctrl)
			mockURLRepository.EXPECT().FindByShortenedString(tt.args.id).Return(tt.want, tt.want1)
			s := &ShortenService{
				urlRepository: mockURLRepository,
			}
			got, got1 := s.GetShortenedURLByID(tt.args.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
