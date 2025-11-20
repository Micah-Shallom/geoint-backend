package llm

import (
	"fmt"

	"github.com/Micah-Shallom/geoint-backend/external/external_models"
)

func (r *RequestObj) LLMCall() (external_models.LLMReportResponse, error) {
	var (
		outBoundResponse external_models.LLMReportResponse
		logger           = r.Logger
		idata            = r.RequestData
		path             = ""
	)

	data, ok := idata.(external_models.LLMReportRequest)
	if !ok {
		logger.Error("gis ", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	err := r.getNewSendRequestObject(data, map[string]string{}, path).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("gis resolve", outBoundResponse, err.Error())
		return outBoundResponse, err
	}

	return outBoundResponse, nil
}
