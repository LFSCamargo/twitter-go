package userModel

import (
	"errors"

	"github.com/LFSCamargo/twitter-go/constants"
	"github.com/LFSCamargo/twitter-go/graph/model"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// User - is the mgm model for the user inside mongo
type User struct {
	mgm.DefaultModel `bson:",inline"`
	Username         string  `json:"username" bson:"username"`
	Email            string  `json:"email" bson:"email"`
	Bio              *string `json:"bio" bson:"bio"`
	Nickname         *string `json:"nickname" bson:"nickname"`
	Picture          *string `json:"picture" bson:"picture"`
	Password         string  `json:"password" bson:"password"`
}

// UpdateProfile - Updates the user information Bio, Nickname and Picture
func UpdateProfile(input *model.UpdateProfileInput, userID string) (*User, error) {
	user := &User{}
	coll := mgm.Coll(user)
	err := coll.FindByID(userID, user)

	if err != nil {
		return nil, errors.New(constants.NotFound)
	}

	if input.Picture != nil {
		user.Picture = input.Picture
	}
	if input.Bio != nil {
		user.Bio = input.Bio
	}
	if input.Nickname != nil {
		user.Nickname = input.Nickname
	}

	updateErr := mgm.Coll(user).Update(user)

	if updateErr != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return user, nil
}

// CreateNewUser - creates a new user inside the database using the mgm
func CreateNewUser(user model.RegisterInput) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedPassword := string(hash)
	dbUser := &User{
		Username: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
	}
	error := mgm.Coll(dbUser).Create(dbUser)

	return dbUser, error
}

// FindByEmail - Finds a user by email
func FindByEmail(email string) (*User, error) {
	user := &User{Email: email}
	coll := mgm.Coll(user)
	err := coll.First(bson.M{"email": email}, user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByID - Finds a user by id
func FindByID(id string) (*User, error) {
	user := &User{}
	coll := mgm.Coll(user)
	err := coll.FindByID(id, user)

	return user, err
}

// AdaptUserModelToGql - Is an adapter to make the output match graphQL
func AdaptUserModelToGql(user *User) *model.User {
	return &model.User{
		Bio:      user.Bio,
		Email:    user.Email,
		Username: user.Username,
		Picture:  user.Picture,
		ID:       user.ID.Hex(),
		Nickname: user.Nickname,
	}
}
