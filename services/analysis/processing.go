package analysis

import (
	"encoding/json"
	"fmt"

	"github.com/Micah-Shallom/geoint-backend/external/external_models"
	"github.com/Micah-Shallom/geoint-backend/external/request"
	"github.com/Micah-Shallom/geoint-backend/internal/models"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage/postgresql"
)

func (w *AnalysisWorker) processGIS(analysis *models.Analysis) error {
	w.Logger.Info("Starting GIS processing for analysis %s", analysis.ID)

	progress := models.AnalysisProgress{
		GISProcessing:    "in_progress",
		VLMAnalysis:      "pending",
		RAGRetrieval:     "pending",
		ReportGeneration: "pending",
	}
	progressJSON, err := json.Marshal(progress)
	if err != nil {
		w.Logger.Error("Failed to marshal GIS progressJSON for analysis %s: %v", analysis.ID, err)
		return fmt.Errorf("failed to marshal GIS progressJSON for analysis: %s", analysis.ID)
	}

	res, err := postgresql.UpdateFields(w.DB.Postgresql, &models.Analysis{}, map[string]any{
		"progress": progressJSON,
	}, "id = ?", analysis.ID)
	if err != nil {
		w.Logger.Error("Failed to update progress for analysis %s: %v", analysis.ID, err)
		return fmt.Errorf("failed to update progress: %v", err)
	}
	if res.RowsAffected == 0 {
		w.Logger.Error("Failed to update progress for analysis %s: %s", analysis.ID, "no rows affected")
		return fmt.Errorf("failed to update progress: no rows affected")
	}

	w.sendWebSocketUpdate(analysis.ID, "processing", progress, "Processing satellite imagery...")

	gisReq := external_models.GISChangeDetectionRequest{
		PastImageURL:    analysis.PastImageryPath,
		PresentImageURL: analysis.PresentImageryPath,
	}

	resp, err := w.ExtReq.SendExternalRequest(request.ProcessGIS, gisReq)
	gisResp, ok := resp.(external_models.GISChangeDetectionResponse)
	if !ok {
		w.Logger.Error("Failed to process GIS for analysis %s: %v", analysis.ID, err)
		return fmt.Errorf("GIS service error: %v", err)
	}

	res, err = postgresql.UpdateFields(w.DB.Postgresql, &models.Analysis{}, map[string]any{
		"change_map_path": gisResp.ChangeMapURL,
	}, "id = ?", analysis.ID)
	if err != nil {
		w.Logger.Error("Failed to update change map path for analysis %s: %v", analysis.ID, err)
		return fmt.Errorf("failed to update change map path: %v", err)
	}
	if res.RowsAffected == 0 {
		w.Logger.Error("Failed to update change map path for analysis %s: %s", analysis.ID, "no rows affected")
		return fmt.Errorf("failed to update change map path: no rows affected")
	}

	analysis.ChangeMapPath = gisResp.ChangeMapURL

	progress.GISProcessing = "completed"
	progressJSON, err = json.Marshal(progress)
	if err != nil {
		w.Logger.Error("Failed to marshal GIS progressJSON for analysis %s: %v", analysis.ID, err)
		return fmt.Errorf("failed to marshal GIS progressJSON for analysis: %s", analysis.ID)
	}

	res, err = postgresql.UpdateFields(w.DB.Postgresql, &models.Analysis{}, map[string]any{
		"progress": progressJSON,
	}, "id = ?", analysis.ID)
	if err != nil {
		w.Logger.Error("Failed to update progress for analysis %s: %v", analysis.ID, err)
		return fmt.Errorf("failed to update progress: %v", err)
	}
	if res.RowsAffected == 0 {
		w.Logger.Error("Failed to update progress for analysis %s: %s", analysis.ID, "no rows affected")
		return fmt.Errorf("failed to update progress: no rows affected")
	}

	w.sendWebSocketUpdate(analysis.ID, "processing", progress, "GIS processing completed")
	w.Logger.Info("GIS processing completed for analysis %s", analysis.ID)

	return nil
}

