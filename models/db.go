package models

import (
	"github.com/cSploit/daemon/models/internal"
	"github.com/jinzhu/gorm"
)

func GetDbInstance() gorm.DB {
	return *internal.Db
}
