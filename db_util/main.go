package db_util

import (
	"github.com/jchavannes/gorm"
)

const (
	DefaultLimit = 25
	MediumLimit  = 100
	LargeLimit   = 500
	ExLargeLimit = 1000
	HugeLimit    = 5000
)

type DB interface {
	Get() (*gorm.DB, error)
}
