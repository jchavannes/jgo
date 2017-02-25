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
	db, err := getDb()
	if err != nil {
		return user, err
	}
	result := db.Find(&user, &user)
	return user, result.Error
}
