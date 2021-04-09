package sherlockneo

// TODO
// - Drop constrains

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

const (
	//ADDRESS of the Neo4j Dockercontainer.
	ADDRESS string = "bolt://10.0.2.15:7687" //"bolt://0.0.0.0:7687"

	//USER will be the user of the db.
	USER string = "neo4j" //Standard username change this in production.

	//PASSWORD will be the password of the neo4j db.
	PASSWORD string = "test" //Standard password change this in production.
)

const (
	//Magicconstant.

	//EMPTYSTRING will return an emptystring.
	EMPTYSTRING string = ""
)

/*
Will return a emptystring.
*/
func getEmptyString() string {
	return EMPTYSTRING
}

/*
GetGraphDBAddress will return the Address of neo4j.
*/
func GetGraphDBAddress() string {
	return ADDRESS
}

/*
getUserName will return the username of the current user.
*/
func getUserName() string {
	return USER
}

/*
getPasswort will return the password of the current user.
*/
func getPasswort() string {
	return PASSWORD
}

/*
GetNewDatabaseConnection will return a new connection/driver for the current neo4j Instance.
*/
func GetNewDatabaseConnection() (neo4j.Driver, error) {
	driver, err := neo4j.NewDriver(
		GetGraphDBAddress(),
		neo4j.BasicAuth(getUserName(),
			getPasswort(),
			getEmptyString()), func(c *neo4j.Config) { c.Encrypted = false })
	if err != nil {
		return nil, fmt.Errorf("an error occrued: %s", err)
	}
	return driver, nil
}

/*
CloseDatabaseConnection will close the drivers to the DB.
Eg. defer localneo.Close()
*/
func CloseDatabaseConnection(driver *neo4j.Driver) {
	(*driver).Close()
}

/*
GetSession will return a session to operate on inorder to send and receive data.
*/
func GetSession(driver neo4j.Driver) (neo4j.Session, error) {
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil, err
	}
	return session, nil
}

/*
CloseSession will close a session to the DB.
Eg. defer localneo.Close()
*/
func CloseSession(session *neo4j.Session) {
	(*session).Close()
}

/*
GetAllNodesAndTheirRelationships will return all nodes with address and the rels to other nodes.
*/
func GetAllNodesAndTheirRelationships(session neo4j.Session, args map[string]interface{}) ([]map[string]string, error) {
	results, err := session.Run(GetAllRels(), args)
	if err != nil {
		return make([]map[string]string, 0), err
	}
	var tojson []map[string]string
	for results.Next() {
		elem := make(map[string]string)
		for _, element := range results.Record().Keys() {
			if value, contains := results.Record().Get(element); contains {
				elem[element] = value.(string)
			}

		}
		tojson = append(tojson, elem)
	}
	defer CloseSession(&session)
	return tojson, nil
}

/*
getFileTypeColor.
*/
func getFileTypeColor(typ string) string {
	switch typ {
	case "Javascript":
		return "#F0B85B"
	case "CSS":
		return "#E891BC"
	case "Image":
		return "#85E196"
	default:
		return "#7CA9EF"
	}
}

/*
getRelationshipsTypeColor.
*/
//nolint: goconst
func getRelationshipsTypeColor(reltype string) string {
	switch reltype {
	case "Requires":
		return "#24A144"
	case "Shows":
		return "#99BA51"
	default:
		return "#08206A"

	}
}

/*
containsNode will check whether or not an node is already in the nodes list.
*/
func containsNode(node string, arr *[]string) bool {
	for _, element := range *arr {
		if element == node {
			return true
		}
	}
	return false
}

