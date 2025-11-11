package vlm

import (
	"github.com/Micah-Shallom/geoint-backend/external"
	"github.com/Micah-Shallom/geoint-backend/utility"
)

type RequestObj struct {
	Name         string
	Path         string
	Method       string
	SuccessCode  int
	RequestData  any
	DecodeMethod string
	Logger       *utility.Logger
	Timeout      bool
}

var (
	JsonDecodeMethod    string = "json"
	PhpSerializerMethod string = "phpserializer"
)

func (r *RequestObj) getNewSendRequestObject(data any, headers map[string]string, urlprefix string) *external.SendRequestObject {
	return external.GetNewSendRequestObject(r.Logger, r.Name, r.Path, r.Method, urlprefix, r.DecodeMethod, headers, r.SuccessCode, data, r.Timeout)
}
