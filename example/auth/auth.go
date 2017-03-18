package auth

import (
	"github.com/jchavannes/jgo/example/db"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
)

func Signup(username string, password string) (db.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, err
	}
	return db.CreateUser(username, string(hashedPassword))
}

func Login(username string, password string) (db.User, error) {
	user, err := db.GetUser(db.User{
		Username: username,
	})
	if err != nil {
		return db.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}
