package domain

import (
	"gopkg.in/mgo.v2/bson"
)

type Video struct {
	Id      bson.ObjectId `json:"id" bson:"_id" `
	VideoId string        `json:"videoid" bson:"videoid"`
	//Transcript xml.Name      `json:"transcript bson:"transcript" xml:"transcript"`
	Content []struct {
		Start string `json:"start" bson:"start" xml:"start,attr"`
		Dur   string `json:"duration" bson:"duration" xml:"dur,attr"`
		Text  string `json:"text" bson:"text" xml:",innerxml"`
	} `json:"content" bson:"content" xml:"text"`
}

type YoutubeSearchResult struct {
	Title     string
	YoutubeId string
}
