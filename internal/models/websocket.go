package models

type Message struct {
	AnalysisID  string `json:"analysis_id"`
	Status      string `json:"status"`
	Progress    any    `json:"progress"`
	CurrentStep string `json:"current_step"`
	UpdatedAt   string `json:"updated_at"`
}
