package ipstack_mocks

import (
	"fmt"

	"github.com/Micah-Shallom/geoint-backend/external/external_models"
	"github.com/Micah-Shallom/geoint-backend/internal/config"
	"github.com/Micah-Shallom/geoint-backend/utility"
)

func IpinfoResolveIp(logger *utility.Logger, idata any) (external_models.IPInfoResponse, error) {
	var (
		key              = config.GetConfig().IPStack.Key
		outBoundResponse external_models.IPInfoResponse
	)

	ip, ok := idata.(string)
	if !ok {
		logger.Error("ipinfo resolve ip", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	// Simulating a response from ipinfo.io
	outBoundResponse.IP = ip
	outBoundResponse.City = "Sample City"
	outBoundResponse.Region = "Sample Region"
	outBoundResponse.Country = "Sample Country"
	outBoundResponse.Location = "40.7128,-74.0060"
	outBoundResponse.Org = "Example ISP"
	outBoundResponse.TimeZone = "AS12345"

	path := "/" + ip + "?token=" + key
	logger.Info("ipinfo resolve ip", ip, path)

	return outBoundResponse, nil
}
