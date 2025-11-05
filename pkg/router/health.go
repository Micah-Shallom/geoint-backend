package router

import (
	"fmt"

	"github.com/Micah-Shallom/geoint-backend/external/request"
	"github.com/Micah-Shallom/geoint-backend/pkg/controller/health"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage"
	"github.com/Micah-Shallom/geoint-backend/utility"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Health(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	health := health.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	healthUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		healthUrl.POST("/health", health.Post)
		healthUrl.GET("/health", health.Get)
	}
	return r
}
