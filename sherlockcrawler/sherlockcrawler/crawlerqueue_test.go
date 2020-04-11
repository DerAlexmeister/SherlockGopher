package sherlockcrawler

import "testing"

func TestQueueHasIDAlreadyInUse(t *testing.T) {
	queue := NewCrawlerQueue()
	task := &CrawlerTaskRequest{}
	id := queue.AppendQueue(task)
	if id == 0 {
		t.Fatal("got a 0 id from the AppendQueue-Method")
	} else if contains := queue.ContainsTaskID(id); !contains {
		t.Fatalf("id is not in the queue altought it should be. ID: %d, Contains: %t", id, contains)
	} else {
		t.Log("successfully found ID in taskqueue")
	}

}

func TestQueueContainsTaskIDFailes(t *testing.T) {
	queue := NewCrawlerQueue()
	task := &CrawlerTaskRequest{}
	id := queue.AppendQueue(task)
	if id == 0 {
		t.Fatal("got a 0 id from the AppendQueue-Method")
	} else if contains := queue.ContainsTaskID(id + 1); contains {
		t.Fatalf("id is not in the queue altought it should be. ID: %d, Contains: %t", id, contains)
	} else {
		t.Log("successfully found ID in taskqueue")
	}

}

func TestQueueAppendFailesBecauseOfNilPointer(t *testing.T) {
	queue := NewCrawlerQueue()
	id := queue.AppendQueue(nil)
	if id == 0 {
		t.Log("got a zero returncode as expected", 0, id)
	} else {
		t.Fatalf("got a non-zero returncode. Expected: %d, Got: %d", 0, id)
	}

}

func TestQueueRemoveTask(t *testing.T) {
	queue := NewCrawlerQueue()
	task := &CrawlerTaskRequest{}
	id := queue.AppendQueue(task)
	if id == 0 {
		t.Fatal("go a 0 id from the AppendQueue-Method")
	} else if contains := queue.ContainsTaskID(id); !contains {
		t.Fatalf("id is not in the queue altought it should be. ID: %d, Contains: %t", id, contains)
	} else {
		t.Log("successfully added task to the taskqueue")
	}
	if works := queue.RemoveFromQueue(id); !works {
		t.Fatalf("couldnt remove task with valid id Expected: %t, Got: %t, ID: %d", true, works, id)
	} else {
		t.Log("successfully removed task")
	}
}

func TestQueueRemoveTaskFailedBecauseOfZeroID(t *testing.T) {
	queue := NewCrawlerQueue()
	task := &CrawlerTaskRequest{}
	id := queue.AppendQueue(task)
	if id == 0 {
		t.Fatal("go a 0 id from the AppendQueue-Method")
	} else if contains := queue.ContainsTaskID(id); !contains {
		t.Fatalf("id is not in the queue altought it should be. ID: %d, Contains: %t", id, contains)
	} else {
		t.Log("successfully added task to the taskqueue")
	}
	if works := queue.RemoveFromQueue(0); works {
		t.Fatalf("could remove task with unvalid id Expected: %t, Got: %t, ID: %d", false, works, id)
	} else {
		t.Log("successfully removed task")
	}
}
