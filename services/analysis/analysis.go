package analysis

import (
	"mime/multipart"

	"github.com/Micah-Shallom/geoint-backend/internal/models"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage"
	"github.com/Micah-Shallom/geoint-backend/utility"
)

func SubmitAnalysis(db *storage.Database, logger *utility.Logger, req models.AnalysisRequest, pastFiles []*multipart.FileHeader, presentFiles []*multipart.FileHeader, supportFiles []*multipart.FileHeader) (models.AnalysisResponse, error) {
	
	return models.AnalysisResponse{}, nil
}
