package mocks

import (
	"fmt"

	gisservicemocks "github.com/Micah-Shallom/geoint-backend/external/mocks/gisservice_mocks"
	"github.com/Micah-Shallom/geoint-backend/external/mocks/ipstack_mocks"
	"github.com/Micah-Shallom/geoint-backend/utility"
)

type ExternalRequest struct {
	Logger     *utility.Logger
	Test       bool
	RequestObj RequestObj
}

type RequestObj struct {
	Name         string
	Path         string
	Method       string
	Headers      map[string]string
	SuccessCode  int
	RequestData  any
	DecodeMethod string
	Logger       *utility.Logger
}

var (
	JsonDecodeMethod    string = "json"
	PhpSerializerMethod string = "phpserializer"
)

func (er ExternalRequest) SendExternalRequest(name string, data any) (any, error) {
	switch name {
	case "ipinfo_resolve_ip":
		return ipstack_mocks.IpinfoResolveIp(er.Logger, data)
	case "process_gis":
		return gisservicemocks.GISMockRequest(er.Logger, data)
	default:
		return nil, fmt.Errorf("request not found")
	}
}
