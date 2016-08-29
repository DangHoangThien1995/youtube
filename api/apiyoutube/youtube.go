package apiyoutube

import (
	"flag"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"youtube/domain"
)

type YT struct {
}

func NewYT() *YT {
	return &YT{}
}

var (
	maxResults = flag.Int64("max-results", 50, "Max YouTube results")
	service    *youtube.Service
	response   *youtube.SearchListResponse
	query      = flag.String("query", "str", "Search term")
	resultType = flag.String("type", "channel", "video")
)

const developerKey = "AIzaSyB9P-_Ep22v_8A_dUD22SpJxxxNkNgIgK4"

func (yt *YT) SearchYoutubeByKey(str string) (result []*domain.YoutubeSearchResult, err error) {
	flag.Parse()

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err = youtube.New(client)
	if err != nil {
		return
	}

	// Make the API call to YouTube.
	call := service.Search.List("id,snippet").
		Type("video").
		ChannelId(str).
		MaxResults(*maxResults).
		VideoCaption("closedCaption")
		//
	response, err = call.Do()
	if err != nil {
		return
	}

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			result = append(result, &domain.YoutubeSearchResult{Title: item.Snippet.Title, YoutubeId: item.Id.VideoId})
		}
	}
	return
}