/*
GetAllNodesAndTheirRelationshipsOptimized will return all nodes with address and the rels to other nodes optimized for  the frontend.
*/
func GetAllNodesAndTheirRelationshipsOptimized(session neo4j.Session, args map[string]interface{}, query string) (map[string][]map[string]string, error) {
	results, err := session.Run(query, args)
	if err != nil {
		return make(map[string][]map[string]string), err
	}
	var existingNodes []string
	var rels, nodes []map[string]string
	for results.Next() {
		handleEntrie := func(entrie string) string {
			if value, contains := (results.Record().Get(entrie)); contains && value != nil {
				return value.(string)
			}
			return ""
		}
		if !containsNode(handleEntrie("Source"), &existingNodes) {
			nodes = append(nodes, map[string]string{
				"id":    handleEntrie("Source"),
				"color": getFileTypeColor(handleEntrie("SourceType")),
			})
			existingNodes = append(existingNodes, handleEntrie("Source"))
		}
		if !containsNode(handleEntrie("Destination"), &existingNodes) {
			nodes = append(nodes, map[string]string{
				"id":    handleEntrie("Destination"),
				"color": getFileTypeColor(handleEntrie("DestinationType")),
			})
			existingNodes = append(existingNodes, handleEntrie("Destination"))
		}
		rels = append(rels, map[string]string{
			"source": handleEntrie("Source"),
			"target": handleEntrie("Destination"),
			"label":  handleEntrie("Relationship"),
			"color":  getRelationshipsTypeColor(handleEntrie("Relationship")),
		})
	}
	InCaseIsEmpty := func(n []map[string]string) []map[string]string {
		if len(n) > 0 {
			return n
		}
		return make([]map[string]string, 0)
	}
	appendNodeAndRel := map[string][]map[string]string{
		"nodes": InCaseIsEmpty(nodes),
		"links": InCaseIsEmpty(rels),
	}
	defer CloseSession(&session)
	return appendNodeAndRel, nil
}

/*
GetAmountOfNodes will return the amount of nodes.
*/
func GetAmountOfNodes(session neo4j.Session, args map[string]interface{}) ([]map[string]int64, error) {
	amountofnodes, err := session.Run(getCountNumberOfNodes(), args)
	if err != nil {
		return make([]map[string]int64, 0), err
	}
	var tojson []map[string]int64
	for amountofnodes.Next() {
		elem := make(map[string]int64)
		for _, element := range amountofnodes.Record().Keys() {
			if value, contains := amountofnodes.Record().Get(element); contains {
				elem[element] = value.(int64)
			}
		}
		tojson = append(tojson, elem)
	}
	defer CloseSession(&session)
	return tojson, nil
}

/*
GetAmountOfRels will return the amout of relationship.
*/
func GetAmountOfRels(session neo4j.Session, args map[string]interface{}) ([]map[string]int64, error) {
	amountofrels, err := session.Run(getCountRelsToNodes(), args)
	if err != nil {
		return make([]map[string]int64, 0), err
	}
	var tojson []map[string]int64
	for amountofrels.Next() {
		elem := make(map[string]int64)
		for _, element := range amountofrels.Record().Keys() {
			if value, contains := amountofrels.Record().Get(element); contains {
				elem[element] = value.(int64)
			}
		}
		tojson = append(tojson, elem)
	}
	defer CloseSession(&session)
	return tojson, nil
}

/*
GetAmountofHTMLNodes will return the amount of html nodes.
*/
func GetAmountofHTMLNodes(session neo4j.Session, args map[string]interface{}) ([]map[string]int64, error) {
	amountofhtmls, err := session.Run(getCountHtmlsNodes(), args)
	if err != nil {
		return make([]map[string]int64, 0), err
	}
	var tojson []map[string]int64
	for amountofhtmls.Next() {
		elem := make(map[string]int64)
		for _, element := range amountofhtmls.Record().Keys() {
			if value, contains := amountofhtmls.Record().Get(element); contains {
				elem[element] = value.(int64)
			}
		}
		tojson = append(tojson, elem)
	}
	defer CloseSession(&session)
	return tojson, nil
}

/*
GetAmountOfStylesheets will return the amount of Stylesheets.
*/
func GetAmountOfStylesheets(session neo4j.Session, args map[string]interface{}) ([]map[string]int64, error) {
	amountofcss, err := session.Run(getCountCSSNodes(), args)
	if err != nil {
		return make([]map[string]int64, 0), err
	}
	var tojson []map[string]int64
	for amountofcss.Next() {
		elem := make(map[string]int64)
		for _, element := range amountofcss.Record().Keys() {
			if value, contains := amountofcss.Record().Get(element); contains {
				elem[element] = value.(int64)
			}
		}
		tojson = append(tojson, elem)
	}
	defer CloseSession(&session)
	return tojson, nil
}

