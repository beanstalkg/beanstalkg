package server

type PriorityQueue interface {
	Init()
	// queue item
	Enqueue(item PriorityQueueItem)
	// remove item from begining
	Dequeue() (item PriorityQueueItem)
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

type Beanstalkg struct {
	tubes map[string]Tube
	tubeCom map[string]chan string
}

func (b *Beanstalkg) Init() {
	// initialize the go routines to handle Heap ops and the comm channel
}

func (b *Beanstalkg) ExecCommand(c Command) string {
	return "USING " + c.Params["tube"]
}