func (w *AnalysisWorker) processVLM(analysisID string) error {
	w.Logger.Info("Starting VLM analysis for analysis %s", analysisID)

	progress := models.AnalysisProgress{
		GISProcessing:    "completed",
		VLMAnalysis:      "in_progress",
		RAGRetrieval:     "pending",
		ReportGeneration: "pending",
	}
	progressJSON, err := json.Marshal(progress)
	if err != nil {
		w.Logger.Error("failed to marshal VLM progressJSON for analysis: %s", analysisID)
		return fmt.Errorf("failed to marshal VLM progressJSON for analysis: %s", analysisID)
	}

	res, err := postgresql.UpdateFields(w.DB.Postgresql, &models.Analysis{}, map[string]any{
		"progress": progressJSON,
	}, "id = ?", analysisID)
	if err != nil {
		w.Logger.Error("failed to update progress for analysis %s: %v", analysisID, err)
		return fmt.Errorf("failed to update progress: %v", err)
	}
	if res.RowsAffected == 0 {
		w.Logger.Error("failed to update progress for analysis %s: %s", analysisID, "no rows affected")
		return fmt.Errorf("failed to update progress: no rows affected")
	}

	w.sendWebSocketUpdate(analysisID, "processing", progress, "Analyzing imagery with VLM...")

	var analysis models.Analysis
	err, _ = postgresql.SelectOneFromDb(w.DB.Postgresql, &analysis, "id = ?", analysisID)
	if err != nil {
		w.Logger.Error("failed to fetch analysis %s: %v", analysisID, err)
		return fmt.Errorf("failed to fetch analysis: %v", err)
	}

	vlmReq := external_models.VLMAnalysisRequest{
		PastImageURL:    analysis.PastImageryPath,
		PresentImageURL: analysis.PresentImageryPath,
		ChangeMapURL:    analysis.ChangeMapPath,
		Prompt:          "You are a military imagery analyst. Analyze these satellite images and identify tactically significant changes. Focus on: new settlements, defensive positions, roads/tracks, obstacles, staging areas. For each change, provide: type, location coordinates, tactical significance.",
	}

	resp, err := w.ExtReq.SendExternalRequest(request.ProcessVLM, vlmReq)
	if err != nil {
		w.Logger.Error("failed to send VLM request for analysis %s: %v", analysis.ID, err)
		return fmt.Errorf("VLM service error: %v", err)
	}

	vlmResp, ok := resp.(external_models.VLMAnalysisResponse)
	if !ok {
		w.Logger.Error("unexpected VLM response type for analysis %s", analysis.ID)
		return fmt.Errorf("unexpected VLM response type")
	}

	vlmJSON, err := json.Marshal(vlmResp)
	if err != nil {
		w.Logger.Error("failed to marshal vlmJSON: %v", err)
		return fmt.Errorf("failed to marshal vlmJSON: %v", err)
	}

	progress.VLMAnalysis = "completed"
	progressJSON, err = json.Marshal(progress)
	if err != nil {
		w.Logger.Error("failed to marshal VLM progress for analysis %s: %v", analysis.ID, err)
		return fmt.Errorf("failed to marshal VLM progress: %v", err)
	}

	res, err = postgresql.UpdateFields(w.DB.Postgresql, &models.Analysis{}, map[string]any{
		"vlm_analysis": vlmJSON,
		"progress":     progressJSON,
	}, "id = ?", analysis.ID)
	if err != nil {
		w.Logger.Error("failed to update VLM analysis and progress for analysis %s: %v", analysis.ID, err)
		return fmt.Errorf("failed to update VLM analysis and progress: %v", err)
	}
	if res.RowsAffected == 0 {
		w.Logger.Error("failed to update VLM analysis and progress for analysis %s: no rows affected", analysis.ID)
		return fmt.Errorf("failed to update VLM analysis and progress: no rows affected")
	}

	w.sendWebSocketUpdate(analysis.ID, "processing", progress, "VLM analysis completed")
	w.Logger.Info("VLM analysis completed for analysis %s", analysis.ID)

	return nil
}

