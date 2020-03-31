package sherlockcrawler

/*
CrawlerQueue will be the queue of the current CrawlerTaskRequest.
*/
type CrawlerQueue struct {
	Queue map[string]*CrawlerTaskRequest
}

/*
getCurrentQueue will return a pointer to the current Queue.
*/
func (que *CrawlerQueue) getCurrentQueue() *(map[string]*CrawlerTaskRequest) {
	return &que.Queue
}

/*
ContainsAddress will check whether or not a addr is allready in use or not.
*/
func (que *CrawlerQueue) ContainsAddress(addr string) bool {
	if _, contains := (*que.getCurrentQueue())[addr]; !contains {
		return false
	}
	return true
}

/*
AppendQueue will append the current queue with a new CrawlerTaskRequest.
*/
func (que *CrawlerQueue) AppendQueue(target string, task *CrawlerTaskRequest) error {
	if !que.ContainsAddress(target) {
		(*que.getCurrentQueue())[target] = task
		return nil //TODO Returntype noch mal überarbeiten weil Error vlt nicht das beste.
	}
	return nil //TODO Returntype noch mal überarbeiten weil Error vlt nicht das beste.
}

/*
RemoveFromQueue will remove a task from the queue by a given address.
*/
func (que *CrawlerQueue) RemoveFromQueue(target string) bool {
	if !que.ContainsAddress(target) {
		delete((*que.getCurrentQueue()), target)
		return true
	}
	return false
}
