package db_util

import (
	"github.com/jchavannes/gorm"
	"github.com/jchavannes/jgo/jerr"
)

func RetryFind(db *gorm.DB, out interface{}, where ...interface{}) error {
	result := db.Find(out, where...)
	if result.Error == nil {
		return nil
	} else if IsInvalidConnectionError(result.Error) {
		// Retry if connection error
		result = db.Find(out, where...)
		if result.Error == nil {
			return nil
		}
	}
	return jerr.Get("error running query", result.Error)
}
