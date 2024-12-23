package auth

import (
	"fmt"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/database"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/errors"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

func Register(user models.User) (string, error) {
	exists, err := database.Exists(models.UserLogin{
		Login:    user.Login,
		Password: user.Password,
	}, false)

	if err != nil {
		return "", err
	}
	if exists {
		return "", &errors.LoginExistsError{Login: user.Login}
	}
	err = database.AddUser(user)
	if err != nil {
		return "", err
	}
	token, err := createToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func Login(credentials models.UserLogin) (string, error) {
	exists, err := database.Exists(credentials, true)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", &errors.AuthError{}
	}
	foundUser, err := database.GetUser(credentials.Login)
	if err != nil {
		return "", err
	}
	token, err := createToken(foundUser)
	if err != nil {
		return "", err
	}
	return token, nil
}

func createToken(user models.User) (string, error) {
	if user.Role == "" {
		user.Role = "customer"
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": user.Role,
		"name": user.Login,
	})
	token, err := claims.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error while signing jwt token: %v", err)
	}
	return token, nil
}
