package llmlservicemocks

import (
	"time"

	"github.com/Micah-Shallom/geoint-backend/external/external_models"
	"github.com/Micah-Shallom/geoint-backend/utility"
)

func MockLLMReportGeneration(logger *utility.Logger, data any) (external_models.LLMReportResponse, error) {
	req := data.(external_models.LLMReportRequest)
	logger.Info("Mock LLM Service - Generating IPB report", req.OperatorQuery.OperationType)

	// Simulate processing delay
	time.Sleep(5 * time.Second)

	response := external_models.LLMReportResponse{}

	response.Report = external_models.IPBReport{
		ExecutiveSummary: "Analysis of River Kaduna crossing sites reveals significant enemy defensive preparation in the target area. Recent imagery analysis identified new settlement construction, defensive positions, and supply routes indicating enemy anticipation of operations in this area. Combined with elevated IED threat and monsoon conditions, deliberate assault operations will require extensive engineer support and route clearance.",

		GeospatialAnalysis: "Target Area: Grid 32N 450000 1250000, River Kaduna crossing zone. Terrain consists of soft alluvial plains with 2-4m river depth during current monsoon season. VLM analysis identified three primary tactical changes: (1) New settlement at coordinates 9.0821°N, 7.5341°E with approximately 15-20 structures, (2) New unpaved track connecting settlement to main road, approximately 2.5km length, (3) Possible defensive position at 9.0798°N, 7.5312°E with cleared perimeter.",

		WeatherLightAnalysis: "Current Period: Monsoon season, heavy rainfall expected. River levels elevated 1.5-2m above dry season baseline. Visibility: Reduced during rainfall. BMNT: 0615 hrs, EENT: 1845 hrs. Moonrise: 2245 hrs (67% illumination). Recommendation: Night operations offer concealment but river crossing more hazardous due to elevated water levels.",

		OCOKA: external_models.OCOKAAnalysis{
			Observation: "Enemy observation posts likely located at new settlement (coordinates 9.0821°N, 7.5341°E) providing overwatch of primary crossing site. Elevated terrain 800m north of river offers 270-degree observation arc covering approach routes and crossing points.",

			CoverConcealment: "Limited cover along river approaches. Vegetation cleared in areas consistent with defensive preparation. Attacking force will require smoke obscuration during approach. Natural depressions offer limited concealment for staging forces.",

			Obstacles: "Primary obstacle: River Kaduna (2-4m depth, 1.5-2kt current). Soft mud banks require engineer bridging. IED threat documented along approach roads per intelligence reporting. Suspected defensive position at coordinates 9.0798°N, 7.5312°E creates engagement area covering likely crossing points.",

			KeyTerrain: "Primary crossing site: Grid 32N 450200 1250300 (historical site, hardest banks). High ground at coordinates 9.0821°N, 7.5341°E (new settlement location) dominates crossing area - must be secured or suppressed prior to crossing operations.",

			AvenuesOfApproach: "Ground Avenue: Axis COBRA from south along Route 4, vulnerable to IED threat and observation from high ground. Air Avenue: LZ HAWK at Grid 32N 449800 1250100 permits helicopter insertion 1.2km from objective, outside direct fire range.",
		},

		ThreatCOAs: external_models.ThreatCOAs{
			MLCOA: "Most Likely COA: Enemy conducts defense from prepared positions utilizing new settlement as observation post and command node. Defensive belt established 100-500m from river with IEDs along approach routes. Upon contact, enemy defends from prepared positions with crew-served weapons, targeting engineer assets during crossing. Reserves positioned at settlement conduct counterattack against bridgehead.",

			MDCOA: "Most Dangerous COA: Enemy has prepared complex ambush incorporating IEDs, indirect fire, and mobile reserves. Allows friendly force to initiate crossing before engaging with massed fires. Simultaneously conducts infiltration of special purpose forces to attack command posts and logistics sites in rear areas.",
		},

		TargetingList: []external_models.TargetingItem{
			{
				TargetNumber:     "T001",
				Description:      "Settlement/Observation Post",
				Location:         "9.0821°N, 7.5341°E",
				Priority:         "HIGH",
				EngagementMethod: "Artillery/CAS",
				Justification:    "Command and observation node identified by VLM analysis, dominates crossing site",
			},
			{
				TargetNumber:     "T002",
				Description:      "Suspected Fighting Position",
				Location:         "9.0798°N, 7.5312°E",
				Priority:         "HIGH",
				EngagementMethod: "Direct Fire/Artillery",
				Justification:    "Defensive position identified by VLM analysis, overlooks river approach",
			},
			{
				TargetNumber:     "T003",
				Description:      "Supply Route",
				Location:         "Between 9.0805°N, 7.5320°E and 9.0818°N, 7.5355°E",
				Priority:         "MEDIUM",
				EngagementMethod: "Interdiction/Mining",
				Justification:    "New track identified by VLM analysis connecting settlement to main road, enemy resupply route",
			},
		},

		InformationCollectionPlan: "Priority Intelligence Requirements (PIR): 1) Confirm enemy strength and disposition at settlement (coordinates 9.0821°N, 7.5341°E). 2) Identify IED locations along Route 4 approach. 3) Determine river current speed and exact depth at primary crossing site. Collection Assets: ISR platform to conduct pattern-of-life analysis on settlement (24-48hr prior). Engineer recon of crossing site (night prior). Route clearance teams to confirm IED threat (H-6 to H-2). Named Areas of Interest: NAI 1 - Settlement area, NAI 2 - Suspected defensive positions, NAI 3 - Crossing site banks.",
	}

	response.Metadata.Model = "llama-3-70b-instruct"
	response.Metadata.TokensGenerated = 3847
	response.Metadata.GenerationTime = "127 seconds"

	logger.Info("Mock LLM Service - Report generation completed")
	return response, nil
}
