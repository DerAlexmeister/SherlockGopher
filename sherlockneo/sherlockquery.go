package sherlockneo

const (
	//Cypherstatements.

	//DROPGRAPH the entire graph.
	dropgraph string = "MATCH (n) DETACH DELETE n"

	//contains will check whether or not a node is already in use.
	contains string = "MATCH (x) WHERE x.Address = \"%s\" RETURN CASE WHEN count(x) = 0 THEN false ELSE true END as contains"

	//constrains will be the constrains for a DB type html.
	constrains string = "CREATE CONSTRAINT ON (c:Website) ASSERT c.Address IS UNIQUE"

	//constrainsjs will be the constrains for a DB type javascript.
	contrainsjs string = "CREATE CONSTRAINT ON (j:Javascript) ASSERT j.Address IS UNIQUE"

	//constrainscc will be the constrains for a DB type css.
	contrainscss string = "CREATE CONSTRAINT ON (s:StyleSheet) ASSERT s.Address IS UNIQUE"

	//constrainsimg will be the constrains for a DB type image.
	contrainsimg string = "CREATE CONSTRAINT ON (i:Image) ASSERT i.Address IS UNIQUE"

	//addnode will add a node into the neo4j DB.
	addnode string = "MERGE (c:Website {%s})"

	//addnodecss will add a node of the type StyleSheet into the db.
	addnodecss string = "MERGE (c:StyleSheet {%s})"

	//addnodejs will add a node of the type javascript into the db.
	addnodejs string = "MERGE (c:Javascript {%s})"

	//addnodeimg will add a node of the type image into the db.
	addnodeimg string = "MERGE (c:Image {%s})"

	//returnall will return all nodes and their relationsships in the db.
	returnall string = "MATCH (n) RETURN properties(n)"

	//returnnode will return all information of a node.
	returnnode string = "MATCH (n) WHERE n.Address=\"%s\" RETURN properties(n)"

	//returnall will return all nodes and their relationsships in the db.
	returnallrels string = "MATCH (n)-[r]->(k) RETURN n.Address as Source, n.FileType as SourceType, Type(r) as Relationship, k.Address as Destination, k.FileType as DestinationType"

	//connectbylink will connect two given nodes by a link relationship.
	connectbylink string = "MERGE (a:Website {Address:\"%s\"}) MERGE (b:Website {Address:\"%s\", Status: \"unverified\", FileType:\"%s\"}) MERGE (a)-[r:Links]->(b)"

	//connectbyRequireCSS will connect a website with a StyleSheet instance by the relationship requires.
	connectbyRequireCSS string = "MERGE (a:Website {Address:\"%s\"}) MERGE (b:StyleSheet {Address:\"%s\", Status: \"unverified\", FileType:\"%s\"}) MERGE (a)-[r:Requires]->(b)"

	//connectbyRequireJS will connect a website with a javascript instance by the relationship requires.
	connectbyRequireJS string = "MERGE (a:Website {Address:\"%s\"}) MERGE (b:Javascript {Address:\"%s\", Status: \"unverified\", FileType:\"%s\"}) MERGE (a)-[r:Requires]->(b)"

	//connectbyRequireShows will connect a website with a image instance by the relationship shows.
	connectbyShows string = "MERGE (a:Website {Address:\"%s\"}) MERGE (b:Image {Address:\"%s\", Status: \"unverified\", FileType:\"%s\"}) MERGE (a)-[r:Shows]->(b)"

	//countnumberofnodes will count the number of nodes.
	countnumberofnodes string = "MATCH (n) RETURN count(n) as amountofnodes"

	//countnumberofrels will count the number of relationships.
	countnumberofrels string = "MATCH ()-[r]->() RETURN count(r) as amountofrels"

	//countnumberofStyleSheets will count the number of CSS files in the db.
	countnumberofstylesheets string = "MATCH (n) WHERE n.FileType = \"CSS\" RETURN count(n) as amountofsheets"

	//countnumberofjavascript will count the number of javascripts in the db.
	countnumberofjavascript string = "MATCH (n) WHERE n.FileType = \"Javascript\" RETURN count(n) as amountofjs"

	//countnumberofimages will count the number of Images in the db.
	countnumberofimages string = "MATCH (n) WHERE n.FileType = \"Image\" RETURN count(n) as amountofimages"

	//countnumberofhtml will count the number of HTML sites in the db.
	countnumberofhtml string = "MATCH (n) WHERE n.FileType = \"HTML\" RETURN count(n) as amountofhtmls"

	//responseTimeInTableAndStatusCode will be for each website the responsetime and the code so easy to put in a table.
	responseTimeInTableAndStatusCode string = "MATCH (n) RETURN n.Address as Address, n.Responsetime as Responsetime, n.Statuscode as Statuscode"

	//verified will be a query to verify a node.
	verified string = "MERGE (n {Address: \"%s\"}) SET n.Status = \"verified\" RETURN n"

	//query to update the properties of a node.
	setproperties string = "MERGE (n {Address: \"%s\" SET c += {props} RETURN n"

	//updateprops will update the properties of a node.
	updateprops string = "MATCH (x) Where x.Address=\"%s\" SET x += {%s}"

	//connector will connect two nodes
	connector string = "MERGE (a%s {Address:\"%s\" %s}) MERGE (b%s {Address:\"%s\" %s}) MERGE (a)-[r:%s]->(b)"
)

