package webserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/DerAlexx/SherlockGopher/sherlockneo/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//nolint: gochecknoinits
func init() {
	gin.SetMode(gin.ReleaseMode)
}

func compareMaps(wantJSON string, have map[string]string) {

}

func TestGraphFetchWholeGraphHighPerformanceV1(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sut := New()
	c.Request = httptest.NewRequest("GET", "/alloptimized", nil)

	//Mocking
	ctrl := gomock.NewController(t)
	mockDriver := mocks.NewMockDriver(ctrl)
	mockRecord := mocks.NewMockRecord(ctrl)
	mockResult := mocks.NewMockResult(ctrl)
	mockSession := mocks.NewMockSession(ctrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil).Times(1)
	mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true).Times(1)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	gomock.InOrder(
		mockRecord.EXPECT().Get("Source").Return("www.example.com/0", true),
		mockRecord.EXPECT().Get("Source").Return("www.example.com/0", true),
		mockRecord.EXPECT().Get("SourceType").Return("#08206A", true),
		mockRecord.EXPECT().Get("Source").Return("www.example.com/0", true),
		mockRecord.EXPECT().Get("Destination").Return("www.example.com/1", true),
		mockRecord.EXPECT().Get("Destination").Return("www.example.com/1", true),
		mockRecord.EXPECT().Get("DestinationType").Return("#08206A", true),
		mockRecord.EXPECT().Get("Destination").Return("www.example.com/1", true),
		mockRecord.EXPECT().Get("Source").Return("www.example.com/0", true),
		mockRecord.EXPECT().Get("Destination").Return("www.example.com/1", true),
		mockRecord.EXPECT().Get("Relationship").Return("Requires", true),
		mockRecord.EXPECT().Get("Relationship").Return("Requires", true))

	mockResult.EXPECT().Record().Return(mockRecord).MinTimes(1)
	mockSession.EXPECT().Close().Times(1)

	sut.Driver = mockDriver
	sut.GraphFetchWholeGraphHighPerformanceV1(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "map[links:[map[color:#24A144 label:Requires source:www.example.com/0 target:www.example.com/1]] nodes:[map[color:#7CA9EF id:www.example.com/0] map[color:#7CA9EF id:www.example.com/1]]]", fmt.Sprint(got))
}

func TestGraphFetchWholeGraphHighPerformanceV1BadSession(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sut := New()
	c.Request = httptest.NewRequest("GET", "/alloptimized", nil)

	//Mocking
	ctrl := gomock.NewController(t)
	mockDriver := mocks.NewMockDriver(ctrl)
	mockSession := mocks.NewMockSession(ctrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, &WebserverTestError{}).Times(1)

	sut.Driver = mockDriver
	sut.GraphFetchWholeGraphHighPerformanceV1(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "500 Internal Server Error", w.Result().Status)
	assert.Equal(t, "A Problem occurred while trying to connect to the Database", got["Message"])
}

func TestGraphNodeDetailsV1(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sut := New()
	c.Request = httptest.NewRequest("POST", "/detailsofnode", bytes.NewReader([]byte("{\n        \"url\": \"www.example.com\"\n    } ")))

	//Mocking
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	first := mockResult.EXPECT().Next().Return(true).Times(1)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil)
	mockRecord.EXPECT().Keys().Return([]string{"Status"}).Times(1)
	mockRecord.EXPECT().Get("Status").Return(map[string]interface{}{"Address": "example.com/", "Status": "verified"}, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(2)
	mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil).Times(1)
	mockSession.EXPECT().Close().Times(2)

	sut.Driver = mockDriver
	sut.GraphNodeDetailsV1(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "map[example.com/:map[Address:example.com/ Status:verified]]", fmt.Sprint(got))
}

func TestGraphNodeDetailsV1BadSession(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sut := New()
	c.Request = httptest.NewRequest("POST", "/detailsofnode", bytes.NewReader([]byte("{\n        \"url\":\"example.com\"    } ")))

	//Mocking
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, &WebserverTestError{})

	sut.Driver = mockDriver
	sut.GraphNodeDetailsV1(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "500 Internal Server Error", w.Result().Status)
	assert.Equal(t, "A Problem occurred while trying to connect to the Database", got["Message"])
}

