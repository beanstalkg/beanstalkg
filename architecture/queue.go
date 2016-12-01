package architecture

import (
	"log"
)

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
	// delete an item and return it by id
	Delete(id string) PriorityQueueItem
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
		log.Println("QUEUE delayed job got ready: ", delayedJob)
		delayedJob = tube.Delayed.Dequeue()
		delayedJob.(*Job).SetState(READY)
		tube.Ready.Enqueue(delayedJob)
	}
	// reserved jobs are put to ready
	reservedJob := tube.Reserved.Peek()
	if reservedJob != nil && reservedJob.Key() <= 0 {
		log.Println("QUEUE found reserved job thats ready: ", reservedJob)
		reservedJob = tube.Reserved.Dequeue()
		reservedJob.(*Job).SetState(READY)
		tube.Ready.Enqueue(reservedJob)
	}
	// ready jobs are sent
	if tube.AwaitingClients.Peek() != nil && tube.Ready.Peek() != nil {
		//log.Println("*********************************************************************")
		availableClientConnection := tube.AwaitingClients.Dequeue()
		readyJob := tube.Ready.Dequeue().(*Job)
		client := availableClientConnection.(*AwaitingClient)
		log.Println("QUEUE sending job to client: ", client.id)
		client.Request.Job = *readyJob
		client.SendChannel <- client.Request
		readyJob.SetState(RESERVED)
		tube.Reserved.Enqueue(readyJob)
	}
}