/*
GetAmountOfJavascriptFiles will return the amount of Javascript files.
*/
func GetAmountOfJavascriptFiles(session neo4j.Session, args map[string]interface{}) ([]map[string]int64, error) {
	amountofjs, err := session.Run(getCountJavascriptNodes(), args)
	if err != nil {
		return make([]map[string]int64, 0), err
	}
	var tojson []map[string]int64
	for amountofjs.Next() {
		elem := make(map[string]int64)
		for _, element := range amountofjs.Record().Keys() {
			if value, contains := amountofjs.Record().Get(element); contains {
				elem[element] = value.(int64)
			}
		}
		tojson = append(tojson, elem)
	}
	defer CloseSession(&session)
	return tojson, nil
}

/*
GetAmountOfImages will return amount of images.
*/
func GetAmountOfImages(session neo4j.Session, args map[string]interface{}) ([]map[string]int64, error) {
	amountofimages, err := session.Run(getCountImageNodes(), args)
	if err != nil {
		return make([]map[string]int64, 0), err
	}
	var tojson []map[string]int64
	for amountofimages.Next() {
		elem := make(map[string]int64)
		for _, element := range amountofimages.Record().Keys() {
			if value, contains := amountofimages.Record().Get(element); contains {
				elem[element] = value.(int64)
			}
		}
		tojson = append(tojson, elem)
	}
	defer CloseSession(&session)
	return tojson, nil
}

/*
GetPerformanceOfSite will return the performance index of each site like address, rTT and statuscode.
*/
func GetPerformanceOfSite(session neo4j.Session, args map[string]interface{}) ([]map[string]string, error) {
	performance, err := session.Run(getResponseTimeInTableAndStatusCode(), args)
	if err != nil {
		return make([]map[string]string, 0), err
	}
	var tojson []map[string]string
	for performance.Next() {
		elem := make(map[string]string)
		for _, element := range performance.Record().Keys() {
			if value, contains := performance.Record().Get(element); contains {
				switch value := value.(type) {
				case int64:
					elem[element] = fmt.Sprintf("%d", value)
				case string:
					elem[element] = value
				case int:
					elem[element] = fmt.Sprintf("%d", value)
				}
			}
		}
		tojson = append(tojson, elem)
	}
	defer CloseSession(&session)
	return tojson, nil
}

/*
GetDetailsOfNode will return information of a node given as parameter.
*/
func GetDetailsOfNode(session neo4j.Session, target string) (map[string]map[string]interface{}, error) {
	details, err := session.Run(fmt.Sprintf(getReturnNode(), target), nil)
	if err != nil {
		return make(map[string]map[string]interface{}), err
	}
	tojson := make(map[string]map[string]interface{})
	for details.Next() {
		for _, element := range details.Record().Keys() {
			if value, contains := details.Record().Get(element); contains {
				tojson[(value.(map[string]interface{}))["Address"].(string)] = value.(map[string]interface{})
			}
		}
	}
	defer CloseSession(&session)
	return tojson, nil
}

/*
DropTable will drop the entire database. BE CAREFUL!
*/
func DropTable(session neo4j.Session) (map[string]interface{}, error) {
	result, err := session.Run(getDropGraph(), nil)
	response := make(map[string]interface{})
	if err != nil {
		response["Message"] = err
		return response, err
	}
	for result.Next() {
		for _, element := range result.Record().Keys() {
			if value, contains := result.Record().Get(element); contains {
				response[element] = value
			}
		}
	}
	defer CloseSession(&session)
	return response, nil
}

