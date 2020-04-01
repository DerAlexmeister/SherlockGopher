package sherlockcrawler

/*
Sherlockdependencies is an type to manage all dependencies of sherlockcrawler.
*/
type Sherlockdependencies struct {
	Webserver func() //Webserver for gRPC
	Analyser  func() //Analyser for gRPC
}

/*
NewSherlockDependencies will return a new sherlockdependencies instance to put it in the dependencies
in a sherlockcrawler object.
*/
func NewSherlockDependencies() *Sherlockdependencies {
	return &Sherlockdependencies{}
}

//TODO getter und setter fehlen noch.
