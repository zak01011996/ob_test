package currfeed

import (
	"fmt"
	"strconv"

	"github.com/zak01011996/bitcoin-ticker/service"
)

// herokuapp API response struct
type herokuResp struct {
	To   string `json:"To"`
	From string `json:"From"`
	Rate string `json:"Rate"` // YEP, rate here comes in string GJ
}

// Currency feed service
// Service: herokuapp.com
// Current API url: http://rate-exchange.herokuapp.com/fetchRate?from=<fromCurrency>&to=<toCurrency>
type HerokuFeed struct {
	service.BaseService
}

// Get rate in given currencies
func (s *HerokuFeed) GetRate(fromCurr string, toCurr string) (res float64, err error) {
	var response herokuResp

	// Compile url
	url := s.GetURL() + fmt.Sprintf("fetchRate?from=%s&to=%s", fromCurr, toCurr)

	// Get data from service
	err = s.GetJSON(url, &response)
	if err != nil {
		return
	}

	// Try to parse result
	res, err = strconv.ParseFloat(response.Rate, 64)
	return
}

// Creates new instance of service
func NewHerokuFeed(srvUrl string) *HerokuFeed {
	return &HerokuFeed{service.BaseService{
		ServiceUrl:  srvUrl,
		ServiceName: "Herokuapp.com",
	}}
}
