package sherlockneo

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"

	"github.com/DerAlexx/SherlockGopher/sherlockneo/mocks"
	"github.com/golang/mock/gomock"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type testError struct{}

func (err *testError) Error() string {
	return "expected test-error"
}

//nolint:dupl
func Test_GetAmountOfImages(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	key := "amountofimages"
	value := int64(23)

	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	mockSession.EXPECT().Run(fmt.Sprint(getCountImageNodes()), nil).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true).Times(1)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockRecord.EXPECT().Keys().Return([]string{key}).Times(1)
	mockRecord.EXPECT().Get(key).Return(value, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(2)
	mockSession.EXPECT().Close().Times(1)

	result, _ := GetAmountOfImages(mockSession, nil)

	if result[0][key] != value {
		t.Errorf("Expected value of key %s to be %d but was %d", key, value, result[0][key])
	}

}

func Test_GetDetailsOfNode(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	first := mockResult.EXPECT().Next().Return(true).Times(1)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockRecord.EXPECT().Keys().Return([]string{"Status"}).Times(1)
	mockRecord.EXPECT().Get("Status").Return(map[string]interface{}{"Address": "example.com/0", "Status": "verified"}, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(2)
	mockSession.EXPECT().Run(gomock.Any(), nil).Return(mockResult, nil).Times(1)
	mockSession.EXPECT().Close().Times(1)

	result, _ := GetDetailsOfNode(mockSession, "example.com/0")
	if result["example.com/0"]["Status"] != "verified" {
		t.Errorf("Expected value of key %s to be %s but was %s", "Address", "verified", result["verified"]["Address"])
	}
}

//nolint:dupl
func Test_GetAmountOfJavascriptFiles(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	key := "amountofjs"
	value := int64(212)

	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	mockSession.EXPECT().Run(fmt.Sprint(getCountJavascriptNodes()), nil).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true).Times(1)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockRecord.EXPECT().Keys().Return([]string{key}).Times(1)
	mockRecord.EXPECT().Get(key).Return(value, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(2)
	mockSession.EXPECT().Close().Times(1)

	result, _ := GetAmountOfJavascriptFiles(mockSession, nil)

	if result[0][key] != value {
		t.Errorf("Expected value of key %s to be %d but was %d", key, value, result[0][key])
	}
}

//nolint:dupl
func Test_GetAmountOfStylesheets(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	key := "amountofsheets"
	value := int64(3)

	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	mockSession.EXPECT().Run(fmt.Sprint(getCountCSSNodes()), nil).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true).Times(1)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockRecord.EXPECT().Keys().Return([]string{key}).Times(1)
	mockRecord.EXPECT().Get(key).Return(value, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(2)
	mockSession.EXPECT().Close().Times(1)

	result, _ := GetAmountOfStylesheets(mockSession, nil)

	if result[0][key] != value {
		t.Errorf("Expected value of key %s to be %d but was %d", key, value, result[0][key])
	}
}

//nolint:dupl
func Test_GetAmountofHTMLNodes(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	key := "amountofhtmls"
	value := int64(10)

	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	mockSession.EXPECT().Run(fmt.Sprint(getCountHtmlsNodes()), nil).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true).Times(1)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockRecord.EXPECT().Keys().Return([]string{key}).Times(1)
	mockRecord.EXPECT().Get(key).Return(value, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(2)
	mockSession.EXPECT().Close().Times(1)

	result, _ := GetAmountofHTMLNodes(mockSession, nil)

	if result[0][key] != value {
		t.Errorf("Expected value of key %s to be %d but was %d", key, value, result[0][key])
	}
}

//nolint:dupl
func Test_GetAmountOfRels(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	key := "amountofrels"
	value := int64(5)

	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	mockSession.EXPECT().Run(fmt.Sprint(getCountRelsToNodes()), nil).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true).Times(1)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockRecord.EXPECT().Keys().Return([]string{key}).Times(1)
	mockRecord.EXPECT().Get(key).Return(value, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(2)
	mockSession.EXPECT().Close().Times(1)

	result, _ := GetAmountOfRels(mockSession, nil)

	if result[0][key] != value {
		t.Errorf("Expected value of key %s to be %d but was %d", key, value, result[0][key])
	}
}

//nolint:dupl
func Test_GetAmountOfNodes(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	key := "amountofnodes"
	value := int64(199)

	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	mockSession.EXPECT().Run(fmt.Sprint(getCountNumberOfNodes()), nil).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true).Times(1)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockRecord.EXPECT().Keys().Return([]string{key}).Times(1)
	mockRecord.EXPECT().Get(key).Return(value, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(2)
	mockSession.EXPECT().Close().Times(1)

	result, _ := GetAmountOfNodes(mockSession, nil)

	if result[0][key] != value {
		t.Errorf("Expected value of key %s to be %d but was %d", key, value, result[0][key])
	}
}

func Test_GetAllNodesAndTheirRelationshipsOptimized(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	mockSession.EXPECT().Run(gomock.Any(), nil).Return(mockResult, nil).Times(1)

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

	result, err := GetAllNodesAndTheirRelationshipsOptimized(mockSession, nil, "MATCH (n)-[r]->(k) RETURN n.Address as Source, n.FileType as SourceType, Type(r) as Relationship, k.Address as Destination, k.FileType as DestinationType")

	if fmt.Sprint(result) != "map[links:[map[color:#24A144 label:Requires source:www.example.com/0 target:www.example.com/1]] nodes:[map[color:#7CA9EF id:www.example.com/0] map[color:#7CA9EF id:www.example.com/1]]]" || err != nil {
		t.Fail()
	}
}

func Test_GetAllNodesAndTheirRelationships(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	key := "n.Address"
	value := "41"

	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	mockSession.EXPECT().Run(fmt.Sprint(GetAllRels()), nil).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true).Times(1)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockRecord.EXPECT().Keys().Return([]string{key}).Times(1)
	mockRecord.EXPECT().Get(key).Return(value, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(2)
	mockSession.EXPECT().Close().Times(1)

	result, _ := GetAllNodesAndTheirRelationships(mockSession, nil)

	if result[0][key] != value {
		t.Errorf("Expected value of key %s to be %s but was %s", key, value, result[0][key])
	}
}

func Test_GetPerformanceOfSite(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	addrKey := "n.Address"
	respKey := "n.Responsetime"
	statusKey := "n.Statuscode"
	addrValue := int64(64)
	respValue := "333"
	statusValue := "293"

	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	mockSession.EXPECT().Run(fmt.Sprint(getResponseTimeInTableAndStatusCode()), nil).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockRecord.EXPECT().Keys().Return([]string{addrKey, respKey, statusKey})
	mockRecord.EXPECT().Get(addrKey).Return(addrValue, true)
	mockRecord.EXPECT().Get(respKey).Return(respValue, true)
	mockRecord.EXPECT().Get(statusKey).Return(statusValue, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(4)
	mockSession.EXPECT().Close().Times(1)

	result, _ := GetPerformanceOfSite(mockSession, nil)

	results := result[0]
	if results[addrKey] != fmt.Sprintf("%d", addrValue) || results[respKey] != respValue || results[statusKey] != statusValue {
		t.Fail()
	}
}

func Test_DropTable(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	key := "amountofnodes"
	value := int64(199)

	mockRecord := mocks.NewMockRecord(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	mockSession.EXPECT().Run(fmt.Sprint(getDropGraph()), nil).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true).Times(1)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockRecord.EXPECT().Keys().Return([]string{key}).Times(1)
	mockRecord.EXPECT().Get(key).Return(value, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(2)
	mockSession.EXPECT().Close().Times(1)

	result, _ := DropTable(mockSession)

	if result[key] != value {
		t.Errorf("Expected value of key %s to be %d but was %d", key, value, result[key])
	}
}

func Test_DropTableError(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockSession := mocks.NewMockSession(mockCtrl)

	mockSession.EXPECT().Run(fmt.Sprint(getDropGraph()), nil).Return(nil, &testError{}).Times(1)

	_, err := DropTable(mockSession)

	if err.Error() != "expected test-error" {
		t.Fail()
	}
}

func Test_getFileTypeColor(t *testing.T) {
	if getFileTypeColor("Javascript") != "#F0B85B" ||
		getFileTypeColor("CSS") != "#E891BC" ||
		getFileTypeColor("Image") != "#85E196" ||
		getFileTypeColor("Test") != "#7CA9EF" {
		t.Fail()
	}
}

func Test_getRelationshipsTypeColor(t *testing.T) {
	if getRelationshipsTypeColor("Requires") != "#24A144" ||
		getRelationshipsTypeColor("Shows") != "#99BA51" ||
		getRelationshipsTypeColor("Test") != "#08206A" {
		t.Fail()
	}
}

func Test_containsNode(t *testing.T) {
}

/*func Test_ConvertNeoDataIntoMap(t *testing.T) {
	result := ConvertNeoDataIntoMap(*NewNeoData("example.com/0", 200, 300, nil, []string{"requires", "links"}, HTML), true)

	if result["Address"] != "example.com/0" ||
		result["Statuscode"] != 200 ||
		result["Responsetime"] != time.Duration(300) ||
		result["Filetype"] != "HTML" ||
		fmt.Sprintf("%s", result["Status"]) != "verified" {
		t.Fail()
	}
}

func Test_ConvertNeoDataIntoMapUnverified(t *testing.T) {
	result := ConvertNeoDataIntoMap(*NewNeoData("example.com/0", 200, 300, nil, []string{"requires", "links"}, HTML), false)

	if result["Address"] != "example.com/0" ||
		result["Statuscode"] != 200 ||
		result["Responsetime"] != time.Duration(300) ||
		result["Filetype"] != "HTML" ||
		fmt.Sprintf("%s", result["Status"]) != "unverified" {
		t.Fail()
	}
}*/

func Test_GetQueryByFiletype(t *testing.T) {
	if getQueryByFiletype(HTML) != getAddNode() ||
		getQueryByFiletype(Image) != getAddImageNode() ||
		getQueryByFiletype(CSS) != getAddStyleSheetNode() ||
		getQueryByFiletype(Javascript) != getAddJavascriptNode() {
		t.Fail()
	}
}

//func Test_CreateAUnverifiedNode(t *testing.T){
//	mockCtrl := gomock.NewController(t)
//
//	defer mockCtrl.Finish()
//
//	args := *NewNeoData("example.com/0", 200, 300, nil, []string{"requires", "links"}, HTML)
//
//	mockSession := mocks.NewMockSession(mockCtrl)
//	mockTransaction := mocks.NewMockTransaction(mockCtrl)
//
//	mockSession.EXPECT().WriteTransaction(gomock.Any()).Return(nil, nil).Times(1)
//	mockSession.EXPECT().Close().Times(1)
//
//	mockTransaction.EXPECT().Run(getQueryByFiletype(args.CleanFiletype()),
//		ConvertNeoDataIntoMap(args, false))
//	CreateAUnverifiedNode(mockSession, NeoData{})
//}

func Test_GetSession(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil)

	result, _ := GetSession(mockDriver)
	if mockSession != result.(neo4j.Session) {
		t.Fail()
	}
}

func Test_GetSessionErr(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, big.ErrNaN{})

	_, result := GetSession(mockDriver)
	if result == nil {
		t.Fail()
	}
}

func Test_ContainsNode(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	input := ""

	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)
	mockRecord := mocks.NewMockRecord(mockCtrl)

	mockSession.EXPECT().Run(gomock.Any(), nil).Return(mockResult, nil).Times(1)

	mockResult.EXPECT().Next().Return(true)

	mockResult.EXPECT().Record().Return(mockRecord).Times(1)
	mockRecord.EXPECT().Get(gomock.Any()).Return(true, true)

	result := ContainsNode(mockSession, input)
	if !result {
		t.Fail()
	}
}

