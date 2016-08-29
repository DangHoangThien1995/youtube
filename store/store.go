package store

import (
	"youtube/api/apiyoutube"
	"youtube/domain"
)

type Store struct {
	Yt *apiyoutube.YT
}

func NewStore() *Store {
	return &Store{
		Yt: apiyoutube.NewYT(),
	}
}

func (this *Store) GetResultVideos(str string) (result []*domain.YoutubeSearchResult, err error) {
	return this.Yt.SearchYoutubeByKey(str)
}
