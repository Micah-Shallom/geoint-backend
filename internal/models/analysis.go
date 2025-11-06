package models

import (
	"mime/multipart"
	"time"
)

type Analysis struct {
    ID                 string   `gorm:"primaryKey"`
    OperationType      string
    AreaOfOperation    string
    MissionDescription string
    Constraints        string
    PastImagery        
    PresentImagery      
    SupportDocuments   []*multipart.FileHeader // paths
    Status             string   
    CreatedAt          time.Time
    CompletedAt        *time.Time
    ReportJSON         string   // final report
}


type AnalysisRequest struct {
	OperationType      string                  `form:"operation_type" json:"operation_type" binding:"required"`
	AreaOfOperation    string                  `form:"area_of_operation" json:"area_of_operation" binding:"required"`
	MissionDescription string                  `form:"mission_description" json:"mission_description" binding:"required"`
	Constraints        string                  `form:"constraints" json:"constraints"`
	PastImagery        *multipart.FileHeader   `form:"past_imagery" json:"past_imagery" binding:"required"`
	PresentImagery     *multipart.FileHeader   `form:"present_imagery" json:"present_imagery" binding:"required"`
	SupportDocuments   []*multipart.FileHeader `form:"support_documents" json:"support_documents"`
}

type AnalysisResponse struct {
	AnalysisID    string    `json:"analysis_id"`
	Status        string    `json:"status"`
	EstimatedTime time.Time `json:"estimated_time"`
	CreatedAt     time.Time `json:"created_at"`
}
