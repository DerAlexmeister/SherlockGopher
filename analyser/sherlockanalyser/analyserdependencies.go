package sherlockanalyser

import crawlerproto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyser"

/*
AnalyserDependency represents all services which are needed to run the analyser.
*/
type AnalyserDependency struct {
	Neo4J   func() //TODO
	Crawler func() crawlerproto.AnalyserInterface
}

/*
NewAnalyserDependencies will return an empty dependency struct.
*/
func NewAnalyserDependencies(crawler crawlerproto.AnalyserInterface) *AnalyserDependency {
	return &AnalyserDependency{
		Crawler: crawler
	}
}