func TestGraphNodeDetailsV1MalformedURL(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sut := New()
	c.Request = httptest.NewRequest("POST", "/detailsofnode", bytes.NewReader([]byte("{\n        \"url\":\"examp\"    } ")))

	//Mocking
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil)

	sut.Driver = mockDriver
	sut.GraphNodeDetailsV1(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "400 Bad Request", w.Result().Status)
	assert.Equal(t, "Malformed URL", got["Message"])
}

func TestGraphFetchWholeGraphV1(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sut := New()
	c.Request = httptest.NewRequest("GET", "/graph/v1/all", nil)

	//Mocking
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	key := "n.Address"
	value := "41"

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)
	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil)
	mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true).Times(1)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockRecord.EXPECT().Keys().Return([]string{key}).Times(1)
	mockRecord.EXPECT().Get(key).Return(value, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(2)
	mockSession.EXPECT().Close().Times(2)

	sut.Driver = mockDriver
	sut.GraphFetchWholeGraphV1(c)

	var got []interface{}
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "[map[n.Address:41]]", fmt.Sprint(got))
}

func TestGraphPerformanceOfSitesV1BadSession(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sut := New()
	c.Request = httptest.NewRequest("GET", "/graph/v1/all", nil)

	//Mocking
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockDriver.EXPECT().Session(gomock.Any()).Return(nil, &WebserverTestError{})

	sut.Driver = mockDriver

	sut.GraphPerformanceOfSitesV1(c)
	var got gin.H

	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "500 Internal Server Error", w.Result().Status)
	assert.Equal(t, "A Problem occurred while trying to connect to the Database", fmt.Sprint(got["Message"]))
	fmt.Println(got)
}

func TestGraphPerformanceOfSitesV1(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sut := New()
	c.Request = httptest.NewRequest("GET", "/graph/v1/all", nil)

	//Mocking
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)
	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)

	addrKey := "n.Address"
	respKey := "n.Responsetime"
	statusKey := "n.Statuscode"
	addrValue := int64(64)
	respValue := "333"
	statusValue := "293"

	mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockRecord.EXPECT().Keys().Return([]string{addrKey, respKey, statusKey})
	mockRecord.EXPECT().Get(addrKey).Return(addrValue, true)
	mockRecord.EXPECT().Get(respKey).Return(respValue, true)
	mockRecord.EXPECT().Get(statusKey).Return(statusValue, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(4)
	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil)
	mockSession.EXPECT().Close().Times(2)
	sut.Driver = mockDriver

	sut.GraphPerformanceOfSitesV1(c)
	var got []interface{}

	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "[map[n.Address:64 n.Responsetime:333 n.Statuscode:293]]", fmt.Sprint(got))
	fmt.Println(got)
}

