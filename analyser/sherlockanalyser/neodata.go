package sherlockanalyser

import (
	"net/http"
	"time"
)

type neoData struct {
	address string //The crawled address
	statusCode int //The status code of the website request
	responseTime time.Duration //The response time of the website request
	responseHeader *http.Header //The header of the website request
	relations []string //The found links
}

func NewNeoData(address string, statusCode int, responseTime time.Duration, responseHeader *http.Header, relations []string) *neoData {
	return &neoData{address: address, statusCode: statusCode, responseTime: responseTime, responseHeader: responseHeader, relations: relations}
}

func (nData *neoData) Relations() []string {
	return nData.relations
}

func (nData *neoData) SetRelations(relations []string) {
	nData.relations = relations
}

func (nData *neoData) ResponseHeader() *http.Header {
	return nData.responseHeader
}

func (nData *neoData) SetResponseHeader(responseHeader *http.Header) {
	nData.responseHeader = responseHeader
}

func (nData *neoData) ResponseTime() time.Duration {
	return nData.responseTime
}

func (nData *neoData) SetResponseTime(responseTime time.Duration) {
	nData.responseTime = responseTime
}

func (nData *neoData) StatusCode() int {
	return nData.statusCode
}

func (nData *neoData) SetStatusCode(statusCode int) {
	nData.statusCode = statusCode
}

func (nData *neoData) Address() string {
	return nData.address
}

func (nData *neoData) SetAddress(address string) {
	nData.address = address
}

