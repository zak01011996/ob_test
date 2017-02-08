package btcfeed

import (
	"fmt"

	"github.com/zak01011996/bitcoin-ticker/service"
)

// bitcoincharts.com API response struct
type bcResp map[string]struct {
	Rate float64 `json:"last"`
}

// Bitcoin feed service
// Service: blockchain.info
// Current API url: https://blockchain.info/ticker
type BCFeed struct {
	service.BaseService
}

// Get rate in given currencies
func (s *BCFeed) GetRate(curr string) (res float64, err error) {
	var response bcResp

	// Compile url
	url := s.GetURL() + "ticker"

	// Get data from service
	err = s.GetJSON(url, &response)
	if err != nil {
		return
	}

	// Try to get rate
	if v, ok := response[curr]; ok {
		res = v.Rate
	} else {
		err = fmt.Errorf(service.EmptyResp, s.GetName())
	}

	return
}

// Creates new instance of service
func NewBCFeed(srvUrl string) *BCFeed {
	return &BCFeed{service.BaseService{
		ServiceUrl:  srvUrl,
		ServiceName: "Blockchain.info",
	}}
}
