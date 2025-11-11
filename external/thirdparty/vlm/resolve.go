package vlm

import (
	"fmt"

	"github.com/Micah-Shallom/geoint-backend/external/external_models"
)

func (r *RequestObj) VLMCall() (external_models.VLMAnalysisResponse, error) {
	var (
		outBoundResponse external_models.VLMAnalysisResponse
		logger           = r.Logger
		idata            = r.RequestData
		path             = ""
	)

	data, ok := idata.(external_models.VLMAnalysisRequest)
	if !ok {
		logger.Error("vlm ", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	err := r.getNewSendRequestObject(data, map[string]string{}, path).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("vlm resolve", outBoundResponse, err.Error())
		return outBoundResponse, err
	}

	return outBoundResponse, nil
}
