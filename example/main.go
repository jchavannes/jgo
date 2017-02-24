package main

import (
	"github.com/jchavannes/jgo/web"
	"golang.org/x/crypto/bcrypt"
	"git.jasonc.me/sandbox/go-lang-idea-plugin/testData/mockSdk-1.1.2/src/pkg/fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"errors"
)

func main() {
	server := web.Server{
		Port: 8248,
		Sessions: true,
		TemplateDir: "templates",
		StaticDir: "pub",
		Routes: []web.Route{{
			Pattern: "/",
			Handler: func(r *web.Request) {
				r.Render("index")
			},
		}, {
			Pattern: "/signup",
			Handler: func(r *web.Request) {
				r.Render("signup")
			},
		}, {
			Pattern: "/signup-submit",
			CsrfProtect: true,
			Handler: func(r *web.Request) {
				username := r.GetFormValue("username")
				password := r.GetFormValue("password")
				user, err := signup(username, password)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
				} else {
					fmt.Printf("User: %#v\n", user)
				}
			},
		}},
	}
	server.Run()
}

type User struct {
	Id           uint `gorm:"primary_key"`
	Username     string `gorm:"unique_index"`
	PasswordHash string
}

func signup(username string, password string) (User, error) {
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
