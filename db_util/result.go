package db_util

import (
	"gorm.io/gorm"
	"github.com/jchavannes/jgo/jerr"
)

func Result(result *gorm.DB) (*gorm.DB, error) {
	if result.Error != nil {
		return result, jerr.Get("error with result", result.Error)
	}
	return result, nil
}
