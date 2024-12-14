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
	exists, err := database.Exists(user, false)
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

func Login(user models.User) (string, error) {
	exists, err := database.Exists(user, true)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", &errors.AuthError{}
	}
	token, err := createToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func createToken(user models.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud":  user.Role,
		"name": user.Login,
	})
	token, err := claims.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error while signing jwt token: %v", err)
	}
	return token, nil
}
