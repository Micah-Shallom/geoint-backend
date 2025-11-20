package vlmservicemocks

import (
	"time"

	"github.com/Micah-Shallom/geoint-backend/external/external_models"
	"github.com/Micah-Shallom/geoint-backend/utility"
)

func MockVLMAnalysis(logger *utility.Logger, data any) (external_models.VLMAnalysisResponse, error) {
	req := data.(external_models.VLMAnalysisRequest)
	logger.Info("Mock VLM Service - Analyzing imagery", req.PastImageURL, req.PresentImageURL)

	// Simulate processing delay
	time.Sleep(5 * time.Second)

	response := external_models.VLMAnalysisResponse{
		Changes: []external_models.VLMChange{
			{
				ID:   "change-1",
				Type: "new_settlement",
				Location: &external_models.Location{
					Lat: 9.0821,
					Lon: 7.5341,
				},
				Description:          "New settlement cluster of approximately 15-20 structures located 800m north of river. Appears to be recent construction based on disturbed earth patterns.",
				TacticalSignificance: "HIGH - Could serve as staging area or observation post overlooking crossing site",
				Confidence:           0.89,
			},
			{
				ID:   "change-2",
				Type: "new_road",
				Coordinates: []external_models.Location{
					{Lat: 9.0805, Lon: 7.5320},
					{Lat: 9.0818, Lon: 7.5355},
				},
				Description:          "New unpaved track connecting settlement to main road. Approximately 2.5km length, 4-5m width.",
				TacticalSignificance: "MEDIUM - Provides supply route to new settlement",
				Confidence:           0.92,
			},
			{
				ID:   "change-3",
				Type: "defensive_position",
				Location: &external_models.Location{
					Lat: 9.0798,
					Lon: 7.5312,
				},
				Description:          "Possible fighting position or bunker. Cleared area with bermed perimeter, approximately 15m diameter.",
				TacticalSignificance: "HIGH - Overlooks river approach from southwest",
				Confidence:           0.76,
			},
		},
		Summary: "Detected 3 tactically significant changes indicating possible defensive preparation along river approach.",
	}

	logger.Info("Mock VLM Service - Analysis completed: %d", len(response.Changes))
	return response, nil
}
