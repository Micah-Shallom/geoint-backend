package gisservicemocks

import (
	"github.com/Micah-Shallom/geoint-backend/external/external_models"
	"github.com/Micah-Shallom/geoint-backend/utility"
)

func GISMockRequest(logger *utility.Logger, data any) (external_models.GISChangeDetectionResponse, error) {
	response := external_models.GISChangeDetectionResponse{
		ChangeMapURL: "https://example.com/changedetection",
		Metadata: external_models.ChangeDetectionMetadata{
			ChangeAreas: []external_models.ChangeArea{
				{
					Type: "Polygon",
					Bbox: []float64{1, 1, 1, 1},
				},
			},
			SignificantChanges:     0,
			TotalChangesDetected:   0,
			ChangeDetectionVersion: "1.0.0",
		},
	}

	return response, nil
}
