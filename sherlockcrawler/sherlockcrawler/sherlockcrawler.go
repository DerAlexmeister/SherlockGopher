package sherlockcrawler

/*
Sherlockcrawler will be the Crawlerservice.
*/
type Sherlockcrawler struct {
	Queue        CrawlerQueue //Queue with all tasks
	Dependencies *Sherlockdependencies
}

/*
InjectDependency will inject the dependencies for the a sherlockcrawler instance
into the actual instance.
*/
func (sherlock Sherlockcrawler) InjectDependency(deps *Sherlockdependencies) {
	sherlock.Dependencies = deps
}