func Test_ContainsNodeFalse(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	input := ""

	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)
	mockRecord := mocks.NewMockRecord(mockCtrl)

	mockSession.EXPECT().Run(gomock.Any(), nil).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true)
	second := mockResult.EXPECT().Next().Return(false)

	gomock.InOrder(first, second)

	mockResult.EXPECT().Record().Return(mockRecord).Times(1)
	mockRecord.EXPECT().Get(gomock.Any()).Return(false, false)

	result := ContainsNode(mockSession, input)
	if result {
		t.Fail()
	}
}

func Test_ConvertNeoDataIntoMap(t *testing.T) {
	data := NeoData{}
	data.SetCrawledLink(NewNeoLink("example.org/1", HTML))
	data.SetRelations([]*NeoLink{NewNeoLink("example.org/2", HTML), NewNeoLink("example.org/0", HTML)})
	data.SetResponseHeader(&http.Header{})
	data.SetResponseTime(123)
	data.SetStatusCode(200)

	resultMap := ConvertNeoDataIntoMap(data)
	assert.Equal(t, fmt.Sprint(resultMap), "map[Address:example.org/1 FileType:HTML Responsetime:123ns Status:verified Statuscode:200]")

}

func Test_ConvertNeoDataIntoMapNoError(t *testing.T) {
	data := NeoData{}
	data.SetCrawledLink(NewNeoLink("example.org/1", HTML))
	data.SetRelations([]*NeoLink{NewNeoLink("example.org/2", HTML), NewNeoLink("example.org/0", HTML)})
	data.SetResponseHeader(&http.Header{})
	data.SetError("Error")
	data.SetResponseTime(123)
	data.SetStatusCode(200)

	resultMap := ConvertNeoDataIntoMap(data)
	assert.Equal(t, fmt.Sprint(resultMap), "map[Address:example.org/1 FileType:HTML Status:unverified Statuscode:0]")

}

