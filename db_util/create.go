package db_util

import "github.com/jchavannes/jgo/jerr"

func Create(_db DB, value interface{}) error {
	db, err := _db.Get()
	if err != nil {
		return jerr.Get("error getting db", err)
	}
	result := db.Create(value)
	if result.Error != nil {
		return jerr.Get("error running query", result.Error)
	}
	return nil
}

func CreateMany(_db DB, objects []interface{}) error {
	db, err := _db.Get()
	if err != nil {
		return jerr.Get("error getting db", err)
	}
	err = bulkInsert(db, objects)
	if err != nil {
		return jerr.Get("error with bulk insert", err)
	}
	return nil
}
