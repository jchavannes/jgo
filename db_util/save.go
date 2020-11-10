package db_util

import (
	"github.com/jchavannes/jgo/jerr"
)

func Save(_db DB, value interface{}) error {
	db, err := _db.Get()
	if err != nil {
		return jerr.Get("error getting db", err)
	}
	result := db.Save(value)
	if result.Error != nil {
		return jerr.Get("error saving value", result.Error)
	}
	return nil
}
