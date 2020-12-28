package replyModel

import (
	"errors"

	"github.com/LFSCamargo/twitter-go/utils/array"

	"github.com/LFSCamargo/twitter-go/constants"
	tweetModel "github.com/LFSCamargo/twitter-go/database/models/tweets"
	userModel "github.com/LFSCamargo/twitter-go/database/models/user"
	"github.com/LFSCamargo/twitter-go/graph/model"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Reply - is the mgm model for the tweet inside mongo
type Reply struct {
	mgm.DefaultModel `bson:",inline"`
	Reference        string   `json:"reference" bson:"reference"`
	Text             string   `json:"text" bson:"text"`
	User             string   `json:"user" bson:"user"`
	Likes            []string `json:"likes" bson:"likes"`
	Active           bool     `json:"active" bson:"active"`
}

// LikeReply - adds another person like to the reply
func LikeReply(tweetID string, userID string) (*Reply, error) {
	reply := &Reply{}
	coll := mgm.Coll(reply)
	findErr := coll.FindByID(tweetID, reply)
	if findErr != nil {
		return nil, errors.New(constants.NotFound)
	}

	_, found := array.FindItem(reply.Likes, userID)

	if found {
		reply.Likes = array.RemoveItem(reply.Likes, userID)
	} else {
		reply.Likes = append(reply.Likes, userID)
	}

	updateErr := mgm.Coll(reply).Update(reply)

	if updateErr != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return GetReply(tweetID)
}

// GetReplies - gets all the replies from a tweet
func GetReplies(limit int, reference string) (*model.RepliesPaginationOutput, error) {
	dbReply, err := tweetModel.GetTweet(reference)

	if err != nil {
		return nil, errors.New(constants.NotFound)
	}

	result := []*Reply{}
	first := int64(limit)
	resultErr := mgm.Coll(&Reply{}).SimpleFind(
		&result,
		bson.M{"reference": dbReply.ID.Hex()},
		&options.FindOptions{
			Limit: &first,
		},
	)
	if resultErr != nil {
		return nil, resultErr
	}
	total := []*Reply{}
	totalerr := mgm.Coll(&Reply{}).SimpleFind(&result, bson.M{})
	if totalerr != nil {
		return nil, totalerr
	}

	totalCount := len(total)
	count := len(result)

	replies := []*model.Reply{}

	for _, reply := range result {
		replies = append(replies, AdaptToGqlReply(reply))
	}

	return &model.RepliesPaginationOutput{
		PageInfo: &model.PageInfo{
			HasNextPage: totalCount > count,
			PageSize:    count,
		},
		Replies: replies,
	}, nil
}

// CreateReply - it creates a reply to a tweet
func CreateReply(input model.CreateTweet, reference string, userID string) (*Reply, error) {
	tweet, err := tweetModel.GetTweet(reference)

	if err != nil {
		return nil, errors.New(constants.NotFound)
	}

	newReply := &Reply{
		Text:      input.Text,
		User:      userID,
		Likes:     []string{},
		Active:    true,
		Reference: tweet.ID.Hex(),
	}

	error := mgm.Coll(newReply).Create(newReply)

	return newReply, error
}

// DeleteReply - it marks a tweet that already got created as active = false
func DeleteReply(id string, user *userModel.User) error {
	reply := &Reply{}
	coll := mgm.Coll(reply)
	findErr := coll.FindByID(id, reply)

	if findErr != nil {
		return errors.New(constants.NotFound)
	}

	if reply.User != user.ID.Hex() {
		return errors.New(constants.InsuficientPermissions)
	}

	reply.Active = false
	updateErr := mgm.Coll(reply).Update(reply)

	if updateErr != nil {
		return errors.New(constants.InternalServerError)
	}

	return nil
}

// GetReply - Finds a reply by it's id
func GetReply(id string) (*Reply, error) {
	reply := &Reply{}
	coll := mgm.Coll(reply)
	findErr := coll.FindByID(id, reply)
	if findErr != nil {
		return nil, errors.New(constants.NotFound)
	}
	return reply, nil
}

// AdaptToGqlReply - Adapter to match the GQL output
func AdaptToGqlReply(reply *Reply) *model.Reply {
	likes := []*model.User{}

	for _, like := range reply.Likes {
		user, _ := userModel.FindByID(like)
		likes = append(likes, &model.User{
			Email:    user.Email,
			ID:       user.ID.Hex(),
			Picture:  user.Picture,
			Username: user.Username,
		})
	}

	user, _ := userModel.FindByID(reply.User)

	return &model.Reply{
		ID:    reply.ID.Hex(),
		Likes: likes,
		User: &model.User{
			Email:    user.Email,
			ID:       user.ID.Hex(),
			Picture:  user.Picture,
			Username: user.Username,
		},
		Text: reply.Text,
	}
}
