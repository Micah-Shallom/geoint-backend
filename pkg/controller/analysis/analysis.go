package analysis

import (
	"net/http"

	"github.com/Micah-Shallom/geoint-backend/external/request"
	"github.com/Micah-Shallom/geoint-backend/internal/models"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage"
	"github.com/Micah-Shallom/geoint-backend/services/analysis"
	"github.com/Micah-Shallom/geoint-backend/utility"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) SubmitAnalysis(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(50 << 20); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	req := models.AnalysisRequest{
		OperationType:      c.PostForm("operation_type"),
		AreaOfOperation:    c.PostForm("area_of_operation"),
		MissionDescription: c.PostForm("mission_description"),
		Constraints:        c.PostForm("constraints"),
	}

	pastFiles := c.Request.MultipartForm.File["past_imagery"]
	presentFiles := c.Request.MultipartForm.File["present_imagery"]
	supportFiles := c.Request.MultipartForm.File["support_documents"]

	if len(pastFiles) == 0 || len(presentFiles) == 0 {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Past imagery and present imagery are required", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if req.OperationType == "" || req.AreaOfOperation == "" || req.MissionDescription == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Required fields are missing", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validator.Struct(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	response, err := analysis.SubmitAnalysis(base.Db, base.Logger, req, pastFiles, presentFiles, supportFiles)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to submit analysis", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Analysis submitted successfully", response)
	c.JSON(http.StatusOK, rd)
}
