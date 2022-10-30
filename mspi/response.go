package mspi

import (
	"github.com/com-gft-tsbo-source/go-common/ms-framework/microservice"
)

// ###########################################################################
// ###########################################################################
// MsPi Response - Device
// ###########################################################################
// ###########################################################################

// PiResponse ...
type PiResponse struct {
	microservice.Response
	Value      string `json:"value"`
	Iterations int    `json:"iterations"`
	Precision  int    `json:"precision"`
}

// ###########################################################################

// InitPiResponse Constructor of a response of ms-pi
func InitPiResponse(r *PiResponse, code int, status string, ms *MsPi) {
	microservice.InitResponseFromMicroService(&r.Response, ms, code, status)
	r.Value = "n/a"
	r.Iterations = 0
	r.Precision = 0
}

// NewPiResponse ...
func NewPiResponse(code int, status string, ms *MsPi) *PiResponse {
	var r PiResponse
	InitPiResponse(&r, code, status, ms)
	return &r
}
