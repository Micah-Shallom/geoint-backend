package external_models

// GIS Service Models
type GISChangeDetectionRequest struct {
	PastImageURL    string `json:"past_image"`
	PresentImageURL string `json:"present_image"`
}

type ChangeArea struct {
	Type string    `json:"type"`
	Bbox []float64 `json:"bbox"`
}

type ChangeDetectionMetadata struct {
	ChangeAreas            []ChangeArea `json:"change_areas"`
	SignificantChanges     int          `json:"significant_changes"`
	TotalChangesDetected   int          `json:"total_changes_detected"`
	ChangeDetectionVersion string       `json:"change_detection_version"`
}

type GISChangeDetectionResponse struct {
	ChangeMapURL string                  `json:"change_map"`
	Metadata     ChangeDetectionMetadata `json:"metadata"`
}

// VLM Service Models
type VLMAnalysisRequest struct {
	PastImageURL    string `json:"past_image"`
	PresentImageURL string `json:"present_image"`
	ChangeMapURL    string `json:"change_map"`
	Prompt          string `json:"prompt"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type VLMChange struct {
	ID                   string     `json:"id"`
	Type                 string     `json:"type"`
	Location             *Location  `json:"location,omitempty"`
	Coordinates          []Location `json:"coordinates,omitempty"`
	Description          string     `json:"description"`
	TacticalSignificance string     `json:"tactical_significance"`
	Confidence           float64    `json:"confidence"`
}

type VLMAnalysisResponse struct {
	Changes []VLMChange `json:"changes"`
	Summary string      `json:"summary"`
}

// LLM Service Models
type LLMReportRequest struct {
	OperatorQuery struct {
		OperationType      string `json:"operation_type"`
		AreaOfOperation    string `json:"area_of_operation"`
		MissionDescription string `json:"mission_description"`
		Constraints        string `json:"constraints"`
	} `json:"operator_query"`
	RAGContext     []string            `json:"rag_context"`
	VLMAnalysis    VLMAnalysisResponse `json:"vlm_analysis"`
	ReportSections []string            `json:"report_sections"`
	SystemPrompt   string              `json:"system_prompt"`
}

type TargetingItem struct {
	TargetNumber     string `json:"target_number"`
	Description      string `json:"description"`
	Location         string `json:"location"`
	Priority         string `json:"priority"`
	EngagementMethod string `json:"engagement_method"`
	Justification    string `json:"justification"`
}

type OCOKAAnalysis struct {
	Observation       string `json:"observation"`
	CoverConcealment  string `json:"cover_concealment"`
	Obstacles         string `json:"obstacles"`
	KeyTerrain        string `json:"key_terrain"`
	AvenuesOfApproach string `json:"avenues_of_approach"`
}

type ThreatCOAs struct {
	MLCOA string `json:"mlcoa"`
	MDCOA string `json:"mdcoa"`
}

type IPBReport struct {
	ExecutiveSummary          string          `json:"executive_summary"`
	GeospatialAnalysis        string          `json:"geospatial_analysis"`
	WeatherLightAnalysis      string          `json:"weather_light_analysis"`
	OCOKA                     OCOKAAnalysis   `json:"ocoka"`
	ThreatCOAs                ThreatCOAs      `json:"threat_coas"`
	TargetingList             []TargetingItem `json:"targeting_list"`
	InformationCollectionPlan string          `json:"information_collection_plan"`
}

type LLMReportResponse struct {
	Report   IPBReport `json:"report"`
	Metadata struct {
		Model           string `json:"model"`
		TokensGenerated int    `json:"tokens_generated"`
		GenerationTime  string `json:"generation_time"`
	} `json:"metadata"`
}
