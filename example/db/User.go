package db

import (
	"errors"
)

type User struct {
	Id           uint `gorm:"primary_key"`
	Username     string `gorm:"unique_index"`
	PasswordHash string
}

func GetUser(user User) (User, error) {
	if user.Id == 0 && user.Username == "" {
		return user, errors.New("Must set either Id or Username.")
	}
	result := find(&user, &user)
	return user, result.Error
}

func CreateUser(username string, hashedPassword string) (User, error) {
	user := User{
		Username: username,
		PasswordHash: string(hashedPassword),
	}
	result := create(&user)
	return user, result.Error
}
