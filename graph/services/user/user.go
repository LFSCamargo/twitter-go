package user

import (
	"errors"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/LFSCamargo/graphql-go-boilerplate/constants"
	userModel "github.com/LFSCamargo/graphql-go-boilerplate/database/models/user"
	"github.com/LFSCamargo/graphql-go-boilerplate/graph/model"
	"github.com/dgrijalva/jwt-go"
)

// GetUserFromToken - Gets the user data from the token string
func GetUserFromToken(tokenStr string) (*userModel.User, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "topsecret"
	}
	byteSecret := []byte(secret)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return byteSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["user_id"].(string)
		dbUser, err := userModel.FindByEmail(email)
		return dbUser, err
	} else {
		return nil, errors.New("Invalid Token")
	}
}

func createToken(email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "topsecret"
	}
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = email
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// LoginUser - it's the function that check if the user exists
func LoginUser(user *model.LoginInput) (*model.TokenOutput, error) {
	dbUser, err := userModel.FindByEmail(user.Email)
	if err != nil {
		return nil, errors.New(constants.WrongEmailOrPassword)
	}

	if err != nil {
		return nil, err
	}

	passwordError := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))

	if passwordError != nil {
		return nil, errors.New(constants.WrongEmailOrPassword)
	}

	token, error := createToken(user.Email)

	if error != nil {
		return nil, error
	}

	return &model.TokenOutput{
		Token: token,
	}, nil
}

// RegisterNewUser - it's the function that creates a new db user and returns a jwt token
func RegisterNewUser(user *model.RegisterInput) (*model.TokenOutput, error) {
	dbUser, _ := userModel.FindByEmail(user.Email)

	if dbUser != nil {
		return nil, errors.New(constants.UserAlreadyRegistered)
	}

	_, err := userModel.CreateNewUser(user)
	if err != nil {
		return nil, err
	}

	token, error := createToken(user.Email)

	if error != nil {
		return nil, error
	}

	return &model.TokenOutput{
		Token: token,
	}, nil
}
