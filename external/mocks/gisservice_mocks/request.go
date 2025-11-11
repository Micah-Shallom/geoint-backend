package gisservicemocks

import (
	"time"

	"github.com/Micah-Shallom/geoint-backend/external/external_models"
	"github.com/Micah-Shallom/geoint-backend/utility"
)

func MockGISChangeDetection(logger *utility.Logger, data any) (external_models.GISChangeDetectionResponse, error) {
	req := data.(external_models.GISChangeDetectionRequest)
	
	logger.Info("Mock GIS Service - Processing change detection of pastimage: %s and present image: %s", req.PastImageURL, req.PresentImageURL)

	// Simulate processing delay
	time.Sleep(5 * time.Second)

	response := external_models.GISChangeDetectionResponse{
		ChangeMapURL: "http://localhost:9000/geointbucket/changemaps/mock_change_map.tif",
	}

	response.Metadata.TotalChangesDetected = 47
	response.Metadata.SignificantChanges = 12
	response.Metadata.ChangeAreas = []external_models.ChangeArea{
		{
			Type: "new_structure",
			Bbox: []float64{450120, 1250340, 450180, 1250390},
		},
		{
			Type: "vegetation_removal",
			Bbox: []float64{450500, 1250100, 450650, 1250200},
		},
		{
			Type: "road_construction",
			Bbox: []float64{450300, 1250500, 450450, 1250550},
		},
	}

	logger.Info("Mock GIS Service - Change detection completed. Total changes: %d", response.Metadata.TotalChangesDetected)
	return response, nil
}
