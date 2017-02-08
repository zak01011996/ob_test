package currfeed

import (
	"fmt"

	"github.com/zak01011996/bitcoin-ticker/service"
)

// mycurrency.net API response struct
type mcnResp struct {
	CurrCode string  `json:"currency_code"`
	Rate     float64 `json:"rate"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
}

// Currency feed service
// Service: mycurrency.net
// Current API url: http://www.mycurrency.net/service/rates
type MCNFeed struct {
	service.BaseService
}

// Get rate in given currencies
func (s *MCNFeed) GetRate(fromCurr string, toCurr string) (res float64, err error) {
	var response []mcnResp

	// Compile url
	url := s.GetURL() + "service/rates"

	// Try to get data
	err = s.GetJSON(url, &response)
	if err != nil {
		return
	}

	// Check response
	if len(response) <= 0 {
		err = fmt.Errorf(service.EmptyResp, s.GetName())
		return
	}

	var fromRate, toRate float64

	// Find and calculate rates
	for _, v := range response {
		if v.CurrCode == fromCurr {
			fromRate = v.Rate
		}

		if v.CurrCode == toCurr {
			toRate = v.Rate
		}
	}

	// Check if zero, who knows what could happen here, but it shouldn't crash our process
	if toRate == 0 {
		err = fmt.Errorf(service.ZeroValue, s.GetName())
	}

	res = toRate / fromRate
	return
}

// Creates new instance of service
func NewMCNFeed(srvUrl string) *MCNFeed {
	return &MCNFeed{service.BaseService{
		ServiceUrl:  srvUrl,
		ServiceName: "MyCurrency.net",
	}}
}
