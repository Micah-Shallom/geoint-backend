package storage

import (
	"gorm.io/gorm"

	"github.com/Micah-Shallom/geoint-backend/utility"
	"github.com/minio/minio-go/v7"
)

type Database struct {
	Postgresql *gorm.DB
	Minio      *minio.Client
}

var (
	DB     *Database = &Database{}
	Logger *utility.Logger
)

func Connection() *Database {
	return DB
}
