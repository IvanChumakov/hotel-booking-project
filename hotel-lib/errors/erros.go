package errors

import "fmt"

type LoginExistsError struct {
	Login string
}

func (e *LoginExistsError) Error() string {
	return fmt.Sprintf("такой логин уже существует %s", e.Login)
}

type AuthError struct{}

func (e *AuthError) Error() string {
	return "такого пользователя не существует"
}
