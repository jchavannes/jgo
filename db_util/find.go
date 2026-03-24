package db_util

import (
	"reflect"

	"github.com/jchavannes/jgo/jerr"
)

func Find(_db DB, out interface{}, where ...interface{}) error {
	db, err := _db.Get()
	if err != nil {
		return jerr.Get("error getting db", err)
	}
	result := db.Find(out, where...)
	if result.Error != nil {
		return jerr.Get("error running query", result.Error)
	}
	if result.RowsAffected == 0 && !isSlice(out) {
		return RecordNotFoundError
	}
	return nil
}

func FindPreload(_db DB, columns []string, out interface{}, where ...interface{}) error {
	db, err := _db.Get()
	if err != nil {
		return jerr.Get("error getting db", err)
	}
	for _, column := range columns {
		db = db.Preload(column)
	}
	result := db.Find(out, where...)
	if result.Error != nil {
		return jerr.Get("error running query", result.Error)
	}
	if result.RowsAffected == 0 && !isSlice(out) {
		return RecordNotFoundError
	}
	return nil
}

func isSlice(v interface{}) bool {
	rt := reflect.TypeOf(v)
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return rt.Kind() == reflect.Slice
}

func First(_db DB, out interface{}, where ...interface{}) error {
	db, err := _db.Get()
	if err != nil {
		return jerr.Get("error getting db", err)
	}
	result := db.First(out, where...)
	if result.Error != nil {
		return jerr.Get("error running query", result.Error)
	}
	return nil
}

func Last(_db DB, out interface{}, where ...interface{}) error {
	db, err := _db.Get()
	if err != nil {
		return jerr.Get("error getting db", err)
	}
	result := db.Last(out, where...)
	if result.Error != nil {
		return jerr.Get("error running query", result.Error)
	}
	return nil
}
