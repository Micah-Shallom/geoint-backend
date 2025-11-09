package analysis

import (
	"fmt"
	"net/http"

	"github.com/Micah-Shallom/geoint-backend/internal/models"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage/postgresql"
	"github.com/Micah-Shallom/geoint-backend/services/websocket"
	"github.com/Micah-Shallom/geoint-backend/utility"
	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, implement proper origin checking
	},
}

func WsService(c *gin.Context, db *storage.Database, logger *utility.Logger, hub *websocket.Hub, analysisID string) error {
	var (
		analysis models.Analysis
	)

	exists := postgresql.CheckExists(db.Postgresql, &analysis, "id = ?", analysisID)
	if !exists {
		return fmt.Errorf("analysis id not found")
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("failed to upgrade connection", err)
		return err
	}

	client := &websocket.Client{
		Hub:        hub,
		Conn:       conn,
		Send:       make(chan []byte, 256),
		AnalysisID: analysisID,
	}

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()

	return nil

}
