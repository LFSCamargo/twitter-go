package reply

import (
	"context"
	"errors"

	"github.com/LFSCamargo/twitter-go/auth"
	"github.com/LFSCamargo/twitter-go/constants"
	replyModel "github.com/LFSCamargo/twitter-go/database/models/reply"
	"github.com/LFSCamargo/twitter-go/graph/model"
)

// CreateTweet - Create a tweet as a reply to another tweet
func CreateTweet(ctx context.Context, input model.CreateTweet, id string) (*model.Reply, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, errors.New(constants.NotLogged)
	}

	reply, err := replyModel.CreateReply(input, id, user.ID.Hex())

	if err != nil {
		return nil, err
	}

	return replyModel.AdaptToGqlReply(reply), nil
}

// DeleteReply - Deletes a reply linked to a tweet
func DeleteReply(ctx context.Context, id string) (*model.MessageOutput, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, errors.New(constants.NotLogged)
	}

	err := replyModel.DeleteReply(id, user)

	if err != nil {
		return nil, err
	}

	return &model.MessageOutput{
		Message: constants.ReplyDeleted,
	}, nil
}

// GetReplies - Gets all the replies linked to a tweet
func GetReplies(ctx context.Context, input *model.PaginationInput, tweetID string) (*model.RepliesPaginationOutput, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, errors.New(constants.NotLogged)
	}

	limit := 10
	if input != nil {
		limit = input.First
	}

	repliesConn, err := replyModel.GetReplies(limit, tweetID)

	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return repliesConn, err
}

// GetReply - Gets the reply
func GetReply(ctx context.Context, replyID string) (*model.Reply, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, errors.New(constants.NotLogged)
	}

	reply, err := replyModel.GetReply(replyID)

	if err != nil {
		return nil, err
	}

	return replyModel.AdaptToGqlReply(reply), nil
}

// LikeReply - likes a reply from a tweet
func LikeReply(ctx context.Context, replyID string) (*model.Reply, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, errors.New(constants.NotFound)
	}

	reply, err := replyModel.LikeReply(replyID, user.ID.Hex())

	if err != nil {
		return nil, err
	}

	return replyModel.AdaptToGqlReply(reply), nil
}
