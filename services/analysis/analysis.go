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
	"github.com/Micah-Shallom/geoint-backend/services/websocket"
	"github.com/Micah-Shallom/geoint-backend/utility"
)

func SubmitAnalysis(db *storage.Database, logger *utility.Logger, req models.AnalysisRequest, pastFiles []*multipart.FileHeader, presentFiles []*multipart.FileHeader, supportFiles []*multipart.FileHeader, hub *websocket.Hub) (models.AnalysisResponse, error) {
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
	logger.Info("Past imagery uploaded to %s", pastImageryURL)

	presentImageryURL, err := minio.UploadImagery(logger, analysisID, "present", presentFiles[0])
	if err != nil {
		return models.AnalysisResponse{}, fmt.Errorf("failed to upload present imagery: %v", err)
	}
	logger.Info("Present imagery uploaded to %s", presentImageryURL)

	for _, file := range supportFiles {
		docURL, err := minio.UploadSupportDocument(logger, analysisID, file)
		if err != nil {
			logger.Warning("Failed to upload support document %s: %v", file.Filename, err)
			continue
		}
		supportDocURLs = append(supportDocURLs, docURL)
	}
	logger.Info("Support documents uploaded: %d", len(supportDocURLs))

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
		logger.Error("Failed to create analysis record: %v", err)
		return models.AnalysisResponse{}, fmt.Errorf("failed to create analysis record: %v", err)
	}
	logger.Info("Analysis record created successfully with ID %s", analysisID)

	go ProcessAnalysisBackground(db, logger, analysisID, hub)

	response := models.AnalysisResponse{
		AnalysisID:    analysisID,
		Status:        "processing",
		EstimatedTime: "10-15 minutes",
		CreatedAt:     analysis.CreatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func ProcessAnalysisBackground(db *storage.Database, logger *utility.Logger, analysisID string, hub *websocket.Hub) {
	logger.Info("Starting background processing for analysis %s", analysisID)

	worker := NewAnalysisWorker(db, logger, hub)
	worker.ProcessAnalysis(analysisID)

	logger.Info("Background processing completed for analysis %s", analysisID)
}
