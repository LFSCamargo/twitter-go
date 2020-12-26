package userModel

import (
	"github.com/LFSCamargo/twitter-go/graph/model"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// User - is the mgm model for the user inside mongo
type User struct {
	mgm.DefaultModel `bson:",inline"`
	Username         string `json:"username" bson:"username"`
	Email            string `json:"email" bson:"email"`
	Picture          string `json:"picture" bson:"picture"`
	Password         string `json:"password" bson:"password"`
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
		Picture:  user.Picture,
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
