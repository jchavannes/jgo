package db

import (
	"errors"
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
