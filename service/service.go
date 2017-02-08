package service

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

// Error templates
var (
	BadStatus = "Bad response from service (%s), status code is: %d"
	EmptyResp = "Empty response (%s)"
	ZeroValue = "Zero rate, WTF?! (%s)"
)

// Base service declaration
type BaseService struct {
	ServiceUrl  string
	ServiceName string
}

// Name of service
func (s *BaseService) GetName() string {
	return s.ServiceName
}

// Get service url
func (s *BaseService) GetURL() string {
	return s.ServiceUrl
}

// Service JSON call and try to unmarshal response
func (s *BaseService) GetJSON(url string, res interface{}) (err error) {

	// Make request to service
	response, err := http.Get(url)
	if err != nil {
		return
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = fmt.Errorf(BadStatus, s.GetName(), response.StatusCode)
		return
	}

	// Initialize decoder
	decoder := json.NewDecoder(response.Body)

	// Try to parse response
	err = decoder.Decode(res)
	return
}

// Service XML call and try to unmarshal response
func (s *BaseService) GetXML(url string, res interface{}) (err error) {

	// Make request to service
	response, err := http.Get(url)
	if err != nil {
		return
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = fmt.Errorf(BadStatus, s.GetName(), response.StatusCode)
		return
	}

	// Initialize decoder
	decoder := xml.NewDecoder(response.Body)

	// Try to parse response
	err = decoder.Decode(res)
	return
}
