package db

import (
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"errors"
)

type Session struct {
	Id       uint `gorm:"primary_key"`
	CookieId string `gorm:"unique_index"`
	UserId   uint
	StartTs  uint
}

func (s *Session) Save() error {
	db, err := getDb()
	if err != nil {
		return err
	}
	result := db.Save(s)
	return result.Error
}

func (s *Session) IsLoggedIn() bool {
	session, err := GetSession(s.CookieId)
	loggedIn := false
	if err == nil {
		_, err := GetUser(User{
			Id: session.UserId,
		})
		if err == nil {
			loggedIn = true
		}
	}
	return loggedIn
}

func (s *Session) Signup(username string, password string) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	user, err := CreateUser(username, string(hashedPassword))
	if err != nil {
		return User{}, errors.New(fmt.Sprintf("Error signing up: %s\n", err))
	}
	session, err := GetSession(s.CookieId)
	if err != nil {
		return User{}, errors.New(fmt.Sprintf("Error getting session: %s\n", err))
	}
	session.UserId = user.Id
	err = session.Save()
	if err != nil {
		return User{}, errors.New(fmt.Sprintf("Error saving session: %s\n", err))
	}
	return user, nil
}

func (s *Session) Login(username string, password string) (User, error) {
	user, err := GetUser(User{Username: username})
	if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return User{}, err
	}

	session, err := GetSession(s.CookieId)
	if err != nil {
		return User{}, errors.New(fmt.Sprintf("Error getting session: %s\n", err))
	}

	session.UserId = user.Id
	err = session.Save()
	if err != nil {
		return User{}, errors.New(fmt.Sprintf("Error saving session: %s\n", err))
	}

	return user, nil
}

func GetSession(cookieId string) (*Session, error) {
	session := &Session{
		CookieId: cookieId,
	}
	db, err := getDb()
	if err != nil {
		return session, err
	}
	result := db.Find(session, session)
	if result.Error != nil && result.Error.Error() == "record not found" {
		result = db.Create(session)
	}
	return session, result.Error
}
