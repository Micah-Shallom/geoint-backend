package router

import (
	"fmt"

	"github.com/Micah-Shallom/geoint-backend/external/request"
	"github.com/Micah-Shallom/geoint-backend/pkg/controller/analysis"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage"
	"github.com/Micah-Shallom/geoint-backend/services/websocket"
	"github.com/Micah-Shallom/geoint-backend/utility"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Analysis(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger, hub *websocket.Hub) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	analysisCtrl := analysis.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq, Hub: hub}
	wsCtrl := analysis.WSController{Hub: hub, Logger: logger, Db: db}

	analysisUrl := r.Group(fmt.Sprintf("%v/analysis", ApiVersion))
	{
		analysisUrl.POST("/submit", analysisCtrl.SubmitAnalysis)
		analysisUrl.GET("/ws/:analysis_id", wsCtrl.HandleWebSocket)
	}
	return r
}
