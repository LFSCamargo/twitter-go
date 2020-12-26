package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/LFSCamargo/twitter-go/auth"
	"github.com/LFSCamargo/twitter-go/graph/generated"
	"github.com/LFSCamargo/twitter-go/graph/model"
	"github.com/LFSCamargo/twitter-go/graph/services/tweets"
	"github.com/LFSCamargo/twitter-go/graph/services/user"
)

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.TokenOutput, error) {
	return user.LoginUser(input)
}

func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (*model.TokenOutput, error) {
	return user.RegisterNewUser(input)
}

func (r *mutationResolver) AddReply(ctx context.Context, input model.CreateTweet) (*model.Reply, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteReply(ctx context.Context, input string) (*model.MessageOutput, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateTweet(ctx context.Context, input model.CreateTweet, tweet string) (*model.Tweet, error) {
	return tweets.CreateTweet(ctx, input)
}

func (r *mutationResolver) DeleteTweet(ctx context.Context, id string) (*model.MessageOutput, error) {
	return tweets.DeleteTweet(ctx, id)
}

func (r *mutationResolver) LikeTweet(ctx context.Context, id string) (*model.Tweet, error) {
	return tweets.LikeTweet(ctx, id)
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, errors.New("Not Logged")
	}
	return &model.User{
		Email:    user.Email,
		Username: user.Username,
		Picture:  user.Picture,
		ID:       user.ID.Hex(),
	}, nil
}

func (r *queryResolver) Replies(ctx context.Context, input *model.PaginationInput, id string) (*model.RepliesPaginationOutput, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Reply(ctx context.Context, id string) (*model.Reply, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Tweets(ctx context.Context, input *model.PaginationInput) (*model.TweetsPaginationOutput, error) {
	return tweets.GetTweets(ctx, input)
}

func (r *queryResolver) Tweet(ctx context.Context, id string) (*model.Tweet, error) {
	return tweets.GetTweet(ctx, id)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
