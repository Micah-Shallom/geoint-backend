package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/Micah-Shallom/geoint-backend/internal/config"
	"github.com/Micah-Shallom/geoint-backend/internal/models/migrations"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage/minio"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage/postgresql"
	"github.com/Micah-Shallom/geoint-backend/pkg/router"
	"github.com/Micah-Shallom/geoint-backend/utility"
	"github.com/go-playground/validator/v10"
)

func main() {

	logger := utility.NewLogger() //Warning !!!!! Do not recreate anywhere 

	configuration := config.Setup(logger, "./app")
	postgresql.ConnectToDatabase(logger, configuration.Database)
	minio.ConnectToMinio(logger, configuration.Minio)

	db := storage.Connection()

	validatorRef := validator.New()
	validatorRef.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "-" {
			return ""
		}
		return name
	})
	utility.RegisterCustomValidations(validatorRef)

	if configuration.Database.Migrate {
		migrations.RunAllMigrations(db)
	}

	r := router.Setup(logger, validatorRef, db, &configuration.App)

	utility.LogAndPrint(logger, fmt.Sprintf("Server is starting at 127.0.0.1:%s", configuration.Server.Port))
	log.Fatal(r.Run(":" + configuration.Server.Port))
}