func Test_ConvertNeoLinkIntoMode(t *testing.T) {
	link := NewNeoLink("example.org/1", HTML)

	resultMap := ConvertNeoLinkIntoNode(link)
	if fmt.Sprint(resultMap) != "map[Address:example.org/1 FileType:HTML]" {
		t.Fail()
	}
}

func Test_stringifymap(t *testing.T) {
	inputMap := make(map[string]interface{})
	inputMap["string"] = "test"
	inputMap["int64"] = int64(1)
	inputMap["uint64"] = uint64(1)
	inputMap["int"] = int(1)
	inputMap["[]string"] = []string{"test"}
	dur, _ := time.ParseDuration("123ms")
	inputMap["time.Duration"] = dur
	inputMap["default"] = true
	stringifymap(inputMap)
	//assert.Equal(t, stringifymap(inputMap), "int: 1, []string: \"test\", time.Duration: 123, string: \"test\", int64: 1, uint64: 1" )
}

func Test_CreateANode(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	data := NeoData{}
	data.SetCrawledLink(NewNeoLink("example.org/1", HTML))
	data.SetError("This is an error")
	data.SetRelations([]*NeoLink{NewNeoLink("example.org/2", HTML), NewNeoLink("example.org/0", HTML)})
	data.SetResponseHeader(&http.Header{})
	data.SetResponseTime(123)
	data.SetStatusCode(200)

	defer mockCtrl.Finish()
	mockSession := mocks.NewMockSession(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession.EXPECT().Run(gomock.Any(), nil).Return(mockResult, nil).Times(1)
	err := CreateANode(mockSession, data)
	if err != nil {
		log.Fatal("CreateANode failed")
	}
}

func Test_CreateANodeError(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	data := NeoData{}
	data.SetCrawledLink(NewNeoLink("example.org/1", HTML))
	data.SetError("This is an error")
	data.SetRelations([]*NeoLink{NewNeoLink("example.org/2", HTML), NewNeoLink("example.org/0", HTML)})
	data.SetResponseHeader(&http.Header{})
	data.SetResponseTime(123)
	data.SetStatusCode(200)

	defer mockCtrl.Finish()
	mockSession := mocks.NewMockSession(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession.EXPECT().Run(gomock.Any(), nil).Return(mockResult, &testError{}).Times(1)

	err := CreateANode(mockSession, data)
	assert.Equal(t, err.Error(), (&testError{}).Error())
}

func Test_getRelByFileType(t *testing.T) {
	assert.Equal(t, getRelByFileType(HTML), "Links")
	assert.Equal(t, getRelByFileType(CSS), "Requires")
	assert.Equal(t, getRelByFileType(Image), "Shows")
	assert.Equal(t, getRelByFileType(Javascript), "Requires")
}

func Test_getNameByFileType(t *testing.T) {
	assert.Equal(t, getNameByFileType(HTML), "Website")
	assert.Equal(t, getNameByFileType(CSS), "StyleSheet")
	assert.Equal(t, getNameByFileType(Image), "Image")
	assert.Equal(t, getNameByFileType(Javascript), "Javascript")
}

func Test_RunConstrains(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()
	mockSession := mocks.NewMockSession(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession.EXPECT().Run(gomock.Any(), nil).Return(mockResult, &testError{}).Times(4)
	mockSession.EXPECT().Close()
	RunConstrains(mockSession)
}

func Test_UpdateProperties(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	input := make(map[string]interface{})
	input["Address"] = ""

	defer mockCtrl.Finish()
	mockSession := mocks.NewMockSession(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockSession.EXPECT().Run(gomock.Any(), nil).Return(mockResult, nil).Times(1)
	mockSession.EXPECT().Close()
	err := UpdateProperties(mockSession, input)
	if err != nil {
		log.Fatal("UpdateProperties failed")
	}
}

func Test_CreateRelationshipsNodeNotContained(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	nl1 := NewNeoLink("example.org/1", HTML)
	nl2 := NewNeoLink("example.org/2", HTML)

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil)
	mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, &testError{}).MinTimes(1)
	mockSession.EXPECT().Close()
	mockResult.EXPECT().Err().Return(&testError{})

	ret := CreateRelationships(mockDriver, *nl1, *nl2)

	assert.Equal(t, ret.Error(), (&testError{}).Error())
}

func Test_CreateRelationshipsBothContained(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	nl1 := NewNeoLink("example.org/1", HTML)
	nl2 := NewNeoLink("example.org/2", HTML)

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockRecord := mocks.NewMockRecord(mockCtrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil)
	mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil).MinTimes(1)

	mockResult.EXPECT().Next().Return(true).MinTimes(2)

	mockResult.EXPECT().Record().Return(mockRecord).MinTimes(1)
	mockRecord.EXPECT().Get(gomock.Any()).Return(true, true).MinTimes(1)
	mockSession.EXPECT().Close()

	ret := CreateRelationships(mockDriver, *nl1, *nl2)

	assert.Equal(t, ret, nil)
}

//nolint: dupl
func Test_CreateRelationshipsOneContained(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	nl1 := NewNeoLink("example.org/1", HTML)
	nl2 := NewNeoLink("example.org/2", HTML)

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockRecord := mocks.NewMockRecord(mockCtrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil)
	mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil).MinTimes(1)

	mockResult.EXPECT().Next().Return(true).MinTimes(2)

	mockResult.EXPECT().Record().Return(mockRecord).MinTimes(1)
	first := mockRecord.EXPECT().Get(gomock.Any()).Return(true, true)
	second := mockRecord.EXPECT().Get(gomock.Any()).Return(false, true)

	gomock.InOrder(first, second)
	mockSession.EXPECT().Close()

	ret := CreateRelationships(mockDriver, *nl1, *nl2)

	assert.Equal(t, ret, nil)
}

