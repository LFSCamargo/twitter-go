package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/LFSCamargo/graphql-go-boilerplate/auth"
	"github.com/LFSCamargo/graphql-go-boilerplate/graph/generated"
	"github.com/LFSCamargo/graphql-go-boilerplate/graph/model"
	"github.com/LFSCamargo/graphql-go-boilerplate/graph/services/user"
)

func (r *mutationResolver) Login(ctx context.Context, input *model.LoginInput) (*model.TokenOutput, error) {
	tokenOut, err := user.LoginUser(input)
	return tokenOut, err
}

func (r *mutationResolver) Register(ctx context.Context, input *model.RegisterInput) (*model.TokenOutput, error) {
	tokenOut, err := user.RegisterNewUser(input)
	return tokenOut, err
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
		ID:       user.ID.String(),
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
