package sherlockneo

import (
	"net/http"
	"time"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

/*
NeoData will be an instance representing the data for neo4j.
*/
type NeoData struct {
	crawledLink    *NeoLink      //The crawled address
	statusCode     int           //The status code of the website request
	responseTime   time.Duration //The response time of the website request
	responseHeader *http.Header  //The header of the website request
	cralwerError   string        //This field is set incase of an error.
	relations      []*NeoLink    //The found links
}

/*
NeoDataInterface will be the Interface implemented by neodata
*/
type NeoDataInterface interface {
	Save(session neo4j.Session)
}

/*
NeoLink will be an instance representing the data for a link.
*/
type NeoLink struct {
	address  string   //The address
	fileType FileType //FileType
}

//FileType will represent the possible types of a node.
type FileType string

const (
	//HTML will be FileType HTML
	HTML FileType = "HTML"

	//Image will be FileType Image
	Image FileType = "Image"

	//CSS will be FileType Stylesheet
	CSS FileType = "CSS"

	//Javascript will be FileType Javascript
	Javascript FileType = "Javascript"
)

/*
NewNeoData will return an new instance of NeoData.
*/
func NewNeoData(crawledLink *NeoLink, statusCode int, responseTime time.Duration, responseHeader *http.Header, cralwerError string, relations []*NeoLink) *NeoData {
	return &NeoData{crawledLink: crawledLink, statusCode: statusCode, responseTime: responseTime, responseHeader: responseHeader, cralwerError: cralwerError, relations: relations}
}

/*
Relations will return an array of relationships.
*/
func (nData *NeoData) Relations() []*NeoLink {
	return nData.relations
}

/*
SetRelations will set the relationships.
*/
func (nData *NeoData) SetRelations(relations []*NeoLink) {
	nData.relations = relations
}

/*
ResponseHeader will return the response header.
*/
func (nData *NeoData) ResponseHeader() *http.Header {
	return nData.responseHeader
}

/*
SetResponseHeader will set the responseHeader.
*/
func (nData *NeoData) SetResponseHeader(responseHeader *http.Header) {
	nData.responseHeader = responseHeader
}

/*
ResponseTime will return the responseTime.
*/
func (nData *NeoData) ResponseTime() time.Duration {
	return nData.responseTime
}

/*
SetResponseTime will set the responseTime.
*/
func (nData *NeoData) SetResponseTime(responseTime time.Duration) {
	nData.responseTime = responseTime
}

/*
StatusCode will return the statusCode.
*/
func (nData *NeoData) StatusCode() int {
	return nData.statusCode
}

/*
SetStatusCode will set the statusCode.
*/
func (nData *NeoData) SetStatusCode(statusCode int) {
	nData.statusCode = statusCode
}

/*
CrawledLink will return the crawled neoLink.
*/
func (nData *NeoData) CrawledLink() *NeoLink {
	return nData.crawledLink
}

/*
SetCrawledLink will set the crawled neoLink.
*/
func (nData *NeoData) SetCrawledLink(crawledLink *NeoLink) {
	nData.crawledLink = crawledLink
}

/*
NewNeoLink will return a new neoLink.
*/
func NewNeoLink(address string, fileType FileType) *NeoLink {
	return &NeoLink{address: address, fileType: fileType}
}

/*
FileType will return the fileType.
*/
func (n *NeoLink) FileType() FileType {
	return n.fileType
}

/*
FileTypeAsString will return the fileType as string.
*/
func (n *NeoLink) FileTypeAsString() string {
	return string(n.fileType)
}

/*
SetFileType will set the fileType.
*/
func (n *NeoLink) SetFileType(fileType FileType) {
	n.fileType = fileType
}

/*
Address will return the address.
*/
func (n *NeoLink) Address() string {
	return n.address
}

/*
SetAddress will set the address.
*/
func (n *NeoLink) SetAddress(address string) {
	n.address = address
}

/*
HasError Will return true incase there is an error in this neodata.
*/
func (nData *NeoData) HasError() bool {
	return nData.cralwerError != ""
}

/*
SetError will set the error field.
*/
func (nData *NeoData) SetError(message string) {
	nData.cralwerError = message
}
