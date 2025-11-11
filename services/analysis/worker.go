package analysis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Micah-Shallom/geoint-backend/external/request"
	"github.com/Micah-Shallom/geoint-backend/internal/models"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage/postgresql"
	"github.com/Micah-Shallom/geoint-backend/services/websocket"
	"github.com/Micah-Shallom/geoint-backend/utility"
)

type AnalysisWorker struct {
	DB     *storage.Database
	Logger *utility.Logger
	Hub    *websocket.Hub
	ExtReq request.ExternalRequest
}

func NewAnalysisWorker(db *storage.Database, logger *utility.Logger, hub *websocket.Hub) *AnalysisWorker {
	return &AnalysisWorker{
		DB:     db,
		Logger: logger,
		Hub:    hub,
		ExtReq: request.ExternalRequest{Logger: logger, Test: true},
	}
}

func (w *AnalysisWorker) ProcessAnalysis(analysisID string) {
	var (
		analysis  models.Analysis
		startTime = time.Now()
	)
	w.Logger.Info("Starting analysis processing for %s", analysisID)

	err, _ := postgresql.SelectOneFromDb(w.DB.Postgresql, &analysis, "id = ?", analysisID)
	if err != nil {
		w.Logger.Error("Failed to fetch analysis: %v", err)
		return
	}

	// GIS Processing
	if err := w.processGIS(&analysis); err != nil {
		w.Logger.Error("GIS processing failed: %v", err)
		w.updateStatus(&analysis, "failed", "GIS processing failed")
		return
	}

	// VLM Analysis
	if err := w.processVLM(analysisID); err != nil {
		w.Logger.Error("VLM analysis failed: %v", err)
		w.updateStatus(&analysis, "failed", "VLM analysis failed")
		return
	}

	// LLM Report Generation
	if err := w.processLLM(analysis.ID); err != nil {
		w.Logger.Error("Report generation failed: %v", err)
		w.updateStatus(&analysis, "failed", "Report generation failed")
		return
	}

	// Calculate processing time
	processingTime := time.Since(startTime)
	completedAt := time.Now()

	// Final update
	progress := models.AnalysisProgress{
		GISProcessing:    "completed",
		VLMAnalysis:      "completed",
		RAGRetrieval:     "completed",
		ReportGeneration: "completed",
	}
	progressJSON, _ := json.Marshal(progress)

	progressUpdate := map[string]any{
		"status":          "completed",
		"progress":        progressJSON,
		"processing_time": fmt.Sprintf("%dm %ds", int(processingTime.Minutes()), int(processingTime.Seconds())%60),
		"completed_at":    completedAt,
	}

	res, err := postgresql.UpdateFields(w.DB.Postgresql, &models.Analysis{}, progressUpdate, "id = ?", analysisID)
	if err != nil {
		w.Logger.Error("Failed to update analysis progress: %v", err)
	}

	if res.RowsAffected == 0 {
		w.Logger.Error("fields not updated")
		return
	}

	// Send final WebSocket update
	w.sendWebSocketUpdate(analysisID, "completed", progress, "Report ready")
	w.Logger.Info("Analysis processing completed for %s in %v", analysisID, processingTime)
}

func (w *AnalysisWorker) sendWebSocketUpdate(analysisID, status string, progress models.AnalysisProgress, currentStep string) {
	message := &models.Message{
		AnalysisID:  analysisID,
		Status:      status,
		Progress:    progress,
		CurrentStep: currentStep,
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	if w.Hub == nil || w.Hub.Broadcast == nil {
		w.Logger.Error("Hub or Hub.Broadcast is nil, skipping WebSocket update")
		return
	}

	w.Logger.Info("Sending WebSocket update for analysis %s: %s - %s", analysisID, status, currentStep)
	w.Hub.Broadcast <- message
}

func (w *AnalysisWorker) updateStatus(analysis *models.Analysis, status string, message string) {
	w.DB.Postgresql.Model(&models.Analysis{}).
		Where("id = ?", analysis.ID).
		Update("status", status)

	w.Logger.Info("Analysis status updated for %s to %s: %s", analysis.ID, status, message)
}
