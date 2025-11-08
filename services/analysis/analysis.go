package analysis

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/Micah-Shallom/geoint-backend/internal/models"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage/minio"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage/postgresql"
	"github.com/Micah-Shallom/geoint-backend/utility"
)

func SubmitAnalysis(db *storage.Database, logger *utility.Logger, req models.AnalysisRequest, pastFiles []*multipart.FileHeader, presentFiles []*multipart.FileHeader, supportFiles []*multipart.FileHeader) (models.AnalysisResponse, error) {
	var (
		analysis       models.Analysis
		supportDocURLs []string
	)

	analysisID := utility.GenerateUUID()

	if err := utility.ValidateImageryFile(pastFiles[0]); err != nil {
		return models.AnalysisResponse{}, err
	}

	if err := utility.ValidateImageryFile(presentFiles[0]); err != nil {
		return models.AnalysisResponse{}, err
	}

	for _, file := range supportFiles {
		if err := utility.ValidateDocumentFile(file); err != nil {
			return models.AnalysisResponse{}, err
		}
	}

	pastImageryURL, err := minio.UploadImagery(logger, analysisID, "past", pastFiles[0])
	if err != nil {
		return models.AnalysisResponse{}, fmt.Errorf("failed to upload past imagery: %v", err)
	}
	logger.Info("Past imagery uploaded", pastImageryURL)

	presentImageryURL, err := minio.UploadImagery(logger, analysisID, "present", presentFiles[0])
	if err != nil {
		return models.AnalysisResponse{}, fmt.Errorf("failed to upload present imagery: %v", err)
	}
	logger.Info("Present imagery uploaded", presentImageryURL)

	for _, file := range supportFiles {
		docURL, err := minio.UploadSupportDocument(logger, analysisID, file)
		if err != nil {
			logger.Warning("Failed to upload support document", file.Filename, err)
			continue
		}
		supportDocURLs = append(supportDocURLs, docURL)
	}
	logger.Info("Support documents uploaded", len(supportDocURLs))

	supportDocsJSON, err := json.Marshal(supportDocURLs)
	if err != nil {
		return models.AnalysisResponse{}, fmt.Errorf("failed to marshal support documents: %v", err)
	}

	progress := models.AnalysisProgress{
		GISProcessing:    "pending",
		VLMAnalysis:      "pending",
		RAGRetrieval:     "pending",
		ReportGeneration: "pending",
	}

	progressJSON, err := json.Marshal(progress)
	if err != nil {
		return models.AnalysisResponse{}, fmt.Errorf("failed to marshal progress: %v", err)
	}

	analysis = models.Analysis{
		ID:                   analysisID,
		OperationType:        req.OperationType,
		AreaOfOperation:      req.AreaOfOperation,
		MissionDescription:   req.MissionDescription,
		Constraints:          req.Constraints,
		PastImageryPath:      pastImageryURL,
		PresentImageryPath:   presentImageryURL,
		SupportDocumentPaths: supportDocsJSON,
		Status:               "processing",
		Progress:             progressJSON,
		CreatedAt:            time.Now(),
	}

	if err := postgresql.CreateOneRecord(db.Postgresql, &analysis); err != nil {
		logger.Error("Failed to create analysis record", err)
		return models.AnalysisResponse{}, fmt.Errorf("failed to create analysis record: %v", err)
	}
	logger.Info("Analysis record created successfully", analysisID)

	go ProcessAnalysisBackground(db, logger, analysisID)

	response := models.AnalysisResponse{
		AnalysisID:    analysisID,
		Status:        "processing",
		EstimatedTime: "10-15 minutes",
		CreatedAt:     analysis.CreatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func ProcessAnalysisBackground(db *storage.Database, logger *utility.Logger, analysisID string) {
	logger.Info("Starting background processing for analysis", analysisID)

	// This is a placeholder for the actual processing pipeline
	// In future implementations, this will:
	// 1. Call GIS Service for change detection
	// 2. Call VLM Service for imagery analysis
	// 3. Call RAG Service for document retrieval
	// 4. Call LLM Service for report generation

	logger.Info("Background processing initiated", analysisID)

	// Update status to indicate processing has started
	updateProgress := models.AnalysisProgress{
		GISProcessing:    "in_progress",
		VLMAnalysis:      "pending",
		RAGRetrieval:     "pending",
		ReportGeneration: "pending",
	}

	progressJSON, _ := json.Marshal(updateProgress)

	db.Postgresql.Model(&models.Analysis{}).
		Where("id = ?", analysisID).
		Update("progress", progressJSON)

	logger.Info("Progress updated for analysis", analysisID)
}
