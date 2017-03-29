package db

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var _db *gorm.DB

func getDb() (*gorm.DB, error) {
	if _db == nil {
		var err error
		_db, err = gorm.Open("sqlite3", "jgo-example.db")
		if err != nil {
			return _db, errors.New("Failed to connect to database\n")
		}
		interfaces := []interface{}{
			User{},
			Session{},
			Message{},
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

func save(value interface{}) *gorm.DB {
	db, _ := getDb()
	if db.Error != nil {
		fmt.Printf("Db error: %s\n", db.Error)
		return db
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
	if result.Error != nil && !result.RecordNotFound() {
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
