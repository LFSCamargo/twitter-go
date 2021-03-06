package tweetModel

import (
	"errors"

	"github.com/LFSCamargo/twitter-go/constants"
	userModel "github.com/LFSCamargo/twitter-go/database/models/user"
	"github.com/LFSCamargo/twitter-go/graph/model"
	"github.com/LFSCamargo/twitter-go/utils/array"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Tweet - is the mgm model for the tweet inside mongo
type Tweet struct {
	mgm.DefaultModel `bson:",inline"`
	Text             string   `json:"text" bson:"text"`
	User             string   `json:"user" bson:"user"`
	Likes            []string `json:"likes" bson:"likes"`
	Active           bool     `json:"active" bson:"active"`
}

// LikeTweet - adds another person like to the tweet
func LikeTweet(tweetID string, userID string) (*Tweet, error) {
	tweet := &Tweet{}
	coll := mgm.Coll(tweet)
	findErr := coll.FindByID(tweetID, tweet)

	if findErr != nil {
		return nil, errors.New(constants.NotFound)
	}

	_, found := array.FindItem(tweet.Likes, userID)

	if found {
		tweet.Likes = array.RemoveItem(tweet.Likes, userID)
	} else {
		tweet.Likes = append(tweet.Likes, userID)
	}

	updateErr := mgm.Coll(tweet).Update(tweet)

	if updateErr != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return GetTweet(tweetID)
}

// CreateNewTweet - it creates a new tweet inside the database
func CreateNewTweet(tweet model.CreateTweet, userID string) (*Tweet, error) {
	newTweet := &Tweet{
		Text:   tweet.Text,
		User:   userID,
		Likes:  []string{},
		Active: true,
	}

	error := mgm.Coll(newTweet).Create(newTweet)

	return newTweet, error
}

// DeleteTweet - it marks a tweet that already got created as active = false
func DeleteTweet(id string, user *userModel.User) error {
	tweet := &Tweet{}
	coll := mgm.Coll(tweet)
	findErr := coll.FindByID(id, tweet)

	if findErr != nil {
		return errors.New(constants.NotFound)
	}

	if tweet.User != user.ID.Hex() {
		return errors.New(constants.InsuficientPermissions)
	}

	tweet.Active = false
	updateErr := mgm.Coll(tweet).Update(tweet)

	if updateErr != nil {
		return errors.New(constants.InternalServerError)
	}

	return nil
}

// GetTweet - gets a tweet from its Mongo Id
func GetTweet(id string) (*Tweet, error) {
	tweet := &Tweet{}
	coll := mgm.Coll(tweet)
	findErr := coll.FindByID(id, tweet)
	if findErr != nil {
		return nil, errors.New(constants.NotFound)
	}
	return tweet, nil
}

// GetTweets - Gets tweets with pagination
func GetTweets(limit int) (*model.TweetsPaginationOutput, error) {
	result := []*Tweet{}
	first := int64(limit)
	resulterr := mgm.Coll(&Tweet{
		Active: true,
	}).SimpleFind(&result, bson.M{"active": true}, &options.FindOptions{
		Limit: &first,
	})
	if resulterr != nil {
		return nil, resulterr
	}

	total := []*Tweet{}
	totalerr := mgm.Coll(&Tweet{
		Active: true,
	}).SimpleFind(&result, bson.M{"active": true})
	if totalerr != nil {
		return nil, totalerr
	}

	totalCount := len(total)
	count := len(result)

	tweets := []*model.Tweet{}

	for _, tweet := range result {
		tweets = append(tweets, AdaptToGqlTweet(tweet))
	}

	return &model.TweetsPaginationOutput{
		PageInfo: &model.PageInfo{
			HasNextPage: totalCount > count,
			PageSize:    count,
		},
		Tweets: tweets,
	}, nil
}

// AdaptToGqlTweet - it adapts the output to the gql output
func AdaptToGqlTweet(tweet *Tweet) *model.Tweet {
	likes := []*model.User{}

	for _, like := range tweet.Likes {
		user, _ := userModel.FindByID(like)
		likes = append(likes, &model.User{
			Email:    user.Email,
			ID:       user.ID.Hex(),
			Picture:  user.Picture,
			Username: user.Username,
		})
	}

	user, _ := userModel.FindByID(tweet.User)

	return &model.Tweet{
		ID:    tweet.ID.Hex(),
		Likes: likes,
		User: &model.User{
			Email:    user.Email,
			ID:       user.ID.Hex(),
			Picture:  user.Picture,
			Username: user.Username,
		},
		Text: tweet.Text,
	}
}