//nolint: dupl
func Test_CreateRelationshipsOtherContained(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	nl1 := NewNeoLink("example.org/1", HTML)
	nl2 := NewNeoLink("example.org/2", HTML)

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockRecord := mocks.NewMockRecord(mockCtrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil)
	mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil).MinTimes(1)

	mockResult.EXPECT().Next().Return(true).MinTimes(2)

	mockResult.EXPECT().Record().Return(mockRecord).MinTimes(1)
	first := mockRecord.EXPECT().Get(gomock.Any()).Return(false, true)
	second := mockRecord.EXPECT().Get(gomock.Any()).Return(true, true)

	gomock.InOrder(first, second)
	mockSession.EXPECT().Close()

	ret := CreateRelationships(mockDriver, *nl1, *nl2)

	assert.Equal(t, ret, nil)
}

func Test_CreateRelationshipsDBError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	nl1 := NewNeoLink("example.org/1", HTML)
	nl2 := NewNeoLink("example.org/2", HTML)

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockDriver.EXPECT().Session(gomock.Any()).Return(nil, &testError{})

	ret := CreateRelationships(mockDriver, *nl1, *nl2)

	assert.Equal(t, ret.Error(), (&testError{}).Error())
}

func Test_SaveErrorCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	nl := NeoLink{}
	nl.SetAddress("example.org")
	nl.SetFileType(HTML)
	sut := NewNeoData(&nl, 200, 100, nil, "", nil)

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil).MinTimes(1)
	mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, &testError{}).MinTimes(1)
	mockSession.EXPECT().Close().MinTimes(1)

	err := sut.Save(mockDriver)

	assert.Equal(t, err.Error(), (&testError{}).Error())

}