/*
ConvertNeoDataIntoMap will turn a NeoData instance into an  map[string]interface{}.
*/
func ConvertNeoDataIntoMap(instance NeoData) map[string]interface{} { //TODO error case.
	param := make(map[string]interface{})
	param["Address"] = instance.crawledLink.Address()
	param["FileType"] = instance.crawledLink.FileTypeAsString()
	if !instance.HasError() {
		param["Statuscode"] = instance.StatusCode()
		param["Responsetime"] = instance.ResponseTime()
		param["Status"] = "verified"
		if instance.ResponseHeader() != nil {
			for key, value := range *instance.ResponseHeader() {
				param[key] = value
			}
		}
	} else {
		param["Statuscode"] = 0
		param["Status"] = "unverified"
	}
	return param
}

/*
ConvertNeoLinkIntoNode will turn a NeoLink into a map[string]interface{} //TODO
*/
func ConvertNeoLinkIntoNode(link *NeoLink) map[string]interface{} {
	param := make(map[string]interface{})
	if link != nil {
		param["Address"] = link.Address()
		param["FileType"] = link.FileType()
		return param
	}
	return nil
}

/*
getQueryByFiletype will return the query depending on the filetype.
*/
func getQueryByFiletype(typ FileType) string {
	switch typ {
	case CSS:
		return getAddStyleSheetNode()
	case Image:
		return getAddImageNode()
	case Javascript:
		return getAddJavascriptNode()
	default:
		return getAddNode()

	}
}

/*
stringifymap will turn a map of elements into a string.
*/
//nolint: gocritic, gosimple
func stringifymap(args map[string]interface{}) string {
	var stringifyedmap string = ""
	for key, value := range args {
		switch value.(type) {
		case int64:
			stringifyedmap += strings.ReplaceAll(key, "-", "") + ": " + fmt.Sprintf("%d", value) + ", "
		case uint64:
			stringifyedmap += strings.ReplaceAll(key, "-", "") + ": " + fmt.Sprintf("%d", value) + ", "
		case int:
			stringifyedmap += strings.ReplaceAll(key, "-", "") + ": " + fmt.Sprintf("%d", value) + ", "
		case string:
			stringifyedmap += strings.ReplaceAll(key, "-", "") + ": \"" + value.(string) + "\", "
		case []string:
			stringifyedmap += strings.ReplaceAll(key, "-", "") + ": \"" + strings.Trim(value.([]string)[0], "\"") + "\", " //TODO array mit mehr elementen ordentlich auspacken.
		case time.Duration:
			stringifyedmap += strings.ReplaceAll(key, "-", "") + ": " + fmt.Sprintf("%d", value.(time.Duration).Milliseconds()) + ", " //TODO time.Duration in zeit nicht int
		default:
			log.WithFields(log.Fields{
				"key":  key,
				"type": reflect.TypeOf(value),
			}).Info("Unknown key type")
		}
	}
	stringifyedmap = strings.TrimRight(stringifyedmap, ", ")
	return stringifyedmap
}

/*
CreateANode will create an node and put it in the neo4j db.
*/
func CreateANode(session neo4j.Session, args NeoData) error {
	_, err := session.Run(fmt.Sprintf(getQueryByFiletype(args.crawledLink.FileType()), stringifymap(ConvertNeoDataIntoMap(args))), nil)
	if err != nil {
		log.Info(err)
		return err
	}
	return nil
}

func getRelByFileType(typ FileType) string {
	switch typ {
	case CSS:
		return "Requires"
	case Image:
		return "Shows"
	case Javascript:
		return "Requires"
	default:
		return "Links"

	}
}

func getNameByFileType(typ FileType) string {
	switch typ {
	case CSS:
		return "StyleSheet"
	case Image:
		return "Image"
	case Javascript:
		return "Javascript"
	default:
		return "Website"

	}
}

