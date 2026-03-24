package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _db *gorm.DB

func getDb() (*gorm.DB, error) {
	if _db == nil {
		var err error
		_db, err = gorm.Open(sqlite.Open("jgo-example.db"), &gorm.Config{})
		if err != nil {
			return _db, errors.New("Failed to connect to database\n")
		}
		err = _db.AutoMigrate(
			&User{},
			&Session{},
			&Message{},
		)
		if err != nil {
			return _db, err
		}
	}
	return _db, nil
}

func save(value interface{}) *gorm.DB {
	db, _ := getDb()
	if db == nil {
		fmt.Printf("Db error: failed to get db\n")
		return &gorm.DB{}
	}
	result := db.Save(value)
	if result.Error != nil {
		fmt.Printf("Db error: %s\n", result.Error)
		return result
	}
	return result
}

func find(out interface{}, where ...interface{}) *gorm.DB {
	db, _ := getDb()
	result := db.Find(out, where...)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Printf("Db error: %s\n", result.Error)
		return result
	}
	return result
}

func create(value interface{}) *gorm.DB {
	db, _ := getDb()
	result := db.Create(value)
	if result.Error != nil {
		fmt.Printf("Db error: %s\n", result.Error)
		return result
	}
	return result
}
