package btcfeed

import (
	"fmt"
	"time"

	"github.com/zak01011996/bitcoin-ticker/service"
)

// coindesk.com API response struct
type cdResp struct {
	Time struct {
		updated time.Time `json:"updatedISO"`
	} `json:"time"`
	Disc string `json:"disclaimer"`
	Bpi  map[string]struct {
		CurrCode string  `json:"code"`
		Rate     float64 `json:"rate_float"`
	} `json:"bpi"`
}

// Bitcoin feed service
// Service: coindesk.com
// Current API url: http://api.coindesk.com/v1/bpi/currentprice.json
type CDFeed struct {
	service.BaseService
}

// Get rate in given currencies
func (s *CDFeed) GetRate(curr string) (res float64, err error) {
	var response cdResp

	// Compile url
	url := s.GetURL() + "v1/bpi/currentprice.json"

	// Get data from service
	err = s.GetJSON(url, &response)
	if err != nil {
		return
	}

	// Try to get rate
	if v, ok := response.Bpi[curr]; ok {
		res = v.Rate
	} else {
		err = fmt.Errorf(service.EmptyResp, s.GetName())
	}

	return
}

// Creates new instance of service
func NewCDFeed(srvUrl string) *CDFeed {
	return &CDFeed{service.BaseService{
		ServiceUrl:  srvUrl,
		ServiceName: "Coindesk.com",
	}}
}
