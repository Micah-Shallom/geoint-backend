package analysis

import (
	"encoding/json"
	"fmt"

	"github.com/Micah-Shallom/geoint-backend/external/external_models"
	"github.com/Micah-Shallom/geoint-backend/external/request"
	"github.com/Micah-Shallom/geoint-backend/internal/models"
)

func (w *AnalysisWorker) processGIS(analysis *models.Analysis) error {
	w.Logger.Info("Starting GIS processing", analysis.ID)

	// Update progres
	progress := models.AnalysisProgress{
		GISProcessing:    "in_progress",
		VLMAnalysis:      "pending",
		RAGRetrieval:     "pending",
		ReportGeneration: "pending",
	}
	progressJSON, _ := json.Marshal(progress)

	w.DB.Postgresql.Model(&models.Analysis{}).
		Where("id = ?", analysis.ID).
		Update("progress", progressJSON)

	w.sendWebSocketUpdate(analysis.ID, "processing", progress, "Processing satellite imagery...")

	gisReq := external_models.GISChangeDetectionRequest{
		PastImageURL:    analysis.PastImageryPath,
		PresentImageURL: analysis.PresentImageryPath,
	}

	resp, err := w.ExtReq.SendExternalRequest(request.ProcessGIS, gisReq)
	gisResp, ok := resp.(external_models.GISChangeDetectionResponse)
	if !ok {
		return fmt.Errorf("GIS service error: %v", err)
	}

	w.DB.Postgresql.Model(&models.Analysis{}).
		Where("id = ?", analysis.ID).
		Update("change_map_path", gisResp.ChangeMapURL)

	analysis.ChangeMapPath = gisResp.ChangeMapURL

	progress.GISProcessing = "completed"
	progressJSON, _ = json.Marshal(progress)

	w.DB.Postgresql.Model(&models.Analysis{}).
		Where("id = ?", analysis.ID).
		Update("progress", progressJSON)

	w.sendWebSocketUpdate(analysis.ID, "processing", progress, "GIS processing completed")
	w.Logger.Info("GIS processing completed", analysis.ID)

	return nil
}
