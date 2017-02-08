package currfeed

import (
	"encoding/xml"
	"fmt"

	"github.com/zak01011996/bitcoin-ticker/service"
)

// mycurrency.net API response struct
type ecbResp struct {
	XMLName  xml.Name `xml:"Envelope"`
	CubeRoot struct {
		CubeDate struct {
			RateList []struct {
				CurrCode string  `xml:"currency,attr"`
				Rate     float64 `xml:"rate,attr"`
			} `xml:"Cube"`
		} `xml:"Cube"`
	} `xml:"Cube"`
}

// Currency feed service
// Service: ecb.europa.eu
// Current API url: http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml
type ECBFeed struct {
	service.BaseService
}

// Get rate in given currencies
func (s *ECBFeed) GetRate(fromCurr string, toCurr string) (res float64, err error) {
	// Get data from service
	var response ecbResp

	// Compile url
	url := s.GetURL() + "stats/eurofxref/eurofxref-daily.xml"

	// Try to get data
	err = s.GetXML(url, &response)
	if err != nil {
		return
	}

	// Get rate list
	data := response.CubeRoot.CubeDate.RateList

	// Check response
	if len(data) <= 0 {
		err = fmt.Errorf(service.EmptyResp, s.GetName())
		return
	}

	var fromRate, toRate float64

	// Find and calculate rates
	for _, v := range data {
		if v.CurrCode == fromCurr {
			fromRate = v.Rate
		}

		if v.CurrCode == toCurr {
			toRate = v.Rate
		}
	}

	// EUR here equals 1
	if fromCurr == EUR_CODE {
		fromRate = 1
	}

	// Check if zero, who knows what could happen here, but it shouldn't crash our process
	if toRate == 0 {
		err = fmt.Errorf(service.ZeroValue, s.GetName())
	}

	res = toRate / fromRate
	return
}

// Creates new instance of service
func NewECBFeed(srvUrl string) *ECBFeed {
	return &ECBFeed{service.BaseService{
		ServiceUrl:  srvUrl,
		ServiceName: "Ecb.europa.eu",
	}}
}
