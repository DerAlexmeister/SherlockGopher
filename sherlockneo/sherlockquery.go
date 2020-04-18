package sherlockneo

const (
	//Cypherstatements.

	//DROPGRAPH the entire graph.
	DROPGRAPH string = "MATCH (n) DETACH DELETE n"

	//CONTAINS will check whether or not a node is already in use.
	CONTAINS string = "MATCH (x) WHERE x.name = \"%s\" RETURN x"

	//ADDNODE will add a node into the neo4j DB.
	ADDNODE string = "UNWIND {props} as prop CREATE (a:Website {address:prop.address, statuscode:prop.statuscode, responsetime:prop.responsetime, Header:prop.header, status:prop.status});"

	//CONSTRAINS will be the constrains for a DB.
	CONSTRAINS string = "CREATE CONSTRAINT ON (c:Website) ASSERT c.address IS UNIQUE"

	//RETURNALL will return all nodes and their relationsships in the db.
	RETURNALL string = "MATCH (n) RETURN n"

	//CONNECTBYLINK will connect two given nodes by a link relationship.
	CONNECTBYLINK string = "MATCH (f:Website), (s:Website) WHERE f.address = \"%s\" AND s.address = \"%s\" MERGE (f)-[:]->(s);"

	//STARTERKIDOFNODE Will return a subset of nodes connected directly to a given node.
	STARTERKIDOFNODE string = "MATCH (a)-[:]->(b) WHERE a.address = \"%s\" RETURN a, b" // TODO Vlt. mit Limit.

	//COUNTNUMBEROFNODES will count the number of nodes.
	COUNTNUMBEROFNODES string = "MATCH (n) RETURN count(n) as count"

	//COUNTNUMBEROFRELS will count the number of relationships.
	COUNTNUMBEROFRELS string = "MATCH ()-[r]->() RETURN count(r) as count"

	//COUNTNUMBEROFSTYLESHEETS will count the number of CSS files in the db.
	COUNTNUMBEROFSTYLESHEETS string = "MATCH (n) WHERE n.Filetype = \"CSS\" RETURN count(n) as count"

	//COUNTNUMBEROFJAVASCRIPT will count the number of javascripts in the db.
	COUNTNUMBEROFJAVASCRIPT string = "MATCH (n) WHERE n.Filetype = \"Javascript\" RETURN count(n) as count"

	//COUNTNUMBEROFIAMGES will count the number of Images in the db.
	COUNTNUMBEROFIAMGES string = "MATCH (n) WHERE n.Filetype = \"Image\" RETURN count(n) as count"

	//COUNTNUMBEROFHTML will count the number of HTML sites in the db.
	COUNTNUMBEROFHTML string = "MATCH (n) WHERE n.Filetype = \"HTML\" RETURN count(n) as count"
)

//TODO GETTER

/*
getStarterKidOfNode will return the query for a small subset of nodes.
*/
func getStarterKidOfNode() string {
	return STARTERKIDOFNODE
}

/*
getDropGraphQuery will return the query for droping the entire graph.
*/
func getDropGraph() string {
	return DROPGRAPH
}