/*
CreateRelationships will create a Link between two nodes.
*/
//nolint: unconvert
func CreateRelationships(driver neo4j.Driver, source NeoLink, target NeoLink) error {
	session, err := GetSession(driver)
	if err != nil {
		return err
	}

	relationships := func(styp FileType, ttyp FileType, source string, target string) (string, error) {
		containsSource := ContainsNode(session, source)
		containsTarget := ContainsNode(session, target)
		switch {
		case containsSource && containsTarget:
			return fmt.Sprintf(getConnector(), "", source, "", "", target, "", getRelByFileType(ttyp)), nil
		case containsSource && !containsTarget:
			return fmt.Sprintf(getConnector(), "", source, "", fmt.Sprintf(":%s", string(getNameByFileType(ttyp))), target, fmt.Sprintf(", FileType: \"%s\"", string(ttyp)), getRelByFileType(ttyp)), nil
		case !containsSource && containsTarget:
			return fmt.Sprintf(getConnector(), fmt.Sprintf(":%s", string(getNameByFileType(styp))), source, fmt.Sprintf(", FileType: \"%s\"", string(styp)), "", target, "", getRelByFileType(ttyp)), nil
		default:
			return fmt.Sprintf(getConnector(), fmt.Sprintf(":%s", string(getNameByFileType(styp))), source, fmt.Sprintf(", FileType: \"%s\"", string(styp)), fmt.Sprintf(":%s", string(ttyp)), target, fmt.Sprintf(", FileType: \"%s\"", string(ttyp)), getRelByFileType(ttyp)), nil
		}
	}
	query, err := relationships(source.FileType(), target.FileType(), source.Address(), target.Address())
	if err != nil {
		return err
	}
	result, err := session.Run(query, nil)
	if err != nil {
		log.Info(result.Err())
		CloseSession(&session)
		return err
	}
	CloseSession(&session)
	return nil
}

/*
Save will save a Neo4J Entrie.
*/
func (nData *NeoData) Save(pDriver neo4j.Driver) error {
	session, err := GetSession(pDriver)
	if err != nil {
		return err
	}

	if !ContainsNode(session, nData.CrawledLink().Address()) {
		err := CreateANode(session, *nData)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Save->CreateANode failed")
			CloseSession(&session)
			return err
		}
		log.WithFields(log.Fields{
			"link": nData.CrawledLink().Address(),
		}).Info("[-] Not contained")
	} else {
		err := UpdateProperties(session, ConvertNeoDataIntoMap(*nData))
		log.WithFields(log.Fields{
			"address": nData.CrawledLink().Address(),
		}).Info("Update")
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Save->UpdateProperties failed")
			CloseSession(&session)
			return err
		}
	}
	for _, entry := range nData.Relations() {
		if entry != nil && nData != nil && nData.CrawledLink() != nil {
			err := CreateRelationships(pDriver, *nData.CrawledLink(), *entry)
			if err != nil {
				log.Fatal("CreateRelationships failed")
			}
		} else {
			log.WithFields(log.Fields{
				"entry": entry,
				"link":  *nData.CrawledLink(),
			}).Info("Save")
		}
	}
	CloseSession(&session)
	return nil
}

/*
UpdateProperties will update the properties of a Neodata.
*/
func UpdateProperties(session neo4j.Session, args map[string]interface{}) error {
	var timeSleep time.Duration = 2
	if args["Address"] == "" {
		log.Info("Empty Address-Field")
	}
	varA := fmt.Sprintf(getUpdatePropsQuery(), args["Address"], stringifymap(args))
	time.Sleep(timeSleep)
	_, err := session.Run(varA, nil)
	if err != nil {
		log.Error("UpdateProperties failed")
		return err
	}
	defer CloseSession(&session)
	return nil
}

/*
RunConstrains will add the constrains into the neo4j db.Filetype
*/
func RunConstrains(session neo4j.Session) {
	var constrains []string = []string{
		getConstrains(),
		getJSConstrain(),
		getStylesheetConstrains(),
		getImageConstrains(),
	}
	for _, element := range constrains {
		if _, err := session.Run(element, nil); err != nil {
			log.Infof("A problem occurred while running the constrains-function: %s  ", err.Error())
		}
	}
	defer CloseSession(&session)
}

/*
ContainsNode will check whether or not a node is in the Database.
*/
func ContainsNode(session neo4j.Session, target string) bool {
	query := fmt.Sprintf(getContains(), target)
	result, err := session.Run(query, nil)
	if err != nil {
		log.Info(err)
		return false
	}
	for result.Next() {
		if value, containskey := result.Record().Get("contains"); containskey {
			return value.(bool)
		}
	}
	return false
}
