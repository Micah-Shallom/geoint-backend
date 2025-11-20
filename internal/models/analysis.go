package models

import (
	"mime/multipart"
	"time"

	"gorm.io/datatypes"
)

type Analysis struct {
	ID                   string         `gorm:"primaryKey;type:uuid" json:"id"`
	OperationType        string         `gorm:"type:varchar(255);not null" json:"operation_type"`
	AreaOfOperation      string         `gorm:"type:varchar(500);not null" json:"area_of_operation"`
	MissionDescription   string         `gorm:"type:text;not null" json:"mission_description"`
	Constraints          string         `gorm:"type:text" json:"constraints"`
	PastImageryPath      string         `gorm:"type:varchar(500)" json:"past_imagery_path"`
	PresentImageryPath   string         `gorm:"type:varchar(500)" json:"present_imagery_path"`
	ChangeMapPath        string         `gorm:"type:varchar(500)" json:"change_map_path"`
	SupportDocumentPaths datatypes.JSON `gorm:"type:jsonb" json:"support_document_paths"`
	Status               string         `gorm:"type:varchar(50);default:'pending'" json:"status"`
	Progress             datatypes.JSON `gorm:"type:jsonb" json:"progress"`
	VLMAnalysis          datatypes.JSON `gorm:"type:jsonb" json:"vlm_analysis"`
	RAGContext           datatypes.JSON `gorm:"type:jsonb" json:"rag_context"`
	ReportJSON           datatypes.JSON `gorm:"type:jsonb" json:"report_json"`
	ProcessingTime       string         `gorm:"type:varchar(50)" json:"processing_time"`
	CreatedAt            time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	CompletedAt          *time.Time     `json:"completed_at"`
}

type AnalysisProgress struct {
	GISProcessing    string `json:"gis_processing"`
	VLMAnalysis      string `json:"vlm_analysis"`
	RAGRetrieval     string `json:"rag_retrieval"`
	ReportGeneration string `json:"report_generation"`
}

type AnalysisResponse struct {
	AnalysisID    string `json:"analysis_id"`
	Status        string `json:"status"`
	EstimatedTime string `json:"estimated_time"`
	CreatedAt     string `json:"created_at"`
}

type AnalysisRequest struct {
	OperationType      string                  `json:"operation_type" validate:"required"`
	AreaOfOperation    string                  `json:"area_of_operation" validate:"required"`
	MissionDescription string                  `json:"mission_description" validate:"required"`
	Constraints        string                  `json:"constraints" validate:"required"`
	PastImagery        []*multipart.FileHeader `json:"past_imagery" validate:"required"`
	PresentImagery     []*multipart.FileHeader `json:"present_imagery" validate:"required"`
	SupportDocuments   []*multipart.FileHeader `json:"support_documents" validate:"required"`
}
