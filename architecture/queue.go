package architecture

import "log"

type PriorityQueue interface {
	Init()
	// queue item
	Enqueue(item PriorityQueueItem)
	// get the highest priority item without removing
	Peek() (item PriorityQueueItem)
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
	Name            string
	Ready           PriorityQueue
	Reserved        PriorityQueue
	Delayed         PriorityQueue
	Buried          PriorityQueue
	AwaitingClients PriorityQueue
}
// Process runs all the necessary operations for upkeep of the tube
// TODO unit test
func (tube *Tube) Process() {
	delayedJob := tube.Delayed.Peek()
	if delayedJob != nil && delayedJob.Key() <= 0 {
		log.Println("delayed job got ready: ", delayedJob)
		delayedJob = tube.Delayed.Dequeue()
		tube.Ready.Enqueue(delayedJob)
	}
	// reserved jobs are put to ready
	reservedJob := tube.Reserved.Peek()
	if reservedJob != nil && reservedJob.Key() <= 0 {
		reservedJob = tube.Reserved.Dequeue()
		tube.Ready.Enqueue(reservedJob)
	}
	// ready jobs are sent
	availableClientConnection := tube.AwaitingClients.Dequeue()
	if (availableClientConnection != nil) {
		readyJob := tube.Ready.Dequeue().(*Job)
		availableClientConnection.(*AwaitingClient).SendChannel <- *readyJob
		tube.Reserved.Enqueue(readyJob)
	}
}
