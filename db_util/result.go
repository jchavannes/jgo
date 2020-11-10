package db_util

import (
	"github.com/jchavannes/gorm"
	"github.com/jchavannes/jgo/jerr"
)

func Result(result *gorm.DB) (*gorm.DB, error) {
	if result.Error != nil {
		return result, jerr.Get("error with result", result.Error)
	}
	return result, nil
}