func (w *AnalysisWorker) processLLM(analysisID string) error {
	w.Logger.Info("Starting report generation for analysis %s", analysisID)

	progress := models.AnalysisProgress{
		GISProcessing:    "completed",
		VLMAnalysis:      "completed",
		RAGRetrieval:     "completed",
		ReportGeneration: "in_progress",
	}
	progressJSON, err := json.Marshal(progress)
	if err != nil {
		w.Logger.Error("failed to marshal progress for analysis %s: %v", analysisID, err)
		return fmt.Errorf("failed to marshal progress: %v", err)
	}

	res, err := postgresql.UpdateFields(w.DB.Postgresql, &models.Analysis{}, map[string]any{
		"progress": progressJSON,
	}, "id = ?", analysisID)
	if err != nil {
		w.Logger.Error("failed to update progress for analysis %s: %v", analysisID, err)
		return fmt.Errorf("failed to update progress: %v", err)
	}
	if res.RowsAffected == 0 {
		w.Logger.Error("failed to update progress for analysis %s: no rows affected", analysisID)
		return fmt.Errorf("failed to update progress: no rows affected")
	}

	w.sendWebSocketUpdate(analysisID, "processing", progress, "Generating IPB report...")

	var analysis models.Analysis
	err, _ = postgresql.SelectOneFromDb(w.DB.Postgresql, &analysis, "id = ?", analysisID)
	if err != nil {
		w.Logger.Error("Failed to re-fetch analysis for LLM processing %s: %v", analysisID, err)
		return fmt.Errorf("failed to re-fetch analysis: %v", err)
	}

	var vlmResp external_models.VLMAnalysisResponse
	if err := json.Unmarshal(analysis.VLMAnalysis, &vlmResp); err != nil {
		w.Logger.Error("failed to unmarshal VLM analysis for analysis %s: %v", analysisID, err)
		return fmt.Errorf("failed to unmarshal VLM analysis: %v", err)
	}

	// Mock RAG context (in production, this would come from RAG service)
	ragContext := []string{
		"IED threat along River Kaduna increased significantly in Q3 2025. Common placement locations include: river approach roads, bridge crossing points, and concealed positions 50-100m from water's edge.",
		"Enemy defensive tactics for river crossings typically establish observation posts 500-1000m from water on high ground. Defensive positions use natural terrain for cover.",
		"River Kaduna depth varies 2-4m during monsoon season (current). Current speed 1.5-2 knots. Banks are soft mud, requiring engineer support for vehicle crossing.",
	}

	llmReq := external_models.LLMReportRequest{
		RAGContext:  ragContext,
		VLMAnalysis: vlmResp,
		ReportSections: []string{
			"executive_summary",
			"geospatial_analysis",
			"weather_light_analysis",
			"ocoka",
			"threat_coas",
			"targeting_list",
			"information_collection_plan",
		},
		SystemPrompt: "You are a military intelligence analyst. Generate a comprehensive Intelligence Preparation of the Battlefield (IPB) report based on the provided information. Be specific, tactical, and actionable.",
	}

	llmReq.OperatorQuery.OperationType = analysis.OperationType
	llmReq.OperatorQuery.AreaOfOperation = analysis.AreaOfOperation
	llmReq.OperatorQuery.MissionDescription = analysis.MissionDescription
	llmReq.OperatorQuery.Constraints = analysis.Constraints

	llmResp, err := w.ExtReq.SendExternalRequest(request.ProcessLLM, llmReq)
	if err != nil {
		w.Logger.Error("LLM service error for analysis %s: %v", analysisID, err)
		return fmt.Errorf("LLM service error: %v", err)
	}

	typedLLMResp, ok := llmResp.(external_models.LLMReportResponse)
	if !ok {
		w.Logger.Error("LLM service error: unexpected response type for analysis %s", analysisID)
		return fmt.Errorf("LLM service error: unexpected response type")
	}

	reportJSON, err := json.Marshal(typedLLMResp.Report)
	if err != nil {
		w.Logger.Error("failed to marshal LLM report for analysis %s: %v", analysisID, err)
		return fmt.Errorf("failed to marshal LLM report: %v", err)
	}

	progress.ReportGeneration = "completed"
	progressJSON, err = json.Marshal(progress)
	if err != nil {
		w.Logger.Error("failed to marshal final progress for analysis %s: %v", analysisID, err)
		return fmt.Errorf("failed to marshal final progress: %v", err)
	}

	res, err = postgresql.UpdateFields(w.DB.Postgresql, &models.Analysis{}, map[string]any{
		"report_json": reportJSON,
		"progress":    progressJSON,
	}, "id = ?", analysisID)
	if err != nil {
		w.Logger.Error("failed to update report and progress for analysis %s: %v", analysisID, err)
		return fmt.Errorf("failed to update report and progress: %v", err)
	}
	if res.RowsAffected == 0 {
		w.Logger.Error("failed to update report and progress for analysis %s: no rows affected", analysisID)
		return fmt.Errorf("failed to update report and progress: no rows affected")
	}

	w.sendWebSocketUpdate(analysisID, "processing", progress, "Report generation completed")
	w.Logger.Info("Report generation completed for analysis %s", analysisID)

	return nil
}
