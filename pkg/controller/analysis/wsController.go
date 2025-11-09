package analysis

import (
	"fmt"
	"net/http"

	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage"
	"github.com/Micah-Shallom/geoint-backend/services/analysis"
	"github.com/Micah-Shallom/geoint-backend/services/websocket"
	"github.com/Micah-Shallom/geoint-backend/utility"
	"github.com/gin-gonic/gin"
)

type WSController struct {
	Hub    *websocket.Hub
	Logger *utility.Logger
	Db     *storage.Database
}

func (base *WSController) HandleWebSocket(c *gin.Context) {
	analysisID := c.Param("analysis_id")
	if analysisID == "" {
		base.Logger.Info("empty analysis id")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "analysis id is empty", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := analysis.WsService(c, base.Db, base.Logger, base.Hub, analysisID); err != nil {
		base.Logger.Info(fmt.Sprintf("error setting up websocket service for analysis: %d", analysisID), err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "error setting up websocket service", err, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("Websocket connection establised")
	rd := utility.BuildSuccessResponse(http.StatusOK, "websocket connection established", nil)
	c.JSON(http.StatusOK, rd)
}
