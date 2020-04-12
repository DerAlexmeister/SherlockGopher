package sherlockanalyser

/*
AnalyserDependency represents all services which are needed to run the analyser.
*/
type AnalyserDependency struct {
	Neo4J   func()
	Crawler func()
}

/*
NewAnalyserDependencies will return an empty dependency struct.
*/
func NewAnalyserDependencies() *AnalyserDependency {
	return &AnalyserDependency{}
}
