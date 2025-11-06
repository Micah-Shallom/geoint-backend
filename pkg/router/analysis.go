package router

import (
	"fmt"

	"github.com/Micah-Shallom/geoint-backend/external/request"
	"github.com/Micah-Shallom/geoint-backend/pkg/controller/analysis"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage"
	"github.com/Micah-Shallom/geoint-backend/utility"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Analysis(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	analysis := analysis.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	analysisUrl := r.Group(fmt.Sprintf("%v/analysis", ApiVersion))
	{
		analysisUrl.POST("/submit", analysis.SubmitAnalysis)
	}
	return r
}
