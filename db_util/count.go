package db_util

import "github.com/jchavannes/jgo/jerr"

func Count(_db DB, where interface{}) (uint, error) {
	db, err := _db.Get()
	if err != nil {
		return 0, jerr.Get("error getting db", err)
	}
	var totalCount uint
	result := db.Model(where).Where(where).Count(&totalCount)
	if result.Error != nil {
		return 0, jerr.Get("error running query", result.Error)
	}
	return totalCount, nil
}
