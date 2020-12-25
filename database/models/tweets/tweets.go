package tweetModel

import (
	"github.com/kamva/mgm/v3"
	"github.com/LFSCamargo/twitter-go/graph/model"
)

// Tweet - is the mgm model for the tweet inside mongo
type Tweet struct {
	mgm.DefaultModel `bson:",inline"`
	ID               string   `json:"id" bson:"id"`
	Text             string   `json:"text" bson:"text"`
	User             string   `json:"user" bson:"user"`
	Likes            []string `json:"likes" bson:"likes"`
}

func CreateNewTweet(input model.)