func TestGraphMetaV1(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sut := New()
	c.Request = httptest.NewRequest("GET", "/graph/v1/all", nil)

	//Mocking
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)
	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil)

	sut.Driver = mockDriver

	//MockImages
	keyImages := "amountofimages"
	valueImages := int64(23)

	imageRun := mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil).Times(1)
	imageNextFirst := mockResult.EXPECT().Next().Return(true).Times(1)
	imageNextSecond := mockResult.EXPECT().Next().Return(false).Times(1)
	imageKeys := mockRecord.EXPECT().Keys().Return([]string{keyImages}).Times(1)
	imageGetKey := mockRecord.EXPECT().Get(keyImages).Return(valueImages, true)
	imageRecordFirst := mockResult.EXPECT().Record().Return(mockRecord)
	imageRecordSecond := mockResult.EXPECT().Record().Return(mockRecord)
	imageClose := mockSession.EXPECT().Close().Times(1)
	gomock.InOrder(imageRun, imageNextFirst, imageRecordFirst, imageKeys, imageRecordSecond, imageGetKey, imageNextSecond, imageClose)

	//MockCSS

	keyCSS := "amountofsheets"
	valueCSS := int64(3)

	cssRun := mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil).Times(1)
	cssNextFirst := mockResult.EXPECT().Next().Return(true).Times(1)
	cssNextSecond := mockResult.EXPECT().Next().Return(false).Times(1)
	cssKeys := mockRecord.EXPECT().Keys().Return([]string{keyCSS}).Times(1)
	cssGetKey := mockRecord.EXPECT().Get(keyCSS).Return(valueCSS, true)
	cssRecordFirst := mockResult.EXPECT().Record().Return(mockRecord)
	cssRecordSecond := mockResult.EXPECT().Record().Return(mockRecord)
	cssClose := mockSession.EXPECT().Close().Times(1)

	gomock.InOrder(cssRun, cssNextFirst, cssRecordFirst, cssKeys, cssRecordSecond, cssGetKey, cssNextSecond, cssClose)

	//MockJavaScript

	keyJS := "amountofjs"
	valueJS := int64(212)

	jsRun := mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil).Times(1)
	jsNextFirst := mockResult.EXPECT().Next().Return(true).Times(1)
	jsNextSecond := mockResult.EXPECT().Next().Return(false).Times(1)
	jsKeys := mockRecord.EXPECT().Keys().Return([]string{keyJS}).Times(1)
	jsGetKey := mockRecord.EXPECT().Get(keyJS).Return(valueJS, true)
	jsRecordFirst := mockResult.EXPECT().Record().Return(mockRecord)
	jsRecordSecond := mockResult.EXPECT().Record().Return(mockRecord)
	jsClose := mockSession.EXPECT().Close().Times(1)

	gomock.InOrder(jsRun, jsNextFirst, jsRecordFirst, jsKeys, jsRecordSecond, jsGetKey, jsNextSecond, jsClose)

	//MockHTML

	keyHTML := "amountofhtmls"
	valueHTML := int64(10)

	htmlRun := mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil).Times(1)
	htmlNextFirst := mockResult.EXPECT().Next().Return(true).Times(1)
	htmlNextSecond := mockResult.EXPECT().Next().Return(false).Times(1)
	htmlKeys := mockRecord.EXPECT().Keys().Return([]string{keyHTML}).Times(1)
	htmlGetKey := mockRecord.EXPECT().Get(keyHTML).Return(valueHTML, true)
	htmlRecordFirst := mockResult.EXPECT().Record().Return(mockRecord)
	htmlRecordSecond := mockResult.EXPECT().Record().Return(mockRecord)
	htmlClose := mockSession.EXPECT().Close().Times(1)

	gomock.InOrder(htmlRun, htmlNextFirst, htmlRecordFirst, htmlKeys, htmlRecordSecond, htmlGetKey, htmlNextSecond, htmlClose)

	//MockRels

	keyRels := "amountofrels"
	valueRels := int64(5)

	relsRun := mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil).Times(1)
	relsNextFirst := mockResult.EXPECT().Next().Return(true).Times(1)
	relsNextSecond := mockResult.EXPECT().Next().Return(false).Times(1)
	relsKeys := mockRecord.EXPECT().Keys().Return([]string{keyRels}).Times(1)
	relsGetKey := mockRecord.EXPECT().Get(keyRels).Return(valueRels, true)
	relsRecordFirst := mockResult.EXPECT().Record().Return(mockRecord)
	relsRecordSecond := mockResult.EXPECT().Record().Return(mockRecord)
	relsClose := mockSession.EXPECT().Close().Times(1)

	gomock.InOrder(relsRun, relsNextFirst, relsRecordFirst, relsKeys, relsRecordSecond, relsGetKey, relsNextSecond, relsClose)

	//MockNodes

	keyNodes := "amountofnodes"
	valueNodes := int64(199)

	nodesRun := mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil).Times(1)
	nodesNextFirst := mockResult.EXPECT().Next().Return(true).Times(1)
	nodesNextSecond := mockResult.EXPECT().Next().Return(false).Times(1)
	nodesKeys := mockRecord.EXPECT().Keys().Return([]string{keyNodes}).Times(1)
	nodesGetKey := mockRecord.EXPECT().Get(keyNodes).Return(valueNodes, true)
	nodesRecordFirst := mockResult.EXPECT().Record().Return(mockRecord)
	nodesRecordSecond := mockResult.EXPECT().Record().Return(mockRecord)
	nodesClose := mockSession.EXPECT().Close().Times(1)

	gomock.InOrder(nodesRun, nodesNextFirst, nodesRecordFirst, nodesKeys, nodesRecordSecond, nodesGetKey, nodesNextSecond, nodesClose)

	mockSession.EXPECT().Close().Times(1)

	sut.GraphMetaV1(c)
	var got []interface{}

	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	//check if string representation of map can change! Tests might fail unexpectedly...
	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "[map[amountofnodes:199] map[amountofrels:5] map[amountofhtmls:10] map[amountofsheets:3] map[amountofjs:212] map[amountofimages:23]]", fmt.Sprint(got))
	fmt.Println(got)
}
