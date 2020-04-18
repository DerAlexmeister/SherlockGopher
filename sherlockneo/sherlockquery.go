package sherlockneo

const (
	//Cypherstatements.

	//DROPGRAPH the entire graph.
	dropgraph string = "MATCH (n) DETACH DELETE n"

	//contains will check whether or not a node is already in use.
	contains string = "MATCH (x) WHERE x.name = \"%s\" RETURN x"

	//addnode will add a node into the neo4j DB.
	addnode string = "UNWIND {props} as prop CREATE (a:Website {address:prop.address, statuscode:prop.statuscode, responsetime:prop.responsetime, Header:prop.header, status:prop.status});"

	//constrains will be the constrains for a DB type html.
	constrains string = "CREATE CONSTRAINT ON (c:Website) ASSERT c.address IS UNIQUE"

	//constrainsjs will be the constrains for a DB type javascript.
	contrainsjs string = "CREATE CONSTRAINT ON (j:Javascript) ASSERT j.address IS UNIQUE"

	//constrainscc will be the constrains for a DB type css.
	contrainscss string = "CREATE CONSTRAINT ON (s:StyleSheet) ASSERT s.address IS UNIQUE"

	//constrainsimg will be the constrains for a DB type image.
	contrainsimg string = "CREATE CONSTRAINT ON (i:Image) ASSERT i.address IS UNIQUE"

	//returnall will return all nodes and their relationsships in the db.
	returnall string = "MATCH (n) RETURN n"

	//connectbylink will connect two given nodes by a link relationship.
	connectbylink string = "MATCH (f:Website), (s:Website) WHERE f.address = \"%s\" AND s.address = \"%s\" MERGE (f)-[:]->(s);"

	//starterkidofnode Will return a subset of nodes connected directly to a given node.
	starterkidofnode string = "MATCH (a)-[:]->(b) WHERE a.address = \"%s\" RETURN a, b" // TODO Vlt. mit Limit.

	//countnumberofnodes will count the number of nodes.
	countnumberofnodes string = "MATCH (n) RETURN count(n) as count"

	//countnumberofrels will count the number of relationships.
	countnumberofrels string = "MATCH ()-[r]->() RETURN count(r) as count"

	//countnumberofstylesheets will count the number of CSS files in the db.
	countnumberofstylesheets string = "MATCH (n) WHERE n.Filetype = \"CSS\" RETURN count(n) as count"

	//countnumberofjavascript will count the number of javascripts in the db.
	countnumberofjavascript string = "MATCH (n) WHERE n.Filetype = \"Javascript\" RETURN count(n) as count"

	//countnumberofimages will count the number of Images in the db.
	countnumberofimages string = "MATCH (n) WHERE n.Filetype = \"Image\" RETURN count(n) as count"

	//countnumberofhtml will count the number of HTML sites in the db.
	countnumberofhtml string = "MATCH (n) WHERE n.Filetype = \"HTML\" RETURN count(n) as count"
)

/*
getdropgraphQuery will return the query for droping the entire graph.
*/
func getDropGraph() string {
	return dropgraph
}

/*
GetContains will return the query to check if a node is contained.
*/
func GetContains() string {
	return contains
}

/*
GetAddNode will return the query to add a node.
*/
func GetAddNode() string {
	return addnode
}

/*
GetConstrains will return the query to the main constrain.
*/
func GetConstrains() string {
	return constrains
}

/*
GetReturnAll will return the query to fetch all entitys in a DB.
*/
func GetReturnAll() string {
	return returnall
}

/*
GetConnectbyLink will return the query to link to nodes together
*/
func GetConnectbyLink() string {
	return connectbylink
}

/*
GetStarterKidOfNode will return the query for a small subset of nodes.
*/
func GetStarterKidOfNode() string {
	return starterkidofnode
}

/*
GetCountNumberOfNodes will return the query to get the number of nodes in the db.
*/
func GetCountNumberOfNodes() string {
	return countnumberofnodes
}

/*
GetCountRelsToNodes will return the query to get the number of relationships in the db.
*/
func GetCountRelsToNodes() string {
	return countnumberofrels
}

/*
GetCountCSSNodes will return the query to get the number of type stylesheets in the db.
*/
func GetCountCSSNodes() string {
	return countnumberofstylesheets
}

/*
GetCountJavascriptNodes will return the query to get the number of type javascripts in the db.
*/
func GetCountJavascriptNodes() string {
	return countnumberofjavascript
}

/*
GetCountImageNodes will return the query to get the number of nodes of type images in the db.
*/
func GetCountImageNodes() string {
	return countnumberofimages
}

/*
GetCountHtmlsNodes will return the query to get the number of nodes of type html in the db.
*/
func GetCountHtmlsNodes() string {
	return countnumberofhtml
}
