package storage

import (
	"gorm.io/gorm"

	"github.com/Micah-Shallom/geoint-backend/utility"
)

type Database struct {
	Postgresql *gorm.DB
}

var (
	DB     *Database = &Database{}
	Logger *utility.Logger
)

func Connection() *Database {
	return DB
}
