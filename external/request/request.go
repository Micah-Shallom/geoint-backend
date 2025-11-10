package request

import (
	"fmt"
	"net/http"

	"github.com/Micah-Shallom/geoint-backend/external/mocks"
	gisservice "github.com/Micah-Shallom/geoint-backend/external/thirdparty/gis_service"
	"github.com/Micah-Shallom/geoint-backend/external/thirdparty/ipstack"
	"github.com/Micah-Shallom/geoint-backend/internal/config"
	"github.com/Micah-Shallom/geoint-backend/utility"
)

type ExternalRequest struct {
	Logger *utility.Logger
	Test   bool
}

var (
	JsonDecodeMethod string = "json"
	IpinfoResolveIp  string = "ipinfo_resolve_ip"
	ProcessGIS       string = "process_gis"
)

func (er ExternalRequest) SendExternalRequest(name string, data any) (any, error) {
	var (
		config = config.GetConfig()
	)
	if !er.Test {
		switch name {
		case IpinfoResolveIp:
			obj := ipstack.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v", config.IPStack.BaseUrl),
				Method:       http.MethodGet,
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.IpinfoResolveIp()
		case ProcessGIS:
			obj := gisservice.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v", config.GIS.BaseURL),
				Method:       http.MethodPost,
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.GISRequest()
		default:
			return nil, fmt.Errorf("request not found")
		}
	} else {
		mer := mocks.ExternalRequest{Logger: er.Logger, Test: true}
		return mer.SendExternalRequest(name, data)
	}
}