/*
getdropgraphQuery will return the query for droping the entire graph.
*/
func getDropGraph() string {
	return dropgraph
}

/*
getContains will return the query to check if a node is contained.
*/
func getContains() string {
	return contains
}

/*
GetAddNode will return the query to add a node.
*/
func getAddNode() string {
	return addnode
}

/*
getAddStyleSheetNode will return the query to add a StyleSheet node.
*/
func getAddStyleSheetNode() string {
	return addnodecss
}

/*
getAddJavascriptNode will return the query to add a javascript node.
*/
func getAddJavascriptNode() string {
	return addnodejs
}

/*
getAddImageNode will return the query to add a StyleSheet node.
*/
func getAddImageNode() string {
	return addnodeimg
}

/*
GetConstrains will return the query to the main constrain.
*/
func getConstrains() string {
	return constrains
}

/*
getImageConstrains will return the constrains for image files.
*/
func getImageConstrains() string {
	return contrainsimg
}

/*
getStylesheetConstrains will return constrains for StyleSheets files.
*/
func getStylesheetConstrains() string {
	return contrainscss
}

/*
getJSConstrain will return the constrains for javascript files.
*/
func getJSConstrain() string {
	return contrainsjs
}

/*
GetReturnAll will return the query to fetch all entitys in a DB.
*/
func getReturnAll() string {
	return returnall
}

/*
GetConnectbyLink will return the query to link to nodes together.
*/
func getConnectbyLink() string {
	return connectbylink
}

/*
GetCountNumberOfNodes will return the query to get the number of nodes in the db.
*/
func getCountNumberOfNodes() string {
	return countnumberofnodes
}

/*
GetCountRelsToNodes will return the query to get the number of relationships in the db.
*/
func getCountRelsToNodes() string {
	return countnumberofrels
}

/*
GetCountCSSNodes will return the query to get the number of type StyleSheets in the db.
*/
func getCountCSSNodes() string {
	return countnumberofstylesheets
}

/*
GetCountJavascriptNodes will return the query to get the number of type javascripts in the db.
*/
func getCountJavascriptNodes() string {
	return countnumberofjavascript
}

/*
GetCountImageNodes will return the query to get the number of nodes of type images in the db.
*/
func getCountImageNodes() string {
	return countnumberofimages
}

/*
GetCountHtmlsNodes will return the query to get the number of nodes of type html in the db.
*/
func getCountHtmlsNodes() string {
	return countnumberofhtml
}

/*
GetResponseTimeInTableAndStatusCode will return the query to get for each address the statuscode and the RTT.
*/
func getResponseTimeInTableAndStatusCode() string {
	return responseTimeInTableAndStatusCode
}

/*
GetAllRels will return the query to get all relationships.
*/
func getAllRels() string {
	return returnallrels
}

/*
getReturnNode will return the query to get information of a node.
*/
func getReturnNode() string {
	return returnnode
}

/*
getVerify will return the query to verfiy a node no matter which type.
*/
func getVerify() string {
	return verified
}

/*
getUpdateProperties will return the query to update the properties of any kind of node.
*/
func getUpdateProperties() string {
	return setproperties
}

/*
getCSSconnection will return the css-connection query.
*/
func getCSSconnection() string {
	return connectbyRequireCSS
}

/*
getJavascriptConnection will return the query for the javascript-connection.
*/
func getJavascriptConnection() string {
	return connectbyRequireJS
}

/*
getShowsConnection will return the query for the image-connection.
*/
func getShowsConnection() string {
	return connectbyShows
}

/*
getUpdatePropsQuery will return the query for an update.
*/
func getUpdatePropsQuery() string {
	return updateprops
}

/*
getConnector will return the query to connect any kind of nodes.
*/
func getConnector() string {
	return connector
}
