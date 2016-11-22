package architecture

type PriorityQueue interface {
	Init()
	// queue item
	Enqueue(item PriorityQueueItem)
	// remove item from begining
	Dequeue() (item PriorityQueueItem)
	// find an item by id in the queue
	Find(id string) (item PriorityQueueItem)
}

type PriorityQueueItem interface {
	Key() int64
	Id() string
}

type Tube struct {
	Name     string
	Queue    PriorityQueue
	Reserved PriorityQueue
	Delayed  PriorityQueue
	Buried   []PriorityQueueItem
}
