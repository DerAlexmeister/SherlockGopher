package sherlockanalyser

/*
AnalyserDependency represents all services which are needed to run the analyser.
*/
type AnalyserDependency struct {
	Neo4J   func() //TODO
	Crawler func() //crawlerproto.AnalyserInterface
}

/*
NewAnalyserDependencies will return an empty dependency struct.
*/
func NewAnalyserDependencies() *AnalyserDependency {
	return &AnalyserDependency{
		//Crawler: crawler
	}
}
