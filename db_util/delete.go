package db_util

import "github.com/jchavannes/jgo/jerr"

func Delete(_db DB, value interface{}) error {
	db, err := _db.Get()
	if err != nil {
		return jerr.Get("error getting db", err)
	}
	result := db.Delete(value, value)
	if result.Error != nil {
		return jerr.Get("error deleting", result.Error)
	}
	return nil
}
