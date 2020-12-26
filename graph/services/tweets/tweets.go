package tweets

import (
	"context"
	"errors"

	"github.com/LFSCamargo/twitter-go/auth"
	"github.com/LFSCamargo/twitter-go/constants"
	tweetModel "github.com/LFSCamargo/twitter-go/database/models/tweets"
	"github.com/LFSCamargo/twitter-go/graph/model"
)

// GetTweet - Gets a tweet from mongo and validate user
func GetTweet(ctx context.Context, id string) (*model.Tweet, error) {
	tweet, err := tweetModel.GetTweet(id)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}
	return tweetModel.AdaptToGqlTweet(tweet), nil
}

// GetTweets - gets the tweets and returns a connection
func GetTweets(ctx context.Context, input *model.PaginationInput) (*model.TweetsPaginationOutput, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, errors.New(constants.NotLogged)
	}
	first := 10
	if input != nil {
		first = input.First
	}
	page, err := tweetModel.GetTweets(first)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}
	return page, nil
}

// DeleteTweet - Logic to delete the tweet
func DeleteTweet(ctx context.Context, id string) (*model.MessageOutput, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, errors.New(constants.NotLogged)
	}
	err := tweetModel.DeleteTweet(id)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}
	return &model.MessageOutput{
		Message: constants.TweetDeleted,
	}, nil
}

// CreateTweet - Logic to create the tweet
func CreateTweet(ctx context.Context, input model.CreateTweet) (*model.Tweet, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, errors.New(constants.NotLogged)
	}
	tweet, tweetErr := tweetModel.CreateNewTweet(input, user.ID.Hex())
	if tweetErr != nil {
		return nil, tweetErr
	}

	return tweetModel.AdaptToGqlTweet(tweet), nil
}
