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

type RoomAlreadyTakenError struct {
	Room int
}

func (e *RoomAlreadyTakenError) Error() string {
	return fmt.Sprintf("комната с номером %d уже забронирована на ваши даты или такой комнаты не существует", e.Room)
}