func Test_SaveErrorContains(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sut := NeoData{
		crawledLink:    NewNeoLink("example.org", HTML),
		statusCode:     200,
		responseTime:   100,
		responseHeader: nil,
		cralwerError:   "",
		relations:      nil,
	}

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil).MinTimes(1)
	firstRun := mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, &testError{})
	secondRun := mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil)
	mockSession.EXPECT().Close().MinTimes(1)

	gomock.InOrder(firstRun, secondRun)

	err := sut.Save(mockDriver)

	assert.Equal(t, err, (nil))

}

func Test_Save(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sut := NeoData{
		crawledLink:    NewNeoLink("example.org", HTML),
		statusCode:     200,
		responseTime:   100,
		responseHeader: nil,
		cralwerError:   "",
		relations:      nil,
	}

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockRecord := mocks.NewMockRecord(mockCtrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil).MinTimes(1)
	first := mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil)
	second := mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, &testError{})
	mockResult.EXPECT().Next().Return(true)
	mockSession.EXPECT().Close().MinTimes(1)

	gomock.InOrder(first, second)

	mockResult.EXPECT().Record().Return(mockRecord).MinTimes(1)
	mockRecord.EXPECT().Get(gomock.Any()).Return(true, true).MinTimes(1)

	err := sut.Save(mockDriver)

	assert.Equal(t, err.Error(), (&testError{}).Error())

}

