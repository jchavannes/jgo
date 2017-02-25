package db

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
)

func Signup(username string, password string) (User, error) {
	db, err := getDb()
	if err != nil {
		return User{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	user := User{
		Username: username,
		PasswordHash: string(hashedPassword),
	}
	result := db.Create(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func Login(username string, password string) (User, error) {
	user, err := GetUser(User{
		Username: username,
	})
	if err != nil {
		return User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return User{}, err
	}
	return user, nil
}

var _db *gorm.DB

func getDb() (*gorm.DB, error) {
	if _db == nil {
		var err error
		_db, err = gorm.Open("sqlite3", "go-example.db")
		if err != nil {
			return _db, errors.New("Failed to connect to database\n")
		}
		interfaces := []interface{}{
			User{},
			Session{},
		}
		for _, iface := range interfaces {
			result := _db.AutoMigrate(iface)
			if result.Error != nil {
				return result, result.Error
			}
		}
	}
	return _db, nil
}
