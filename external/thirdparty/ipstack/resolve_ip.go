package ipstack

import (
	"fmt"

	"github.com/Micah-Shallom/geoint-backend/external/external_models"
	"github.com/Micah-Shallom/geoint-backend/internal/config"
)

func (r *RequestObj) IpinfoResolveIp() (external_models.IPInfoResponse, error) {
	var (
		key              = config.GetConfig().IPStack.Key
		outBoundResponse external_models.IPInfoResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	ip, ok := idata.(string)
	if !ok {
		logger.Error("ipinfo resolve ip", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	path := "/" + ip + "?token=" + key

	logger.Info("ipinfo resolve ip", ip)
	err := r.getNewSendRequestObject(nil, map[string]string{}, path).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("ipinfo resolve ip", outBoundResponse, err.Error())
		return outBoundResponse, err
	}

	return outBoundResponse, nil
}