func Test_SaveWithRelations(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sut := NeoData{
		crawledLink:    NewNeoLink("example.org", HTML),
		statusCode:     200,
		responseTime:   100,
		responseHeader: nil,
		cralwerError:   "",
		relations:      []*NeoLink{NewNeoLink("www.example.org", HTML), NewNeoLink("test.com", HTML)},
	}

	mockDriver := mocks.NewMockDriver(mockCtrl)
	mockSession := mocks.NewMockSession(mockCtrl)
	mockResult := mocks.NewMockResult(mockCtrl)
	mockRecord := mocks.NewMockRecord(mockCtrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil).MinTimes(1)
	mockResult.EXPECT().Next().Return(true).MinTimes(1)
	mockSession.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, nil).MinTimes(1)
	mockSession.EXPECT().Close().MinTimes(1)

	mockResult.EXPECT().Record().Return(mockRecord).MinTimes(1)
	mockRecord.EXPECT().Get(gomock.Any()).Return(true, true).MinTimes(1)
	err := sut.Save(mockDriver)

	assert.Equal(t, err, nil)

}

func Test_GetNewDatabaseConnection(t *testing.T) {
	driver, err := GetNewDatabaseConnection()

	assert.NotEqual(t, driver, nil)
	assert.Equal(t, err, nil)

}

func TestQuery(t *testing.T) {
	if getReturnAll() != "MATCH (n) RETURN properties(n)" {
		t.Fatal("Error in Query returnall")
	}
	if getConnectbyLink() != "MERGE (a:Website {Address:\"%s\"}) MERGE (b:Website {Address:\"%s\", Status: \"unverified\", FileType:\"%s\"}) MERGE (a)-[r:Links]->(b)" {
		t.Fatal("Error in Query connectbylink")
	}

	if getVerify() != "MERGE (n {Address: \"%s\"}) SET n.Status = \"verified\" RETURN n" {
		t.Fatal("Error in Query verified")
	}

	if getUpdateProperties() != "MERGE (n {Address: \"%s\" SET c += {props} RETURN n" {
		t.Fatal("Error in Query setproperties")
	}

	if getCSSconnection() != "MERGE (a:Website {Address:\"%s\"}) MERGE (b:StyleSheet {Address:\"%s\", Status: \"unverified\", FileType:\"%s\"}) MERGE (a)-[r:Requires]->(b)" {
		t.Fatal("Error in Query connectbyRequireCSS")
	}

	if getJavascriptConnection() != "MERGE (a:Website {Address:\"%s\"}) MERGE (b:Javascript {Address:\"%s\", Status: \"unverified\", FileType:\"%s\"}) MERGE (a)-[r:Requires]->(b)" {
		t.Fatal("Error in Query connectbyRequireJS")
	}

	if getShowsConnection() != "MERGE (a:Website {Address:\"%s\"}) MERGE (b:Image {Address:\"%s\", Status: \"unverified\", FileType:\"%s\"}) MERGE (a)-[r:Shows]->(b)" {
		t.Fatal("Error in Query connectbyShows")
	}

}
